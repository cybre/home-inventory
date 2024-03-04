package utils

import "unicode"

func FirstLetterUppercase(s string) string {
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])

	return string(r)
}
