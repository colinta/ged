package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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

// cliOptions holds parsed CLI flags separate from rule arguments.
type cliOptions struct {
	inputFiles []string // --input=file
	writeback  bool     // --write: overwrite input file in place
	writeTo    string   // --write-to=file: explicit output file
	ruleArgs   []string // remaining args that are rules
}

// parseCliOptions separates CLI flags from rule arguments.
// Flags:
//
//	--input=FILE or --input FILE   read from FILE instead of stdin (repeatable)
//	--write                        write output back to input file (requires --input)
//	--write-to=FILE or --write-to FILE   write output to FILE
//	--                             everything after -- is a rule argument (not a flag)
//
// Everything else is treated as a rule argument.
func parseCliOptions(args []string) (*cliOptions, error) {
	opts := &cliOptions{}
	bareRules := false

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if bareRules {
			opts.ruleArgs = append(opts.ruleArgs, arg)
			continue
		}

		if arg == "--" {
			bareRules = true
			continue
		}

		if arg == "--write" {
			opts.writeback = true
			continue
		}

		if strings.HasPrefix(arg, "--write-to=") {
			opts.writeTo = arg[len("--write-to="):]
			continue
		}
		if arg == "--write-to" {
			i++
			if i >= len(args) {
				return nil, fmt.Errorf("--write-to requires a filename")
			}
			opts.writeTo = args[i]
			continue
		}

		if strings.HasPrefix(arg, "--input=") {
			opts.inputFiles = append(opts.inputFiles, arg[len("--input="):])
			continue
		}
		if arg == "--input" {
			i++
			if i >= len(args) {
				return nil, fmt.Errorf("--input requires a filename")
			}
			opts.inputFiles = append(opts.inputFiles, args[i])
			continue
		}

		opts.ruleArgs = append(opts.ruleArgs, arg)
	}

	// Validate flag combinations
	if opts.writeback && len(opts.inputFiles) == 0 {
		return nil, fmt.Errorf("--write requires --input")
	}
	if opts.writeback && opts.writeTo != "" {
		return nil, fmt.Errorf("--write and --write-to are mutually exclusive")
	}
	if opts.writeTo != "" && len(opts.inputFiles) > 1 {
		return nil, fmt.Errorf("--write-to cannot be used with multiple input files")
	}

	return opts, nil
}

// run executes ged with the given arguments and I/O streams.
// This is separated from main() for testability.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	opts, err := parseCliOptions(args)
	if err != nil {
		return err
	}

	if len(opts.ruleArgs) < 1 {
		return fmt.Errorf("usage: ged <rule> [rule...]")
	}

	// Parse all rules, handling { } blocks for conditionals.
	allParsed, err := parser.ParseArgs(opts.ruleArgs)
	if err != nil {
		return fmt.Errorf("error parsing rules: %w", err)
	}

	// No input files — use stdin/stdout as before
	if len(opts.inputFiles) == 0 {
		output := stdout
		if opts.writeTo != "" {
			f, err := os.Create(opts.writeTo)
			if err != nil {
				return fmt.Errorf("error creating output file: %w", err)
			}
			defer f.Close()
			output = f
		}
		return processStream(allParsed, stdin, output)
	}

	// Process each input file
	for _, inputFile := range opts.inputFiles {
		if err := processFile(allParsed, inputFile, stdout, opts); err != nil {
			return err
		}
	}
	return nil
}

// processFile handles a single input file with the given options.
func processFile(allParsed []any, inputFile string, stdout io.Writer, opts *cliOptions) error {
	f, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("error opening %s: %w", inputFile, err)
	}
	defer f.Close()

	if opts.writeback {
		return writeBack(allParsed, f, inputFile)
	}

	output := stdout
	if opts.writeTo != "" {
		outFile, err := os.Create(opts.writeTo)
		if err != nil {
			return fmt.Errorf("error creating output file: %w", err)
		}
		defer outFile.Close()
		output = outFile
	}

	return processStream(allParsed, f, output)
}

// writeBack processes a file and writes the result back to the same path.
// It uses a temporary file + rename for safety.
func writeBack(allParsed []any, input io.Reader, filePath string) error {
	// Read and process into memory
	lines, err := readLines(input)
	if err != nil {
		return err
	}

	lines, err = applyDocRules(allParsed, lines)
	if err != nil {
		return err
	}

	// Write to a temp file in the same directory, then rename
	dir := filepath.Dir(filePath)
	tmp, err := os.CreateTemp(dir, ".ged-*")
	if err != nil {
		return fmt.Errorf("error creating temp file: %w", err)
	}
	tmpName := tmp.Name()

	// Write output to temp file
	for _, line := range lines {
		if _, err := fmt.Fprintln(tmp, line); err != nil {
			tmp.Close()
			os.Remove(tmpName)
			return fmt.Errorf("error writing temp file: %w", err)
		}
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("error closing temp file: %w", err)
	}

	// Preserve original file permissions
	info, err := os.Stat(filePath)
	if err == nil {
		os.Chmod(tmpName, info.Mode())
	}

	// Atomic rename
	if err := os.Rename(tmpName, filePath); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("error replacing file: %w", err)
	}

	return nil
}

// processStream applies rules to input and writes to output.
// This is the main processing function used for both stdin and file inputs.
func processStream(allParsed []any, input io.Reader, output io.Writer) error {
	// Build a list of DocumentRules.
	// Consecutive LineRules are wrapped in an ApplyAllRule.
	var docRules []rule.DocumentRule
	var pendingLineRules []rule.LineRule

	for _, parsed := range allParsed {
		switch r := parsed.(type) {
		case rule.LineRule:
			pendingLineRules = append(pendingLineRules, r)
		case rule.DocumentRule:
			if len(pendingLineRules) > 0 {
				docRules = append(docRules, rule.NewApplyAllRule(pendingLineRules))
				pendingLineRules = nil
			}
			docRules = append(docRules, r)
		default:
			return fmt.Errorf("unknown rule type from parser: %T", parsed)
		}
	}

	// If there are no document rules, stream line-by-line.
	if len(docRules) == 0 {
		pipeline := engine.NewPipeline(pendingLineRules...)
		scanner := bufio.NewScanner(input)
		ctx := &rule.LineContext{}

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
				fmt.Fprintln(output, result)
			}
		}
		return scanner.Err()
	}

	// Document rules exist — buffer all input.
	if len(pendingLineRules) > 0 {
		docRules = append(docRules, rule.NewApplyAllRule(pendingLineRules))
	}

	lines, err := readLines(input)
	if err != nil {
		return err
	}

	for _, dr := range docRules {
		lines, err = dr.ApplyDocument(lines)
		if err != nil {
			return fmt.Errorf("error applying rules: %w", err)
		}
	}

	for _, line := range lines {
		fmt.Fprintln(output, line)
	}

	return nil
}

// readLines reads all lines from a reader.
func readLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}
	return lines, nil
}

// applyDocRules buffers input lines and applies parsed rules as document rules.
// Used by writeBack where we need full document processing.
func applyDocRules(allParsed []any, lines []string) ([]string, error) {
	var docRules []rule.DocumentRule
	var pendingLineRules []rule.LineRule

	for _, parsed := range allParsed {
		switch r := parsed.(type) {
		case rule.LineRule:
			pendingLineRules = append(pendingLineRules, r)
		case rule.DocumentRule:
			if len(pendingLineRules) > 0 {
				docRules = append(docRules, rule.NewApplyAllRule(pendingLineRules))
				pendingLineRules = nil
			}
			docRules = append(docRules, r)
		}
	}
	if len(pendingLineRules) > 0 {
		docRules = append(docRules, rule.NewApplyAllRule(pendingLineRules))
	}

	for _, dr := range docRules {
		var err error
		lines, err = dr.ApplyDocument(lines)
		if err != nil {
			return nil, fmt.Errorf("error applying rules: %w", err)
		}
	}

	return lines, nil
}
