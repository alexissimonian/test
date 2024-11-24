package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alexissimonian/test/bootdev/gator/internal/database"
	"github.com/google/uuid"
)

func handlerReset(s *state, c command) error {
	if len(c.args) != 0 {
		return fmt.Errorf("Reset takes no arguments. Got %v\n", len(c.args))
	}

	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Could not restet db users: %v\n", err)
	}
    
    fmt.Printf("Successfully restet all users !\n")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Register expects one argument. Got: %v\n", len(cmd.args))
	}

	username := cmd.args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      username,
	})

	if err != nil {
		return fmt.Errorf("Something went wrong registering your user: %v\n", err)
	}

	err = s.config.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Something went wrong setting user in config: %v\n", err)
	}

	fmt.Printf("User: %v registered with id: %v\n", user.Name, user.ID)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Login expects at least one argument")
	}

	if len(cmd.args) > 1 {
		return fmt.Errorf("Too many arguments for login command, expected 1, got %v",
			len(cmd.args))
	}

	username := cmd.args[0]
    _, err := s.db.GetUser(context.Background(), username)
    if err != nil {
        return fmt.Errorf("Could not find user: %v\n", username)
    }

	if err := s.config.SetUser(username); err != nil {
		return fmt.Errorf("Error applying username : %v", err)
	}

	fmt.Printf("Logged in as user: %v\n", username)

	return nil
}

func handlerUsers(s *state, c command) error {
	if len(c.args) != 0 {
		return fmt.Errorf("Users expects no arguments. Got %v\n", len(c.args))
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Something went wrong fetching users: %v\n", err)
	}

	connectedUserName := s.config.CurrentUserName

	for _, user := range users {
		if connectedUserName == user.Name {
			fmt.Printf("%v (current)\n", user.Name)
		} else {
			fmt.Printf("%v\n", user.Name)
		}
	}
	return nil
}
