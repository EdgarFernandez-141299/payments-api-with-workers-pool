package utils

import (
	"encoding/base64"
)

func Base64Decode(input string) string {
	decoded, err := base64.RawStdEncoding.DecodeString(input)
	if err != nil {
		return ""
	}
	return string(decoded)
}
