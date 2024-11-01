package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ecastellanosr/rssagg/internal/config"
	"github.com/ecastellanosr/rssagg/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	// the CLI tool needs a minimum of 2 arguments as it need the name of the command and an argument for that command.
	if len(os.Args) < 1 {
		// if it does not, then return an error and stop the process.
		fmt.Printf("no command was given\n")
		os.Exit(1)
	}
	// take the arguments of the CLI execute
	command_name := os.Args[1]

	arguments := os.Args[2:]
	// read your configuration file
	current_config, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	//take the url of the database from the config file
	dburl := current_config.Db_url
	// open the database
	db, err := sql.Open("postgres", dburl)
	if err != nil {
		fmt.Println(err)
	}
	dbQueries := database.New(db)
	//update the current state with the configuration file data
	state_current := state{
		db:           dbQueries,
		config_state: &current_config,
	}
	//Commands
	cmds := commands{
		command: map[string]func(*state, command) error{},
	}
	//register the commands that can be used
	cmds.register("login", middlewareLoggedIn(loginhandler))
	cmds.register("register", register)
	cmds.register("reset", reset)
	cmds.register("getusers", GetUsers)
	cmds.register("agg", agg)
	cmds.register("addfeed", middlewareLoggedIn(addfeed))
	cmds.register("feeds", feeds)
	cmds.register("follow", middlewareLoggedIn(follow))
	cmds.register("following", middlewareLoggedIn(following))
	cmds.register("unfollow", middlewareLoggedIn(unfollow))
	cmds.register("browse", middlewareLoggedIn(browse))
	//current command that is taking place
	command := command{
		name:      command_name,
		arguments: arguments,
	}
	//run the command
	err = cmds.run(&state_current, command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// for now it gets updated in all the commands
	current_config.SetUser()
}
