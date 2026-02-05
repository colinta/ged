package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/colinta/ged/internal/parser"
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run executes ged with the given arguments and I/O streams.
// This is separated from main() for testability.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ged <rule>")
	}

	// Parse the rule from command line
	ruleStr := args[0]
	rule, err := parser.ParseRule(ruleStr)
	if err != nil {
		return fmt.Errorf("error parsing rule: %w", err)
	}

	// Read stdin line by line
	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		line := scanner.Text()

		// Apply the rule
		results, err := rule.Apply(line)
		if err != nil {
			return fmt.Errorf("error applying rule: %w", err)
		}

		// Print each result line
		for _, result := range results {
			fmt.Fprintln(stdout, result)
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	return nil
}
