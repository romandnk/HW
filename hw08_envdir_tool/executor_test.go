package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	testCase := []struct {
		name         string
		cmdArgs      []string
		expectedCode int
	}{
		{
			name:         "without errors",
			cmdArgs:      []string{"echo", "Hello, World!"},
			expectedCode: 0,
		},
		{
			name:         "process has not started",
			cmdArgs:      []string{"random_command"},
			expectedCode: -1,
		},
		{
			name:         "directory does not exist",
			cmdArgs:      []string{"ls", "/any_directory"},
			expectedCode: 2,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			actualCode := RunCmd(tc.cmdArgs, nil)
			require.Equal(t, tc.expectedCode, actualCode)
		})
	}
}
