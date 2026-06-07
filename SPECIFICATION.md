# Vetol Specification

## Overview

Vetol is a CLI tool that analyzes bash command strings using Abstract Syntax Tree (AST) parsing and performs strict security validation to detect forbidden commands regardless of how they are nested or obfuscated.

**Module**: `github.com/tf63/vetol`
**Language**: Go 1.26
**CLI Command**: `vetol`

## Purpose

Provide robust and strict security validation for bash commands that cannot be bypassed through command chaining, substitution, or other complex syntax patterns.

**Problem Being Solved:**
Simple string matching-based validation is insufficient and can be bypassed:

- Command chaining: `ls && rm -rf /` - bypasses simple "rm" prefix check
- Command substitution: `cat "$(rm -rf /)"` - hides forbidden command inside variable expansion
- Pipelines and other shell constructs - can conceal forbidden operations

**Solution:**

- Parse bash command strings into a complete Abstract Syntax Tree (AST) using `mvdan.cc/sh/v3/syntax`
- Traverse the entire AST to detect all commands, including those nested in:
  - Command chains (`&&`, `||`, `;`)
  - Pipelines (`|`)
  - Command substitution (`$()`, backticks)
  - Other shell constructs
- Perform strict validation against allowlist or denylist rules at the AST level
- Reject any command string containing forbidden commands, regardless of nesting depth or context

## Tools

### Dependencies

- **mvdan.cc/sh/v3/syntax**: Bash command parser and AST builder
- Go 1.26 or later

## CLI

### Command: vetol

The main entry point that provides bash command validation functionality.

```bash
vetol [OPTIONS] <COMMAND_STRING>
```

### Subcommands

- `check` - Analyze and validate a bash command string

## Usage

### Check Command with Allowlist Mode

Parse and validate a bash command string against an allowed command list.

```bash
# Single commands
vetol check --mode allowlist --rules ls,cat,grep "ls -la /tmp"

# Multi-command sequences
vetol check -m allowlist -r "docker compose exec app go fmt,ls,cat" "docker compose exec app go fmt ./..."
```

### Check Command with Denylist Mode

Parse and validate a bash command string against a forbidden command list.

```bash
# Single commands
vetol check --mode denylist --rules rm,dd,mkfs "cat /etc/passwd"

# Multi-command sequences
vetol check -m denylist -r "docker compose exec app rm,rm -rf" "docker compose exec app rm -rf /"
```

### Rule Format

Rules are comma-separated, where each rule can be:

- **Single command**: `ls`, `cat`, `rm`, `dd`
- **Multi-command sequence**: `docker compose exec app rm`, `docker compose exec app go fmt`

Rules use **prefix matching**: a rule matches a command sequence if the command sequence starts with all elements of the rule.

**Prefix Matching Examples:**

- Rule `echo` matches: `echo`, `echo Hello`, `echo 'Hello World'` ✓
- Rule `ls -la` matches: `ls -la`, `ls -la /tmp` ✓
- Rule `docker compose exec app rm` matches: `docker compose exec app rm -rf /` ✓
- Rule `docker compose exec app rm` does NOT match: `docker rm -rf /` ✗
- Rule `docker compose exec app rm` does NOT match: `docker compose exec app go fmt` ✗
- Rule `rm` does NOT match: `rmdir` ✗

### Options

- `--mode <mode>`, `-m <mode>`: Security validation mode (required)
  - `allowlist`: Only allow explicitly permitted commands
  - `denylist`: Allow all commands except explicitly forbidden ones
- `--rules <RULES>`, `-r <RULES>`: Comma-separated list of rules for validation (required)
  - Each rule can be a single command or a space-separated sequence of commands
  - In allowlist mode: list of allowed command/sequences
  - In denylist mode: list of forbidden command/sequences
- `<COMMAND_STRING>`: The bash command string to validate (positional argument)

## Config

### Configuration Methods

Rules can be specified via one of two methods (not both):

1. **Command-line options**: `--mode` and `--rules`
2. **Configuration file**: `--config` (JSON format only)

If `--config` is specified, `--mode` and `--rules` must NOT be provided. Specifying both will result in an error.

### Configuration File Format (JSON)

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

### Usage Examples

**Using command-line options:**

```bash
vetol check -m denylist -r "rm,dd" "cat /etc/passwd"
```

**Using configuration file:**

```bash
vetol check --config rules.json "docker compose exec app rm -rf /"
```

**Error cases:**

```bash
# ERROR: Cannot mix --config with --mode/--rules
vetol check --config rules.json -m denylist "echo test"
```

### Mode Options

- **allowlist**: Only allow explicitly permitted commands/sequences. Commands/sequences not in the allowlist are rejected.
- **denylist**: Allow all commands/sequences except those explicitly forbidden. Commands/sequences in the denylist are rejected.

## Architecture

### Directory Structure

```
vetol/
├── cmd/
│   └── main.go              # CLI entry point
├── internal/
│   ├── parser/              # Bash AST parsing logic
│   └── validator/           # Security validation logic
├── pkg/
│   └── rules/               # Rule management and configuration
├── go.mod
└── go.sum
```

## Development Roadmap

- [ ] Setup go.mod with mvdan.cc/sh/v3/syntax dependency
- [ ] Implement AST parser wrapper
- [ ] Implement allowlist/denylist validator
- [ ] Create CLI command structure
- [ ] Add configuration file support
- [ ] Write comprehensive tests
- [ ] Add documentation and examples
