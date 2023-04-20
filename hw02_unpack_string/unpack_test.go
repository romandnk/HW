package hw02unpackstring

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input         string
		expected      string
		expectedError error
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde", expectedError: nil},
		{input: "abcd", expected: "abcd", expectedError: nil},
		{input: "aaa0b", expected: "aab", expectedError: nil},
		{input: "", expected: "", expectedError: nil},
		{input: "a0", expected: "", expectedError: nil},
		{input: "aasd0f0ghm0", expected: "aasgh", expectedError: nil},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc", expectedError: nil},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`, expectedError: nil},
		{input: `qwe\45`, expected: `qwe44444`, expectedError: nil},
		{input: `qwe\\5`, expected: `qwe\\\\\`, expectedError: nil},
		{input: `qwe\\\3`, expected: `qwe\3`, expectedError: nil},
		{input: `\\a\\`, expected: `\a\`, expectedError: nil},
		{input: `\0\\`, expected: `0\`, expectedError: nil},
		{input: `\56a`, expected: `555555a`, expectedError: nil},
		{input: "45", expected: "", expectedError: ErrDigitOnTheFirstPlace},
		{input: "3abc", expected: "", expectedError: ErrDigitOnTheFirstPlace},
		{input: `qw\ne`, expected: "", expectedError: ErrEscapeLetter},
		{input: `t5b8\h`, expected: "", expectedError: ErrEscapeLetter},
		{input: "aaa12b", expected: "", expectedError: ErrNumberInString},
		{input: "a55b", expected: "", expectedError: ErrNumberInString},
		{input: `abc\\333`, expected: "", expectedError: ErrNumberInString},
		{input: `\\67`, expected: "", expectedError: ErrNumberInString},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.Equal(t, tc.expected, result)
			require.Equal(t, tc.expectedError, err)
		})
	}
}
