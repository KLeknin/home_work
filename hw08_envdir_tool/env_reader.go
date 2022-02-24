package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func fileFirstLine(fileName string) (value string, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer func(e *error) {
		err := file.Close()
		if *e == nil && err != nil {
			*e = err
		}
	}(&err)

	reader := bufio.NewReader(file)
	value, err = reader.ReadString('\n')
	if errors.Is(err, io.EOF) {
		err = nil
	}
	if err != nil {
		return "", err
	}
	if value[len(value)-1] == '\n' {
		value = value[:len(value)-1]
	}
	return
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Place your code here
	fi, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := Environment{}
	var envValue EnvValue
	for _, f := range fi {
		if strings.Index(f.Name(), "=") > 0 {
			return nil, fmt.Errorf("errror file name %s contains \"=\"", f.Name())
		}
		if f.IsDir() {
			continue
		}

		envValue = EnvValue{}
		envValue.NeedRemove = f.Size() == 0
		if f.Size() > 0 {
			fullFileName := dir + string(os.PathSeparator) + f.Name()
			envValue.Value, err = fileFirstLine(fullFileName)
			if err != nil {
				return nil, fmt.Errorf("errror reading file %s: %w", f.Name(), err)
			}
		}
		envValue.Value = strings.TrimRight(envValue.Value, " \t\r")
		envValue.Value = strings.Replace(envValue.Value, "\x00", "\n", -1)
		env[f.Name()] = envValue
	}
	return env, nil
}
