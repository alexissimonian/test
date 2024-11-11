package main

import (
	"fmt"
)

func callbackMapb(cfg *config, args ...string) error {
	if cfg.previousLocationAreaURL == nil {
		return fmt.Errorf("You're on the first page")
	}
	pokeClient := cfg.pokeapiHttpClient
	locationAreas, err := pokeClient.ListLocationAreas(cfg.previousLocationAreaURL)
	if err != nil {
		return err
	}

	for _, area := range locationAreas.Results {
		fmt.Println(area.Name)
	}

	cfg.nextLocationAreaURL = locationAreas.Next
	cfg.previousLocationAreaURL = locationAreas.Previous
	return nil
}
