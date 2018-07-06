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
			Key:   "test/when_command_ends_with_zero_exit_code",
			Value: []byte("worker-1"),
		}, nil)
		if err != nil {
			t.Error(err)
			return
		}

		testable := New("true", []string{}, "test/when_command_ends_with_zero_exit_code", "worker-1", time.Second)

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
			Key:   "test/when_command_ends_with_non_zero_exit_code",
			Value: []byte("worker-1"),
		}, nil)
		if err != nil {
			t.Error(err)
			return
		}

		testable := New("false", []string{}, "test/when_command_ends_with_non_zero_exit_code", "worker-1", time.Second)

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
