package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func Read() (Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("Error getting configPath: %v", err)
	}
	gatorConfigFile, err := os.Open(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("Error opening the gatorconfig file: %v", err)
	}

	gatorConfig, err := io.ReadAll(gatorConfigFile)
	if err != nil {
		return Config{}, fmt.Errorf("Error reading config file: %v", err)
	}

	config := Config{}
	err = json.Unmarshal(gatorConfig, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Error parsing config into a struct: %v", err)
	}
	return config, nil
}

func (c *Config) SetUser(username string) error {
	if len(username) < 1 {
		return fmt.Errorf("Invalid username. Must be at least one character.")
	}

	c.CurrentUserName = username
    file, err := json.MarshalIndent(c, "", " ")
    if err != nil {
        return fmt.Errorf("Error converting config struct to file bytes: %v", err)
    }
    
    filePath, err := getConfigFilePath()
    if err != nil {
        return fmt.Errorf("Error getting configPath: %v", err)
    }
    err = os.WriteFile(filePath, file, 644)
    return nil
}

func getConfigFilePath() (string, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error getting home directory: %v", err)
	}
	configPath := homeDirectory + "/.gatorconfig.json"
	return configPath, nil
}
