package main

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
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
