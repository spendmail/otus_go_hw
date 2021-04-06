package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

var (
	ErrUnreadableDir  = errors.New("unable to read the dir")
	ErrUnreadableFile = errors.New("unable to read the file")
	ErrUnreadableLine = errors.New("unable to read the file's line")
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

	content, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, ErrUnreadableDir
	}

	for _, fileInfo := range content {
		filePath := path.Join(dir, fileInfo.Name())

		file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
		if err != nil {
			return nil, ErrUnreadableFile
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		firstLine, err := reader.ReadString('\n')
		if err != nil {
			return nil, ErrUnreadableLine
		}

		//@Todo: если файл полностью пустой (длина - 0 байт), то envdir удаляет переменную окружения с именем S
		//@Todo: имя S не должно содержать =
		//@Todo: пробелы и табуляция в конце T удаляются
		//@Todo: терминальные нули (0x00) заменяются на перевод строки (\\n);
		//@Todo: if trim(line) != "" then set

		environmentMap[fileInfo.Name()] = EnvValue{Value: firstLine, NeedRemove: false}
	}

	return environmentMap, nil
}
