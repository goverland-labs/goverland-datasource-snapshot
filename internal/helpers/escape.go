package helpers

import (
	"encoding/json"
	"strings"
)

func EscapeIllegalCharacters(data string) string {
	result := strings.ReplaceAll(data, "\u0000", "")
	result = strings.ReplaceAll(result, "\\u0000", "")

	return result
}

func EscapeIllegalCharactersJson(data json.RawMessage) json.RawMessage {
	return json.RawMessage(EscapeIllegalCharacters(string(data)))
}
