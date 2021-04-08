package main

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("dir is not readable test", func(t *testing.T) {
		dir := "/tmp/unreadable_dir"

		_ = os.Mkdir(dir, 0o300)

		_, err := ReadDir(dir)

		require.Truef(t, errors.Is(err, ErrDirNotReadable), "actual error %q", err)

		_ = os.Remove(dir)
	})

	t.Run("file is not readable test", func(t *testing.T) {
		dir := "/tmp/readable_dir"
		file := "TEST"

		_ = os.Mkdir(dir, 0o755)
		filePath := path.Join(dir, file)
		f, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o300)
		_, _ = f.WriteString("VALUE\n")
		_ = f.Close()

		_, err := ReadDir(dir)

		require.Truef(t, errors.Is(err, ErrFileNotReadable), "actual error %q", err)

		_ = os.Remove(filePath)
		_ = os.Remove(dir)
	})

	t.Run("unacceptable character test", func(t *testing.T) {
		dir := "/tmp/readable_dir"
		file := "TEST"

		_ = os.Mkdir(dir, 0o755)
		filePath := path.Join(dir, file)
		f, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
		_, _ = f.WriteString("VALUE=\n")
		_ = f.Close()

		_, err := ReadDir(dir)

		require.Truef(t, errors.Is(err, ErrUnacceptableCharacter), "actual error %q", err)

		_ = os.Remove(filePath)
		_ = os.Remove(dir)
	})
}
