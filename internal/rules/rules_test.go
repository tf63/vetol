package rules

import (
	"testing"
)

// TestRuleMatchesCommandOnly tests Rule.Matches with command-only rules.
func TestRuleMatchesCommandOnly(t *testing.T) {
	tests := []struct {
		name     string
		rule     Rule
		commands []string
		want     bool
	}{
		{
			name:     "exact command match",
			rule:     Rule{Command: "ls"},
			commands: []string{"ls"},
			want:     true,
		},
		{
			name:     "command with arguments",
			rule:     Rule{Command: "ls"},
			commands: []string{"ls", "-la", "/tmp"},
			want:     true,
		},
		{
			name:     "different command",
			rule:     Rule{Command: "ls"},
			commands: []string{"cat", "file.txt"},
			want:     false,
		},
		{
			name:     "partial command match fails",
			rule:     Rule{Command: "ls"},
			commands: []string{"lsof"},
			want:     false,
		},
		{
			name:     "empty commands",
			rule:     Rule{Command: "ls"},
			commands: []string{},
			want:     false,
		},
		{
			name:     "multi-word command",
			rule:     Rule{Command: "docker compose"},
			commands: []string{"docker", "compose", "exec", "app", "ls"},
			want:     true,
		},
		{
			name:     "multi-word command with complete match",
			rule:     Rule{Command: "docker compose up -d"},
			commands: []string{"docker", "compose", "up", "-d"},
			want:     true,
		},
		{
			name:     "multi-word command partial match fails",
			rule:     Rule{Command: "docker compose"},
			commands: []string{"docker", "run", "ls"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rule.Matches(tt.commands)
			if got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRuleMatchesWithIncludeShortFlags tests Include matching with short flags.
func TestRuleMatchesWithIncludeShortFlags(t *testing.T) {
	tests := []struct {
		name     string
		rule     Rule
		commands []string
		want     bool
	}{
		{
			name:     "short flag single character present",
			rule:     Rule{Command: "grep", Include: []string{"-r"}},
			commands: []string{"grep", "-r"},
			want:     true,
		},
		{
			name:     "short flag in combined flags",
			rule:     Rule{Command: "grep", Include: []string{"-r"}},
			commands: []string{"grep", "-rl"},
			want:     true,
		},
		{
			name:     "short flag in combined flags with different order",
			rule:     Rule{Command: "grep", Include: []string{"-r"}},
			commands: []string{"grep", "-lr"},
			want:     true,
		},
		{
			name:     "short flag multiple characters all present",
			rule:     Rule{Command: "ls", Include: []string{"-la"}},
			commands: []string{"ls", "-la"},
			want:     true,
		},
		{
			name:     "short flag multiple characters in single option",
			rule:     Rule{Command: "ls", Include: []string{"-la"}},
			commands: []string{"ls", "-lah"},
			want:     true,
		},
		{
			name:     "short flag multiple characters spread across options",
			rule:     Rule{Command: "ls", Include: []string{"-la"}},
			commands: []string{"ls", "-l", "-a"},
			want:     true,
		},
		{
			name:     "short flag missing character",
			rule:     Rule{Command: "ls", Include: []string{"-la"}},
			commands: []string{"ls", "-l"},
			want:     false,
		},
		{
			name:     "short flag not present at all",
			rule:     Rule{Command: "grep", Include: []string{"-r"}},
			commands: []string{"grep", "-l"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rule.Matches(tt.commands)
			if got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRuleMatchesWithIncludeLongFlags tests Include matching with long flags.
func TestRuleMatchesWithIncludeLongFlags(t *testing.T) {
	tests := []struct {
		name     string
		rule     Rule
		commands []string
		want     bool
	}{
		{
			name:     "long flag with value matches",
			rule:     Rule{Command: "ls", Include: []string{"--color"}},
			commands: []string{"ls", "--color=auto"},
			want:     true,
		},
		{
			name:     "long flag exact match",
			rule:     Rule{Command: "ls", Include: []string{"--color=auto"}},
			commands: []string{"ls", "--color=auto"},
			want:     true,
		},
		{
			name:     "long flag not present",
			rule:     Rule{Command: "ls", Include: []string{"--color"}},
			commands: []string{"ls", "-la"},
			want:     false,
		},
		{
			name:     "multiple long flags all present",
			rule:     Rule{Command: "cmd", Include: []string{"--verbose", "--debug"}},
			commands: []string{"cmd", "--verbose", "--debug"},
			want:     true,
		},
		{
			name:     "multiple long flags one missing",
			rule:     Rule{Command: "cmd", Include: []string{"--verbose", "--debug"}},
			commands: []string{"cmd", "--verbose"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rule.Matches(tt.commands)
			if got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRuleMatchesWithIncludeNoPrefix tests Include matching with non-prefix patterns.
func TestRuleMatchesWithIncludeNoPrefix(t *testing.T) {
	tests := []struct {
		name     string
		rule     Rule
		commands []string
		want     bool
	}{
		{
			name:     "non-prefix pattern exact match",
			rule:     Rule{Command: "cat", Include: []string{"file.txt"}},
			commands: []string{"cat", "file.txt"},
			want:     true,
		},
		{
			name:     "non-prefix pattern not present",
			rule:     Rule{Command: "cat", Include: []string{"file.txt"}},
			commands: []string{"cat", "other.txt"},
			want:     false,
		},
		{
			name:     "non-prefix pattern among multiple arguments",
			rule:     Rule{Command: "cat", Include: []string{"file.txt"}},
			commands: []string{"cat", "file1.txt", "file.txt", "file3.txt"},
			want:     true,
		},
		{
			name:     "multiple non-prefix patterns all present",
			rule:     Rule{Command: "diff", Include: []string{"file1.txt", "file2.txt"}},
			commands: []string{"diff", "file1.txt", "file2.txt"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rule.Matches(tt.commands)
			if got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRuleMatchesWithExclude tests Exclude matching.
func TestRuleMatchesWithExclude(t *testing.T) {
	tests := []struct {
		name     string
		rule     Rule
		commands []string
		want     bool
	}{
		{
			name:     "exclude short flag present",
			rule:     Rule{Command: "grep", Include: []string{"-r"}, Exclude: []string{"-l"}},
			commands: []string{"grep", "-r", "-l"},
			want:     false,
		},
		{
			name:     "exclude short flag not present",
			rule:     Rule{Command: "grep", Include: []string{"-r"}, Exclude: []string{"-l"}},
			commands: []string{"grep", "-r"},
			want:     true,
		},
		{
			name: "exclude long flag exact match",
			rule: Rule{
				Command: "cmd",
				Include: []string{"--verbose"},
				Exclude: []string{"--quiet"},
			},
			commands: []string{"cmd", "--verbose", "--quiet"},
			want:     false,
		},
		{
			name: "exclude long flag not present",
			rule: Rule{
				Command: "cmd",
				Include: []string{"--verbose"},
				Exclude: []string{"--quiet"},
			},
			commands: []string{"cmd", "--verbose"},
			want:     true,
		},
		{
			name:     "no include only exclude",
			rule:     Rule{Command: "cmd", Exclude: []string{"-d"}},
			commands: []string{"cmd", "-a", "-b"},
			want:     true,
		},
		{
			name:     "no include only exclude with excluded flag",
			rule:     Rule{Command: "cmd", Exclude: []string{"-d"}},
			commands: []string{"cmd", "-d"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rule.Matches(tt.commands)
			if got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRuleMatchesComplexScenarios tests Rule.Matches with complex scenarios.
func TestRuleMatchesComplexScenarios(t *testing.T) {
	tests := []struct {
		name     string
		rule     Rule
		commands []string
		want     bool
	}{
		{
			name: "rule with both include and exclude all pass",
			rule: Rule{
				Command: "grep",
				Include: []string{"-r", "pattern"},
				Exclude: []string{"-q"},
			},
			commands: []string{"grep", "-r", "pattern", "file.txt"},
			want:     true,
		},
		{
			name: "rule with both include and exclude include fails",
			rule: Rule{
				Command: "grep",
				Include: []string{"-r", "pattern"},
				Exclude: []string{"-q"},
			},
			commands: []string{"grep", "-r"},
			want:     false,
		},
		{
			name:     "rule with both include and exclude exclude fails",
			rule:     Rule{Command: "grep", Include: []string{"-r"}, Exclude: []string{"-q"}},
			commands: []string{"grep", "-r", "-q"},
			want:     false,
		},
		{
			name:     "rule with command and options",
			rule:     Rule{Command: "docker run", Include: []string{"-it"}},
			commands: []string{"docker", "run", "-it", "ubuntu"},
			want:     true,
		},
		{
			name:     "rule short flag character in longer option",
			rule:     Rule{Command: "tar", Include: []string{"-z"}},
			commands: []string{"tar", "-xzf", "file.tar.gz"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rule.Matches(tt.commands)
			if got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGrepFlagOrdering tests grep -lr and -rl flag ordering equivalence.
func TestGrepFlagOrdering(t *testing.T) {
	tests := []struct {
		name     string
		rule     Rule
		commands []string
		want     bool
	}{
		{
			name:     "grep -lr pattern matches rule with -lr",
			rule:     Rule{Command: "grep", Include: []string{"-lr"}},
			commands: []string{"grep", "-lr", "pattern", "dir"},
			want:     true,
		},
		{
			name:     "grep -rl pattern matches rule with -lr",
			rule:     Rule{Command: "grep", Include: []string{"-lr"}},
			commands: []string{"grep", "-rl", "pattern", "dir"},
			want:     true,
		},
		{
			name:     "grep -lr pattern matches rule with -rl",
			rule:     Rule{Command: "grep", Include: []string{"-rl"}},
			commands: []string{"grep", "-lr", "pattern", "dir"},
			want:     true,
		},
		{
			name:     "grep -rl pattern matches rule with -rl",
			rule:     Rule{Command: "grep", Include: []string{"-rl"}},
			commands: []string{"grep", "-rl", "pattern", "dir"},
			want:     true,
		},
		{
			name:     "grep -r -l pattern matches rule with -lr",
			rule:     Rule{Command: "grep", Include: []string{"-lr"}},
			commands: []string{"grep", "-r", "-l", "pattern", "dir"},
			want:     true,
		},
		{
			name:     "grep -l -r pattern matches rule with -rl",
			rule:     Rule{Command: "grep", Include: []string{"-rl"}},
			commands: []string{"grep", "-l", "-r", "pattern", "dir"},
			want:     true,
		},
		{
			name:     "grep -rl pattern and -v excluded",
			rule:     Rule{Command: "grep", Include: []string{"-rl"}, Exclude: []string{"-v"}},
			commands: []string{"grep", "-rl", "pattern", "dir"},
			want:     true,
		},
		{
			name:     "grep -rl pattern but -v excluded fails",
			rule:     Rule{Command: "grep", Include: []string{"-rl"}, Exclude: []string{"-v"}},
			commands: []string{"grep", "-rlv", "pattern", "dir"},
			want:     false,
		},
		{
			name:     "grep -lr with -v in separate option matches",
			rule:     Rule{Command: "grep", Include: []string{"-lr"}},
			commands: []string{"grep", "-lr", "-v", "pattern", "dir"},
			want:     true,
		},
		{
			name:     "grep -lrn matches rule with -l and -r and -n",
			rule:     Rule{Command: "grep", Include: []string{"-l", "-r", "-n"}},
			commands: []string{"grep", "-lrn", "pattern", "dir"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rule.Matches(tt.commands)
			if got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRuleMatchesEdgeCases tests Rule.Matches with edge cases.
func TestRuleMatchesEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		rule     Rule
		commands []string
		want     bool
	}{
		{
			name:     "empty rule command",
			rule:     Rule{Command: ""},
			commands: []string{"ls"},
			want:     false,
		},
		{
			name:     "single character short flag",
			rule:     Rule{Command: "ls", Include: []string{"-l"}},
			commands: []string{"ls", "-l"},
			want:     true,
		},
		{
			name:     "single character short flag not present",
			rule:     Rule{Command: "ls", Include: []string{"-l"}},
			commands: []string{"ls", "-a"},
			want:     false,
		},
		{
			name:     "long flag with multiple equals signs",
			rule:     Rule{Command: "cmd", Include: []string{"--key=value=extra"}},
			commands: []string{"cmd", "--key=value=extra"},
			want:     true,
		},
		{
			name:     "long flag prefix match with equals in value",
			rule:     Rule{Command: "cmd", Include: []string{"--color"}},
			commands: []string{"cmd", "--color=auto=extra"},
			want:     true,
		},
		{
			name:     "multiple options with same flag character",
			rule:     Rule{Command: "grep", Include: []string{"-r"}},
			commands: []string{"grep", "-r", "-r", "pattern"},
			want:     true,
		},
		{
			name:     "exclude with short flag in combined option",
			rule:     Rule{Command: "grep", Exclude: []string{"-v"}},
			commands: []string{"grep", "-rvl"},
			want:     false,
		},
		{
			name:     "include long flag partial command match",
			rule:     Rule{Command: "docker compose", Include: []string{"--version"}},
			commands: []string{"docker", "compose", "--version"},
			want:     true,
		},
		{
			name:     "command exact match with trailing options",
			rule:     Rule{Command: "cd"},
			commands: []string{"cd", "/tmp"},
			want:     true,
		},
		{
			name:     "rule command longer than actual command",
			rule:     Rule{Command: "docker compose up"},
			commands: []string{"docker", "compose"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rule.Matches(tt.commands)
			if got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
