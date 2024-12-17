package main

import (
	"database/sql"
	"log"
	"os"
	"sync/atomic"

	"github.com/alexissimonian/test/bootdev/chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	database *database.Queries
	platform string
	serverSecret string
}

func loadConfig(c *apiConfig) {
	loadEnvVariables()
	loadServerSecret(c)
	loadPlatform(c)
	loadDatabase(c)
}

func loadEnvVariables(){
	err := godotenv.Load()
	if err != nil {
		log.Panicf("Something went wrong loading environment variables: %v\n", err)
	}
}

func loadServerSecret(c *apiConfig){
	c.serverSecret = os.Getenv("SERVER_SECRET")
}

func loadDbUrl() string {	
	dbUrl := os.Getenv("DB_URL")
	return dbUrl
}

func loadPlatform(c *apiConfig) {
	c.platform = os.Getenv("PLATFORM")
}

func loadDatabase(c *apiConfig) {
	dburl := loadDbUrl()

	openedDBConnection, err := sql.Open("postgres", dburl)
	if err != nil {
		log.Panicf("Something went wrong opening database connection: %v\n", err)
	}

	c.database = database.New(openedDBConnection)
}