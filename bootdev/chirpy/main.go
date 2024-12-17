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
	serveMux.HandleFunc("POST /api/chirps", cfg.createChirpHandler)
	serveMux.HandleFunc("GET /api/chirps", cfg.getChirpsHandler)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirpHandler)
	serveMux.HandleFunc("POST /api/users", cfg.createUserHandler)
	serveMux.HandleFunc("POST /api/login", cfg.loginHandler)
	server := &http.Server{}
	server.Addr = ":8080"
	server.Handler = serveMux
	server.ListenAndServe()
}
