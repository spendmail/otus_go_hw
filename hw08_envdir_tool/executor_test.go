package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("executing test", func(t *testing.T) {
		dir, _ := os.Getwd()
		args := []string{"/bin/bash", path.Join(dir, "testdata/echo.sh"), "arg1=1", "arg2=2"}

		env := make(Environment)

		exitCode := RunCmd(args, env)

		require.Equal(t, exitCode, 0)
	})

	t.Run("exit status 0 test", func(t *testing.T) {
		dir, _ := os.Getwd()
		args := []string{"/bin/bash", path.Join(dir, "testdata/exit0.sh")}

		env := make(Environment)

		exitCode := RunCmd(args, env)

		require.Equal(t, 0, exitCode)
	})

	t.Run("exit status 1 test", func(t *testing.T) {
		dir, _ := os.Getwd()
		args := []string{"/bin/bash", path.Join(dir, "testdata/exit1.sh")}

		env := make(Environment)

		exitCode := RunCmd(args, env)

		require.Equal(t, 1, exitCode)
	})
}
