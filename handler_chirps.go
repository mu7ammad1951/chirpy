package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/mu7ammad1951/chirpy-boot/internal/auth"
	"github.com/mu7ammad1951/chirpy-boot/internal/database"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {

	var responseData []database.Chirp
	var err error
	if req.URL.Query().Has("author_id") {
		queryUserID, err := uuid.Parse(req.URL.Query().Get("author_id"))
		if err != nil {
			log.Printf("invalid author_id: %v", err)
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		responseData, err = cfg.dbQueries.GetChirpsByUserID(req.Context(), queryUserID)
		if err != nil {
			log.Printf("error fetching chirps: %v", err)
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		responseData, err = cfg.dbQueries.GetChirps(req.Context())
		if err != nil {
			log.Printf("error fetching chirps: %v", err)
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
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

	switch req.URL.Query().Get("sort") {
	case "asc":
		sort.Slice(formattedResponseData, func(i, j int) bool {
			return formattedResponseData[i].CreatedAt.Before(formattedResponseData[j].CreatedAt)
		})
	case "desc":
		sort.Slice(formattedResponseData, func(i, j int) bool {
			return formattedResponseData[i].CreatedAt.After(formattedResponseData[j].CreatedAt)
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

	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("error getting bearer token: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "permission denied")
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.secretString)
	if err != nil {
		log.Printf("error validating token: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "permission denied")
		return
	}

	res, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleanedChirp,
		UserID: userID,
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

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		log.Printf("error parsing chirpID: %v\n", err)
		respondWithError(w, http.StatusBadRequest, "invalid chirp id")
		return
	}

	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("error getting bearer token: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "permission denied")
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.secretString)
	if err != nil {
		log.Printf("error validating token: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "permission denied")
		return
	}

	chirpInfo, err := cfg.dbQueries.GetChirpByID(req.Context(), chirpID)
	if err != nil {
		log.Printf("error retrieving chirp: %v\n", err)
		respondWithError(w, http.StatusNotFound, "could not find chirp")
		return
	}

	if chirpInfo.UserID != userID {
		log.Printf("unauthorized: %v\n", err)
		respondWithError(w, http.StatusForbidden, "you are not authorized to delete this chirp")
		return
	}

	err = cfg.dbQueries.DeleteChirpByID(req.Context(), chirpID)
	if err != nil {
		log.Printf("chirp not found: %v", err)
		respondWithError(w, http.StatusNotFound, "chirp does not exist")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
