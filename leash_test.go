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
				t.Error(err)
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
				t.Error("error expected")
				return
			}
		case <-time.After(5 * time.Second):
			t.Error("timeout exceeded")
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
}
