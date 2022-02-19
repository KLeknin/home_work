package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type sParams struct {
	tstName, toFile string
	offset, limit   int64
	resultData      string
}

func TestCopy(t *testing.T) {
	// Place your code here.
	tempDir, err := os.MkdirTemp("", "tmp")
	check(err)
	defer os.RemoveAll(tempDir)

	tstData := "1234567890"
	tstFileName := tempDir + string(os.PathSeparator) + "tstFile.txt"

	tstFile, err := os.Create(tstFileName)
	check(err)

	fmt.Fprint(tstFile, tstData)
	tstFile.Close()

	tstParams := []sParams{
		{"Full copy", "FullCopy.txt", 0, 0, tstData},
		{"Full copy, oversize limit", "FullCopy.txt", 0, 100500, tstData},
		{"Head copy", "HeadCopy.txt", 0, 3, "123"},
		{"Tail copy", "TailCopy.txt", 7, 0, "890"},
		{"Middle copy", "MiddleCopy.txt", 5, 2, "67"},
	}

	for _, tstParam := range tstParams {
		t.Run(tstParam.tstName, func(t *testing.T) {
			resultFileName := tempDir + string(os.PathSeparator) + tstParam.toFile

			err = Copy(tstFileName, resultFileName, tstParam.offset, tstParam.limit)
			check(err)

			resultBytes, err := os.ReadFile(resultFileName)
			check(err)

			resultStr := string(resultBytes)
			require.Equal(t, tstParam.resultData, resultStr)
		})
	}

	tstParam := tstParams[0]

	t.Run("Error File exist", func(t *testing.T) {
		err = Copy(tstFileName, tstFileName, tstParam.offset, tstParam.limit)
		require.Equal(t, ErrFileExist, err)
	})

	t.Run("Error Unsupported file", func(t *testing.T) {
		err = Copy(tempDir, tstFileName, tstParam.offset, tstParam.limit)
		require.Equal(t, ErrUnsupportedFile, err)
	})

	t.Run("Error Offset Exceeds File Size", func(t *testing.T) {
		resultFileName := tempDir + string(os.PathSeparator) + "_" + tstParam.toFile
		err = Copy(tstFileName, resultFileName, 100500, tstParam.limit)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})
}
