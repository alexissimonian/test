package main

import (
	"database/sql"
	"log"
	"sync/atomic"

	"github.com/alexissimonian/test/bootdev/chirpy/internal/config"
	"github.com/alexissimonian/test/bootdev/chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	config config.Config
	database *database.Queries
}

func loadConfig(c *apiConfig) {
	loadConfigFile(c)
	loadDatabase(c)
}

func loadConfigFile(c *apiConfig) {
	config, err := config.Read()
	if err != nil {
		log.Panicf("Error reading config: %v\n", err)
	}
	c.config = config
}

func loadDatabase(c *apiConfig) {
	dburl := c.config.DbURL

	openedDBConnection, err := sql.Open("postgres", dburl)
	if err != nil {
		log.Panicf("Something went wrong opening database connection: %v\n", err)
	}

	c.database = database.New(openedDBConnection)
}