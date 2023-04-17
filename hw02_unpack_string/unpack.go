package hw02unpackstring

import (
	"errors"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var (
		sRune      = []rune(s)
		result     []rune
		prev       rune
		checkSlash bool
	)
	// check if digit is on the first place
	if len(sRune) > 0 && unicode.IsDigit(sRune[0]) {
		return "", ErrInvalidString
	}
	for i := range sRune {
		cur := sRune[i]
		switch {
		case unicode.IsDigit(cur): // check if current symbol is digit
			if checkSlash {
				result = append(result, cur)
				prev = cur
				checkSlash = false
			} else {
				numRepetition := int(cur - '0')
				if numRepetition == 0 {
					result = result[:len(result)-1]
				} else {
					for j := 0; j < numRepetition-1; j++ {
						result = append(result, prev)
					}
					prev = cur
				}
			}
		case checkSlash && !unicode.IsLetter(cur): // check if there were two backslash
			result = append(result, '\\')
			checkSlash = false
			prev = cur
		case string(cur) == "\\": // make flag "checkSlash" true if current symbol is a backslash
			checkSlash = true
		default:
			// check if before letter was a backslash
			if checkSlash {
				return "", ErrInvalidString
			}
			result = append(result, cur)
			prev = cur
		}
	}
	return string(result), nil
}
