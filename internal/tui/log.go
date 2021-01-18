package tui

import (
	"fmt"
	"os"
)

var Verbose bool

func Debug(format string, a ...interface{}) {
	if Verbose {
		format = fmt.Sprintf("DEBUG: %s\n", format)
		fmt.Fprintf(os.Stderr, ApplyStylef(Gray, format, a...))
	}
}

func Warn(format string, a ...interface{}) {
	format = fmt.Sprintf("%s %s\n", ApplyStyle(Yellow, "WARNING:"), format)
	fmt.Fprintf(os.Stderr, format, a...)
}

func Print(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func FatalIfError(message string, err error) {
	if err != nil {
		if message != "" {
			fmt.Fprintln(os.Stderr, message)
		}
		fmt.Fprintf(os.Stderr, ApplyStylef(Red, "Error: %v\n", err))
		os.Exit(255)
	}
}

func Fatal(message string) {
	fmt.Fprintf(os.Stderr, ApplyStylef(Red, "Error: %v\n", message))
	os.Exit(255)
}
