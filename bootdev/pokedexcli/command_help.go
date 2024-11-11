package main

import "fmt"

func callbackHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to Pokedex help menu :")
	fmt.Println("Here are available commands:")
    for _, command := range getCommands() {
        fmt.Printf("- %s : %s\n", command.name, command.description)
    }
    
    return nil
}
