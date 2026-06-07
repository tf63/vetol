package rules

import (
	"encoding/json"
	"fmt"
	"os"
)

// ConfigFile represents the structure of a JSON configuration file.
type ConfigFile struct {
	Mode  string   `json:"mode"`
	Rules []string `json:"rules"`
}

// LoadConfigFromFile loads a configuration from a JSON file.
func LoadConfigFromFile(filePath string) (Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var configFile ConfigFile
	if err := json.Unmarshal(data, &configFile); err != nil {
		return Config{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	mode := Mode(configFile.Mode)
	if mode != ModeWhitelist && mode != ModeBlacklist {
		return Config{}, fmt.Errorf("invalid mode: %s", configFile.Mode)
	}

	return NewConfig(mode, configFile.Rules), nil
}
