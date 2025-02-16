package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mu7ammad1951/chirpy-boot/internal/auth"
)

const (
	maxExpirationInSeconds = 3600
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	reqJSON := struct {
		UserRequest
		ExpiresIn *int `json:"expires_in_seconds"`
	}{}
	err := decoder.Decode(&reqJSON)
	if err != nil {
		log.Printf("error decoding request")
		respondWithError(w, http.StatusInternalServerError, "")
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(req.Context(), reqJSON.Email)
	if err != nil {
		log.Printf("error retrieving user")
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password")
		return
	}

	err = auth.CheckPasswordHash(reqJSON.Password, user.HashedPassword)
	if err != nil {
		log.Printf("wrong password")
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password")
		return
	}

	expiresIn := maxExpirationInSeconds

	if reqJSON.ExpiresIn != nil {
		expiresIn = *reqJSON.ExpiresIn
	}

	if expiresIn > maxExpirationInSeconds {
		expiresIn = maxExpirationInSeconds
	}

	tokenString, err := auth.MakeJWT(user.ID, cfg.secretString, time.Duration(expiresIn)*time.Second)
	if err != nil {
		log.Printf("error creating JWT: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		UserResponse
		Token string `json:"token"`
	}{
		UserResponse: UserResponse{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: tokenString,
	})

}
