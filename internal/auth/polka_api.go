package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authValue := headers.Get("Authorization")
	if authValue == "" {
		return "", errors.New("missing header")
	}

	authTokenString, found := strings.CutPrefix(authValue, "ApiKey ")
	if !found {
		return "", errors.New("malformed header")
	}
	authTokenString = strings.TrimSpace(authTokenString)
	if authTokenString == "" {
		return "", errors.New("missing token string")
	}
	return authTokenString, nil
}
