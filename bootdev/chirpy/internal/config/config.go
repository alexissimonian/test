package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func Read() (Config, error){
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("Error getting configPath: %v", err)
	}
	chirpyConfigFile, err := os.Open(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("Error opening the chirpyconfig file: %v", err)
	}
	
	chirpyConfig, err := io.ReadAll(chirpyConfigFile)
	if err != nil {
		return Config{}, fmt.Errorf("Error reading config file: %v", err)
	}

	config := Config{}
	err = json.Unmarshal(chirpyConfig, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Error parsing config into a struct: %v", err)
	}
	return config, nil
}

func getConfigFilePath() (string, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Error getting home directory: %v", err)
	}
	configPath := homeDirectory + "/.chirpy.json"
	return configPath, nil
}
