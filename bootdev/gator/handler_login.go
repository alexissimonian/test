package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Login expects at least one argument")
	}

	if len(cmd.args) > 1 {
		return fmt.Errorf("Too many arguments for login command, expected 1, got %v",
			len(cmd.args))
	}

	username := cmd.args[0]
	if err := s.config.SetUser(username); err != nil {
		return fmt.Errorf("Error applying username : %v", err)
	}

	fmt.Printf("User with username: %v, has been set !\n", username)

	return nil
}
