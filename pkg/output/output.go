package output

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

// Verbosity levels
const (
	LevelQuiet   = 0 // Suppress all non-essential output
	LevelNormal  = 1 // Default output level
	LevelVerbose = 2 // Show additional details
	LevelDebug   = 3 // Show debug information
)

var (
	verbosity           = LevelNormal
	stdout    io.Writer = os.Stdout
	stderr    io.Writer = os.Stderr
)

// SetVerbosity sets the global verbosity level
func SetVerbosity(level int) {
	if level < LevelQuiet {
		level = LevelQuiet
	}
	if level > LevelDebug {
		level = LevelDebug
	}
	verbosity = level
}

// GetVerbosity returns the current verbosity level
func GetVerbosity() int {
	return verbosity
}

// IsQuiet returns true if verbosity is at quiet level (0)
func IsQuiet() bool {
	return verbosity == LevelQuiet
}

// IsNormal returns true if verbosity is at least normal level (1)
func IsNormal() bool {
	return verbosity >= LevelNormal
}

// IsVerbose returns true if verbosity is at least verbose level (2)
func IsVerbose() bool {
	return verbosity >= LevelVerbose
}

// IsDebug returns true if verbosity is at debug level (3)
func IsDebug() bool {
	return verbosity >= LevelDebug
}

// Print prints to stdout if not in quiet mode
func Print(format string, a ...interface{}) {
	if IsNormal() {
		fmt.Fprintf(stdout, format, a...)
	}
}

// Println prints a line to stdout if not in quiet mode
func Println(a ...interface{}) {
	if IsNormal() {
		fmt.Fprintln(stdout, a...)
	}
}

// PrintVerbose prints only when verbosity is >= 2
func PrintVerbose(format string, a ...interface{}) {
	if IsVerbose() {
		fmt.Fprintf(stdout, format, a...)
	}
}

// PrintDebug prints debug information to stderr when verbosity is >= 3
func PrintDebug(format string, a ...interface{}) {
	if IsDebug() {
		fmt.Fprintf(stderr, "[DEBUG] "+format, a...)
	}
}

// PrintError always prints errors regardless of verbosity
func PrintError(format string, a ...interface{}) {
	fmt.Fprintf(stderr, format, a...)
}

// Success prints a success message (green) if not in quiet mode
func Success(format string, a ...interface{}) {
	if IsNormal() {
		green := color.New(color.FgGreen, color.Bold)
		green.Fprintf(stdout, format, a...)
	}
}

// Warning prints a warning message (yellow) if not in quiet mode
func Warning(format string, a ...interface{}) {
	if IsNormal() {
		yellow := color.New(color.FgYellow)
		yellow.Fprintf(stdout, format, a...)
	}
}

// Info prints an info message (cyan) if not in quiet mode
func Info(format string, a ...interface{}) {
	if IsNormal() {
		cyan := color.New(color.FgCyan, color.Bold)
		cyan.Fprintf(stdout, format, a...)
	}
}

// Header prints a header box if not in quiet mode
func Header(title string) {
	if IsNormal() {
		cyan := color.New(color.FgCyan, color.Bold)
		fmt.Println()
		cyan.Println("╔════════════════════════════════════════════════════════════╗")
		cyan.Printf("║          %-47s ║\n", title)
		cyan.Println("╚════════════════════════════════════════════════════════════╝")
		fmt.Println()
	}
}

// Section prints a section header if not in quiet mode
func Section(emoji, title string) {
	if IsNormal() {
		yellow := color.New(color.FgYellow, color.Bold)
		yellow.Printf("%s %s\n", emoji, title)
	}
}
