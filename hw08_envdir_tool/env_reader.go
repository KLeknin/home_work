package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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
	if err == io.EOF {
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
		if f.IsDir() {
			continue
		}

		envValue = EnvValue{}
		if f.Size() > 0 {
			fullFileName := dir + string(os.PathSeparator) + f.Name()
			envValue.Value, err = fileFirstLine(fullFileName)
			if err != nil {
				return nil, fmt.Errorf("errror reading file %s: %v", f.Name(), err)
			}
		}
		envValue.NeedRemove = envValue.Value == ""
		env[f.Name()] = envValue
	}
	return env, nil
}
