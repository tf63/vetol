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
	mode := fs.String("mode", "", "Security validation mode (allowlist or denylist)")
	modeShort := fs.String("m", "", "Short flag for mode")
	rulesStr := fs.String("rules", "", "Comma-separated list of rules")
	rulesShort := fs.String("r", "", "Short flag for rules")
	configPath := fs.String("config", "", "Path to configuration file")

	if err := fs.Parse(os.Args[1:]); err != nil {
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
		if modeValue != rules.ModeAllowlist && modeValue != rules.ModeDenylist {
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
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Vetol - Security command validator

Usage: vetol [OPTIONS] <COMMAND_STRING>

Options:
  --mode, -m <mode>        Security validation mode (allowlist or denylist)
  --rules, -r <RULES>      Comma-separated list of rules
  --config <PATH>          Path to configuration file (JSON)
  --help, -h              Show this help message

Arguments:
  <COMMAND_STRING>         The bash command string to validate

Configuration Methods (use one):
  1. With --mode and --rules flags:
     vetol -m allowlist -r "ls,cat,grep" "ls -la /tmp"

  2. With --config file:
     vetol --config rules.json "docker compose exec app rm -rf /"

Mode Options:
  allowlist               Only allow explicitly permitted commands
  denylist               Allow all commands except explicitly forbidden ones

Rule Format:
  Single command:         ls, cat, rm
  Command sequence:       docker compose exec app rm, docker run rm -rf
  Comma-separated:        ls,cat,grep or "docker compose exec app,cat"

Examples:
  vetol -m allowlist -r "ls,cat,grep" "ls -la /tmp"
  vetol -m denylist -r "rm,dd" "cat /etc/passwd"
  vetol --config rules.json "docker compose exec app rm -rf /"
  vetol --help             Show this help message
`)
}
