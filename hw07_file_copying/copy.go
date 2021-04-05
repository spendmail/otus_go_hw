package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	step                         = int64(1024 * 1024)
	ErrUnsupportedFile           = errors.New("unsupported file")
	ErrOffsetExceedsFileSize     = errors.New("offset exceeds file size")
	ErrSourceFileNotFound        = errors.New("source file is not found")
	ErrSourceFileNotReadable     = errors.New("source file is not readable")
	ErrSourceFileStatNotReadable = errors.New("source file stat is not readable")
	ErrDestinationFileCreate     = errors.New("destination file can't be created")
	ErrSourceFileSeek            = errors.New("source file can't be sought")
	ErrFileCopy                  = errors.New("source file can't be copied into destination file")
)

// Copying "limit" bytes from "fromPath" to "toPath" with offset "offset".
func Copy(fromPath, toPath string, offset, limit int64) error {
	// Attempting to open a source file
	srcFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrSourceFileNotFound
		}
		return ErrSourceFileNotReadable
	}
	defer srcFile.Close()

	// Getting file statistics
	srcFileInfo, err := srcFile.Stat()
	if err != nil {
		return ErrSourceFileStatNotReadable
	}

	// Getting file size
	srcFileSize := srcFileInfo.Size()
	if srcFileSize == 0 {
		return ErrUnsupportedFile
	}

	// Attempting to create a destination file
	dstFile, err := os.Create(toPath)
	if err != nil {
		return ErrDestinationFileCreate
	}
	defer dstFile.Close()

	// Checking the offset
	if offset >= srcFileSize {
		return ErrOffsetExceedsFileSize
	} else if offset < 0 {
		offset = 0
	}

	// Setting the file cursor
	_, err = srcFile.Seek(offset, 0)
	if err != nil {
		return ErrSourceFileSeek
	}

	// Checking the limit
	if limit <= 0 {
		limit = srcFileSize
	} else if limit > srcFileSize-offset {
		limit = srcFileSize - offset
	}

	// Defines a progress bar
	bar := pb.StartNew(int(limit))

	// Trying to copy n byte to destination file in the loop by step bytes per time
	for limit > 0 {
		if step >= limit {
			step = limit
		}

		// Copying step bytes into a destination file
		_, err = io.CopyN(dstFile, srcFile, step)
		if err != nil {
			return ErrFileCopy
		}

		// Adding values into a progress bar
		bar.Add(int(step))

		limit -= step
	}

	// Finishing progress bar's work
	bar.FinishPrint("Done!")

	return nil
}
