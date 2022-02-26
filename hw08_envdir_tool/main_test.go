package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvPrepare(t *testing.T) {
	// envPrepare(env []string) Environment
	t.Run("a=1;b=2", func(t *testing.T) {
		sEnv := []string{"a=1", "b=2"}
		rEnv := Environment{
			"a": {"1", false},
			"b": {"2", false},
		}
		require.Equal(t, rEnv, envPrepare(sEnv))
	})
	t.Run("nil", func(t *testing.T) {
		rEnv := Environment{}
		require.Equal(t, rEnv, envPrepare(nil))
	})
	t.Run("empty", func(t *testing.T) {
		var sEnv []string
		rEnv := Environment{}
		require.Equal(t, rEnv, envPrepare(sEnv))
	})
}
