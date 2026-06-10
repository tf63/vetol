package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tf63/vetol/internal/validator"
	"github.com/tf63/vetol/pkg/io"
	"github.com/tf63/vetol/pkg/logger"
)

func main() {
	// Check for help flag
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h" || os.Args[1] == "help") {
		printUsage()
		os.Exit(0)
	}

	// Check for help flag in arguments
	for _, arg := range os.Args[1:] {
		if arg == "--help" || arg == "-h" {
			printUsage()
			os.Exit(0)
		}
	}

	fs := flag.NewFlagSet("vetol", flag.ContinueOnError)
	configPath := fs.String("config", "", "Path to configuration file")

	if err := fs.Parse(os.Args[1:]); err != nil {
		logger.Error("failed to parse flags", "error", err)
		os.Exit(1)
	}

	// Validate --config is provided
	if *configPath == "" {
		logger.Error("--config is required")
		printUsage()
		os.Exit(1)
	}

	// Get positional argument (command string)
	args := fs.Args()
	if len(args) == 0 {
		logger.Error("command string is required")
		printUsage()
		os.Exit(1)
	}

	commandStr := args[0]

	// Load configuration
	cfg, err := io.LoadConfigFromFile(*configPath)
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Validate command
	val := validator.NewValidator(&cfg)
	result, err := val.Validate(commandStr)
	if err != nil {
		logger.Error("failed to validate command", "error", err)
		os.Exit(1)
	}

	// Output result
	if result.Valid {
		fmt.Println("ALLOW")
		os.Exit(0)
	} else {
		fmt.Println("DENY")
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Vetol - Vet agent tools

Usage: vetol --config <PATH> <COMMAND_STRING>

Options:
  --config <PATH>          Path to configuration file (JSON) - REQUIRED
  --help, -h              Show this help message

Arguments:
  <COMMAND_STRING>         The bash command string to validate

Configuration File Format:
{
  "mode": "allowlist",
  "rules": [
    {
      "command": "grep",
      "include": ["-r"],
      "exclude": ["-rule"]
    }
  ]
}

Mode Options:
  allowlist               Only allow explicitly permitted commands
  denylist               Allow all commands except explicitly forbidden ones

Rule Fields:
  command                 Command name (prefix matching)
  include                 Array of allowed options or patterns
  exclude                 Array of forbidden options or patterns

Flag Types:
  -X                      Short flag (character containment matching)
  --flag                  Long flag (exact match)
  flag                    Option without prefix (exact match)

Examples:
  vetol --config vetol.json "grep -r pattern file.txt"
  vetol --config vetol.json "docker compose exec app ls -la"
  vetol --help             Show this help message
`)
}
