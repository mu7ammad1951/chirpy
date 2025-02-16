package main

import (
	"log"
	"net/http"
	"time"

	"github.com/mu7ammad1951/chirpy-boot/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	clientRefreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("error getting beared token: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}

	databaseRefreshToken, err := cfg.dbQueries.GetRefreshToken(req.Context(), clientRefreshToken)
	if err != nil {
		log.Printf("error retrieving refresh token: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}

	if (databaseRefreshToken.ExpiresAt.Compare(time.Now().UTC())) < 0 || (databaseRefreshToken.RevokedAt.Valid) {
		log.Printf("expired token or revoked token: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}

	token, err := auth.MakeJWT(databaseRefreshToken.UserID, cfg.secretString)
	if err != nil {
		log.Printf("error creating JWT: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	clientRefreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("error getting beared token: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}

	err = cfg.dbQueries.UpdateRefreshToken(req.Context(), clientRefreshToken)
	if err != nil {
		log.Printf("failed to revoke token: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "failed to revoke refresh token")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
