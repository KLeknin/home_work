package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func envIsSame(a, b Environment) bool {
	var s string
	for s = range a {
		if a[s] != b[s] {
			return false
		}
	}
	for s = range b {
		if a[s] != b[s] {
			return false
		}
	}
	return true
}

func TestReadDir(t *testing.T) {
	// Place your code here
	tempDir, err := os.MkdirTemp("", "tmp")
	check(err)
	defer os.RemoveAll(tempDir)

	tstEnv := Environment{
		"first":   {"first1", false},
		"SECOND":  {"SECOND2", false},
		"Third":   {"Third3", false},
		"fourTH":  {"fourTH4", false},
		"fifth":   {"\nfifth", true},
		"sixth":   {"\n", true},
		"seventh": {"", true},
	}

	var tstFile *os.File
	defer func() {
		if tstFile != nil {
			tstFile.Close()
		}
	}()
	for envValue := range tstEnv {
		tstFileName := tempDir + string(os.PathSeparator) + envValue
		tstFile, err = os.Create(tstFileName)
		check(err)
		_, err = fmt.Fprint(tstFile, tstEnv[envValue].Value)
		check(err)
		tstFile.Close()
	}
	emptyValue := EnvValue{"", true}
	tstEnv["fifth"] = emptyValue
	tstEnv["sixth"] = emptyValue

	t.Run("ReadDir Test", func(t *testing.T) {
		dir, err := ReadDir(tempDir)
		require.Nilf(t, err, "error: %v", err)
		require.Truef(t, envIsSame(dir, tstEnv), "environments is not same: \n%v\n%v", dir, tstEnv)
	})

	//todo "errror file name %s contains \"=\""
}
