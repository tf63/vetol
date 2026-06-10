package validator

import (
	"testing"

	"github.com/tf63/vetol/internal/rules"
)

func TestNewValidator(t *testing.T) {
	config := rules.Config{
		Mode: rules.ModeAllowlist,
		Rules: []rules.Rule{
			{Command: "ls"},
		},
	}
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
		config                rules.Config
		commandStr            string
		expectedValid         bool
		expectedViolatedCount int
		wantErr               bool
	}{
		{
			name: "allowlist mode - allowed single command",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "ls"},
				},
			},
			commandStr:            "ls",
			expectedValid:         true,
			expectedViolatedCount: 0,
			wantErr:               false,
		},
		{
			name: "allowlist mode - allowed command with args",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "ls", Include: []string{"-la", "/tmp"}},
				},
			},
			commandStr:            "ls -la /tmp",
			expectedValid:         true,
			expectedViolatedCount: 0,
			wantErr:               false,
		},
		{
			name: "allowlist mode - forbidden command",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "ls"},
				},
			},
			commandStr:            "cat file.txt",
			expectedValid:         false,
			expectedViolatedCount: 1,
			wantErr:               false,
		},
		{
			name: "allowlist mode - multiple commands piped",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "ls"},
					{Command: "grep", Include: []string{"test"}},
				},
			},
			commandStr:            "ls | grep test",
			expectedValid:         true,
			expectedViolatedCount: 0,
			wantErr:               false,
		},
		{
			name: "allowlist mode - one of multiple piped commands not allowed",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "ls"},
				},
			},
			commandStr:            "ls | grep test",
			expectedValid:         false,
			expectedViolatedCount: 1,
			wantErr:               false,
		},
		{
			name: "denylist mode - forbidden command",
			config: rules.Config{
				Mode: rules.ModeDenylist,
				Rules: []rules.Rule{
					{Command: "cat", Include: []string{"file.txt"}},
				},
			},
			commandStr:            "cat file.txt",
			expectedValid:         false,
			expectedViolatedCount: 1,
			wantErr:               false,
		},
		{
			name: "denylist mode - allowed command",
			config: rules.Config{
				Mode: rules.ModeDenylist,
				Rules: []rules.Rule{
					{Command: "cat"},
				},
			},
			commandStr:            "ls",
			expectedValid:         true,
			expectedViolatedCount: 0,
			wantErr:               false,
		},
		{
			name: "denylist mode - multiple forbidden commands",
			config: rules.Config{
				Mode: rules.ModeDenylist,
				Rules: []rules.Rule{
					{Command: "cat", Include: []string{"file.txt"}},
				},
			},
			commandStr:            "ls && cat file.txt",
			expectedValid:         false,
			expectedViolatedCount: 1,
			wantErr:               false,
		},
		{
			name: "empty command string",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "ls"},
				},
			},
			commandStr:            "",
			expectedValid:         true,
			expectedViolatedCount: 0,
			wantErr:               false,
		},
		{
			name: "invalid syntax - parsing error",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "ls"},
				},
			},
			commandStr:            "ls (",
			expectedValid:         false,
			expectedViolatedCount: 0,
			wantErr:               true,
		},
		{
			name: "command substitution with allowed commands",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "echo"},
					{Command: "ls"},
				},
			},
			commandStr:            "echo $(ls)",
			expectedValid:         true,
			expectedViolatedCount: 0,
			wantErr:               false,
		},
		{
			name: "command substitution with forbidden command",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "echo"},
				},
			},
			commandStr:            "echo $(cat file.txt)",
			expectedValid:         false,
			expectedViolatedCount: 1,
			wantErr:               false,
		},
		{
			name: "allowlist with prefix matching - command with quoted args",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "echo"},
				},
			},
			commandStr:            "echo 'Hello World'",
			expectedValid:         true,
			expectedViolatedCount: 0,
			wantErr:               false,
		},
		{
			name: "allowlist with prefix matching - command with multiple args",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "ls", Include: []string{"-la"}},
				},
			},
			commandStr:            "ls -la /tmp",
			expectedValid:         true,
			expectedViolatedCount: 0,
			wantErr:               false,
		},
		{
			name: "denylist with prefix matching - forbidden command with args",
			config: rules.Config{
				Mode: rules.ModeDenylist,
				Rules: []rules.Rule{
					{Command: "cat"},
				},
			},
			commandStr:            "cat -n file.txt",
			expectedValid:         false,
			expectedViolatedCount: 1,
			wantErr:               false,
		},
		{
			name: "denylist with prefix matching - different command not blocked",
			config: rules.Config{
				Mode: rules.ModeDenylist,
				Rules: []rules.Rule{
					{Command: "cat"},
				},
			},
			commandStr:            "catalog /tmp",
			expectedValid:         true,
			expectedViolatedCount: 0,
			wantErr:               false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator(&tt.config)

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
		config             rules.Config
		commandStr         string
		shouldHaveViolated bool
	}{
		{
			name: "result fields for valid command",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "ls"},
				},
			},
			commandStr:         "ls",
			shouldHaveViolated: false,
		},
		{
			name: "result fields for invalid command",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "ls"},
				},
			},
			commandStr:         "cat file.txt",
			shouldHaveViolated: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator(&tt.config)

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
		config            rules.Config
		commandStr        string
		expectedViolated  bool
		expectedSequences int
	}{
		{
			name: "two commands with && both allowed",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "ls"},
					{Command: "cat", Include: []string{"file.txt"}},
				},
			},
			commandStr:        "ls && cat file.txt",
			expectedViolated:  false,
			expectedSequences: 2,
		},
		{
			name: "two commands with && first forbidden",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "cat"},
				},
			},
			commandStr:        "ls && cat file.txt",
			expectedViolated:  true,
			expectedSequences: 2,
		},
		{
			name: "pipe with multiple commands",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "ls"},
					{Command: "grep", Include: []string{"test"}},
					{Command: "awk", Include: []string{"{print}"}},
				},
			},
			commandStr:        "ls | grep test | awk '{print}'",
			expectedViolated:  false,
			expectedSequences: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator(&tt.config)

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
		name        string
		config      rules.Config
		commandStr  string
		expectedErr bool
	}{
		{
			name: "whitespace only",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "ls"},
				},
			},
			commandStr:  "   ",
			expectedErr: false,
		},
		{
			name: "command with tabs",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "ls"},
				},
			},
			commandStr:  "ls\t-la",
			expectedErr: false,
		},
		{
			name: "command with newlines",
			config: rules.Config{
				Mode: rules.ModeAllowlist,
				Rules: []rules.Rule{
					{Command: "ls"},
				},
			},
			commandStr:  "ls\n-la",
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator(&tt.config)

			_, err := v.Validate(tt.commandStr)

			if (err != nil) != tt.expectedErr {
				t.Errorf("Validate() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}
