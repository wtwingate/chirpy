package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	const root = "."
	const port = "8080"

	cfg := newApiConfig()
	mux := http.NewServeMux()

	fileServer := http.StripPrefix("/app", http.FileServer(http.Dir(root)))
	mux.Handle("/app/*", cfg.middlewareMetrics(fileServer))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/metrics", cfg.handlerMetrics)
	mux.HandleFunc("/api/reset", cfg.handlerMetricsReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", root, port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plan; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

type apiConfig struct {
	fileserverHits int
}

func newApiConfig() *apiConfig {
	cfg := apiConfig{
		fileserverHits: 0,
	}
	return &cfg
}

func (cfg *apiConfig) middlewareMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsReport() []byte {
	report := fmt.Sprintf("Hits: %v", strconv.Itoa(cfg.fileserverHits))
	return []byte(report)
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(cfg.metricsReport())
}

func (cfg *apiConfig) handlerMetricsReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits = 0
	w.Write(cfg.metricsReport())
}
