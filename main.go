package main

import (
	"fmt"
	"os"

	"github.com/ecastellanosr/rssagg/internal/config"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("not enough arguments were provided\n")
		os.Exit(1)
	}
	command_name := os.Args[1]
	argument := os.Args[2]

	current_config, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	state_current := state{
		config_state: &current_config,
	}
	cmds := commands{
		command: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)

	command := command{
		name:      command_name,
		arguments: []string{argument},
	}
	err = cmds.run(&state_current, command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if command.name == "login" {
		current_config.Current_user_name = state_current.config_state.Current_user_name
		current_config.SetUser()
	}
}
