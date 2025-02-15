package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, fmt.Sprintf(`
	<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
	</html>`,
		cfg.fileserverHits.Load()))

}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		log.Printf("invalid request: permission denied")
		respondWithError(w, http.StatusForbidden, "Permission denied")
		return
	}

	err := cfg.dbQueries.ResetUsers(req.Context())
	if err != nil {
		log.Printf("error resetting table: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "")
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("State Reset"))
	cfg.fileserverHits.Store(0)
}
