# Contributing to Vetol

Thank you for your interest in contributing to Vetol! This document provides guidelines and instructions for setting up your development environment.

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
│       └── main.go              # CLI entry point
├── internal/
│   ├── parser/                  # Bash AST parsing logic
│   ├── rules/                   # Rule matching and validation (internal)
│   │   ├── rules.go             # Rule matching logic
│   │   ├── config.go            # Configuration validation
│   │   └── *_test.go            # Tests (100% coverage)
│   └── validator/               # High-level validation orchestration
├── pkg/
│   ├── io/                      # Configuration file I/O
│   │   └── config.go            # JSON config loading
│   └── logger/                  # Logging utilities (public API)
│       └── logger.go            # Error, Warn, Info, Debug functions
├── testdata/                    # Test configuration files
├── tests/
│   └── test.sh                  # Integration tests (97 test cases)
├── go.mod
└── go.sum
```

### Module Responsibilities

- **cmd/vetol**: CLI entry point - loads config and orchestrates validation
- **pkg/io**: Public configuration file parsing
- **pkg/logger**: Public logging utilities
- **internal/rules**: Rule matching and validation (not exported)
- **internal/parser**: Bash AST parsing (not exported)
- **internal/validator**: Validation orchestration (not exported)
