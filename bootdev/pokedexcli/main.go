package main

import (
	"time"

	"github.com/alexissimonian/test/bootdev/pokedexcli/internal/pokeapi"
)

type config struct {
	pokeapiHttpClient       pokeapi.Client
	nextLocationAreaURL     *string
	previousLocationAreaURL *string
	caughtPokemons          map[string]pokeapi.PokemonResponse
}

func main() {
	config := config{
		pokeapiHttpClient: pokeapi.NewClient(time.Hour),
		caughtPokemons:    make(map[string]pokeapi.PokemonResponse),
	}
	startREPL(&config)
}
