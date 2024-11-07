package main

import (
	"fmt"
	"log"
)

func callbackMapb(cfg *config) error {
    pokeClient := cfg.httpClient
	locationAreas, err := pokeClient.ListLocationAreas(cfg.previousLocationAreaURL)
	if err != nil {
		log.Fatal(err)
	}

	for _, area := range locationAreas.Results {
		fmt.Println(area.Name)
	}

    cfg.nextLocationAreaURL = locationAreas.Next
    cfg.previousLocationAreaURL = locationAreas.Previous
	return nil
}
