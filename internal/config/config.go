// /internal/config/config.go
package config

import (
	"os"
	"encoding/json"
	"path/filepath"
	"fmt"
)

type Config struct {
	DBURL string `json:"DBURL"`
	CurrentUserName string `json:"CurrentUserName"`
}

func GetConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	const configFileName = ".gatorconfig.json"
	configFilePath := filepath.Join(homeDir, configFileName)

	return configFilePath
}

func Read(path string) (Config, error) {
	data, err := os.ReadFile(path)
    if err != nil {
        return Config{}, fmt.Errorf("reading config file: %w", err)
    }

	var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return Config{}, fmt.Errorf("parsing config json: %w", err)
    }

    return cfg, nil

}

func Write(path string, cfg Config) error {
	data, err := json.Marshal(cfg)
    if err != nil {
        return fmt.Errorf("marshaling config: %w", err)
    }

    if err := os.WriteFile(path, data, 0644); err != nil {
        return fmt.Errorf("writing config file: %w", err)
    }

    return nil
}

func (cfg *Config) SetUser(username string, configFilePath string) error {
	cfg.CurrentUserName = username

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(configFilePath, data, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}
	fmt.Printf("User '%s' logged in successfully.\n", username)

	return nil
}