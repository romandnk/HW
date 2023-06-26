package main

import (
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmdArgs []string, env Environment) (returnCode int) {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...) //nolint:gosec

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	envs := setEnvs(env)

	cmd.Env = envs

	if err := cmd.Run(); err != nil {
		return cmd.ProcessState.ExitCode()
	}

	return 0
}

func setEnvs(env Environment) []string {
	globalEnv := make(map[string]string, len(os.Environ()))
	resultEnvs := make([]string, 0, len(os.Environ())+len(env))

	for _, val := range os.Environ() {
		envKeyVal := strings.Split(val, "=")
		globalEnv[envKeyVal[0]] = envKeyVal[1]
	}

	for key, val := range env {
		if val.NeedRemove {
			delete(globalEnv, key)
			continue
		}
		globalEnv[key] = val.Value
	}

	for key, val := range globalEnv {
		var keyValueEnv strings.Builder
		keyValueEnv.WriteString(key + "=" + val)
		resultEnvs = append(resultEnvs, keyValueEnv.String())
	}

	return resultEnvs
}
