package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/crisp-coder/chirpy/internal/config"
	"github.com/crisp-coder/chirpy/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DB_URL)
	if err != nil {
		fmt.Println(err)
	}
	dbQueries := database.New(db)

	api_cfg := apiConfig{
		db: dbQueries,
	}

	logFile, err := setupLogging("application.log")
	if err != nil {
		log.Fatalf("Failed to set up logging: %v", err)
	}
	defer func() {
		err := logFile.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error closing log file")
		}
	}()
	log.Println("log start")

	// Redirect for /
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/", http.StatusFound)
	})

	mux.Handle("GET /app/", http.StripPrefix("/app", api_cfg.appHandler()))
	mux.HandleFunc("GET /admin/metrics", api_cfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", api_cfg.resetHandler)
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("POST /api/users", api_cfg.registerUserHandler)
	mux.HandleFunc("POST /api/chirps", api_cfg.addChirpHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Run server
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
