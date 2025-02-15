package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mu7ammad1951/chirpy-boot/internal/auth"
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
		respondWithError(w, http.StatusInternalServerError, "incorrect email or password")
		return
	}

	err = auth.CheckPasswordHash(reqJSON.Password, user.HashedPassword)
	if err != nil {
		log.Printf("wrong password")
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password")
		return
	}

	respondWithJSON(w, http.StatusOK, UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}
