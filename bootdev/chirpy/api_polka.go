package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/alexissimonian/test/bootdev/chirpy/internal/auth"
	"github.com/google/uuid"
)

type polkaWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) polkaWebhookHandler(rw http.ResponseWriter, r *http.Request) {
	apikey, err := auth.GetApiKey(r.Header)
	if err != nil {
		log.Printf("Apikey not found: %v\n", err)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Unauthorized request"))
		return
	}

	if apikey != cfg.polkaApiKey {
		log.Printf("Incorrect API key: %v\n", err)
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Unauthorized request"))
		return
	}

	request := polkaWebhookRequest{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)
	if err != nil {
		log.Printf("Problem when decoding Polka webhook request: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Couldn't decode request"))
		return
	}

	if request.Event != "user.upgraded" {
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	user, err := cfg.database.GetUserById(r.Context(), request.Data.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			log.Printf("User: %v not found: %v\n", request.Data.UserID, err)
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("User not found"))
			return
		}

		log.Printf("Something went wrong retreiving user: %v, err: %v\n", request.Data.UserID, err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Could not retreive user"))
		return
	}

	err = cfg.database.UpgradeUser(r.Context(), user.ID)
	if err != nil {
		log.Printf("Error when upgrading user: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error when upgrading user"))
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
