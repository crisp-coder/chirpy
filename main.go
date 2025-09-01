package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

func setupLogging(logFilePath string) (*os.File, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	return file, nil
}

type apiConfig struct {
	FileserverHits atomic.Int32
}

func main() {
	// Log file setup
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

	// API metric config setup
	cfg := apiConfig{}

	// Redirect for /
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("root redirect firing...\n")
		http.Redirect(w, r, "/app/", http.StatusFound)
	})

	// Handlers
	mux.Handle("GET /app/", http.StripPrefix("/app", middlewareLog(cfg.middlewareMetricsInc(http.FileServer(http.Dir("public"))))))
	mux.HandleFunc("GET /healthz", readinessHandler)
	mux.HandleFunc("GET /metrics", cfg.countHandler)
	mux.HandleFunc("POST /reset", cfg.resetHandler)

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

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Println(err)
	}
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("resetHandler firing...\n")
	cfg.FileserverHits.Store(0)
}

func (cfg *apiConfig) countHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "Hits: %d", cfg.FileserverHits.Load())
	if err != nil {
		fmt.Println(err)
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
