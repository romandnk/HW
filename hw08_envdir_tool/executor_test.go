package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmdWithoutErrorsWithoutEnvVariables(t *testing.T) {
	cmdArgs := []string{"echo", "Hello, World!"}
	env := Environment{}

	returnCode := RunCmd(cmdArgs, env)

	require.Equal(t, 0, returnCode)
}

func TestRunCmdCommandError(t *testing.T) {
	cmdArgs := []string{"invalid-command"}
	env := Environment{}

	returnCode := RunCmd(cmdArgs, env)

	require.Equal(t, 1, returnCode)
}

func TestRunCmdWithoutErrorsWithEnvVariables(t *testing.T) {
	cmdArgs := []string{"echo", "Hello, World!"}
	env := Environment{
		"TEST_KEY": EnvValue{
			Value:      "test_value",
			NeedRemove: false,
		},
	}

	returnCode := RunCmd(cmdArgs, env)

	require.Equal(t, 0, returnCode)
}
