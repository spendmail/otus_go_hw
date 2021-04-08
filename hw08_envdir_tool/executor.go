package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmdArgs []string, env Environment) (returnCode int) {
	// Handling given env variables.
	for key, value := range env {
		if value.NeedRemove {
			// Either removing variable if, "NeedRemove" flag has been set.
			err := os.Unsetenv(key)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// Or setting/creating env variable as well.
			err := os.Setenv(key, value.Value)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Preparing the command.
	cmdName := cmdArgs[0]
	args := cmdArgs[1:]
	cmd := exec.Command(cmdName, args...)

	// exec.Command doesn't set env variables, therefore setting them manually.
	cmd.Env = os.Environ()

	// Launching the command.
	stdoutStderr, err := cmd.CombinedOutput()

	// Retrieving exit status code from the error.
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}

	// Printing the output.
	fmt.Println(string(stdoutStderr))

	return 0
}
