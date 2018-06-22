package main

import (
	"fmt"
	"os"
	"time"

	"github.com/iqoption/consul-leash"
)

var (
	path                string
	stoppingDurationStr string
	value               string
)

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString(fmt.Sprintf("Usage: %s", os.Args[0]))
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}

	path = os.Getenv("LEASH_KEY_PATH")
	if path == "" {
		os.Stderr.WriteString("Env variable LEASH_KEY_PATH is required")
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}

	value = os.Getenv("LEASH_KEY_VALUE")
	if value == "" {
		os.Stderr.WriteString("Env variable LEASH_KEY_VALUE is required")
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}

	stoppingDurationStr = os.Getenv("LEASH_STOPPING_DURATION")
	if stoppingDurationStr == "" {
		stoppingDurationStr = "10s"
	}

	stoppingDuration, err := time.ParseDuration(stoppingDurationStr)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Env variable LEASH_STOPPING_DURATION has invalid value (%s)", err.Error()))
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}

	l := leash.New(os.Args[1], os.Args[2:], path, value, stoppingDuration)

	l.Run()

	err = <-l.DoneChan
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
}
