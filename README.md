# Vetol

A CLI tool that analyzes bash command strings using Abstract Syntax Tree (AST) parsing and performs strict security validation to detect forbidden commands regardless of how they are nested or obfuscated.

## Overview

**Vetol** solves the security problem of simple string matching-based command validation, which can be bypassed through command chaining, substitution, or other complex syntax patterns. Instead, Vetol parses the complete bash AST and validates all commands including those nested in pipes, substitutions, chains, and subshells.

### Problem

Simple string matching is insufficient:

- Command chaining: `ls && rm -rf /` - bypasses simple "rm" prefix check
- Command substitution: `cat "$(rm -rf /)"` - hides forbidden command inside variable expansion
- Pipelines and shell constructs - can conceal forbidden operations

### Solution

Vetol uses AST-based parsing to detect all commands, including those nested in:

- Command chains (`&&`, `||`, `;`)
- Pipelines (`|`)
- Command substitution (`$()`, backticks)
- Subshells and other shell constructs

## Installation

### Using `go install`

```bash
go install github.com/tf63/vetol/cmd/vetol@latest
```

Install a specific version:

```bash
go install github.com/tf63/vetol/cmd/vetol@v1.0.0
```

Ensure `$GOPATH/bin` is in your `$PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Usage

### Whitelist Mode

Only allow explicitly permitted commands:

```bash
# Single command
vetol check --mode whitelist --rules ls,cat,grep "ls -la /tmp"

# Multi-command sequence
vetol check -m whitelist -r "docker compose exec app go fmt,ls,cat" "docker compose exec app go fmt ./..."
```

### Blacklist Mode

Allow all commands except explicitly forbidden ones:

```bash
# Single command
vetol check --mode blacklist --rules rm,dd,mkfs "cat /etc/passwd"

# Multi-command sequence
vetol check -m blacklist -r "docker compose exec app rm,rm -rf" "docker compose exec app rm -rf /"
```

### Using Configuration File

Rules can be loaded from a JSON configuration file:

```bash
vetol check --config rules.json "docker compose exec app rm -rf /"
```

**Configuration file format:**

```json
{
  "mode": "blacklist",
  "rules": [
    "docker compose exec app rm",
    "docker compose exec app mkfs",
    "docker run rm -rf",
    "rm -rf",
    "dd"
  ]
}
```

## Command Line Options

- `--mode <mode>`, `-m <mode>`: Security validation mode (required unless using `--config`)
  - `whitelist`: Only allow explicitly permitted commands
  - `blacklist`: Allow all commands except explicitly forbidden ones
- `--rules <RULES>`, `-r <RULES>`: Comma-separated list of rules (required unless using `--config`)
  - Each rule can be a single command or a space-separated sequence of commands
- `--config <PATH>`: Path to a JSON configuration file (alternative to `--mode` and `--rules`)
- `<COMMAND_STRING>`: The bash command string to validate (positional argument)

## Features

✓ **AST-based parsing**: Comprehensive detection of nested and obfuscated commands
✓ **Whitelist/Blacklist modes**: Flexible security policy configuration
✓ **Complex syntax support**: Handles pipes, substitutions, chains, redirects, and subshells
✓ **JSON configuration**: Load rules from configuration files
✓ **Command sequences**: Support for multi-command rules (e.g., `docker compose exec app rm`)

## Examples

### Example 1: Detect Hidden Commands

```bash
# This would pass simple string matching but Vetol detects it:
vetol check -m whitelist -r "echo,cat" 'echo $(rm -rf /)'
# Output: INVALID (because 'rm' is nested in command substitution)
```

### Example 2: Validate Complex Pipelines

```bash
vetol check -m whitelist -r "pwd,grep" "pwd | grep test"
# Output: VALID

vetol check -m whitelist -r "pwd,grep" "pwd | grep test | rm"
# Output: INVALID (because 'rm' is in the pipeline)
```

### Example 3: Command Sequences

```bash
# Allow 'docker compose exec app go' but not 'docker compose exec app rm'
vetol check -m whitelist -r "docker compose exec app go" "docker compose exec app go fmt"
# Output: VALID

vetol check -m whitelist -r "docker compose exec app go" "docker compose exec app rm"
# Output: INVALID
```

## Development

### Prerequisites

- Go 1.26 or later

### Running Tests

```bash
bash tests/test.sh
```

Or run Go tests:

```bash
go test ./...
```

### Build

```bash
go build -o vetol ./cmd
```

## Architecture

```
vetol/
├── cmd/
│   └── vetol/
│       └── main.go          # CLI entry point
├── internal/
│   ├── logger/              # Structured logging utilities
│   ├── parser/              # Bash AST parsing logic
│   └── validator/           # Security validation logic
├── pkg/
│   └── rules/               # Rule management and configuration
├── testdata/                # Test configuration files
├── tests/
│   └── test.sh              # Integration tests
├── go.mod
└── go.sum
```

## Dependencies

- **mvdan.cc/sh/v3/syntax**: Bash command parser and AST builder
- Go 1.26 standard library

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
