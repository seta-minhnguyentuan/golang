package utils

import "encoding/json"

func DecodeJSON[T any](b []byte, out *T) error {
	return json.Unmarshal(b, out)
}
