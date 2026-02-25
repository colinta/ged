package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/colinta/ged/internal/engine"
	"github.com/colinta/ged/internal/parser"
	"github.com/colinta/ged/internal/rule"
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
		return fmt.Errorf("usage: ged <rule> [rule...]")
	}

	// Parse all rules, handling { } blocks for conditionals.
	allParsed, err := parser.ParseArgs(args)
	if err != nil {
		return fmt.Errorf("error parsing rules: %w", err)
	}

	// Build a list of DocumentRules.
	// Consecutive LineRules are wrapped in an ApplyAllRule.
	var docRules []rule.DocumentRule
	var pendingLineRules []rule.LineRule

	for _, parsed := range allParsed {
		switch r := parsed.(type) {
		case rule.LineRule:
			pendingLineRules = append(pendingLineRules, r)
		case rule.DocumentRule:
			// Flush any pending line rules into an ApplyAllRule first
			if len(pendingLineRules) > 0 {
				docRules = append(docRules, rule.NewApplyAllRule(pendingLineRules))
				pendingLineRules = nil
			}
			docRules = append(docRules, r)
		default:
			return fmt.Errorf("unknown rule type from parser: %T", parsed)
		}
	}

	// If there are no document rules, stream stdin line-by-line.
	// This avoids buffering and works with infinite streams (e.g. tail -f).
	if len(docRules) == 0 {
		pipeline := engine.NewPipeline(pendingLineRules...)
		scanner := bufio.NewScanner(stdin)
		ctx := &rule.LineContext{}

		// Call Setup on any rules that need it
		for _, lr := range pendingLineRules {
			if s, ok := lr.(rule.SetupRule); ok {
				s.Setup(ctx)
			}
		}

		for scanner.Scan() {
			ctx.LineNum++
			results, err := pipeline.Process(scanner.Text(), ctx)
			if err != nil {
				return fmt.Errorf("error applying rules: %w", err)
			}
			if ctx.Printing == rule.PrintOff {
				continue
			}
			for _, result := range results {
				fmt.Fprintln(stdout, result)
			}
		}
		return scanner.Err()
	}

	// Document rules exist — flush any trailing line rules and buffer all input.
	if len(pendingLineRules) > 0 {
		docRules = append(docRules, rule.NewApplyAllRule(pendingLineRules))
	}

	scanner := bufio.NewScanner(stdin)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	for _, dr := range docRules {
		var err error
		lines, err = dr.ApplyDocument(lines)
		if err != nil {
			return fmt.Errorf("error applying rules: %w", err)
		}
	}

	for _, line := range lines {
		fmt.Fprintln(stdout, line)
	}

	return nil
}
