package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexissimonian/test/bootdev/chirpy/internal/auth"
	"github.com/alexissimonian/test/bootdev/chirpy/internal/database"
	"github.com/google/uuid"
)

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createUserResponse struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type loginResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
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

	if len(request.Email) == 0 || len(request.Password) == 0 {
		log.Println("Incorect request. Please provide email + password")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Incorect request. Please provide email + password"))
		return
	}

	passwordHash, err := auth.HashPassword(request.Password)
	if err != nil {
		log.Printf("Error hashing password: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	user, err := cfg.database.CreateUser(r.Context(), database.CreateUserParams{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        request.Email,
		PasswordHash: passwordHash,
	})

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			log.Printf("Something went wrong creating user: %v\n", err)
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("User already exists"))
			return
		}

		log.Printf("Something went wrong creating user: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	userCreated := createUserResponse{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
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

func (cfg *apiConfig) loginHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	request := loginUserRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		log.Printf("Error parsing request: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := cfg.database.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("Error: User not found."))
			return
		}

		log.Printf("Something went wrong getting user: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	if auth.CheckPasswordHash(request.Password, user.PasswordHash) != true {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Incorrect email or password"))
		return
	}

	expiresInSeconds, err := time.ParseDuration("3600s")
	if err != nil {
		log.Printf("Error parsing basic duration for expiration token %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error parsing default duration for expiration token"))
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.serverSecret, expiresInSeconds)
	if err != nil {
		log.Printf("Error generating token for user : %v\n", user.Email)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error generating identification token"))
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error when generating refresh token: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error generating refresh token"))
		return
	}

	refreshTokenExpirationDuration, err := time.ParseDuration("1440h")
	if err != nil {
		log.Printf("Problem parsing basic token duration: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error setting expiration date on refresh token"))
		return
	}

	refreshToken, err = cfg.database.CreateRefeshToken(r.Context(), database.CreateRefeshTokenParams{
		Token:     refreshToken,
		ExpiresAt: time.Now().UTC().Add(refreshTokenExpirationDuration),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		RevokedAt: sql.NullTime{},
	})

	loginResponse := loginResponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	}

	responseData, err := json.Marshal(&loginResponse)
	if err != nil {
		log.Printf("Something went wrong encoding user into json: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Write(responseData)
}

func (cfg *apiConfig) resetUsers(r *http.Request) error {
	if cfg.platform != "dev" {
		fmt.Println("platform: ", cfg.platform)
		return fmt.Errorf("Cannot reset users in prod !")
	}

	err := cfg.database.ResetUsers(r.Context())
	if err != nil {
		log.Printf("Something went wrong resetting users: %v\n", err)
	}

	return nil
}

type updateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) updateUserHandler(rw http.ResponseWriter, r *http.Request) {
	userIDString := r.Header.Get("userID")
	request := updateUserRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		log.Printf("Error decoding request: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error when decoding request"))
		return
	}

	if len(request.Email) == 0 || len(request.Password) == 0 {
		log.Printf("Invalid email: %v\n, or password: %v\n", request.Email, request.Password)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Please provide valid email or password"))
	}

	passwordHash, err := auth.HashPassword(request.Password)
	if err != nil {
		log.Printf("Error hashing password: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		log.Printf("Error when parsing userID into UUID: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error when parsing user id"))
		return
	}

	user, err := cfg.database.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:           userID,
		Email:        request.Email,
		PasswordHash: passwordHash,
		UpdatedAt:    time.Now().UTC(),
	})

	updateUserResponse := updateUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	updateUserResponseData, err := json.Marshal(&updateUserResponse)
	if err != nil {
		log.Printf("Error when encoding response: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Error when encoding response"))
		return
	}

	rw.Write(updateUserResponseData)
}
