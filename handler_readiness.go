package main

import (
	"io"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	// Add header
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "200 OK\n")
}
