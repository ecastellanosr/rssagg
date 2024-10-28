package main

import (
	"context"
	"fmt"

	"github.com/ecastellanosr/rssagg/internal/database"
)

type authhandler func(s *state, cmd command, user database.User) error

func middlewareLoggedIn(handler authhandler) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config_state.Current_user_name)
		if err != nil {
			return fmt.Errorf("current user does not exist")
		}
		return handler(s, cmd, user)
	}

}
