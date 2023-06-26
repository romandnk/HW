package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDirWithoutErrors(t *testing.T) {
	path := "./testdata/env"
	actual, err := ReadDir(path)

	expected := Environment{
		"BAR": EnvValue{
			Value:      "newBar",
			NeedRemove: true,
		},
		"EMPTY": EnvValue{
			Value:      "",
			NeedRemove: true,
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
		"SPACESRIGHT": EnvValue{
			Value:      "         SPACE",
			NeedRemove: false,
		},
		"SMALL": EnvValue{
			Value:      "    small",
			NeedRemove: false,
		},
	}

	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestReadDirWithError(t *testing.T) {
	path := "./testdata/env_errors"
	_, err := ReadDir(path)

	require.ErrorIs(t, ErrEqualSign, err)
}
