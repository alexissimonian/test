package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alexissimonian/test/bootdev/chirpy/internal/auth"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		fmt.Printf("value: %v\n", cfg.fileServerHits.Load())
		next.ServeHTTP(rw, r)
	})
}

func (cfg *apiConfig) middlewareLoggedInUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Printf("Error getting bearer token : %v\n", err)
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized request"))
			return
		}

		userID, err := auth.ValidateJWT(token, cfg.serverSecret)
		if err != nil {
			log.Printf("Unauthorized token: %v\n", err)
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized request"))
            return
		}

        r.Header.Del("userID")
		r.Header.Add("userID", userID.String())
		next.ServeHTTP(rw, r)
	})
}
