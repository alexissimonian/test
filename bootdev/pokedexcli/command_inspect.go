package main

import "fmt"

func callbackInspect(cfg *config, args ...string) error {
    if len(args) != 1 {
        return fmt.Errorf("No or incorrect number of pokemon names")
    }

    pokemonName := args[0]
    pokemon, ok := cfg.caughtPokemons[pokemonName]
    if !ok {
        return fmt.Errorf("You have not caught %v yet", pokemonName)
    }

    fmt.Printf("Name: %v\n", pokemon.Name)
    fmt.Printf("Height: %v\n", pokemon.Height)
    fmt.Printf("Weight: %v\n", pokemon.Weight)
    fmt.Printf("Stats:\n")
    for _, stat := range pokemon.Stats {
        fmt.Printf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat)
    }
    fmt.Printf("Types:\n")
    for _, typ := range pokemon.Types {
        fmt.Printf("  -%v\n", typ.Type.Name)
    }
    return nil
}
