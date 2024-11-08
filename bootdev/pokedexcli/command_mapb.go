package main

import (
	"fmt"
	"log"
)

func callbackMapb(cfg *config) error {
    if cfg.previousLocationAreaURL == nil {
        return fmt.Errorf("You're on the first page")
    }
    pokeClient := cfg.pokeapiHttpClient
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
