package main

import (
	"github.com/alexissimonian/test/bootdev/pokedexcli/pokeapi"
)

type config struct {
	httpClient              pokeapi.Client
	nextLocationAreaURL     *string
	previousLocationAreaURL *string
}

func main() {
    config := config{
        httpClient: pokeapi.NewClient(),
    }
	startREPL(&config)
}
