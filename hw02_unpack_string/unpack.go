package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrStringStartWithDigit = errors.New("string starts with a digit")

var ErrStringContainsDoubleDigit = errors.New("string contains a double digit")

func Unpack(s string) (string, error) {
	runes := []rune(s)
	length := len(runes)

	// If string is empty, then return ""
	if len(runes) == 0 {
		return s, nil
	}

	// If the string starts with a digit, then return an error
	if unicode.IsDigit(runes[0]) {
		return "", ErrStringStartWithDigit
	}

	var stringBuilder strings.Builder

	for current, next := 0, 1; next < length; current, next = current+1, next+1 {
		if unicode.IsDigit(runes[current]) {
			// In case of double digit, return an error
			if unicode.IsDigit(runes[next]) {
				return "", ErrStringContainsDoubleDigit
			}
			continue
		}

		if unicode.IsDigit(runes[next]) {
			// No need to check on error, because of the rune is always a single digit
			digit, _ := strconv.Atoi(string(runes[next]))
			stringBuilder.WriteString(strings.Repeat(string(runes[current]), digit))
		} else {
			stringBuilder.WriteRune(runes[current])
		}
	}

	// Processing the last rune
	if !unicode.IsDigit(runes[length-1]) {
		stringBuilder.WriteRune(runes[length-1])
	}

	return stringBuilder.String(), nil
}
