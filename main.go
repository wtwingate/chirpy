package main

import (
	"log"
	"net/http"

	"github.com/wtwingate/chirpy/internal/database"
)

type apiConfig struct {
	db     *database.DB
	fsHits int
}

func newApiConfig(dbPath string) *apiConfig {
	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatal("could not establish database connection: ", err)
	}

	cfg := apiConfig{
		db:     db,
		fsHits: 0,
	}
	return &cfg
}

func main() {
	const root = "."
	const port = "8080"

	cfg := newApiConfig("./database.json")
	mux := http.NewServeMux()

	fsHandler := cfg.middlewareMetrics(http.StripPrefix("/app", http.FileServer(http.Dir(root))))

	mux.Handle("/app/*", fsHandler)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("/api/reset", cfg.handlerMetricsReset)
	mux.HandleFunc("POST /api/chirps", cfg.handlerPostChirps)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", root, port)
	log.Fatal(srv.ListenAndServe())
}
