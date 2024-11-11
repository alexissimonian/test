package main

import (
	"fmt"
	"strings"
)

func callbackPokedex(cfg *config, args ...string) error {
	if len(cfg.caughtPokemons) == 0 {
		return fmt.Errorf("No pokemons caught yet")
	}

	for k := range cfg.caughtPokemons {
		fmt.Printf("- %v\n",
			strings.ToUpper(string(k[0]))+string(k[1:]),
		)
	}
	return nil
}
