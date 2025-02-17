package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mu7ammad1951/chirpy-boot/internal/auth"
	"github.com/mu7ammad1951/chirpy-boot/internal/database"
)

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	reqJSON := UserRequest{}

	err := decoder.Decode(&reqJSON)
	if err != nil {
		log.Printf("error decoding request: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "")
		return
	}
	hashedPassword, err := auth.HashPassword(reqJSON.Password)
	if err != nil {
		log.Printf("error hashing password")
		respondWithError(w, http.StatusInternalServerError, "error creating user, try again")
		return
	}
	userData, err := cfg.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          reqJSON.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("error creating user: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "")
		return
	}

	respondWithJSON(w, http.StatusCreated, UserResponse{
		ID:        userData.ID,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
		Email:     userData.Email,
	})
}

func (cfg *apiConfig) handlerUserUpdate(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var reqJSON UserRequest
	if err := decoder.Decode(&reqJSON); err != nil {
		log.Printf("error decoding request: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}
	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("error getting token: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}
	userID, err := auth.ValidateJWT(accessToken, cfg.secretString)
	if err != nil {
		log.Printf("error validating JWT: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}

	hashedPassword, err := auth.HashPassword(reqJSON.Password)
	if err != nil {
		log.Printf("error hashing password: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "")
		return
	}

	userInfo, err := cfg.dbQueries.UpdatePasswordEmail(req.Context(), database.UpdatePasswordEmailParams{
		ID:             userID,
		Email:          reqJSON.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("error updating password or email: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "")
		return
	}

	respondWithJSON(w, http.StatusOK, UserResponse{
		ID:        userInfo.ID,
		CreatedAt: userInfo.CreatedAt,
		UpdatedAt: userInfo.UpdatedAt,
		Email:     userInfo.Email,
	})

}
