package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) appHandler() http.Handler {
	fileserver := http.FileServer(http.Dir("public"))
	app_handler := middlewareLog(cfg.middlewareMetricsInc(fileserver))
	return app_handler
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		sendErrorResponse(w, err.Error())
		return
	}

	if len(params.Body) > 140 {
		sendChirpTooLong(w)
		return
	}

	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	cleaned_body := StripBadWords(params.Body, "****", badWords)
	sendCleanedResponse(w, cleaned_body)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Println(err)
	}
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	body := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`,
		cfg.FileserverHits.Load())

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	bytes, err := w.Write([]byte(body))

	if err != nil {
		fmt.Printf("%d bytes written\n", bytes)
		fmt.Println(err)
	}
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits.Store(0)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
