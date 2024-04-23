package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config describes the configuration for Alpacon-CLI
type Config struct {
	ServerAddress string `json:"server_address"`
	Token         string `json:"token"`
	ExpiresAt     string `json:"expires_at"`
}

const (
	ConfigFileName = "config.json"
	ConfigFileDir  = ".alpacon"
)

func CreateConfig(serverAddress string, token string, expiresAt string) error {
	config := Config{
		ServerAddress: serverAddress,
		Token:         token,
		ExpiresAt:     expiresAt,
	}

	return saveConfig(&config)
}

func saveConfig(config *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ConfigFileDir)
	if err = os.MkdirAll(configDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	configFile := filepath.Join(configDir, ConfigFileName)
	file, err := os.Create(configFile)
	if err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err = encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config to JSON: %v", err)
	}

	return nil
}

func DeleteConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ConfigFileDir)
	configFile := filepath.Join(configDir, ConfigFileName)

	err = os.Remove(configFile)
	if err != nil {
		return fmt.Errorf("failed to delete config file: %v", err)
	}

	return nil
}

func LoadConfig() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("failed to get user home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ConfigFileDir)
	configFile := filepath.Join(configDir, ConfigFileName)

	file, err := os.Open(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, fmt.Errorf("config file does not exist: %v", configFile)
		}
		return Config{}, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return Config{}, fmt.Errorf("failed to decode config file: %v", err)
	}

	return config, nil
}
