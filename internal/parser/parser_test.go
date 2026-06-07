package parser

import (
	"slices"
	"testing"
)

func TestNewParser(t *testing.T) {
	p := NewParser()
	if p == nil {
		t.Error("NewParser should not return nil")
	}
}

func TestExtractCommandSequences(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    [][]string
		wantErr bool
	}{
		{
			name:    "simple command",
			command: "ls",
			want:    [][]string{{"ls"}},
			wantErr: false,
		},
		{
			name:    "command with arguments",
			command: "ls -la /tmp",
			want:    [][]string{{"ls", "-la", "/tmp"}},
			wantErr: false,
		},
		{
			name:    "double quoted string",
			command: "echo \"hello\"",
			want:    [][]string{{"echo", "hello"}},
			wantErr: false,
		},
		{
			name:    "single quoted string",
			command: "echo 'world'",
			want:    [][]string{{"echo", "world"}},
			wantErr: false,
		},
		{
			name:    "command substitution",
			command: "echo $(ls -la)",
			want:    [][]string{{"echo"}, {"ls", "-la"}},
			wantErr: false,
		},
		{
			name:    "pipe command",
			command: "ls | grep test",
			want:    [][]string{{"ls"}, {"grep", "test"}},
			wantErr: false,
		},
		{
			name:    "multiple commands with &&",
			command: "ls && cat file.txt",
			want:    [][]string{{"ls"}, {"cat", "file.txt"}},
			wantErr: false,
		},
		{
			name:    "invalid syntax",
			command: "ls (",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty string",
			command: "",
			want:    [][]string{},
			wantErr: false,
		},
		{
			name:    "command with mixed quotes",
			command: "echo 'hello' \"world\"",
			want:    [][]string{{"echo", "hello", "world"}},
			wantErr: false,
		},
		{
			name:    "process substitution",
			command: "cat <(ls -la)",
			want:    [][]string{{"cat"}, {"ls", "-la"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			got, err := p.ExtractCommandSequences(tt.command)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractCommandSequences() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf(
					"ExtractCommandSequences() got %d sequences, want %d",
					len(got),
					len(tt.want),
				)
				return
			}

			for i, seq := range got {
				if !slices.Equal(seq.Commands, tt.want[i]) {
					t.Errorf(
						"ExtractCommandSequences()[%d] got %v, want %v",
						i,
						seq.Commands,
						tt.want[i],
					)
				}
			}
		})
	}
}
