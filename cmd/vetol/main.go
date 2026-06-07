package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tf63/vetol/internal/logger"
	"github.com/tf63/vetol/internal/validator"
	"github.com/tf63/vetol/pkg/rules"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Check for help flag
	if os.Args[1] == "--help" || os.Args[1] == "-h" || os.Args[1] == "help" {
		printUsage()
		os.Exit(0)
	}

	subcommand := os.Args[1]

	if subcommand != "check" {
		logger.Error("unknown subcommand", "subcommand", subcommand)
		printUsage()
		os.Exit(1)
	}

	// Check for help flag in check subcommand
	for _, arg := range os.Args[2:] {
		if arg == "--help" || arg == "-h" {
			printCheckUsage()
			os.Exit(0)
		}
	}

	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	mode := fs.String("mode", "", "Security validation mode (whitelist or blacklist)")
	modeShort := fs.String("m", "", "Short flag for mode")
	rulesStr := fs.String("rules", "", "Comma-separated list of rules")
	rulesShort := fs.String("r", "", "Short flag for rules")
	configPath := fs.String("config", "", "Path to configuration file")

	if err := fs.Parse(os.Args[2:]); err != nil {
		logger.Error("failed to parse flags", "error", err)
		os.Exit(1)
	}

	// Handle short flags
	if *modeShort != "" {
		mode = modeShort
	}
	if *rulesShort != "" {
		rulesStr = rulesShort
	}

	// Get positional argument (command string)
	args := fs.Args()
	if len(args) == 0 {
		logger.Error("command string is required")
		printUsage()
		os.Exit(1)
	}

	commandStr := args[0]

	// Validate configuration options
	if *configPath != "" && (*mode != "" || *rulesStr != "") {
		logger.Error("cannot mix --config with --mode/--rules")
		printUsage()
		os.Exit(1)
	}

	if *configPath == "" && (*mode == "" || *rulesStr == "") {
		logger.Error("either --config or both --mode and --rules are required")
		printUsage()
		os.Exit(1)
	}

	// Validate mode
	if *configPath == "" {
		modeValue := rules.Mode(*mode)
		if modeValue != rules.ModeWhitelist && modeValue != rules.ModeBlacklist {
			logger.Error("invalid mode", "mode", *mode)
			os.Exit(1)
		}
	}

	// Load configuration
	var cfg rules.Config
	var err error

	if *configPath != "" {
		cfg, err = rules.LoadConfigFromFile(*configPath)
		if err != nil {
			logger.Error("failed to load config", "error", err)
			os.Exit(1)
		}
	} else {
		rulesList := strings.Split(*rulesStr, ",")
		modeValue := rules.Mode(*mode)
		cfg = rules.NewConfig(modeValue, rulesList)
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
		if len(result.ViolatedCommands) > 0 {
			logger.Error("violated commands", "commands", result.ViolatedCommands)
		}
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Vetol - Security command validator

Usage: vetol [COMMAND] [OPTIONS]

Commands:
  check                    Analyze and validate a bash command string
  help, -h, --help        Show this help message

Examples:
  vetol check -m whitelist -r "ls,cat,grep" "ls -la /tmp"
  vetol check -m blacklist -r "rm,dd" "cat /etc/passwd"
  vetol check --config rules.json "docker compose exec app rm -rf /"
  vetol --help             Show this help message

Use 'vetol check --help' for more information about the check command.
`)
}

func printCheckUsage() {
	fmt.Fprintf(os.Stderr, `Usage: vetol check [OPTIONS] <COMMAND_STRING>

Options:
  --mode, -m <mode>        Security validation mode (whitelist or blacklist)
  --rules, -r <RULES>      Comma-separated list of rules
  --config <PATH>          Path to configuration file (JSON)
  --help, -h              Show this help message

Arguments:
  <COMMAND_STRING>         The bash command string to validate

Configuration Methods (use one):
  1. With --mode and --rules flags:
     vetol check -m whitelist -r "ls,cat,grep" "ls -la /tmp"

  2. With --config file:
     vetol check --config rules.json "docker compose exec app rm -rf /"

Mode Options:
  whitelist               Only allow explicitly permitted commands
  blacklist               Allow all commands except explicitly forbidden ones

Rule Format:
  Single command:         ls, cat, rm
  Command sequence:       docker compose exec app rm, docker run rm -rf
  Comma-separated:        ls,cat,grep or "docker compose exec app,cat"

Examples:
  vetol check -m whitelist -r "ls,cat,grep" "ls -la /tmp"
  vetol check -m blacklist -r "rm,dd" "cat /etc/passwd"
  vetol check --config rules.json "docker compose exec app rm -rf /"
`)
}
