package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/alexissimonian/test/bootdev/chirpy/internal/database"
	"github.com/google/uuid"
)

type ChirpCreateRequest struct {
	Body string `json:"body"`
}

type ChirpCreateResponse struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type ChirpResponse struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
}

func (cfg *apiConfig) createChirpHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	request := ChirpCreateRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		log.Printf("Error parsing request: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = validateChirp(&request)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}

	userUUID, err := uuid.Parse(r.Header.Get("userID"))
	if err != nil {
        log.Printf("Invalid userID %v, userid: %v\n", err, userUUID)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error: Invalid user id."))
		return
	}

	chirp, err := cfg.database.CreateChirp(r.Context(), database.CreateChirpParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Body: request.Body,
		UserID: userUUID,
	})

	chirpResponse := ChirpCreateResponse{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}

	responseData, err := json.Marshal(&chirpResponse)
	if err != nil {
		log.Printf("Something went wrong encoding chirp into json: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write(responseData)
}

func (cfg *apiConfig) getChirpsHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	chirps, err := cfg.database.GetChirps(r.Context())
	if err != nil {
		log.Printf("Something went wrong getting chirps: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	chirpsResponse := []ChirpResponse{}
	for _, chrip := range chirps {
		chirpsResponse = append(chirpsResponse, ChirpResponse{
			ID: chrip.ID,
			CreatedAt: chrip.CreatedAt,
			UpdatedAt: chrip.UpdatedAt,
			Body: chrip.Body,
		})
	}
	chirpsResponseData, err := json.Marshal(&chirpsResponse)
	if err != nil {
		log.Printf("Something went wrong encoding chirps into json: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(chirpsResponseData)
}

func (cfg *apiConfig) getChirpHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Error: Invalid chirp id."))
		return
	}

	chirp, err := cfg.database.GetChirp(r.Context(), chirpID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("Error: Chirp not found."))
			return
		}

		log.Printf("Something went wrong getting chirp: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	chirpResponse := ChirpResponse{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
	}

	chirpResponseData, err := json.Marshal(&chirpResponse)
	if err != nil {
		log.Printf("Something went wrong encoding chirp into json: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(chirpResponseData)
}

func validateChirp(r *ChirpCreateRequest) error {
	err := validateChirpLength(r)
	if err != nil {
		log.Printf("Error validating chirp length: %v\n", err)
		return err
	}

	currateChirpContent(r)
	return nil
}

func validateChirpLength(r *ChirpCreateRequest) error {
	if len(r.Body) < 1 {
		return fmt.Errorf("Incorect request. No property body found")
	}

	if len(r.Body) > 140 {
		return fmt.Errorf("Chirp is too long")
	}

	return nil
}

func currateChirpContent(r *ChirpCreateRequest) {
	bannedWords := [...]string{"kerfuffle", "fornax", "sharbert"}
	for _, word := range bannedWords {
		if strings.Contains(strings.ToLower(r.Body), word) {
			regexpPattern := fmt.Sprintf("(?i)%v", word)
			regexp := regexp.MustCompile(regexpPattern)
			r.Body = regexp.ReplaceAllString(r.Body, "****")
		}
	}
}
