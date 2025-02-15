package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/mu7ammad1951/chirpy-boot/internal/database"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {

	responseData, err := cfg.dbQueries.GetChirps(req.Context())
	if err != nil {
		log.Printf("error fetching chirps: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var formattedResponseData []ChirpResponse

	for _, chirpResponse := range responseData {
		formattedResponseData = append(formattedResponseData, ChirpResponse{
			ID:        chirpResponse.ID,
			CreatedAt: chirpResponse.CreatedAt,
			UpdatedAt: chirpResponse.UpdatedAt,
			Body:      chirpResponse.Body,
			UserID:    chirpResponse.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, formattedResponseData)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, req *http.Request) {

	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		log.Printf("error parsing chirpID: %v\n", err)
		respondWithError(w, http.StatusBadRequest, "invalid chirp id")
		return
	}
	responseData, err := cfg.dbQueries.GetChirpByID(req.Context(), chirpID)
	if err != nil {
		log.Printf("error retrieving chirp: %v\n", err)
		respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}
	respondWithJSON(w, http.StatusOK, ChirpResponse{
		ID:        responseData.ID,
		CreatedAt: responseData.CreatedAt,
		UpdatedAt: responseData.UpdatedAt,
		Body:      responseData.Body,
		UserID:    responseData.UserID,
	})

}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {

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
		return
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
