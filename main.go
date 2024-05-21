package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/wtwingate/chirpy/internal/database"
)

const dbPath = "./database.json"

type apiConfig struct {
	DB        *database.DB
	jwtSecret string
	apiKey    string
	fsHits    int
}

func newApiConfig(dbPath, jwtSecret, apiKey string) *apiConfig {
	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatal("could not establish database connection: ", err)
	}

	cfg := apiConfig{
		DB:        db,
		jwtSecret: jwtSecret,
		apiKey:    apiKey,
		fsHits:    0,
	}
	return &cfg
}

func main() {
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	apiKey := os.Getenv("API_KEY")

	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *debug {
		os.Remove(dbPath)
	}

	const root = "."
	const port = "8080"

	cfg := newApiConfig(dbPath, jwtSecret, apiKey)
	mux := http.NewServeMux()

	fsHandler := cfg.middlewareMetrics(http.StripPrefix("/app", http.FileServer(http.Dir(root))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", cfg.handlerReset)

	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirpByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handlerDeleteChirp)

	mux.HandleFunc("POST /api/users", cfg.handlerNewUser)
	mux.HandleFunc("PUT /api/users", cfg.handlerUpdateUser)

	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)

	mux.HandleFunc("POST /api/polka/webhooks", cfg.handlerPolkaWebhooks)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", root, port)
	log.Fatal(srv.ListenAndServe())
}
