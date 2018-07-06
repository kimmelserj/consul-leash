package leash

import (
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
)

func TestWorker_Run(t *testing.T) {
	t.Run("when command ends with zero exit code", func(t *testing.T) {
		cc, err := api.NewClient(api.DefaultConfig())
		if err != nil {
			t.Error(err)
			return
		}

		_, err = cc.KV().Put(&api.KVPair{
			Key:   "test_worker_run/when_command_ends_with_zero_exit_code",
			Value: []byte("worker-1"),
		}, nil)
		if err != nil {
			t.Error(err)
			return
		}

		testable := New("true", []string{}, "test_worker_run/when_command_ends_with_zero_exit_code", "worker-1", time.Second)

		err = testable.Run()
		if err != nil {
			t.Error(err)
			return
		}

		select {
		case err := <-testable.DoneChan:
			if err != nil {
				t.Errorf("expected nil, got '%s'", err.Error())
				return
			}
		case <-time.After(5 * time.Second):
			t.Error("timeout exceeded")
			return
		}
	})

	t.Run("when command ends with non-zero exit code", func(t *testing.T) {
		cc, err := api.NewClient(api.DefaultConfig())
		if err != nil {
			t.Error(err)
			return
		}

		_, err = cc.KV().Put(&api.KVPair{
			Key:   "test_worker_run/when_command_ends_with_non_zero_exit_code",
			Value: []byte("worker-1"),
		}, nil)
		if err != nil {
			t.Error(err)
			return
		}

		testable := New("false", []string{}, "test_worker_run/when_command_ends_with_non_zero_exit_code", "worker-1", time.Second)

		err = testable.Run()
		if err != nil {
			t.Error(err)
			return
		}

		select {
		case err := <-testable.DoneChan:
			if err == nil {
				t.Error("expected 'exit status 1', got nil")
				return
			}
			if err.Error() != "exit status 1" {
				t.Errorf("expected 'exit status 1', got '%v'", err.Error())
				return
			}
		case <-time.After(5 * time.Second):
			t.Error("timeout exceeded")
			return
		}
	})

	t.Run("when run unknown command", func(t *testing.T) {
		cc, err := api.NewClient(api.DefaultConfig())
		if err != nil {
			t.Error(err)
			return
		}

		_, err = cc.KV().Put(&api.KVPair{
			Key:   "test_worker_run/when_run_unknown_command",
			Value: []byte("worker-1"),
		}, nil)
		if err != nil {
			t.Error(err)
			return
		}

		testable := New("glory-of-satan", []string{}, "test_worker_run/when_run_unknown_command", "worker-1", time.Second)

		err = testable.Run()
		if err != nil {
			t.Error(err)
			return
		}

		select {
		case err := <-testable.DoneChan:
			if err == nil {
				t.Error(`expected 'exec: "glory-of-satan": executable file not found in $PATH', got nil`)
				return
			}
			if err.Error() != `exec: "glory-of-satan": executable file not found in $PATH` {
				t.Errorf(`expected 'exec: "glory-of-satan": executable file not found in $PATH', got '%v'`, err.Error())
				return
			}
		case <-time.After(5 * time.Second):
			t.Error("timeout exceeded")
			return
		}
	})

	t.Run("when pass empty key path", func(t *testing.T) {
		testable := New("false", []string{}, "", "worker-1", time.Second)

		err := testable.Run()
		if err == nil {
			t.Error("expected 'Must specify a single key to watch', got nil")
			return
		}
		if err.Error() != "Must specify a single key to watch" {
			t.Errorf("expected 'Must specify a single key to watch', got '%v'", err.Error())
			return
		}
	})
}

func TestWorker_StopCommand(t *testing.T) {
	t.Run("when command is not stuck", func(t *testing.T) {
		cc, err := api.NewClient(api.DefaultConfig())
		if err != nil {
			t.Error(err)
			return
		}

		_, err = cc.KV().Put(&api.KVPair{
			Key:   "test_worker_stop_command/when_command_is_not_stuck",
			Value: []byte("worker-1"),
		}, nil)
		if err != nil {
			t.Error(err)
			return
		}

		testable := New("sleep", []string{"10"}, "test_worker_stop_command/when_command_is_not_stuck", "worker-1", time.Second)

		err = testable.Run()
		if err != nil {
			t.Error(err)
			return
		}

		time.Sleep(100 * time.Millisecond)

		err = testable.StopCommand()
		if err != nil {
			t.Error(err)
			return
		}

		select {
		case err := <-testable.DoneChan:
			if err == nil {
				t.Error("expected 'signal: terminated', got nil")
				return
			}
			if err.Error() != "signal: terminated" {
				t.Errorf("expected 'signal: terminated', got '%v'", err.Error())
				return
			}
		case <-time.After(5 * time.Second):
			t.Error("timeout exceeded")
			return
		}
	})

	t.Run("when stop not started command", func(t *testing.T) {
		cc, err := api.NewClient(api.DefaultConfig())
		if err != nil {
			t.Error(err)
			return
		}

		_, err = cc.KV().Put(&api.KVPair{
			Key:   "test_worker_stop_command/when_stop_not_started_command",
			Value: []byte("worker-1"),
		}, nil)
		if err != nil {
			t.Error(err)
			return
		}

		testable := New("false", []string{}, "test_worker_stop_command/when_stop_not_started_command", "worker-1", time.Second)

		err = testable.StopCommand()
		if err != ErrCommandNotStarted {
			t.Errorf("expected '%s', got '%s'", ErrCommandNotStarted.Error(), err.Error())
			return
		}
	})
}
