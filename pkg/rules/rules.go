package rules

import "strings"

// Mode represents the validation mode: whitelist or blacklist.
type Mode string

const (
	ModeWhitelist Mode = "whitelist"
	ModeBlacklist Mode = "blacklist"
)

// Rule represents a single rule which can be a command or a sequence of commands.
type Rule struct {
	Commands []string
}

// NewRule creates a new rule from a string representation.
// The string can be a single command or a space-separated sequence of commands.
func NewRule(ruleStr string) Rule {
	commands := strings.Fields(ruleStr)
	return Rule{Commands: commands}
}

// Matches checks if the provided command sequence matches this rule.
func (r *Rule) Matches(commands []string) bool {
	if len(r.Commands) != len(commands) {
		return false
	}

	for i, cmd := range r.Commands {
		if commands[i] != cmd {
			return false
		}
	}
	return true
}

// Config represents the validation configuration.
type Config struct {
	Mode  Mode
	Rules []Rule
}

// NewConfig creates a new configuration with the given mode and rules.
func NewConfig(mode Mode, ruleStrings []string) Config {
	rules := make([]Rule, len(ruleStrings))
	for i, ruleStr := range ruleStrings {
		rules[i] = NewRule(ruleStr)
	}
	return Config{
		Mode:  mode,
		Rules: rules,
	}
}

// IsValid checks if the provided command sequence is valid according to the rules.
func (c *Config) IsValid(commands []string) bool {
	switch c.Mode {
	case ModeWhitelist:
		return c.isWhitelisted(commands)
	case ModeBlacklist:
		return !c.isBlacklisted(commands)
	default:
		return false
	}
}

// isWhitelisted checks if the command sequence is in the whitelist.
func (c *Config) isWhitelisted(commands []string) bool {
	for _, rule := range c.Rules {
		if rule.Matches(commands) {
			return true
		}
	}
	return false
}

// isBlacklisted checks if the command sequence is in the blacklist.
func (c *Config) isBlacklisted(commands []string) bool {
	for _, rule := range c.Rules {
		if rule.Matches(commands) {
			return true
		}
	}
	return false
}
