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

	subcommand := os.Args[1]

	if subcommand != "check" {
		logger.Error("unknown subcommand", "subcommand", subcommand)
		printUsage()
		os.Exit(1)
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
		fmt.Println("VALID")
		os.Exit(0)
	} else {
		fmt.Println("INVALID")
		if len(result.ViolatedCommands) > 0 {
			logger.Error("violated commands", "commands", result.ViolatedCommands)
		}
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: vetol check [OPTIONS] <COMMAND_STRING>

Options:
  --mode, -m <mode>        Security validation mode (whitelist or blacklist)
  --rules, -r <RULES>      Comma-separated list of rules
  --config <PATH>          Path to configuration file (JSON)

Examples:
  vetol check -m whitelist -r "ls,cat,grep" "ls -la /tmp"
  vetol check -m blacklist -r "rm,dd" "cat /etc/passwd"
  vetol check --config rules.json "docker compose exec rm -rf /"
`)
}
