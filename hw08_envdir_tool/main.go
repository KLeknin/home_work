package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func envPrepare(env []string) Environment {
	result := make(Environment, len(env))
	for _, envStr := range env {
		ss := strings.SplitN(envStr, "=", 2)
		if len(ss) == 0 {
			continue
		}
		envName := ss[0]
		envValue := EnvValue{}
		if len(ss) > 1 {
			envValue.Value = ss[1]
		}
		result[envName] = envValue
	}
	return result
}

func main() {
	// go-envdir /path/to/env/dir command arg1 arg2
	if len(os.Args) < 2 {
		log.Fatalf("not enough params")
	}
	envDir := os.Args[1]
	envNew, err := ReadDir(envDir)
	check(err)

	env := envPrepare(os.Environ())

	for keyName := range envNew {
		if envNew[keyName].NeedRemove {
			delete(env, keyName)
			continue
		}
		newValue := EnvValue{}
		newValue.Value = envNew[keyName].Value
		env[keyName] = newValue
	}

	RunCmd(os.Args[2:], env)
}

func check(e error) {
	if e != nil {
		fmt.Printf("Error: %s\n", e.Error())
	}
}
