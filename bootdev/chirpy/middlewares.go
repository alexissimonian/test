package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		fmt.Printf("value: %v\n", cfg.fileServerHits.Load())
		next.ServeHTTP(rw, r)
	})
}

