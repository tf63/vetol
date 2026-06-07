package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"mvdan.cc/sh/v3/syntax"
)

func main() {
	flag.Parse()

	input := flag.Arg(0)
	if input == "" {
		printUsage()
		os.Exit(1)
	}

	// Example 1: Parse and print AST
	fmt.Println("=== Example 1: Parse and Print AST ===")
	parseAndPrintAST(input)

	fmt.Println("\n=== Example 2: Extract Commands ===")
	extractCommands(input)

	fmt.Println("\n=== Example 3: Walk through AST ===")
	walkAST(input)

	fmt.Println("\n=== Example 4: Get Word Positions ===")
	getWordPositions(input)
}

// parseAndPrintAST parses the input and prints the AST structure
func parseAndPrintAST(input string) {
	file, err := syntax.NewParser().Parse(strings.NewReader(input), "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		return
	}

	fmt.Printf("Parsed file has %d statements\n", len(file.Stmts))
	for i, stmt := range file.Stmts {
		fmt.Printf("Statement %d: %T\n", i, stmt.Cmd)
		printNode(stmt.Cmd, 1)
	}
}

// extractCommands extracts all commands from the input
func extractCommands(input string) {
	file, err := syntax.NewParser().Parse(strings.NewReader(input), "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		return
	}

	var commands []string
	syntax.Walk(file, func(node syntax.Node) bool {
		switch cmd := node.(type) {
		case *syntax.CallExpr:
			if len(cmd.Args) > 0 {
				word := cmd.Args[0]
				if len(word.Parts) > 0 {
					if lit, ok := word.Parts[0].(*syntax.Lit); ok {
						commands = append(commands, lit.Value)
					}
				}
			}
		}
		return true
	})

	fmt.Printf("Found %d commands:\n", len(commands))
	for i, cmd := range commands {
		fmt.Printf("  %d. %s\n", i+1, cmd)
	}
}

// walkAST walks through all nodes in the AST
func walkAST(input string) {
	file, err := syntax.NewParser().Parse(strings.NewReader(input), "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		return
	}

	nodeCount := 0
	syntax.Walk(file, func(node syntax.Node) bool {
		nodeCount++
		fmt.Printf("  Node %d: %T\n", nodeCount, node)
		return true
	})
	fmt.Printf("Total nodes: %d\n", nodeCount)
}

// getWordPositions gets the positions of words in the input
func getWordPositions(input string) {
	file, err := syntax.NewParser().Parse(strings.NewReader(input), "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		return
	}

	var words []*syntax.Word
	syntax.Walk(file, func(node syntax.Node) bool {
		if word, ok := node.(*syntax.Word); ok {
			words = append(words, word)
		}
		return true
	})

	fmt.Printf("Found %d words:\n", len(words))
	for i, word := range words {
		// Get the word text from the original input
		start := uint(word.Pos().Offset())
		end := uint(word.End().Offset())
		if start < uint(len(input)) && end <= uint(len(input)) {
			wordText := input[start:end]
			fmt.Printf("  %d. Position: %d-%d, Text: %s\n", i+1, start, end, wordText)
		}
	}
}

// printNode prints the AST node structure recursively
func printNode(node syntax.Node, indent int) {
	if node == nil {
		return
	}

	prefix := strings.Repeat("  ", indent)

	switch n := node.(type) {
	case *syntax.CallExpr:
		fmt.Printf("%sCallExpr: %d args\n", prefix, len(n.Args))
		for i, arg := range n.Args {
			fmt.Printf("%s  Arg %d:\n", prefix, i)
			printNode(arg, indent+2)
		}

	case *syntax.Word:
		fmt.Printf("%sWord: %d parts\n", prefix, len(n.Parts))
		for i, part := range n.Parts {
			fmt.Printf("%s  Part %d:\n", prefix, i)
			printNode(part, indent+2)
		}

	case *syntax.Lit:
		fmt.Printf("%sLit: %q\n", prefix, n.Value)

	case *syntax.DblQuoted:
		fmt.Printf("%sDblQuoted: %d parts\n", prefix, len(n.Parts))
		for i, part := range n.Parts {
			fmt.Printf("%s  Part %d:\n", prefix, i)
			printNode(part, indent+2)
		}

	case *syntax.SglQuoted:
		fmt.Printf("%sSglQuoted: %q\n", prefix, n.Value)

	case *syntax.Redirect:
		fmt.Printf("%sRedirect: Op=%v\n", prefix, n.Op)
		if n.N != nil {
			fmt.Printf("%s  N:\n", prefix)
			printNode(n.N, indent+2)
		}
		fmt.Printf("%s  Target:\n", prefix)
		printNode(n.Word, indent+2)

	case *syntax.BinaryCmd:
		fmt.Printf("%sBinaryCmd: Op=%v\n", prefix, n.Op)
		fmt.Printf("%s  Left:\n", prefix)
		printNode(n.X, indent+2)
		fmt.Printf("%s  Right:\n", prefix)
		printNode(n.Y, indent+2)

	case *syntax.Subshell:
		fmt.Printf("%sSubshell: %d statements\n", prefix, len(n.Stmts))
		for i, stmt := range n.Stmts {
			fmt.Printf("%s  Stmt %d:\n", prefix, i)
			printNode(stmt.Cmd, indent+2)
		}

	case *syntax.Block:
		fmt.Printf("%sBlock: %d statements\n", prefix, len(n.Stmts))
		for i, stmt := range n.Stmts {
			fmt.Printf("%s  Stmt %d:\n", prefix, i)
			printNode(stmt.Cmd, indent+2)
		}

	case *syntax.IfClause:
		fmt.Printf("%sIfClause:\n", prefix)
		if len(n.Cond) > 0 {
			fmt.Printf("%s  Condition:\n", prefix)
			for _, stmt := range n.Cond {
				printNode(stmt.Cmd, indent+2)
			}
		}
		if len(n.Then) > 0 {
			fmt.Printf("%s  Then:\n", prefix)
			for _, stmt := range n.Then {
				printNode(stmt.Cmd, indent+2)
			}
		}

	case io.WriterTo:
		fmt.Printf("%s%T (custom node)\n", prefix, n)

	default:
		fmt.Printf("%s%T\n", prefix, n)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Demo tool for mvdan.cc/sh/v3/syntax package

Usage: demo <COMMAND>

Examples:
  demo "ls -la /tmp"
  demo "echo hello | grep hello"
  demo "if [ -f file ]; then cat file; fi"
`)
}
