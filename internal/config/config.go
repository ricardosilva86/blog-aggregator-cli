package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	h, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("error getting user's home directory: %w", err)
	}
	f, err := os.Open(filepath.Join(h, ".gatorconfig.json"))
	if err != nil {
		return Config{}, fmt.Errorf("error opening file: %w", err)
	}
	jsonContent, err := ioutil.ReadAll(f)
	if err != nil {
		return Config{}, fmt.Errorf("error reading file: %w", err)
	}
	defer f.Close()

	var config Config
	err = json.Unmarshal(jsonContent, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error decoding json: %w", err)
	}
	return config, nil
}

func (c *Config) SetUser(name string) error {
	config, err := Read()
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	config.CurrentUserName = name

	jsonContent, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("error encoding json: %w", err)
	}

	h, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting user's home directory: %w", err)
	}
	f, err := os.OpenFile(
		filepath.Join(h, ".gatorconfig.json"),
		os.O_WRONLY|os.O_TRUNC,
		0644)
	if err != nil {
		return fmt.Errorf("error getting user's home directory: %w", err)
	}
	defer f.Close()

	_, err = f.Write(jsonContent)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
