package main

import (
	"context"
	"fmt"

	"github.com/alexissimonian/test/bootdev/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, c command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Error getting loggedin user: %v\n", err)
		}

		return handler(s, c, user)
	}
}
