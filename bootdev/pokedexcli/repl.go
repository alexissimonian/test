package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startREPL(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		cleaned := cleanInput(text)

		if len(cleaned) == 0 {
			continue
		}

		userCommand := cleaned[0]
		availableCommands := getCommands()
		command, ok := availableCommands[userCommand]
		if !ok {
			fmt.Println("invalid command")
			continue
		}

        args := []string{}
		if len(cleaned) > 1 {
			args = cleaned[1:]
		}

		err := command.callback(cfg, args...)
		if err != nil {
			fmt.Println(err)
		}
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Opens the help menu",
			callback:    callbackHelp,
		},
		"exit": {
			name:        "exit",
			description: "Turn off Pokedex",
			callback:    callbackExit,
		},
		"map": {
			name:        "map",
			description: "Show next page locations for Pokemons",
			callback:    callbackMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Show last page locations for Pokemons",
			callback:    callbackMapb,
		},
		"explore": {
			name:        "explore {location_area}",
			description: "Show pokemons on given area",
			callback:    callbackExplore,
		},
		"catch": {
			name:        "catch {pokemon_name}",
			description: "Try catching a pokemon",
			callback:    callbackCatch,
		},
		"inspect": {
			name:        "inspect {pokemon_name}",
			description: "Inspect a pokemon",
			callback:    callbackInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "See all pokemon you've caught!",
			callback:    callbackPokedex,
		},
	}
}

func cleanInput(input string) []string {
	lowerInput := strings.ToLower(input)
	words := strings.Fields(lowerInput)
	return words
}
