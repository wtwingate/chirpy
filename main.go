package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/wtwingate/chirpy/internal/database"
)

const dbPath = "./database.json"

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
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *debug {
		os.Remove(dbPath)
	}

	const root = "."
	const port = "8080"

	cfg := newApiConfig(dbPath)
	mux := http.NewServeMux()

	fsHandler := cfg.middlewareMetrics(http.StripPrefix("/app", http.FileServer(http.Dir(root))))

	mux.Handle("/app/*", fsHandler)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("/api/reset", cfg.handlerMetricsReset)
	mux.HandleFunc("POST /api/chirps", cfg.handlerNewChirp)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirpByID)
	mux.HandleFunc("POST /api/users", cfg.handlerNewUser)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", root, port)
	log.Fatal(srv.ListenAndServe())
}
