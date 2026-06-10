package rules

// Config represents the validation configuration.
type Config struct {
	Mode  Mode
	Rules []Rule
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
