package rules

import (
	"testing"
)

// TestConfigIsValidAllowlist tests Config.IsValid in allowlist mode.
func TestConfigIsValidAllowlist(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		commands []string
		want     bool
	}{
		{
			name: "single matching rule",
			config: Config{
				Mode: ModeAllowlist,
				Rules: []Rule{
					{Command: "ls"},
				},
			},
			commands: []string{"ls"},
			want:     true,
		},
		{
			name: "no matching rule",
			config: Config{
				Mode: ModeAllowlist,
				Rules: []Rule{
					{Command: "ls"},
				},
			},
			commands: []string{"cat"},
			want:     false,
		},
		{
			name: "multiple rules with one match",
			config: Config{
				Mode: ModeAllowlist,
				Rules: []Rule{
					{Command: "ls"},
					{Command: "cat"},
					{Command: "grep"},
				},
			},
			commands: []string{"cat", "file.txt"},
			want:     true,
		},
		{
			name: "multiple rules no match",
			config: Config{
				Mode: ModeAllowlist,
				Rules: []Rule{
					{Command: "ls"},
					{Command: "cat"},
				},
			},
			commands: []string{"rm", "file.txt"},
			want:     false,
		},
		{
			name: "rule with include matches",
			config: Config{
				Mode: ModeAllowlist,
				Rules: []Rule{
					{Command: "grep", Include: []string{"-r"}},
				},
			},
			commands: []string{"grep", "-r", "pattern"},
			want:     true,
		},
		{
			name: "rule with include does not match",
			config: Config{
				Mode: ModeAllowlist,
				Rules: []Rule{
					{Command: "grep", Include: []string{"-r"}},
				},
			},
			commands: []string{"grep", "-l"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.IsValid(tt.commands)
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestConfigIsValidDenylist tests Config.IsValid in denylist mode.
func TestConfigIsValidDenylist(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		commands []string
		want     bool
	}{
		{
			name: "single matching denied rule",
			config: Config{
				Mode: ModeDenylist,
				Rules: []Rule{
					{Command: "rm"},
				},
			},
			commands: []string{"rm", "file.txt"},
			want:     false,
		},
		{
			name: "single non-matching rule allowed",
			config: Config{
				Mode: ModeDenylist,
				Rules: []Rule{
					{Command: "rm"},
				},
			},
			commands: []string{"ls"},
			want:     true,
		},
		{
			name: "multiple rules one matches denied",
			config: Config{
				Mode: ModeDenylist,
				Rules: []Rule{
					{Command: "rm"},
					{Command: "dd"},
				},
			},
			commands: []string{"dd", "if=/dev/zero"},
			want:     false,
		},
		{
			name: "multiple rules none match allowed",
			config: Config{
				Mode: ModeDenylist,
				Rules: []Rule{
					{Command: "rm"},
					{Command: "dd"},
				},
			},
			commands: []string{"ls"},
			want:     true,
		},
		{
			name: "rule with include matches denied",
			config: Config{
				Mode: ModeDenylist,
				Rules: []Rule{
					{Command: "grep", Include: []string{"-r"}},
				},
			},
			commands: []string{"grep", "-r", "pattern"},
			want:     false,
		},
		{
			name: "rule with include does not match allowed",
			config: Config{
				Mode: ModeDenylist,
				Rules: []Rule{
					{Command: "grep", Include: []string{"-r"}},
				},
			},
			commands: []string{"grep", "-l"},
			want:     true,
		},
		{
			name: "rule with exclude matches denied",
			config: Config{
				Mode: ModeDenylist,
				Rules: []Rule{
					{Command: "rm", Exclude: []string{"-i"}},
				},
			},
			commands: []string{"rm", "file.txt"},
			want:     false,
		},
		{
			name: "rule with exclude does not match allowed",
			config: Config{
				Mode: ModeDenylist,
				Rules: []Rule{
					{Command: "rm", Exclude: []string{"-i"}},
				},
			},
			commands: []string{"rm", "-i", "file.txt"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.IsValid(tt.commands)
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestConfigIsValidInvalidMode tests Config.IsValid with invalid mode.
func TestConfigIsValidInvalidMode(t *testing.T) {
	config := Config{
		Mode: Mode("invalid"),
		Rules: []Rule{
			{Command: "ls"},
		},
	}
	got := config.IsValid([]string{"ls"})
	want := false
	if got != want {
		t.Errorf("IsValid() with invalid mode = %v, want %v", got, want)
	}
}

// TestConfigIsValidEmptyRules tests Config.IsValid with empty rules.
func TestConfigIsValidEmptyRules(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		commands []string
		want     bool
	}{
		{
			name: "allowlist with empty rules denies all",
			config: Config{
				Mode:  ModeAllowlist,
				Rules: []Rule{},
			},
			commands: []string{"ls"},
			want:     false,
		},
		{
			name: "allowlist with empty rules denies complex command",
			config: Config{
				Mode:  ModeAllowlist,
				Rules: []Rule{},
			},
			commands: []string{"docker", "compose", "up", "-d"},
			want:     false,
		},
		{
			name: "denylist with empty rules allows all",
			config: Config{
				Mode:  ModeDenylist,
				Rules: []Rule{},
			},
			commands: []string{"rm", "-rf", "/"},
			want:     true,
		},
		{
			name: "denylist with empty rules allows simple command",
			config: Config{
				Mode:  ModeDenylist,
				Rules: []Rule{},
			},
			commands: []string{"ls"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.IsValid(tt.commands)
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
