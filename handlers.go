package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/crisp-coder/chirpy/internal/data_models"
	"github.com/crisp-coder/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) AppHandler() http.Handler {
	fileserver := http.FileServer(http.Dir("public"))
	app_handler := middlewareLog(cfg.middlewareMetricsInc(fileserver))
	return app_handler
}

func (cfg *apiConfig) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	temp_user := data_models.User{}
	err := decoder.Decode(&temp_user)

	if err != nil {
		log.Println("error reading request body: %w", err)
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     temp_user.Email,
	})

	if err != nil {
		log.Println("error creating user: %w", err)
	}

	created_user := data_models.User{
		Id:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	}

	dat, err := json.Marshal(created_user)
	if err != nil {
		log.Println("Error marshalling json: %w", err)
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(dat)
	if err != nil {
		log.Println("Error writing response: %w", err)
	}
}

func (cfg *apiConfig) GetChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	chirpID_str := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpID_str)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, err.Error())
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			sendChirpNotFoundResponse(w)
			return
		}
		log.Println(err)
		sendErrorResponse(w, err.Error())
		return
	}

	data_model_chirp := data_models.Chirp{}
	data_model_chirp.Id = chirp.ID
	data_model_chirp.CreatedAt = chirp.CreatedAt
	data_model_chirp.UpdatedAt = chirp.UpdatedAt
	data_model_chirp.Body = chirp.Body
	data_model_chirp.UserID = chirp.UserID
	sendChirpResponse(w, data_model_chirp)
}

func (cfg *apiConfig) GetChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, err.Error())
		return
	}

	data_model_chirps := make([]data_models.Chirp, len(chirps))
	for i := range chirps {
		data_model_chirps[i].Id = chirps[i].ID
		data_model_chirps[i].CreatedAt = chirps[i].CreatedAt
		data_model_chirps[i].UpdatedAt = chirps[i].UpdatedAt
		data_model_chirps[i].Body = chirps[i].Body
		data_model_chirps[i].UserID = chirps[i].UserID
	}
	sendChirpsResponse(w, data_model_chirps)
}

func (cfg *apiConfig) PostChirpsHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	chirp := data_models.Chirp{}
	err := decoder.Decode(&chirp)

	if err != nil {
		log.Println(err)
		sendErrorResponse(w, err.Error())
		return
	}

	if len(chirp.Body) > 140 {
		sendChirpTooLong(w)
		return
	}

	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	cleaned_body := StripBadWords(chirp.Body, "****", badWords)

	saved_chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      cleaned_body,
		UserID:    chirp.UserID,
	})

	if err != nil {
		log.Println(err)
		sendErrorResponse(w, err.Error())
		return
	}

	sendCreatedChirpResponse(w, data_models.Chirp{
		Id:        saved_chirp.ID,
		CreatedAt: saved_chirp.CreatedAt,
		UpdatedAt: saved_chirp.UpdatedAt,
		Body:      saved_chirp.Body,
		UserID:    saved_chirp.UserID,
	})
}

func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Println(err)
	}
}

func (cfg *apiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
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

func (cfg *apiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits.Store(0)
	err := cfg.db.Reset(r.Context())
	if err != nil {
		log.Println("Error resetting database.")
	}
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
