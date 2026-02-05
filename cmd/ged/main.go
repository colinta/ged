package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/colinta/ged/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: ged <rule>")
		os.Exit(1)
	}

	// Parse the rule from command line
	ruleStr := os.Args[1]
	rule, err := parser.ParseRule(ruleStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing rule: %v\n", err)
		os.Exit(1)
	}

	// Read stdin line by line
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		// Apply the rule
		results, err := rule.Apply(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error applying rule: %v\n", err)
			os.Exit(1)
		}

		// Print each result line
		for _, result := range results {
			fmt.Println(result)
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(1)
	}
}
