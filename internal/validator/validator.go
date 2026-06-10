package validator

import (
	"fmt"
	"strings"

	"github.com/tf63/vetol/internal/parser"
	"github.com/tf63/vetol/internal/rules"
)

// ValidationResult represents the result of a validation check.
type ValidationResult struct {
	Valid             bool
	ViolatedCommands  []string
	AllowedCommands   []string
	ForbiddenCommands []string
}

// Validator performs security validation on bash command strings.
type Validator struct {
	parser *parser.Parser
	config *rules.Config
}

// NewValidator creates a new Validator with the given configuration.
func NewValidator(config *rules.Config) *Validator {
	return &Validator{
		parser: parser.NewParser(),
		config: config,
	}
}

// Validate checks if a bash command string is valid according to the configuration.
func (v *Validator) Validate(commandStr string) (ValidationResult, error) {
	sequences, err := v.parser.ExtractCommandSequences(commandStr)
	if err != nil {
		return ValidationResult{}, fmt.Errorf("failed to parse command: %w", err)
	}

	result := ValidationResult{
		Valid:             true,
		ViolatedCommands:  []string{},
		AllowedCommands:   []string{},
		ForbiddenCommands: []string{},
	}

	for _, seq := range sequences {
		isValid := v.config.IsValid(seq.Commands)

		switch v.config.Mode {
		case rules.ModeAllowlist:
			if !isValid {
				result.Valid = false
				result.ViolatedCommands = append(
					result.ViolatedCommands,
					formatCommands(seq.Commands),
				)
			}
		case rules.ModeDenylist:
			if !isValid {
				result.Valid = false
				result.ViolatedCommands = append(
					result.ViolatedCommands,
					formatCommands(seq.Commands),
				)
			}
		}
	}

	return result, nil
}

// formatCommands joins command slices into a single string.
func formatCommands(commands []string) string {
	if len(commands) == 0 {
		return ""
	}

	var result strings.Builder
	for i, cmd := range commands {
		if i > 0 {
			result.WriteString(" ")
		}
		result.WriteString(cmd)
	}
	return result.String()
}
