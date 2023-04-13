package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidString = errors.New("invalid string")
)

func Unpack(s string) (string, error) {
	var (
		temp   strings.Builder
		sRune  = []rune(s)
		result strings.Builder
		cur    rune
		prev   rune
	)

	for i := 0; i < len(sRune); i++ {
		cur = sRune[i]
		// check if current symbol is digit
		if digit, err := strconv.Atoi(string(cur)); err == nil {
			if prev == 0 {
				return "", ErrInvalidString
			}
			// check if previous symbol is digit
			if _, err := strconv.Atoi(string(prev)); err == nil {
				return "", ErrInvalidString
			}
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
