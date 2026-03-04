package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/colinta/ged/internal/diff"
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

// colorMode represents the user's color preference.
type colorMode int

const (
	colorAuto colorMode = iota // detect from terminal
	colorOn                    // always use colors
	colorOff                   // never use colors
)

// cliOptions holds parsed CLI flags separate from rule arguments.
type cliOptions struct {
	inputFiles []string  // --input=file
	writeback  bool      // --write: overwrite input file in place
	writeTo    string    // --write-to=file: explicit output file
	diffMode   bool      // --diff: show diff instead of output
	color      colorMode // --color / --no-color
	ruleArgs   []string  // remaining args that are rules
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

		if arg == "--diff" {
			opts.diffMode = true
			continue
		}

		if arg == "--color" {
			opts.color = colorOn
			continue
		}

		if arg == "--no-color" {
			opts.color = colorOff
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
	if opts.diffMode && opts.writeback {
		return nil, fmt.Errorf("--diff and --write are mutually exclusive")
	}
	if opts.diffMode && opts.writeTo != "" {
		return nil, fmt.Errorf("--diff and --write-to are mutually exclusive")
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

	// Resolve color: auto-detect from stdout if possible
	useColor := resolveColor(opts.color, stdout)

	// Diff mode: compare original vs transformed
	if opts.diffMode {
		return runDiff(allParsed, opts, stdin, stdout, useColor)
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

// resolveColor determines whether to use colors based on mode and output target.
func resolveColor(mode colorMode, output io.Writer) bool {
	switch mode {
	case colorOn:
		return true
	case colorOff:
		return false
	default: // colorAuto
		// Check if output is an *os.File pointing to a character device (terminal).
		// This avoids the golang.org/x/term dependency — os.File.Stat() returns
		// a FileInfo whose Mode() includes ModeCharDevice for TTYs.
		if f, ok := output.(*os.File); ok {
			info, err := f.Stat()
			if err == nil {
				return info.Mode()&os.ModeCharDevice != 0
			}
		}
		return false
	}
}

// runDiff processes input in diff mode: shows changes instead of writing output.
func runDiff(allParsed []any, opts *cliOptions, stdin io.Reader, stdout io.Writer, useColor bool) error {
	if len(opts.inputFiles) == 0 {
		// stdin mode: read all, transform, diff
		original, err := readLines(stdin)
		if err != nil {
			return err
		}
		transformed, err := applyDocRules(allParsed, original)
		if err != nil {
			return err
		}
		return writeDiff(original, transformed, "", stdout, useColor)
	}

	// File mode: diff each file
	for _, inputFile := range opts.inputFiles {
		f, err := os.Open(inputFile)
		if err != nil {
			return fmt.Errorf("error opening %s: %w", inputFile, err)
		}
		original, err := readLines(f)
		f.Close()
		if err != nil {
			return err
		}
		transformed, err := applyDocRules(allParsed, original)
		if err != nil {
			return err
		}
		if err := writeDiff(original, transformed, inputFile, stdout, useColor); err != nil {
			return err
		}
	}
	return nil
}

// writeDiff computes and prints a diff between original and transformed lines.
func writeDiff(original, transformed []string, filename string, output io.Writer, useColor bool) error {
	changes := diff.Compute(original, transformed)
	if !diff.HasChanges(changes) {
		return nil // no output if nothing changed
	}

	// Print header if we have a filename
	if filename != "" {
		header := fmt.Sprintf("--- %s\n+++ %s", filename, filename)
		if useColor {
			header = "\033[1m" + header + "\033[0m"
		}
		fmt.Fprintln(output, header)
	}

	lines := diff.Format(changes, useColor)
	for _, line := range lines {
		fmt.Fprintln(output, line)
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
