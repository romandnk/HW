package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDirWithoutErrors(t *testing.T) {
	testCases := []struct {
		path     string
		expected Environment
	}{
		{
			path: "./testdata/env",
			expected: Environment{
				"BAR": EnvValue{
					Value:      "bar",
					NeedRemove: false,
				},
				"EMPTY": EnvValue{
					Value:      "",
					NeedRemove: false,
				},
				"FOO": EnvValue{
					Value:      "   foo\nwith new line",
					NeedRemove: false,
				},
				"HELLO": EnvValue{
					Value:      `"hello"`,
					NeedRemove: false,
				},
				"UNSET": EnvValue{
					Value:      "",
					NeedRemove: true,
				},
			},
		},
		{
			path: "./testdata/my_test_env",
			expected: Environment{
				"SMALL": EnvValue{
					Value:      "    small",
					NeedRemove: false,
				},
				"SPACESRIGHT": EnvValue{
					Value:      "         SPACE",
					NeedRemove: false,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run("successful", func(t *testing.T) {
			actualEnv, err := ReadDir(tc.path)

			require.NoError(t, err)
			require.Equal(t, tc.expected, actualEnv)
		})
	}
}

func TestReadDirWithError(t *testing.T) {
	path := "./testdata/env_errors"
	_, err := ReadDir(path)

	require.ErrorIs(t, ErrEqualSign, err)
}
