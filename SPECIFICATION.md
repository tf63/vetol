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
vetol --config <PATH> <COMMAND_STRING>
```

### Options

- `--config <PATH>`: Path to configuration file (JSON format) - REQUIRED
- `<COMMAND_STRING>`: The bash command string to validate (positional argument)

## Configuration

### Configuration File Format (JSON)

Configuration is provided via a JSON file with the following structure:

```json
{
  "mode": "allowlist",
  "rules": [
    {
      "command": "grep",
      "include": ["-r"],
      "exclude": []
    },
    {
      "command": "ls",
      "include": ["-la"],
      "exclude": ["-i"]
    },
    {
      "command": "docker compose"
    }
  ]
}
```

#### Rule Fields

- **command** (string, required): The command name or prefix (space-separated for multi-word commands like `docker compose`)
- **include** (array of strings, optional): Patterns or flags that MUST be present for the rule to match
  - If specified, all patterns in the include list must match
  - If empty or omitted, no include constraints apply
- **exclude** (array of strings, optional): Patterns or flags that MUST NOT be present for the rule to match
  - If any pattern in the exclude list matches, the rule does not match
  - If empty or omitted, no exclude constraints apply

#### Pattern Matching Rules

**Short Flags** (starting with `-` but not `--`):

- All characters in the pattern must be present in the command
- Characters can be combined in a single option or spread across multiple options
- Examples:
  - Pattern `-r` matches: `-r`, `-rl`, `-r -l`
  - Pattern `-la` matches: `-la`, `-l -a`

**Long Flags** (starting with `--`):

- Exact match or prefix match with `=` for values
- Examples:
  - Pattern `--color` matches: `--color`, `--color=auto`
  - Pattern `--color=auto` matches only: `--color=auto`

**Non-Flag Patterns** (no prefix):

- Exact match required
- Examples:
  - Pattern `file.txt` matches only: `file.txt`

#### Matching Logic

1. **Command Matching**: Command sequence must start with the rule's command
2. **Include Matching**: If include patterns are specified, ALL must match
3. **Exclude Matching**: If exclude patterns are specified, NONE must match
4. A rule matches only if command matches AND include matches AND exclude does NOT match

### Usage Examples

**Allowlist with include/exclude:**

```bash
vetol --config allowlist.json "grep -r pattern /tmp"
```

Where `allowlist.json` contains:

```json
{
  "mode": "allowlist",
  "rules": [
    {
      "command": "grep",
      "include": ["-r"]
    },
    {
      "command": "ls",
      "include": ["-la"],
      "exclude": ["-i"]
    }
  ]
}
```

**Denylist with complex rules:**

```bash
vetol --config denylist.json "docker compose exec app rm -rf /"
```

Where `denylist.json` contains:

```json
{
  "mode": "denylist",
  "rules": [
    {
      "command": "rm",
      "exclude": ["-i"]
    },
    {
      "command": "dd"
    },
    {
      "command": "docker compose exec app rm"
    }
  ]
}
```

### Mode Options

- **allowlist**: Only allow explicitly permitted commands. Commands not matching any rule are rejected.
- **denylist**: Allow all commands except those explicitly forbidden. Commands matching any rule are rejected.

## Architecture

### Directory Structure

```
vetol/
├── cmd/
│   └── vetol/
│       └── main.go              # CLI entry point
├── internal/
│   ├── parser/                  # Bash AST parsing logic
│   ├── rules/                   # Rule management (internal use only)
│   │   ├── rules.go             # Rule struct and matching logic
│   │   ├── rules_test.go        # Rule tests
│   │   ├── config.go            # Config struct and validation methods
│   │   └── config_test.go       # Config tests
│   └── validator/               # Security validation logic
├── pkg/
│   ├── io/
│   │   └── config.go            # Configuration file loading (JSON)
│   └── logger/                  # Logging utilities
├── tests/
│   └── test.sh                  # Integration tests
├── testdata/                    # Test configuration files
├── go.mod
└── go.sum
```

### Module Organization

- **cmd/vetol**: Entry point that orchestrates configuration loading and validation
- **pkg/io**: Handles JSON configuration file parsing
  - `LoadConfigFromFile()`: Loads configuration from JSON files
- **pkg/logger**: Logging utilities (public API)
  - `Error()`: Log error-level messages
  - `Warn()`: Log warning-level messages
  - `Info()`: Log informational messages
  - `Debug()`: Log debug-level messages
- **internal/rules**: Core rule matching and validation logic (internal use)
  - `Rule.Matches()`: Checks if a command matches a rule
  - `Config.IsValid()`: Validates a command against a configuration
- **internal/parser**: AST parsing for bash commands
- **internal/validator**: High-level command validation orchestration

## Implementation Status

### Completed Features

- [x] Setup go.mod with mvdan.cc/sh/v3/syntax dependency
- [x] Implement AST parser wrapper for bash command parsing
- [x] Implement allowlist/denylist validator with include/exclude patterns
- [x] Create CLI command structure with configuration file support
- [x] Add JSON configuration file support
- [x] Implement comprehensive rule matching system
  - [x] Command prefix matching for single and multi-word commands
  - [x] Short flag pattern matching (character containment)
  - [x] Long flag pattern matching (prefix with = support)
  - [x] Non-flag pattern exact matching
  - [x] Include and exclude constraints
- [x] Write comprehensive tests (100% coverage)
- [x] Add documentation and specifications
