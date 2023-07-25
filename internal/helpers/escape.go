package helpers

import (
	"strings"
)

func EscapeIllegalCharacters(data string) string {
	return strings.ReplaceAll(data, "\u0000", "")
}

func EscapeIllegalCharactersBytes(data []byte) []byte {
	return []byte(EscapeIllegalCharacters(string(data)))
}
