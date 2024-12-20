package main

import (
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	cfg := apiConfig{}
	loadConfig(&cfg)
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	serveMux.HandleFunc("GET /api/healthz", readinessHandler)
	serveMux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	serveMux.HandleFunc("POST /admin/reset", cfg.resetAppHandler)
	serveMux.HandleFunc("POST /api/chirps", cfg.middlewareLoggedInUser(cfg.createChirpHandler))
	serveMux.HandleFunc("GET /api/chirps", cfg.getChirpsHandler)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirpHandler)
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.middlewareLoggedInUser(cfg.deleteChirpHandler))
	serveMux.HandleFunc("POST /api/users", cfg.createUserHandler)
	serveMux.HandleFunc("PUT /api/users", cfg.middlewareLoggedInUser(cfg.updateUserHandler))
	serveMux.HandleFunc("POST /api/login", cfg.loginHandler)
	serveMux.HandleFunc("POST /api/refresh", cfg.refreshTokenHandler)
	serveMux.HandleFunc("POST /api/revoke", cfg.revokeTokenHandler)
	serveMux.HandleFunc("POST /api/polka/webhooks", cfg.polkaWebhookHandler)
	server := &http.Server{}
	server.Addr = ":8080"
	server.Handler = serveMux
	server.ListenAndServe()
}
