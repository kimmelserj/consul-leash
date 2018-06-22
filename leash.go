package leash

import (
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
)

type Worker struct {
	sync.Mutex
	DoneChan chan error

	args            []string
	command         string
	path            string
	stoppingTimeout time.Duration
	value           string

	cmd      *exec.Cmd
	started  bool
	stopping bool
}

func New(command string, args []string, path string, value string, stoppingTimeout time.Duration) *Worker {
	return &Worker{
		DoneChan: make(chan error),

		args:            args,
		command:         command,
		path:            path,
		stoppingTimeout: stoppingTimeout,
		value:           value,
	}
}

func (w *Worker) Run() error {
	wp, err := watch.Parse(map[string]interface{}{
		"type": "key",
		"key":  w.path,
	})
	if err != nil {
		return err
	}

	wp.Handler = func(index uint64, data interface{}) {
		w.Lock()
		defer w.Unlock()

		if data == nil {
			return
		}

		kvPair := data.(*api.KVPair)
		if !w.stopping {
			if !w.started && string(kvPair.Value) == w.value {
				go w.runCommand()
				w.started = true
			} else if w.started && string(kvPair.Value) != w.value {
				w.StopCommand()
			}
		}
	}

	go wp.Run("")

	return nil
}

func (w *Worker) StopCommand() error {
	defer close(w.DoneChan)

	w.stopping = true

	if !w.started {
		return nil
	}

	go func() {
		time.Sleep(w.stoppingTimeout)

		w.cmd.Process.Signal(syscall.SIGKILL)
		os.Exit(1)
	}()

	return nil
}

func (w *Worker) runCommand() {
	w.cmd = exec.Command(w.command, w.args...)
	w.cmd.Stdout = os.Stdout
	w.cmd.Stderr = os.Stderr

	err := w.cmd.Start()
	if err != nil {
		w.DoneChan <- err
	}

	go func() {
		sc := make(chan os.Signal)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM) //@todo need all signals
		for s := range sc {
			w.cmd.Process.Signal(s)
		}
	}()

	err = w.cmd.Wait()
	if err != nil {
		w.DoneChan <- err
	}

	close(w.DoneChan)
}
