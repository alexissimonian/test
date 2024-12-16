package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type ChirpCreateRequest struct {
	Body string `json:"body"`
	UserID string `json:"user_id"`
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
		
		rw.Write([]byte(err.Error()))
	}
}

func validateChirp(r *ChirpCreateRequest) error {
	err := validateChirpLength(r)
	if err != nil {
		log.Printf("Error validating chirp length: %v\n", err)
		return err
	}

	
	bannedWords := [...]string{"kerfuffle", "fornax", "sharbert"}
	for _, word := range bannedWords {
		if strings.Contains(strings.ToLower(request.Body), word) {
			regexpPattern := fmt.Sprintf("(?i)%v", word)
			regexp := regexp.MustCompile(regexpPattern)
			request.Body = regexp.ReplaceAllString(request.Body, "****")
		}
	}

	curratedResponse := chirpCleanResponse{CleanedBody: request.Body}
	data, err := json.Marshal(&curratedResponse)
	if err != nil {
		log.Printf("Something went wrong when encoding response for currated content: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Write(data)
	return
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