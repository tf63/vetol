# Vetol

A CLI tool for validating shell commands using Bash AST analysis, suitable for AI agent hooks and other security-sensitive execution environments.

## Overview

Regex-based command validation can be bypassed through shell syntax such as command substitution (`echo $(rm -rf /)`) or command chaining (`ls && rm -rf /`).

For example, AI agents that rely on pattern matching to filter dangerous commands may fail to detect commands hidden inside nested shell constructs.

**Vetol** parses the Bash Abstract Syntax Tree (AST) and validates every command in the tree, including those nested in pipelines, substitutions, chains, and subshells.

## Installation

### Using `go install`

```bash
go install github.com/tf63/vetol/cmd/vetol
```

## Usage

**Vetol** only validates the provided command string and never executes it.

```bash
$ vetol check -m denylist -r rm "ls"
VALID
$ echo $?
0

$ vetol check -m denylist -r rm "rm -rf /"
INVALID
$ echo $?
1
```

### Allowlist Mode

Only allow explicitly permitted commands:

```bash
# Single command
vetol check --mode allowlist --rules ls,cat,grep "ls -la /tmp"

# Multi-command sequence
vetol check -m allowlist -r "docker compose exec app go fmt,ls,cat" "docker compose exec app go fmt ./..."
```

### Denylist Mode

Allow all commands except explicitly forbidden ones:

```bash
# Single command
vetol check --mode denylist --rules rm,dd,mkfs "cat /etc/passwd"

# Multi-command sequence
vetol check -m denylist -r "docker compose exec app rm,rm -rf" "docker compose exec app rm -rf /"
```

### Using Configuration File

Rules can be loaded from a JSON configuration file:

```bash
vetol check --config rules.json "docker compose exec app rm -rf /"
```

**Configuration file format:**

```json
{
  "mode": "denylist",
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
  - `allowlist`: Only allow explicitly permitted commands
  - `denylist`: Allow all commands except explicitly forbidden ones
- `--rules <RULES>`, `-r <RULES>`: Comma-separated list of rules (required unless using `--config`)
  - Each rule can be a single command or a space-separated sequence of commands
- `--config <PATH>`: Path to a JSON configuration file (alternative to `--mode` and `--rules`)
- `<COMMAND_STRING>`: The bash command string to validate (positional argument)

## Features

✓ **AST-based parsing**: Improved detection of nested and obfuscated commands
✓ **Allowlist/Denylist modes**: Flexible security policy configuration
✓ **Prefix matching**: Rules match command sequences that start with the rule pattern (e.g., rule `echo` matches `echo`, `echo arg1`, `echo arg1 arg2`, etc.)
✓ **Complex syntax support**: Handles pipes, substitutions, chains, redirects, and subshells
✓ **JSON configuration**: Load rules from configuration files
✓ **Command sequences**: Support for multi-command rules (e.g., `docker compose exec app rm`)

## Limitations

Vetol validates command structure through Bash AST analysis.

It does not:

- Analyze command arguments semantically
- Detect malicious behavior inside allowed binaries
- Execute commands
- Provide sandboxing or isolation

## Examples

### Example 1: Simple Commands

```bash
# Allowlist mode: allow only specific commands
vetol check -m allowlist -r "cat,ls" "cat /etc/passwd"
# Output: VALID

vetol check -m allowlist -r "cat,ls" "cat /etc/passwd && rm file.txt"
# Output: INVALID

# Denylist mode: forbid specific commands
vetol check -m denylist -r "rm" "cat /etc/passwd"
# Output: VALID

vetol check -m denylist -r "rm" "rm -rf /"
# Output: INVALID
```

### Example 2: Commands with Options

```bash
# Allowlist specific command with options
vetol check -m allowlist -r "docker" "docker ps"
# Output: VALID

# Forbid specific command options
vetol check -m denylist -r "docker run" "docker build -t image ."
# Output: VALID

vetol check -m denylist -r "docker run" "docker run -it ubuntu"
# Output: INVALID
```

### Example 3: Multiple Commands

```bash
# Allow multi-command sequences
vetol check -m allowlist -r "pwd,ls,cat" "pwd && ls -la && cat file.txt"
# Output: VALID

# Forbid multi-command sequences
vetol check -m denylist -r "rm,dd" "ls && pwd && echo done"
# Output: VALID

vetol check -m denylist -r "rm,dd" "pwd && rm -rf /"
# Output: INVALID
```

### Example 4: Obfuscated Commands

```bash
# Command substitution: simple string matching would miss this
vetol check -m allowlist -r "echo,cat" 'echo $(rm -rf /)'
# Output: INVALID

# Command chaining: prefix matching would miss this
vetol check -m denylist -r "rm" "ls && rm -rf /"
# Output: INVALID

# Pipelines: checks all piped commands
vetol check -m allowlist -r "pwd,grep" "pwd | grep test | rm"
# Output: INVALID
```

## Dependencies

- [**mvdan.cc/sh/v3/syntax**](https://github.com/mvdan/sh): Bash command parser and AST builder
- Go 1.26 standard library

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
go build -o vetol ./cmd/vetol
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

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
