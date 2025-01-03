package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/alexissimonian/test/bootdev/chirpy/internal/database"
	"github.com/google/uuid"
)

type chirpCreateRequest struct {
	Body string `json:"body"`
}

type chirpCreateResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
}

func (cfg *apiConfig) createChirpHandler(rw http.ResponseWriter, r *http.Request) {
	request := chirpCreateRequest{}
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
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Body:      request.Body,
		UserID:    userUUID,
	})

	chirpResponse := chirpCreateResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	responseData, err := json.Marshal(&chirpResponse)
	if err != nil {
		log.Printf("Something went wrong encoding chirp into json: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(responseData)
}

func (cfg *apiConfig) getChirpsHandler(rw http.ResponseWriter, r *http.Request) {
	authorIDString := r.URL.Query().Get("author_id")
	sortQueryParameter := r.URL.Query().Get("sort")
	var authorID uuid.UUID
	var err error
	if len(authorIDString) > 0 {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			log.Printf("Error parsing author id into uuid: %v\n", err)
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("Error author_id incorrect format"))
			return
		}
	}

	var chirps []database.Chirp
	if authorID != uuid.Nil {
		chirps, err = cfg.database.GetChirpsByUser(r.Context(), authorID)
	} else {
		chirps, err = cfg.database.GetChirps(r.Context())
	}

	if sortQueryParameter == "asc" {
		sortChirpsAsc(chirps)
	} else if sortQueryParameter == "desc" {
		sortChirpsDesc(chirps)
	}

	if err != nil {
		log.Printf("Something went wrong getting chirps: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	chirpsResponse := []chirpResponse{}
	for _, chrip := range chirps {
		chirpsResponse = append(chirpsResponse, chirpResponse{
			ID:        chrip.ID,
			CreatedAt: chrip.CreatedAt,
			UpdatedAt: chrip.UpdatedAt,
			Body:      chrip.Body,
		})
	}
	chirpsResponseData, err := json.Marshal(&chirpsResponse)
	if err != nil {
		log.Printf("Something went wrong encoding chirps into json: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(chirpsResponseData)
}

func (cfg *apiConfig) getChirpHandler(rw http.ResponseWriter, r *http.Request) {
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

	chirpResponse := chirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
	}

	chirpResponseData, err := json.Marshal(&chirpResponse)
	if err != nil {
		log.Printf("Something went wrong encoding chirp into json: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(chirpResponseData)
}

func validateChirp(r *chirpCreateRequest) error {
	err := validateChirpLength(r)
	if err != nil {
		log.Printf("Error validating chirp length: %v\n", err)
		return err
	}

	currateChirpContent(r)
	return nil
}

func validateChirpLength(r *chirpCreateRequest) error {
	if len(r.Body) < 1 {
		return fmt.Errorf("Incorect request. No property body found")
	}

	if len(r.Body) > 140 {
		return fmt.Errorf("Chirp is too long")
	}

	return nil
}

func currateChirpContent(r *chirpCreateRequest) {
	bannedWords := [...]string{"kerfuffle", "fornax", "sharbert"}
	for _, word := range bannedWords {
		if strings.Contains(strings.ToLower(r.Body), word) {
			regexpPattern := fmt.Sprintf("(?i)%v", word)
			regexp := regexp.MustCompile(regexpPattern)
			r.Body = regexp.ReplaceAllString(r.Body, "****")
		}
	}
}

func (cfg *apiConfig) deleteChirpHandler(rw http.ResponseWriter, r *http.Request) {
	userIDString := r.Header.Get("userID")
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		log.Printf("Error parsing user id into uuid: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error parsing user id"))
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Error parsing chirpId from request path: %v\n", err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Chirp id not recognised"))
		return
	}

	chirp, err := cfg.database.GetChirp(r.Context(), chirpID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			log.Printf("Chirp not found")
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("Not found"))
			return
		}

		log.Printf("Error when getting chirp: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error when getting chirp"))
		return
	}

	if userID != chirp.UserID {
		log.Printf("Forbidden deletion tentative. User: %v tried to delete chirp: %v\n", userID, chirp.ID)
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte("Forbidden request"))
		return
	}

	err = cfg.database.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirp.ID,
		UserID: userID,
	})

	if err != nil {
		log.Printf("Something went wrong deleting chirp: %v, %v\n", chirp.ID, err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error deleting chirp"))
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}


func sortChirpsAsc(chirps []database.Chirp) {
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})
}

func sortChirpsDesc(chirps []database.Chirp) {
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
	})
}
