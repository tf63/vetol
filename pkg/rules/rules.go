package rules

import "strings"

// Mode represents the validation mode: allowlist or denylist.
type Mode string

const (
	ModeAllowlist Mode = "allowlist"
	ModeDenylist  Mode = "denylist"
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

// Matches checks if the provided command sequence matches this rule using prefix matching.
// The rule matches if the command sequence starts with all elements of the rule.
// For example, rule "echo" matches ["echo"], ["echo", "arg1"], ["echo", "arg1", "arg2"], etc.
func (r *Rule) Matches(commands []string) bool {
	// Prefix matching: rule length must be <= command length
	if len(r.Commands) > len(commands) {
		return false
	}

	// Check if all rule commands match the beginning of the command sequence
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
	case ModeAllowlist:
		return c.isAllowlisted(commands)
	case ModeDenylist:
		return !c.isDenylisted(commands)
	default:
		return false
	}
}

// isAllowlisted checks if the command sequence is in the allowlist.
func (c *Config) isAllowlisted(commands []string) bool {
	for _, rule := range c.Rules {
		if rule.Matches(commands) {
			return true
		}
	}
	return false
}

// isDenylisted checks if the command sequence is in the denylist.
func (c *Config) isDenylisted(commands []string) bool {
	for _, rule := range c.Rules {
		if rule.Matches(commands) {
			return true
		}
	}
	return false
}
