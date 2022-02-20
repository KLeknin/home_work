package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRunCmd(t *testing.T) {
	// Place your code here
	require.Equal(t, 0, RunCmd([]string{"go", "version"}, Environment{}))

}
