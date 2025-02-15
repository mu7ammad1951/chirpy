package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mu7ammad1951/chirpy-boot/internal/database"
)

type ChirpRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	decoder := json.NewDecoder(req.Body)
	var chirpData ChirpRequest
	err := decoder.Decode(&chirpData)
	if err != nil {
		log.Printf("error decoding request: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	cleanedChirp, err := validateAndCleanChirp(chirpData.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	res, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleanedChirp,
		UserID: chirpData.UserID,
	})
	if err != nil {
		log.Printf("error creating chirp: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, ChirpResponse{
		ID:        res.ID,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		Body:      res.Body,
		UserID:    res.UserID,
	})

}

func validateAndCleanChirp(chirp string) (string, error) {
	if len(chirp) > 140 {
		log.Printf("bad request: chirp length > 140")
		return "", errors.New("chirp is too long - max char: 140")
	}

	return filter(chirp), nil
}

func filter(profaneString string) string {
	chirpWords := strings.Split(profaneString, " ")
	for i, word := range chirpWords {
		if profanity(word) {
			chirpWords[i] = "****"
		}
	}
	return strings.Join(chirpWords, " ")
}

func profanity(word string) bool {
	profaneWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	return profaneWords[strings.ToLower(word)]
}
