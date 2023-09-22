package helpers

import "encoding/json"

func Unmarshal[T any](v T, data []byte) (T, error) {
	if string(data) == "null" {
		return v, nil
	}

	err := json.Unmarshal(data, &v)

	return v, err
}
