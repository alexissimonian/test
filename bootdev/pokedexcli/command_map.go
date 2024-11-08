package main

import (
	"fmt"
	"log"
)

func callbackMap(cfg *config) error {
    pokeClient := cfg.pokeapiHttpClient
	locationAreas, err := pokeClient.ListLocationAreas(cfg.nextLocationAreaURL)
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
