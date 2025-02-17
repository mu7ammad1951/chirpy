package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mu7ammad1951/chirpy-boot/internal/auth"
	"github.com/mu7ammad1951/chirpy-boot/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	reqJSON := UserRequest{}
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

	tokenString, err := auth.MakeJWT(user.ID, cfg.secretString)
	if err != nil {
		log.Printf("error creating JWT: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("error creating refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = cfg.dbQueries.AddRefreshToken(req.Context(), database.AddRefreshTokenParams{
		Token:  refreshTokenString,
		UserID: user.ID,
	})
	if err != nil {
		log.Printf("error adding refresh token to database: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		UserResponse
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		UserResponse: UserResponse{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        tokenString,
		RefreshToken: refreshTokenString,
	})

}
