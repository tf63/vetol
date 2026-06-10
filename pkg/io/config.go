package io

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tf63/vetol/pkg/rules"
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
func LoadConfigFromFile(filePath string) (rules.Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return rules.Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var configFile ConfigFile
	if err := json.Unmarshal(data, &configFile); err != nil {
		return rules.Config{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	mode := rules.Mode(configFile.Mode)
	if mode != rules.ModeAllowlist && mode != rules.ModeDenylist {
		return rules.Config{}, fmt.Errorf("invalid mode: %s", configFile.Mode)
	}

	rulesSlice := make([]rules.Rule, len(configFile.Rules))
	for i, rf := range configFile.Rules {
		rulesSlice[i] = rules.Rule{
			Command: rf.Command,
			Include: rf.Include,
			Exclude: rf.Exclude,
		}
	}

	return rules.Config{
		Mode:  mode,
		Rules: rulesSlice,
	}, nil
}
