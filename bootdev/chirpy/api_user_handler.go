package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexissimonian/test/bootdev/chirpy/internal/database"
	"github.com/google/uuid"
)

type createUserRequest struct {
	Email string `json:"email"`
}

type createUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUserHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	request := createUserRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		log.Printf("Error parsing request: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(request.Email) == 0 {
		log.Println("Incorect request. No property email found")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := cfg.database.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     request.Email,
	})

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			log.Printf("Something went wrong creating user: %v\n", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("Something went wrong creating user: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	userCreated := createUserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	responseData, err := json.Marshal(&userCreated)
	if err != nil {
		log.Printf("Something went wrong encoding user into json: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write(responseData)
}