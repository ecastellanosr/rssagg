package main

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strings"
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

func loginhandler(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("no arguments in command line")
	}
	if len(cmd.arguments) >= 2 {
		return fmt.Errorf("handlerLogin can't take more than one user")
	}
	username := cmd.arguments[0]
	firstletter := []byte{username[0]}

	if bytes.ContainsAny(firstletter, "1234567890") {
		return fmt.Errorf("username can't start with a number")
	}
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("current user does not exist")
	}
	s.config_state.Current_user_name = user.Name
	fmt.Printf("Current User: %v\n", user.Name)
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
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("no arguments in command line")
	}

	// if it has multiple arguments, return too many arguments error
	if len(cmd.arguments) >= 2 {
		return fmt.Errorf("register can't take more than one user")
	}

	//get username from arguments and take the first letter
	username := cmd.arguments[0]
	firstletter := []byte{username[0]}

	if bytes.ContainsAny(firstletter, "1234567890") {
		return fmt.Errorf("username can't start with a number")
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
func reset(s *state, cmd command) error {

	// if it has multiple arguments, return too many arguments error
	if len(cmd.arguments) >= 1 {
		return fmt.Errorf("reset does not take any argument")
	}
	err := s.db.ResetTable(context.Background())
	if err != nil {
		return fmt.Errorf("could not remove rows,%w", err)
	}

	fmt.Printf("Reset properly executed\n")

	return nil
}

func GetUsers(s *state, cmd command) error {

	// if it has multiple arguments, return too many arguments error
	if len(cmd.arguments) >= 1 {
		return fmt.Errorf("GetUsers does not take any argument")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("there were no users,%w", err)
	}
	for _, user := range users {
		if user == s.config_state.Current_user_name {
			fmt.Printf("%v (current)\n", user)
		} else {
			fmt.Printf("%v\n", user)
		}
	}
	return nil
}
func agg(s *state, cmd command) error {
	if cmd.arguments == nil {
		return fmt.Errorf("no arguments in command line")
	}
	if len(cmd.arguments) >= 2 {
		return fmt.Errorf("agg can't take more than one link")
	}
	url_command := cmd.arguments[0]
	parsedurl, err := url.Parse(url_command)
	if err != nil {
		// if there is no error, then the name is in the database. return error
		return fmt.Errorf("invalid URL: %w", err)
	}
	path := parsedurl.Path
	if !strings.HasSuffix(path, ".xml") {
		return fmt.Errorf("not an XML feed")
	}
	fmt.Println(s.config_state.Current_user_name)
	rssfeed, err := fetchFeed(context.Background(), url_command)
	if err != nil {
		return err
	}
	fmt.Println(rssfeed)
	return nil
}
func addfeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("too few arguments, addfeed needs name and url")
	}
	if len(cmd.arguments) > 2 {
		return fmt.Errorf("too many arguments, addfeed only takes name and url")
	}
	feed_name := cmd.arguments[0]
	feed_url := cmd.arguments[1]
	_, err := url.Parse(feed_url)
	if err != nil {
		return fmt.Errorf("invalid url, %w", err)
	}
	user_id := user.ID
	feedparams := database.CreatefeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feed_name,
		Url:       feed_url,
		UserID:    user_id,
	}
	feed, err := s.db.Createfeed(context.Background(), feedparams)
	if err != nil {
		return fmt.Errorf("error while creating feed row, %w", err)
	}
	feedfollowparams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user_id,
		FeedID:    feed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), feedfollowparams)
	if err != nil {
		return fmt.Errorf("error while creating feedfollow row, %w", err)
	}
	fmt.Printf("Feed Name: %v\nFeed URL: %v\n", feed.Name, feed.Url)
	return nil
}

func feeds(s *state, cmd command) error {
	if len(cmd.arguments) >= 1 {
		return fmt.Errorf("this command does not take an argument")
	}

	feeds, err := s.db.Feeds(context.Background())
	if err != nil {
		return fmt.Errorf("error while gathering the feeds, %v", err)
	}
	for _, feed := range feeds {
		fmt.Printf("Feed URL:%v\n Feed Name:%v\n User: %v\n", feed.Url, feed.Name, feed.Name_2)
	}
	return nil
}
func follow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("no arguments were passed, follow needs a url")
	}
	if len(cmd.arguments) > 1 {
		return fmt.Errorf("too many arguments, addfeed only takes a url")
	}

	feed_url := cmd.arguments[0]
	_, err := url.Parse(feed_url)
	if err != nil {
		return fmt.Errorf("invalid url, %w", err)
	}

	feed, err := s.db.GetFeed(context.Background(), feed_url)
	if err != nil {
		return fmt.Errorf("this url is not in the feed list, add the url to the feed")
	}

	feedfollowparams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), feedfollowparams)
	if err != nil {
		return fmt.Errorf("error while creating feedfollow row, %w", err)
	}
	fmt.Printf("Feed Name: %v\nUser Name: %v\n", feed.Name, user.Name)
	return nil
}
func following(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) >= 1 {
		return fmt.Errorf("this command does not take an argument")
	}

	feedforuser, err := s.db.GetFeedFollowsforUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error while creating feedfollowUser row, %w", err)
	}
	feedlist := []string{}
	for i, row := range feedforuser {
		if i != len(feedforuser)-1 {
			row.FeedName += ","
		}
		feedlist = append(feedlist, row.FeedName)
	}
	fmt.Println(feedlist, user.Name)
	return nil
}
func unfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("no arguments were passed, command needs a url")
	}
	if len(cmd.arguments) > 1 {
		return fmt.Errorf("too many arguments, command only takes a url")
	}
	feed, err := s.db.GetFeed(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("there is no feed with that URL, %w", err)
	}
	deleteFFParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.db.DeleteFeedFollow(context.Background(), deleteFFParams)
	if err != nil {
		return fmt.Errorf("error while deleting the feed follow, %w", err)
	}
	fmt.Printf("Feed (%v) Follow was successfully deleted\n", feed.Name)
	return nil
}
