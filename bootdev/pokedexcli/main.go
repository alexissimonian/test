package main

import (
	"github.com/alexissimonian/test/bootdev/pokedexcli/internal/pokeapi"
)

type config struct {
	pokeapiHttpClient              pokeapi.Client
	nextLocationAreaURL     *string
	previousLocationAreaURL *string
}

func main() {
    config := config{
        pokeapiHttpClient: pokeapi.NewClient(),
    }
	startREPL(&config)
}
