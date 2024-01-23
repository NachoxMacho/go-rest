package api

import (
	"fmt"
	"net/http"
)

type APIConfig struct {
	FileServerHits int64
}

func (cfg *APIConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHits += 1
		next.ServeHTTP(w, r)
	})
}
func (cfg *APIConfig) Reset(w http.ResponseWriter, r *http.Request) {
	cfg.FileServerHits = 0
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *APIConfig) Metrics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited " + fmt.Sprint(cfg.FileServerHits) + " times!</p></body></html>"))
	w.Header().Set("Content-Type", "text/html")
	// w.WriteHeader(200)
}
