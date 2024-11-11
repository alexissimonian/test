package main

import (
	"fmt"
)

func callbackExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("No or incorrect number of area provided.")
	}
    area := args[0]
    fmt.Printf("Exploring %v...\n", area)
    pokeClient := cfg.pokeapiHttpClient
    locationArea, err := pokeClient.GetLocationArea(area)
    if err != nil {
        return err
    }

    for _, pokemonEncounter := range locationArea.PokemonEncounters {
        fmt.Println(pokemonEncounter.Pokemon.Name)
    }
	return nil
}
