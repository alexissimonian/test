package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/alexissimonian/test/bootdev/gator/internal/config"
	"github.com/alexissimonian/test/bootdev/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db     *database.Queries
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

	dbURL := currentState.config.DbURL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Something went wrong opening database connection: %v\n", err)
	}

	currentState.db = database.New(db)

	allCommands := commands{
		commands: make(map[string]func(*state, command) error),
	}

	allCommands.register("login", handlerLogin)
	allCommands.register("register", handlerRegister)
	allCommands.register("reset", handlerReset)
	allCommands.register("users", handlerUsers)
	allCommands.register("agg", handlerAggregator)
	allCommands.register("addFeed", handlerAddFeed)

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
