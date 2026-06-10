package rules

import (
	"slices"
	"strings"
)

// Mode represents the validation mode: allowlist or denylist.
type Mode string

const (
	ModeAllowlist Mode = "allowlist"
	ModeDenylist  Mode = "denylist"
)

// Rule represents a single rule with command name, include, and exclude patterns.
type Rule struct {
	Command string
	Include []string
	Exclude []string
}

// Matches checks if the provided command sequence matches this rule.
// It first checks command prefix matching (split by spaces), then validates include/exclude patterns.
func (r *Rule) Matches(commands []string) bool {
	if len(commands) == 0 || r.Command == "" {
		return false
	}

	// Split rule command by spaces to support multi-word commands (e.g., "docker compose")
	ruleCommands := strings.Fields(r.Command)

	// Check if input commands start with rule commands (exact match for each command element)
	if len(ruleCommands) > len(commands) {
		return false
	}

	for i, ruleCmd := range ruleCommands {
		if commands[i] != ruleCmd {
			return false
		}
	}

	// Extract options (everything after the rule commands)
	var options []string
	if len(commands) > len(ruleCommands) {
		options = commands[len(ruleCommands):]
	}

	// Check include conditions
	if len(r.Include) > 0 {
		if !matchesInclude(options, r.Include) {
			return false
		}
	}

	// Check exclude conditions
	if len(r.Exclude) > 0 {
		if matchesExclude(options, r.Exclude) {
			return false
		}
	}

	return true
}

// matchesInclude checks if all include patterns match the options.
func matchesInclude(options []string, includes []string) bool {
	for _, inc := range includes {
		if !optionMatches(options, inc) {
			return false
		}
	}
	return true
}

// matchesExclude checks if any exclude pattern matches the options.
func matchesExclude(options []string, excludes []string) bool {
	for _, exc := range excludes {
		if optionMatches(options, exc) {
			return true
		}
	}
	return false
}

// optionMatches checks if a pattern matches any option.
// - Short flags (starting with -): all characters must be present
// - Long flags (starting with --): exact match or prefix match with =
// - Others: exact match
func optionMatches(options []string, pattern string) bool {
	if strings.HasPrefix(pattern, "--") {
		// Long flag: exact match or prefix match with = (e.g., --color matches --color=auto)
		for _, opt := range options {
			if opt == pattern || strings.HasPrefix(opt, pattern+"=") {
				return true
			}
		}
		return false
	} else if strings.HasPrefix(pattern, "-") && len(pattern) > 1 {
		// Short flag: all characters in pattern must be present
		shortChars := pattern[1:] // Get characters after the '-'

		// First, check if any single option contains all characters
		for _, opt := range options {
			if strings.HasPrefix(opt, "-") && !strings.HasPrefix(opt, "--") {
				// This is a short flag, check if all characters are contained
				optChars := opt[1:]
				allCharactersFound := true
				for _, ch := range shortChars {
					if !strings.ContainsRune(optChars, ch) {
						allCharactersFound = false
						break
					}
				}
				if allCharactersFound {
					return true
				}
			}
		}

		// If no single option has all characters, check if characters are spread across options
		for _, ch := range shortChars {
			found := false
			for _, opt := range options {
				if strings.HasPrefix(opt, "-") && !strings.HasPrefix(opt, "--") {
					optChars := opt[1:]
					if strings.ContainsRune(optChars, ch) {
						found = true
						break
					}
				}
			}
			if !found {
				return false
			}
		}
		return true
	}
	// Neither short nor long flag: exact match
	return slices.Contains(options, pattern)
}
