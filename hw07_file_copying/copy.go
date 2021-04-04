package main

import (
	"errors"
	"io"
	"os"
)

var (
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

	// Trying to copy n byte to destination file
	if _, err := io.CopyN(dstFile, srcFile, limit); err != nil {
		return ErrFileCopy
	}

	return nil
}
