package api

import (
	"net/http"

	_ "github.com/lib/pq"
)

func MakeServer(api_cfg *ApiConfig) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/", http.StatusFound)
	})

	mux.Handle("GET /app/", http.StripPrefix("/app", api_cfg.AppHandler()))
	mux.HandleFunc("GET /admin/metrics", api_cfg.MetricsHandler)
	mux.HandleFunc("POST /admin/reset", api_cfg.ResetHandler)
	mux.HandleFunc("GET /api/healthz", api_cfg.ReadinessHandler)
	mux.HandleFunc("POST /api/users", api_cfg.PostUsersHandler)
	mux.HandleFunc("PUT /api/users", api_cfg.PutUsersHandler)
	mux.HandleFunc("POST /api/login", api_cfg.PostLoginHandler)
	mux.HandleFunc("POST /api/refresh", api_cfg.PostRefreshHandler)
	mux.HandleFunc("POST /api/revoke", api_cfg.PostRevokeHandler)
	mux.HandleFunc("POST /api/chirps", api_cfg.PostChirpsHandler)
	mux.HandleFunc("GET /api/chirps", api_cfg.GetChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", api_cfg.GetChirpByIDHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return &server
}
