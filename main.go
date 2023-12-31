package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/palantir/pkg/retry"
	"github.com/spf13/pflag"
)

// set by goreleaser
var version = "unknown"

func main() {
	// handle negative
	maxAttempts := pflag.Uint("max-attempts", 0, "upper limit of number of attempts. 0 indicates no limit.")
	maxBackoff := pflag.Duration("max-backoff", 0, "upper limit of backoff duration. 0 indicates no limit.")
	// handle <=0 values
	initialBackoff := pflag.Duration("initial-backoff", time.Second, "initial backoff duration.")
	// handle <=0 values
	multiplier := pflag.Float64("multiplier", 2, "multiplier to apply after each failed attempt.")
	randomisation := pflag.Float64("randomisation", 0, "randomisation to apply to the multiplication of each backoff")
	showVersion := pflag.BoolP("version", "v", false, "print the version ")

	pflag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	cmd := pflag.Args()
	if len(cmd) == 0 {
		pflag.Usage()
		os.Exit(3)
	}

	getpgid, err := syscall.Getpgid(os.Getpid())
	if err != nil {
		fmt.Println("error getting pgid of retry:", err.Error())
		os.Exit(4)
	}

	var lastExitCode int
	err = retry.Do(context.Background(), func() error {
		cmd2 := exec.Command(cmd[0], cmd[1:]...)
		cmd2.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
			Pgid:    getpgid,
		}
		cmd2.Stdout = os.Stdout
		cmd2.Stderr = os.Stderr
		// confirm how exit code handling works
		err := cmd2.Run()
		if err != nil {
			fmt.Println("error: ", err.Error())
		}
		lastExitCode = cmd2.ProcessState.ExitCode()
		return err
	},
		retry.WithMaxAttempts(int(*maxAttempts)),
		retry.WithMaxBackoff(*maxBackoff),
		retry.WithInitialBackoff(*initialBackoff),
		retry.WithMultiplier(*multiplier),
		retry.WithRandomizationFactor(*randomisation))

	if err != nil {
		fmt.Println("error: ", err.Error())
	}
	os.Exit(lastExitCode)
}
