package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func configPath(segments ...string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("couldn't open home directory, %w", err)
	}
	segments = append([]string{home, "proyects", "go_proyects", "rssaggregator"}, segments...)
	return filepath.Join(segments...), nil
}

func Read() ([]config, error) {

	config_file, err := configPath(".gatorconfig.json")
	if err != nil {
		return nil, err
	}
	jsonfile, err := os.ReadFile(config_file)
	if err != nil {
		return nil, fmt.Errorf("couldn't read json file, %w", err)
	}
	var config []config
	json.Unmarshal(jsonfile, &config)
	return config, nil

}

func (c *config) SetUser(user string) error {
	c.Current_user_name = user
	config_file, err := configPath(".gatorconfig.json")
	if err != nil {
		return err
	}
	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("couldn't Marshal config data, %v", err)
	}
	err = os.WriteFile(config_file, data, 0644)
	if err != nil {
		return fmt.Errorf("couldn't write to configuration file, %v", err)
	}
	return nil
}
