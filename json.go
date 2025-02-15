package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type response_err struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, errorStatus int, errorString string) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	res, err := json.Marshal(response_err{Error: errorString})
	if err != nil {
		log.Printf("error marshalling 'error': %v", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(errorStatus)
	w.Write(res)
}

func respondWithJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	res, err := json.Marshal(v)
	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "")
		return
	}
	w.WriteHeader(status)
	w.Write(res)
}
