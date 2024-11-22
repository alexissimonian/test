package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alexissimonian/test/bootdev/gator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Register expects one argument. Got: %v", len(cmd.args))
	}

    username := cmd.args[0]

    user, err := s.db.GetUser(context.Background(), username)

    arg := database.CreateUserParams{
        ID: uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Name: cmd.args[0],
    } 
}
