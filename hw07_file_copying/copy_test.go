package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	fromPathTest = "testdata/input.txt"
	toPathTest   = "tmp/output.txt"
)

func TestCopyErrOffsetExceedsFileSize(t *testing.T) {
	testCases := []struct {
		name          string
		fromPath      string
		toPath        string
		offset        int64
		limit         int64
		expectedError error
		accessRules   os.FileMode
	}{
		{
			name:          "offset more than file size",
			fromPath:      fromPathTest,
			toPath:        toPathTest,
			offset:        6618,
			limit:         0,
			expectedError: ErrOffsetExceedsFileSize,
			accessRules:   0o777,
		},
		{
			name:          "no permission to the file",
			fromPath:      fromPathTest,
			toPath:        toPathTest,
			offset:        0,
			limit:         0,
			expectedError: ErrNoPermission,
			accessRules:   0o000,
		},
		{
			name:          "paths are equal",
			fromPath:      toPathTest,
			toPath:        "tmp/../tmp/output.txt",
			offset:        0,
			limit:         0,
			expectedError: ErrEqualPaths,
			accessRules:   0o777,
		},
		{
			name:          "offset less than 0",
			fromPath:      fromPathTest,
			toPath:        toPathTest,
			offset:        -1,
			limit:         0,
			expectedError: ErrNegativeOffset,
			accessRules:   0o000,
		},
		{
			name:          "limit less than 0",
			fromPath:      fromPathTest,
			toPath:        toPathTest,
			offset:        0,
			limit:         -1,
			expectedError: ErrNegativeLimit,
			accessRules:   0o000,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if err := os.Chmod(tc.fromPath, tc.accessRules); err != nil {
				t.Error(err)
			}

			err := Copy(tc.fromPath, tc.toPath, tc.offset, tc.limit)

			require.Error(t, err)
			require.EqualError(t, err, tc.expectedError.Error())
		})
	}
	if err := os.Chmod(fromPathTest, 0o777); err != nil {
		t.Error(err)
	}
	_ = os.Remove(toPathTest)
}
