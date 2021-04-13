package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("dir is not readable test", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "temporary")
		if err != nil {
			log.Fatal(err)
		}

		defer os.RemoveAll(dir)

		err = os.Chmod(dir, 0o300)
		if err != nil {
			log.Fatal(err)
		}

		_, err = ReadDir(dir)

		require.Truef(t, errors.Is(err, ErrDirNotReadable), "actual error %q", err)
	})

	t.Run("file is not readable test", func(t *testing.T) {
		content := []byte("VALUE")
		dir, err := ioutil.TempDir("", "temporary")
		if err != nil {
			log.Fatal(err)
		}

		defer os.RemoveAll(dir)

		filePath := filepath.Join(dir, "NAME")
		if err := ioutil.WriteFile(filePath, content, 0o300); err != nil {
			log.Fatal(err)
		}

		_, err = ReadDir(dir)

		require.Truef(t, errors.Is(err, ErrFileNotReadable), "actual error %q", err)
	})

	t.Run("unacceptable character test", func(t *testing.T) {
		content := []byte("VALUE=\n")
		dir, err := ioutil.TempDir("", "temporary")
		if err != nil {
			log.Fatal(err)
		}

		defer os.RemoveAll(dir)

		filePath := filepath.Join(dir, "NAME")
		if err := ioutil.WriteFile(filePath, content, 0o666); err != nil {
			log.Fatal(err)
		}

		_, err = ReadDir(dir)

		require.Truef(t, errors.Is(err, ErrUnacceptableCharacter), "actual error %q", err)
	})
}
