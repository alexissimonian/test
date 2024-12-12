package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/alexissimonian/test/bootdev/chirpy/internal/config"
	"github.com/alexissimonian/test/bootdev/chirpy/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	cfg := apiConfig{}
	loadConfig(&cfg)
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	serveMux.HandleFunc("GET /api/healthz", readinessHandler)
	serveMux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	serveMux.HandleFunc("POST /admin/reset", cfg.resetMetricsHandler)
	serveMux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)
	serveMux.HandleFunc("POST /api/users", createUserHandler)
	server := &http.Server{}
	server.Addr = ":8080"
	server.Handler = serveMux
	server.ListenAndServe()
}

func readinessHandler(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	_, err := rw.Write([]byte("OK"))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) metricsHandler(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Add("Content-Type", "text/html")
	responseContent := fmt.Sprintf(`
    <html>
        <body>
            <h1>Welcome, Chirpy Admin</h1>
            <p>Chirpy has been visited %d times!</p>
        </body>
    </html>
    `, cfg.fileServerHits.Load())
	_, err := rw.Write([]byte(responseContent))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) resetMetricsHandler(rw http.ResponseWriter, _ *http.Request) {
	cfg.fileServerHits.Store(0)
	rw.Write([]byte("Hits reset to zero!"))
}

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

type createUserRequest struct {
	Email string `json:"email"`
}

type createUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func createUserHandler(rw http.ResponseWriter, r *http.Request) {
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

	config, err := config.Read()
	if err != nil {
		log.Printf("Error reading config: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	dburl := config.DbURL

	openedDBConnection, err := sql.Open("postgres", dburl)
	if err != nil {
		log.Printf("Something went wrong opening database connection: %v\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	db := database.New(openedDBConnection)

	user, err := db.CreateUser(r.Context(), database.CreateUserParams{
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
