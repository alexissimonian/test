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
	serveMux.HandleFunc("POST /admin/reset", cfg.resetMetricsHandler)
	serveMux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)
	serveMux.HandleFunc("POST /api/users", cfg.createUserHandler)
	server := &http.Server{}
	server.Addr = ":8080"
	server.Handler = serveMux
	server.ListenAndServe()
}
