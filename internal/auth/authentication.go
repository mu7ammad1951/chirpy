package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, errors.New("invalid token: expired")
	}

	userID, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(userID)
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
