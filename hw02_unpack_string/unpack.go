package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var (
		temp   strings.Builder
		result strings.Builder
		prev   rune
	)

	for _, cur := range s {
		// check if current symbol is digit
		if unicode.IsDigit(cur) {
			if prev == 0 {
				return "", ErrInvalidString
			}
			// check if previous symbol is digit
			if unicode.IsDigit(prev) {
				return "", ErrInvalidString
			}
			// get number of repetitions
			digit := int(cur - '0')
			// trim last symbol if number of repetition is 0
			if digit == 0 {
				temp.WriteString(result.String()[:len(result.String())-1])
				result.Reset()
				result.WriteString(temp.String())
				temp.Reset()
			} else {
				result.WriteString(strings.Repeat(string(prev), digit-1))
				prev = cur
			}
		} else {
			prev = cur
			result.WriteRune(prev)
		}
	}
	return result.String(), nil
}
