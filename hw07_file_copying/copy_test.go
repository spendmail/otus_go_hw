package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTmpFile() (*os.File, error) {
	f, err := os.CreateTemp("/tmp/", "")
	if err != nil {
		return nil, err
	}

	if _, err := f.Write([]byte("short content")); err != nil {
		return nil, err
	}

	if err := f.Close(); err != nil {
		return nil, err
	}

	return f, nil
}

func TestCopy(t *testing.T) {
	t.Run("unsupporting file error", func(t *testing.T) {
		err := Copy("/dev/urandom", "/tmp/out", 0, 0)

		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual error %q", err)
	})

	t.Run("offset exceeds file size error", func(t *testing.T) {
		f, err := createTmpFile()
		if err != nil {
			t.Fatal(err)
		}

		err = Copy(f.Name(), "/tmp/out", 1024*1024*1024, 0)

		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual error %q", err)
	})

	t.Run("source file is not found error", func(t *testing.T) {
		err := Copy("not_existent_file", "/tmp/out", 0, 0)

		require.Truef(t, errors.Is(err, ErrSourceFileNotFound), "actual error %q", err)
	})

	t.Run("source file is not readable error", func(t *testing.T) {
		f, err := createTmpFile()
		if err != nil {
			t.Fatal(err)
		}

		err = os.Chmod(f.Name(), 0o200)
		if err != nil {
			t.Fatal(err)
		}

		err = Copy(f.Name(), "/tmp/out", 0, 0)

		require.Truef(t, errors.Is(err, ErrSourceFileNotReadable), "actual error %q", err)
	})

	t.Run("destination file can't be created error", func(t *testing.T) {
		src, err := createTmpFile()
		if err != nil {
			t.Fatal(err)
		}

		dst, err := createTmpFile()
		if err != nil {
			t.Fatal(err)
		}

		err = os.Chmod(dst.Name(), 0o400)
		if err != nil {
			t.Fatal(err)
		}

		err = Copy(src.Name(), dst.Name(), 0, 0)

		require.Truef(t, errors.Is(err, ErrDestinationFileCreate), "actual error %q", err)
	})
}
