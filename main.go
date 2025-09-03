package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
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

	cfg := apiConfig{}

	// Redirect for /
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/", http.StatusFound)
	})

	mux.Handle("GET /app/", http.StripPrefix("/app", cfg.appHandler()))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.resetHandler)

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
