package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const CONFIG_FILE_NAME = ".chirpyconfig.json"

type Config struct {
	DB_URL     string
	JWT_SECRET string
	POLKA_KEY  string
}

func Read() (Config, error) {
	cfg_path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	bytes, err := os.ReadFile(cfg_path)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) SetUser(username string) error {
	cfg_path, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error getting home directory: %w", err)
	}

	bytes, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("error marshalling config: %w", err)
	}

	err = os.WriteFile(cfg_path, bytes, 0666)
	if err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}
	return nil
}

func getConfigFilePath() (string, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {

		return "", fmt.Errorf("error getting path to home directory: %w", err)
	}
	return home_dir + "/" + CONFIG_FILE_NAME, nil
}
