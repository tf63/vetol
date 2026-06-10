package rules

import (
	"encoding/json"
	"fmt"
	"os"
)

// ConfigFile represents the structure of a JSON configuration file.
type ConfigFile struct {
	Mode  string     `json:"mode"`
	Rules []RuleFile `json:"rules"`
}

// RuleFile represents a single rule in the JSON configuration file.
type RuleFile struct {
	Command string   `json:"command"`
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
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
	if mode != ModeAllowlist && mode != ModeDenylist {
		return Config{}, fmt.Errorf("invalid mode: %s", configFile.Mode)
	}

	rules := make([]Rule, len(configFile.Rules))
	for i, rf := range configFile.Rules {
		rules[i] = Rule(rf)
	}

	return Config{
		Mode:  mode,
		Rules: rules,
	}, nil
}
