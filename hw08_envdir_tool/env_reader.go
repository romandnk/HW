package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ErrEqualSign = errors.New("name of file cannot contain sign '='")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envs := make(Environment)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error opening file: path (%s), err: %w", path, err)
		}

		if info.IsDir() {
			return nil
		}

		fileName := strings.ToUpper(info.Name())
		if strings.ContainsRune(fileName, '=') {
			return ErrEqualSign
		}

		if info.Size() == 0 {
			envs[fileName] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("cannot open file: %s", path)
		}

		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			firstLine := strings.TrimRight(scanner.Text(), " ")

			if strings.ContainsRune(firstLine, '\x00') {
				firstLine = strings.ReplaceAll(firstLine, "\x00", "\n")
			}
			envs[fileName] = EnvValue{
				Value:      firstLine,
				NeedRemove: false,
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return envs, nil
}
