package auth

import (
	"errors"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const costToHashPassword = 10

func GetBearerToken(headers http.Header) (tokenString string, err error) {
	authValue := headers.Get("Authorization")
	if authValue == "" {
		return "", errors.New("missing header")
	}

	authTokenString, found := strings.CutPrefix(authValue, "Bearer ")
	if !found {
		return "", errors.New("malformed header")
	}
	authTokenString = strings.TrimSpace(authTokenString)
	if authTokenString == "" {
		return "", errors.New("missing token string")
	}
	return authTokenString, nil
}

func HashPassword(password string) (string, error) {
	if len(password) < 1 {
		return "", errors.New("password too short")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), costToHashPassword)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func CheckPasswordHash(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return err
	}
	return nil
}
