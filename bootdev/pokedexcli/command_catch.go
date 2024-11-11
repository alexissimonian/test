package main

import (
	"fmt"
	"math/rand"
	"strings"
)

func callbackCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("No or incorrect number of Pokemon provided")
	}

	pokemonName := args[0]
	pokemon, err := cfg.pokeapiHttpClient.GetPokemon(pokemonName)
	if err != nil {
		return fmt.Errorf("Problem getting pokemon: %v, err: %v", pokemonName, err)
	}

	fmt.Printf("Throwing pokeball at %v...\n", pokemonName)
	pokemonXP := pokemon.BaseExperience
	treshold := 50
	randNumber := rand.Intn(pokemonXP)
	if randNumber > treshold {
		return fmt.Errorf("Failed to catch %v\n",
			strings.ToUpper(string(pokemonName[0]))+string(pokemonName[1:]),
		)
	}

	fmt.Printf("%v was caught!\n",
		strings.ToUpper(string(pokemonName[0]))+string(pokemonName[1:]),
	)

    cfg.caughtPokemons[pokemonName] = pokemon
	return nil
}
