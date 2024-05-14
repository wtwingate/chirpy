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

	fileSrv := http.StripPrefix("/app", http.FileServer(http.Dir(root)))
	mux.Handle("/app/*", cfg.middlewareMetrics(fileSrv))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", cfg.handlerMetrics)
	mux.HandleFunc("/reset", cfg.handlerMetricsReset)

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
	cfg.fileserverHits++
	return next
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
