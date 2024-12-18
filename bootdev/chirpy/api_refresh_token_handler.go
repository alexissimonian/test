package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexissimonian/test/bootdev/chirpy/internal/auth"
	"github.com/alexissimonian/test/bootdev/chirpy/internal/database"
)

type refreshTokenResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) refreshTokenHandler(rw http.ResponseWriter, r *http.Request) {
	r.Header.Add("Content-Type", "application/json")
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting authoriztion from header: %v\n", err)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Unauthorized request"))
		return
	}

	refreshToken, err := cfg.database.GetRefreshToken(r.Context(), bearerToken)
	if refreshToken.RevokedAt.Valid {
		log.Printf("Revoked refresh token")
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Unauthorized request"))
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now().UTC()) {
		log.Printf("Expired refresh token")
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Unauthorized request"))
		return
	}

	expirationDurationToken, err := time.ParseDuration("3600s")
	if err != nil {
		log.Printf("Error parsing default expiration duration for token: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error parsing default expiration date for token."))
		return
	}

	token, err := auth.MakeJWT(refreshToken.UserID, cfg.serverSecret, expirationDurationToken)
	if err != nil {
		log.Printf("Error generating token for user : %v\n", refreshToken.UserID)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error generating identification token"))
		return
	}

	refreshTokenResponse := refreshTokenResponse{Token: token}
	refreshTokenResponseData, err := json.Marshal(&refreshTokenResponse)
	if err != nil {
		log.Printf("Error encoding token into a json response: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error encoding response"))
		return
	}

	rw.Write(refreshTokenResponseData)
}

func (cfg *apiConfig) revokeTokenHandler(rw http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting authoriztion from header: %v\n", err)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Unauthorized request"))
		return
	}

	err = cfg.database.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		Token:     bearerToken,
		RevokedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		UpdatedAt: time.Now().UTC(),
	})

    if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized request"))
			return
		}
        log.Printf("Error revoking a refresh token: %v\n", err)
        rw.WriteHeader(http.StatusInternalServerError)
        rw.Write([]byte("Problem revoking refresh token"))
        return
    }

    rw.WriteHeader(http.StatusNoContent)
}
