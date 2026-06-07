package parser

import (
	"fmt"
	"strings"

	"mvdan.cc/sh/v3/syntax"
)

// CommandSequence represents a sequence of commands found in the AST.
type CommandSequence struct {
	Commands []string
}

// Parser is responsible for parsing bash command strings into ASTs and extracting commands.
type Parser struct{}

// NewParser creates a new Parser instance.
func NewParser() *Parser {
	return &Parser{}
}

// ExtractCommandSequences parses a bash command string and extracts all command sequences.
func (p *Parser) ExtractCommandSequences(commandStr string) ([]CommandSequence, error) {
	reader := strings.NewReader(commandStr)
	file, err := syntax.NewParser().Parse(reader, "")
	if err != nil {
		return nil, fmt.Errorf("failed to parse command: %w", err)
	}

	var sequences []CommandSequence
	syntax.Walk(file, func(node syntax.Node) bool {
		switch n := node.(type) {
		case *syntax.CallExpr:
			p.extractFromCallExpr(n, &sequences)
		case *syntax.CmdSubst:
			p.extractFromCmdSubst(n, &sequences)
		case *syntax.ProcSubst:
			p.extractFromProcSubst(n, &sequences)
		}
		return true
	})

	return sequences, nil
}

// extractFromCallExpr extracts command information from a CallExpr node.
func (p *Parser) extractFromCallExpr(call *syntax.CallExpr, sequences *[]CommandSequence) {
	if call == nil || len(call.Args) == 0 {
		return
	}

	var commands []string

	// Extract all values from all arguments
	for _, arg := range call.Args {
		for _, part := range arg.Parts {
			switch v := part.(type) {
			case *syntax.Lit:
				commands = append(commands, v.Value)
			case *syntax.DblQuoted:
				// Extract content from double quoted strings
				if v != nil {
					for _, p := range v.Parts {
						if lit, ok := p.(*syntax.Lit); ok {
							commands = append(commands, lit.Value)
						}
					}
				}
			case *syntax.SglQuoted:
				// Extract content from single quoted strings
				if v != nil {
					commands = append(commands, v.Value)
				}
			}
		}
	}

	if len(commands) > 0 {
		*sequences = append(*sequences, CommandSequence{Commands: commands})
	}
}

// extractFromCmdSubst extracts commands from a command substitution.
func (p *Parser) extractFromCmdSubst(cmdSubst *syntax.CmdSubst, sequences *[]CommandSequence) {
	if cmdSubst == nil {
		return
	}

	// Recursively parse command substitutions
	for _, stmt := range cmdSubst.Stmts {
		if stmt != nil && stmt.Cmd != nil {
			if callExpr, ok := stmt.Cmd.(*syntax.CallExpr); ok {
				p.extractFromCallExpr(callExpr, sequences)
			}
		}
	}
}

// extractFromProcSubst extracts commands from a process substitution.
func (p *Parser) extractFromProcSubst(procSubst *syntax.ProcSubst, sequences *[]CommandSequence) {
	if procSubst == nil {
		return
	}

	// Recursively parse process substitutions
	for _, stmt := range procSubst.Stmts {
		if stmt != nil && stmt.Cmd != nil {
			if callExpr, ok := stmt.Cmd.(*syntax.CallExpr); ok {
				p.extractFromCallExpr(callExpr, sequences)
			}
		}
	}
}
