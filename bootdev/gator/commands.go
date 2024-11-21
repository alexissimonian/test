package main

import (
	"fmt"
	"os"
)

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) error {
    if len(name) == 0 {
        return fmt.Errorf("Invalid name, must be at least one character.")
    }
    c.commands[name] = f
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.commands[cmd.name]
	if !ok {
		fmt.Printf("Invalid command: %v\n", cmd.name)
        os.Exit(1)
	}
    
    err := handler(s, cmd)
    if err != nil {
        return fmt.Errorf("Error running command: %v: %v\n", cmd.name, err)
    }
	return nil
}
