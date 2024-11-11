package main

import (
	"fmt"
)

func callbackMap(cfg *config, args ...string) error {
    pokeClient := cfg.pokeapiHttpClient
	locationAreas, err := pokeClient.ListLocationAreas(cfg.nextLocationAreaURL)
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
