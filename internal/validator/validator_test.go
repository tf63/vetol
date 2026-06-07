package validator

import (
	"testing"

	"github.com/tf63/vetol/pkg/rules"
)

func TestNewValidator(t *testing.T) {
	config := rules.NewConfig(rules.ModeWhitelist, []string{"ls"})
	v := NewValidator(&config)

	// Verify that NewValidator creates a valid Validator
	if v.config != &config {
		t.Error("NewValidator should initialize config correctly")
	}

	if v.parser == nil {
		t.Error("NewValidator should initialize parser")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name                  string
		mode                  rules.Mode
		allowedRules          []string
		commandStr            string
		expectedValid         bool
		expectedViolatedCount int
		wantErr               bool
	}{
		{
			name:                  "whitelist mode - allowed single command",
			mode:                  rules.ModeWhitelist,
			allowedRules:          []string{"ls"},
			commandStr:            "ls",
			expectedValid:         true,
			expectedViolatedCount: 0,
			wantErr:               false,
		},
		{
			name:                  "whitelist mode - allowed command with args",
			mode:                  rules.ModeWhitelist,
			allowedRules:          []string{"ls -la /tmp"},
			commandStr:            "ls -la /tmp",
			expectedValid:         true,
			expectedViolatedCount: 0,
			wantErr:               false,
		},
		{
			name:                  "whitelist mode - forbidden command",
			mode:                  rules.ModeWhitelist,
			allowedRules:          []string{"ls"},
			commandStr:            "rm file.txt",
			expectedValid:         false,
			expectedViolatedCount: 1,
			wantErr:               false,
		},
		{
			name:                  "whitelist mode - multiple commands piped",
			mode:                  rules.ModeWhitelist,
			allowedRules:          []string{"ls", "grep test"},
			commandStr:            "ls | grep test",
			expectedValid:         true,
			expectedViolatedCount: 0,
			wantErr:               false,
		},
		{
			name:                  "whitelist mode - one of multiple piped commands not allowed",
			mode:                  rules.ModeWhitelist,
			allowedRules:          []string{"ls"},
			commandStr:            "ls | grep test",
			expectedValid:         false,
			expectedViolatedCount: 1,
			wantErr:               false,
		},
		{
			name:                  "blacklist mode - forbidden command",
			mode:                  rules.ModeBlacklist,
			allowedRules:          []string{"rm file.txt"},
			commandStr:            "rm file.txt",
			expectedValid:         false,
			expectedViolatedCount: 1,
			wantErr:               false,
		},
		{
			name:                  "blacklist mode - allowed command",
			mode:                  rules.ModeBlacklist,
			allowedRules:          []string{"rm"},
			commandStr:            "ls",
			expectedValid:         true,
			expectedViolatedCount: 0,
			wantErr:               false,
		},
		{
			name:                  "blacklist mode - multiple forbidden commands",
			mode:                  rules.ModeBlacklist,
			allowedRules:          []string{"rm file.txt"},
			commandStr:            "ls && rm file.txt",
			expectedValid:         false,
			expectedViolatedCount: 1,
			wantErr:               false,
		},
		{
			name:                  "empty command string",
			mode:                  rules.ModeWhitelist,
			allowedRules:          []string{"ls"},
			commandStr:            "",
			expectedValid:         true,
			expectedViolatedCount: 0,
			wantErr:               false,
		},
		{
			name:                  "invalid syntax - parsing error",
			mode:                  rules.ModeWhitelist,
			allowedRules:          []string{"ls"},
			commandStr:            "ls (",
			expectedValid:         false,
			expectedViolatedCount: 0,
			wantErr:               true,
		},
		{
			name:                  "command substitution with allowed commands",
			mode:                  rules.ModeWhitelist,
			allowedRules:          []string{"echo", "ls"},
			commandStr:            "echo $(ls)",
			expectedValid:         true,
			expectedViolatedCount: 0,
			wantErr:               false,
		},
		{
			name:                  "command substitution with forbidden command",
			mode:                  rules.ModeWhitelist,
			allowedRules:          []string{"echo"},
			commandStr:            "echo $(rm file.txt)",
			expectedValid:         false,
			expectedViolatedCount: 2,
			wantErr:               false,
		},
		{
			name:                  "whitelist with prefix matching - command with quoted args",
			mode:                  rules.ModeWhitelist,
			allowedRules:          []string{"echo"},
			commandStr:            "echo 'Hello World'",
			expectedValid:         true,
			expectedViolatedCount: 0,
			wantErr:               false,
		},
		{
			name:                  "whitelist with prefix matching - command with multiple args",
			mode:                  rules.ModeWhitelist,
			allowedRules:          []string{"ls -la"},
			commandStr:            "ls -la /tmp",
			expectedValid:         true,
			expectedViolatedCount: 0,
			wantErr:               false,
		},
		{
			name:                  "blacklist with prefix matching - forbidden command with args",
			mode:                  rules.ModeBlacklist,
			allowedRules:          []string{"rm"},
			commandStr:            "rm -rf /",
			expectedValid:         false,
			expectedViolatedCount: 1,
			wantErr:               false,
		},
		{
			name:                  "blacklist with prefix matching - different command not blocked",
			mode:                  rules.ModeBlacklist,
			allowedRules:          []string{"rm"},
			commandStr:            "rmdir /tmp",
			expectedValid:         true,
			expectedViolatedCount: 0,
			wantErr:               false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := rules.NewConfig(tt.mode, tt.allowedRules)
			v := NewValidator(&config)

			result, err := v.Validate(tt.commandStr)

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if result.Valid != tt.expectedValid {
				t.Errorf("Validate() Valid = %v, want %v", result.Valid, tt.expectedValid)
			}

			if len(result.ViolatedCommands) != tt.expectedViolatedCount {
				t.Errorf(
					"Validate() ViolatedCommands count = %d, want %d",
					len(result.ViolatedCommands),
					tt.expectedViolatedCount,
				)
			}
		})
	}
}

func TestValidateResultFields(t *testing.T) {
	tests := []struct {
		name               string
		mode               rules.Mode
		allowedRules       []string
		commandStr         string
		shouldHaveViolated bool
	}{
		{
			name:               "result fields for valid command",
			mode:               rules.ModeWhitelist,
			allowedRules:       []string{"ls"},
			commandStr:         "ls",
			shouldHaveViolated: false,
		},
		{
			name:               "result fields for invalid command",
			mode:               rules.ModeWhitelist,
			allowedRules:       []string{"ls"},
			commandStr:         "rm file.txt",
			shouldHaveViolated: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := rules.NewConfig(tt.mode, tt.allowedRules)
			v := NewValidator(&config)

			result, err := v.Validate(tt.commandStr)
			if err != nil {
				t.Fatalf("Validate() unexpected error: %v", err)
			}

			// Check that ViolatedCommands is not nil
			if result.ViolatedCommands == nil {
				t.Error("ViolatedCommands should not be nil")
			}

			// Check that AllowedCommands is not nil
			if result.AllowedCommands == nil {
				t.Error("AllowedCommands should not be nil")
			}

			// Check that ForbiddenCommands is not nil
			if result.ForbiddenCommands == nil {
				t.Error("ForbiddenCommands should not be nil")
			}

			if tt.shouldHaveViolated && len(result.ViolatedCommands) == 0 {
				t.Error("Expected ViolatedCommands to have entries")
			}

			if !tt.shouldHaveViolated && len(result.ViolatedCommands) > 0 {
				t.Error("Expected ViolatedCommands to be empty")
			}
		})
	}
}

func TestFormatCommands(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		want     string
	}{
		{
			name:     "empty command slice",
			commands: []string{},
			want:     "",
		},
		{
			name:     "single command",
			commands: []string{"ls"},
			want:     "ls",
		},
		{
			name:     "multiple commands",
			commands: []string{"ls", "-la", "/tmp"},
			want:     "ls -la /tmp",
		},
		{
			name:     "commands with special characters",
			commands: []string{"echo", "hello world", "test"},
			want:     "echo hello world test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatCommands(tt.commands)
			if got != tt.want {
				t.Errorf("formatCommands() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValidateMultipleSequences(t *testing.T) {
	tests := []struct {
		name              string
		mode              rules.Mode
		allowedRules      []string
		commandStr        string
		expectedViolated  bool
		expectedSequences int
	}{
		{
			name:              "two commands with && both allowed",
			mode:              rules.ModeWhitelist,
			allowedRules:      []string{"ls", "cat file.txt"},
			commandStr:        "ls && cat file.txt",
			expectedViolated:  false,
			expectedSequences: 2,
		},
		{
			name:              "two commands with && first forbidden",
			mode:              rules.ModeWhitelist,
			allowedRules:      []string{"cat"},
			commandStr:        "ls && cat file.txt",
			expectedViolated:  true,
			expectedSequences: 2,
		},
		{
			name:              "pipe with multiple commands",
			mode:              rules.ModeWhitelist,
			allowedRules:      []string{"ls", "grep test", "awk {print}"},
			commandStr:        "ls | grep test | awk '{print}'",
			expectedViolated:  false,
			expectedSequences: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := rules.NewConfig(tt.mode, tt.allowedRules)
			v := NewValidator(&config)

			result, err := v.Validate(tt.commandStr)
			if err != nil {
				t.Fatalf("Validate() unexpected error: %v", err)
			}

			if result.Valid == tt.expectedViolated {
				t.Errorf(
					"Validate() Valid = %v, expected violated = %v",
					result.Valid,
					tt.expectedViolated,
				)
			}
		})
	}
}

func TestValidateEdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		mode         rules.Mode
		allowedRules []string
		commandStr   string
		expectedErr  bool
	}{
		{
			name:         "whitespace only",
			mode:         rules.ModeWhitelist,
			allowedRules: []string{"ls"},
			commandStr:   "   ",
			expectedErr:  false,
		},
		{
			name:         "command with tabs",
			mode:         rules.ModeWhitelist,
			allowedRules: []string{"ls"},
			commandStr:   "ls\t-la",
			expectedErr:  false,
		},
		{
			name:         "command with newlines",
			mode:         rules.ModeWhitelist,
			allowedRules: []string{"ls"},
			commandStr:   "ls\n-la",
			expectedErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := rules.NewConfig(tt.mode, tt.allowedRules)
			v := NewValidator(&config)

			_, err := v.Validate(tt.commandStr)

			if (err != nil) != tt.expectedErr {
				t.Errorf("Validate() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}
