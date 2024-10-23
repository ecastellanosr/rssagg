package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func configPath(segments ...string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("couldn't open home directory, %w", err)
	}
	segments = append([]string{home, "projects", "go_projects", "rssaggregator"}, segments...)
	return filepath.Join(segments...), nil
}

func Read() (Config, error) {
	var empty_config = Config{}
	config_file, err := configPath(".gatorconfig.json")
	if err != nil {
		return empty_config, err
	}
	jsonfile, err := os.ReadFile(config_file)
	if err != nil {
		return empty_config, fmt.Errorf("couldn't read json file, %w", err)
	}
	var config Config
	json.Unmarshal(jsonfile, &config)
	return config, nil

}

func (c *Config) SetUser() error {
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
