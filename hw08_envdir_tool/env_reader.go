package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var (
	replacingCharacterByte     = byte(0x00)
	linesDelimiter             = '\n'
	linesDelimiterByte         = byte(linesDelimiter)
	linesDelimiterUint8        = uint8(linesDelimiter)
	unacceptableCharacterUint8 = uint8('=')
	trimCharacters             = "\t "
)

var (
	ErrFileNotReadable       = errors.New("unable to read file")
	ErrDirNotReadable        = errors.New("unable to read dir")
	ErrUnacceptableCharacter = errors.New("the file has an unacceptable characters")
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	environmentMap := make(Environment)

	// Reading the dir.
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDirNotReadable, dir)
	}

	// Directory files loop.
	for _, fileInfo := range fileInfos {
		filePath := path.Join(dir, fileInfo.Name())
		fileSize := fileInfo.Size()

		// If file is empty, setting the flag "NeedRemove".
		if fileSize == 0 {
			environmentMap[fileInfo.Name()] = EnvValue{Value: "", NeedRemove: true}
			continue
		}

		// Opening the file.
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFileNotReadable, filePath)
		}

		defer file.Close()

		// Reading the file.
		fileBytes := make([]byte, fileSize)
		_, err = file.Read(fileBytes)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFileNotReadable, filePath)
		}

		// Taking the first line.
		for i, v := range fileBytes {
			// If the content has the unacceptable character, returning with error.
			if v == unacceptableCharacterUint8 {
				return nil, fmt.Errorf("%w: %v", ErrUnacceptableCharacter, filePath)
			}

			// Taking the sub-slice before the delimiter.
			if v == linesDelimiterUint8 {
				fileBytes = fileBytes[0:i]
				break
			}
		}

		// Replace \0 by "linesDelimiterByte".
		fileBytes = bytes.ReplaceAll(fileBytes, []byte{replacingCharacterByte}, []byte{linesDelimiterByte})

		// Trim from the right side.
		value := strings.TrimRight(string(fileBytes), trimCharacters)

		environmentMap[fileInfo.Name()] = EnvValue{Value: value, NeedRemove: false}
	}

	return environmentMap, nil
}
