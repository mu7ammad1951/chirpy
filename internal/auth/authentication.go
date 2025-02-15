package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const costToHashPassword = 10

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
