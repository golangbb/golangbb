package helpers

import (
	"os"
	"testing"
"gotest.tools/v3/assert")

func TestGetEnv(t *testing.T) {
	t.Run("should return env var from os", func(t *testing.T) {
		envVarKey := "envVarKey"
		envVarValue := "envVarValue"
		fallbackValue:= "fallbackValue"

		os.Setenv(envVarKey, envVarValue)
		grabbedEnvVarValue := GetEnv(envVarKey, fallbackValue)

		assert.Equal(t, grabbedEnvVarValue, envVarValue)
		assert.Assert(t, grabbedEnvVarValue != fallbackValue)

		os.Unsetenv(envVarKey)
	})

	t.Run("should return fallback value when env var not set", func(t *testing.T) {
		envVarKey := "envVarKey"
		envVarValue := "envVarValue"
		fallbackValue:= "fallbackValue"

		os.Unsetenv(envVarKey)
		grabbedEnvVarValue := GetEnv(envVarKey, fallbackValue)

		assert.Equal(t, grabbedEnvVarValue, fallbackValue)
		assert.Assert(t, grabbedEnvVarValue != envVarValue)
	})
}
