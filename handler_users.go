package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	reqJSON := struct {
		Email string `json:"email"`
	}{}

	err := decoder.Decode(&reqJSON)
	if err != nil {
		log.Printf("error decoding request: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "")
		return
	}

	userData, err := cfg.dbQueries.CreateUser(req.Context(), reqJSON.Email)
	if err != nil {
		log.Printf("error creating user: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "")
		return
	}

	respondWithJSON(w, http.StatusCreated, struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}{
		ID:        userData.ID,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
		Email:     userData.Email,
	})
}
