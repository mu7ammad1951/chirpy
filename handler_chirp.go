package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type chirp struct {
	Body string `json:"body"`
}

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	decoder := json.NewDecoder(req.Body)
	var chirpData chirp
	err := decoder.Decode(&chirpData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "")
		return
	}

	if len(chirpData.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirpData.Body = filter(chirpData.Body)

	respondWithJSON(w, http.StatusOK, struct {
		CleanedBody string `json:"cleaned_body"`
	}{
		CleanedBody: chirpData.Body,
	})
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
