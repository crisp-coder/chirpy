package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/crisp-coder/chirpy/internal/auth"
	"github.com/crisp-coder/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) AppHandler() http.Handler {
	fileserver := http.FileServer(http.Dir("public"))
	app_handler := middlewareLog(cfg.middlewareMetricsInc(fileserver))
	return app_handler
}

func (cfg *ApiConfig) PostUsersHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	temp_user := User{}
	err := decoder.Decode(&temp_user)

	if err != nil {
		log.Println("error reading request body: %w", err)
		sendErrorResponse(w, "error reading request")
		return
	}

	hashed_password, err := auth.HashPassword(temp_user.Password)
	if err != nil {
		log.Println("Error hashing password: %w", err)
		sendErrorResponse(w, "error hashing password")
		return
	}

	user, err := cfg.Db.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          temp_user.Email,
		HashedPassword: hashed_password,
	})

	if err != nil {
		log.Println("error creating user: %w", err)
		sendErrorResponse(w, "error creating user")
		return
	}

	api_user := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Password:  temp_user.Password,
	}

	dat, err := json.Marshal(api_user)
	if err != nil {
		log.Println("Error marshalling json: %w", err)
		sendErrorResponse(w, "error marshalling json")
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(dat)
	if err != nil {
		log.Println("Error writing response: %w", err)
		sendErrorResponse(w, "error writing response")
		return
	}
}

func (cfg *ApiConfig) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	loginParams := LoginParams{}
	err := decoder.Decode(&loginParams)
	if err != nil {
		log.Println("error decoding login parameters: %w", err)
		sendErrorResponse(w, err.Error())
		return
	}

	user, err := cfg.Db.GetUserByEmail(r.Context(), loginParams.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			sendUserNotFoundResponse(w)
		} else {
			log.Println(err)
		}
		return
	}

	err = auth.CheckPasswordHash(loginParams.Password, user.HashedPassword)
	if err != nil {
		log.Println("failed to match password: %w", err)
		sendIncorrectPasswordResponse(w)
		return
	}

	jwtExpiry := loginParams.ExpiresInSeconds
	if jwtExpiry > 3600 || jwtExpiry == 0 {
		jwtExpiry = 3600
	}

	token, err := auth.MakeJWT(user.ID, cfg.JWT_SECRET, time.Duration(jwtExpiry)*time.Second)
	if err != nil {
		log.Println("error creating jwt: %w", err)
		sendErrorResponse(w, "error creating jwt")
	}

	api_user := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Password:  loginParams.Password,
		Token:     token,
	}

	sendLoginAccepted(w, api_user)
}

func (cfg *ApiConfig) GetChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	chirpID_str := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpID_str)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, err.Error())
		return
	}
	chirp, err := cfg.Db.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			sendChirpNotFoundResponse(w)
		} else {
			log.Println(err)
			sendErrorResponse(w, err.Error())
		}
		return
	}

	api_Chirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	sendChirpResponse(w, api_Chirp)
}

func (cfg *ApiConfig) GetChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.Db.GetChirps(r.Context())
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, err.Error())
		return
	}

	api_Chirp := make([]Chirp, len(chirps))
	for i := range chirps {
		api_Chirp[i].ID = chirps[i].ID
		api_Chirp[i].CreatedAt = chirps[i].CreatedAt
		api_Chirp[i].UpdatedAt = chirps[i].UpdatedAt
		api_Chirp[i].Body = chirps[i].Body
		api_Chirp[i].UserID = chirps[i].UserID
	}
	sendChirpsResponse(w, api_Chirp)
}

func (cfg *ApiConfig) PostChirpsHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)

	if err != nil {
		log.Println("error posting chirp: %w", err)
		sendErrorResponse(w, "error posting chirp")
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

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println("error posting chirp: %w", err)
		sendErrorResponse(w, "error posting chirp")
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.JWT_SECRET)
	if err != nil {
		log.Println("error posting chirp: %w", err)
		sendErrorResponse(w, "error posting chirp")
		return
	}

	user, err := cfg.Db.GetUserByID(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			sendUserNotFoundResponse(w)
			return
		}
		log.Println("error posting chirp: %w", err)
		sendErrorResponse(w, "error posting chirp")
		return
	}

	saved_chirp, err := cfg.Db.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      cleaned_body,
		UserID:    user.ID,
	})

	if err != nil {
		log.Println("error posting chirp: %w", err)
		sendErrorResponse(w, "error posting chirp")
		return
	}

	sendCreatedChirpResponse(w, Chirp{
		ID:        saved_chirp.ID,
		CreatedAt: saved_chirp.CreatedAt,
		UpdatedAt: saved_chirp.UpdatedAt,
		Body:      saved_chirp.Body,
		UserID:    saved_chirp.UserID,
	})
}

func (cfg *ApiConfig) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Println(err)
	}
}

func (cfg *ApiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
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

func (cfg *ApiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits.Store(0)
	err := cfg.Db.ResetUsers(r.Context())
	if err != nil {
		log.Println("error resetting database.")
		sendErrorResponse(w, "error resetting database")
	}
	err = cfg.Db.ResetChirps(r.Context())
	if err != nil {
		log.Println("error resetting database.")
		sendErrorResponse(w, "error resetting database")
	}
}

func (cfg *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
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
