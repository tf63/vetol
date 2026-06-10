package logger

import (
	"fmt"
	"os"
)

// Error logs an error-level message.
func Error(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, "ERROR: %s", msg)
	if len(args) > 0 {
		fmt.Fprintf(os.Stderr, " %v", args)
	}
	fmt.Fprintf(os.Stderr, "\n")
}

// Warn logs a warning-level message.
func Warn(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, "WARN: %s", msg)
	if len(args) > 0 {
		fmt.Fprintf(os.Stderr, " %v", args)
	}
	fmt.Fprintf(os.Stderr, "\n")
}

// Info logs an info-level message.
func Info(msg string, args ...any) {
	fmt.Printf("INFO: %s", msg)
	if len(args) > 0 {
		fmt.Printf(" %v", args)
	}
	fmt.Printf("\n")
}

// Debug logs a debug-level message.
func Debug(msg string, args ...any) {
	fmt.Printf("DEBUG: %s", msg)
	if len(args) > 0 {
		fmt.Printf(" %v", args)
	}
	fmt.Printf("\n")
}
