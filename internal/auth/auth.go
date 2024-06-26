package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Example header:
// Authorization: ApiKey {insert api_key here}
func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")

	if val == "" {
		return "", errors.New("invalid authentication info")
	}

	vals := strings.Split(val, " ")

	if len(vals) != 2 {
		return "", errors.New("invalid auth header")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first part of auth header")
	}

	return vals[1], nil
}
