package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ecastellanosr/rssagg/internal/config"
	"github.com/ecastellanosr/rssagg/internal/database"
	"github.com/google/uuid"
)

type state struct {
	db           *database.Queries
	config_state *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	command map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if cmd.arguments == nil {
		return fmt.Errorf("no arguments in command line")
	}
	if len(cmd.arguments) >= 2 {
		return fmt.Errorf("handlerLogin can't take more than one user")
	}
	username := cmd.arguments[0]
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		// if there is no error, then the name is in the database. return error
		return fmt.Errorf("there is no username with the name: %v", username)
	}
	s.config_state.Current_user_name = username
	fmt.Println("Username has been set")
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.command[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	method := c.command[cmd.name]
	err := method(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func register(s *state, cmd command) error {
	// if register is called but no argument return error
	if cmd.arguments == nil {
		return fmt.Errorf("no arguments in command line")
	}
	// if it has multiple arguments, return too many arguments error
	if len(cmd.arguments) >= 2 {
		return fmt.Errorf("register can't take more than one user")
	}
	//Query to see if there is an existing user with that name, if there is return an error
	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err == nil {
		// if there is no error, then the name is in the database. return error
		return fmt.Errorf("there's already a user with this name. Try another name")
	}
	//user parameters for the database
	user_params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	}
	//create user from the database queries
	user, err := s.db.CreateUser(context.Background(), user_params)
	if err != nil {
		return fmt.Errorf("could not create user, %w", err)
	}
	//print the logs
	fmt.Printf("User ID: %v\n created at: %v\n updated at: %v\n User Name: %v\n", user.ID, user.CreatedAt, user.UpdatedAt, user.Name)
	//update current state user
	s.config_state.Current_user_name = cmd.arguments[0]
	return nil
}
