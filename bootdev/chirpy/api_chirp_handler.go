package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type chirpRequest struct {
	Body string `json:"body"`
}

type chirpErrorResponse struct {
	Error string `json:"error"`
}

type chirpCleanResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

func validateChirpHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	request := chirpRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		log.Printf("Error parsing request: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(request.Body) == 0 {
		log.Println("Incorect request. No property body found")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(request.Body) > 140 {
		errorResponse := chirpErrorResponse{Error: "Chirp is too long"}
		data, err := json.Marshal(&errorResponse)
		if err != nil {
			log.Printf("Something went wrong encoding error into json: %v\n", err)
			rw.WriteHeader(http.StatusInternalServerError)
		}
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write(data)
		return
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