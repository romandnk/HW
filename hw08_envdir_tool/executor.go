package main

import (
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmdArgs []string, env Environment) (returnCode int) {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...) //nolint:gosec

	envSlice := make([]string, 0, len(env))

	for key, value := range env {
		var keyValueEnv strings.Builder
		if value.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				return 126
			}
			if value.Value != "" {
				keyValueEnv.WriteString(key + "=" + value.Value)
			}
			continue
		}
		keyValueEnv.WriteString(key + "=" + value.Value)
		envSlice = append(envSlice, keyValueEnv.String())
	}

	cmd.Env = append(os.Environ(), envSlice...)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return 1
	}

	return 0
}
