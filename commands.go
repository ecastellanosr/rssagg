package main

import (
	"fmt"

	"github.com/ecastellanosr/rssagg/internal/config"
)

type state struct {
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
	s.config_state.Current_user_name = cmd.arguments[0]
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
