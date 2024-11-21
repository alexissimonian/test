package main

import (
	"fmt"
	"os"

	"github.com/alexissimonian/test/bootdev/gator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	currentConfig, err := config.Read()
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	currentState := state{
		config: &currentConfig,
	}

	allCommands := commands{
		commands: make(map[string]func(*state, command) error),
	}

	allCommands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Please provide a command argument.")
		os.Exit(1)
	}

	userCommandName := os.Args[1]
	userCommand := command{
		name: userCommandName,
		args: os.Args[2:],
	}

	err = allCommands.run(&currentState, userCommand)
	if err != nil {
		fmt.Printf("Error handling your command : %v", err)
		os.Exit(1)
	}
}
