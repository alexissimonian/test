package main

import (
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
	config, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v\n", err)
		return
	}

	c.config = config
}