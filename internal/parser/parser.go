// Package parser handles parsing of ged command syntax.
package parser

import (
	"fmt"
	"strings"

	"github.com/colinta/ged/internal/rule"
	"github.com/dlclark/regexp2"
)

// ParseRule parses a rule string and returns the appropriate Rule.
// It handles delimiter detection and dispatches to command-specific parsers.
// Returns either a rule.LineRule or rule.DocumentRule (as any).
func ParseRule(input string) (any, error) {
	// Word commands must be checked first — "sort" starts with 's',
	// which would otherwise match the substitution command.
	if input == "sort" {
		return rule.NewSortRule(), nil
	}
	if input == "reverse" {
		return rule.NewReverseRule(), nil
	}
	if strings.HasPrefix(input, "join") {
		return parseJoin(input)
	}
	if strings.HasPrefix(input, "!between") || strings.HasPrefix(input, "between") {
		return parseBetween(input)
	}
	if strings.HasPrefix(input, "!if") || strings.HasPrefix(input, "if") {
		return parseIf(input)
	}
	if strings.HasPrefix(input, "on") {
		return parseControl(input, "on")
	}
	if strings.HasPrefix(input, "off") {
		return parseControl(input, "off")
	}
	if strings.HasPrefix(input, "after") {
		return parseControl(input, "after")
	}
	if strings.HasPrefix(input, "toggle") {
		return parseControl(input, "toggle")
	}

	if len(input) < 2 {
		return nil, fmt.Errorf("invalid rule: too short")
	}

	command := input[0]
	delimiter := input[1]
	rest := input[2:]

	// Split by delimiter, respecting backslash escapes
	parts, err := splitByDelimiter(rest, delimiter)
	if err != nil {
		return nil, err
	}

	// Quote delimiters mean literal matching — escape regex metacharacters
	if (delimiter == '`' || delimiter == '\'' || delimiter == '"') && len(parts) > 0 {
		parts[0] = regexp2.Escape(parts[0])
	}

	if command == 'p' && delimiter == ':' {
		return parsePrintLineNum(parts)
	} else if command == 'p' {
		return parsePrint(parts)
	} else if command == 'd' && delimiter == ':' {
		return parseDeleteLineNum(parts)
	} else if command == 'd' {
		return parseDelete(parts)
	} else if command == 's' && delimiter == ':' {
		return parseSubstitutionLineNum(parts)
	} else if command == 's' {
		return parseSubstitution(parts)
	} else {
		return nil, fmt.Errorf("unknown command: %c", command)
	}
}

// parseFlags reads a flags string and returns the corresponding RuleOptions.
// Supported flags:
//
//	g — global replacement (SubstitutionRule only)
//	i — case-insensitive matching
func parseFlags(flags string) []rule.RuleOption {
	var opts []rule.RuleOption
	if strings.Contains(flags, "g") {
		opts = append(opts, rule.WithGlobal())
	}
	if strings.Contains(flags, "i") {
		opts = append(opts, rule.WithIgnoreCase())
	}
	return opts
}

// flagsFromParts extracts flags from the trailing element of a parts slice.
// For commands like p/pat/ and d/pat/, flags are in parts[1].
// For substitution s/pat/repl/flags, flags are in parts[2].
// Returns the options parsed from the given index, or nil if index is out of range.
func flagsFromParts(parts []string, flagIndex int) []rule.RuleOption {
	if flagIndex < len(parts) {
		return parseFlags(parts[flagIndex])
	}
	return nil
}

// parseJoin handles "join" (bare) and "join/sep/" syntax.
func parseJoin(input string) (*rule.JoinRule, error) {
	if input == "join" {
		return rule.NewJoinRule(""), nil
	}

	if len(input) < 5 {
		return nil, fmt.Errorf("invalid join syntax: %q", input)
	}

	delimiter := input[4]
	rest := input[5:]

	parts, err := splitByDelimiter(rest, delimiter)
	if err != nil {
		return nil, err
	}

	if len(parts) < 1 {
		return nil, fmt.Errorf("join requires a separator")
	}

	return rule.NewJoinRule(parts[0]), nil
}

// splitByDelimiter splits a string by delimiter, respecting backslash escapes.
// The delimiter at the end is required (trailing part can be empty for flags).
// Returns the parts with escape sequences processed.
func splitByDelimiter(input string, delimiter byte) ([]string, error) {
	var parts []string
	var current strings.Builder

	i := 0
	for i < len(input) {
		ch := input[i]

		if ch == '\\' && i+1 < len(input) {
			next := input[i+1]
			if next == delimiter {
				// Escaped delimiter - write the delimiter itself
				current.WriteByte(delimiter)
				i += 2
				continue
			} else if next == '\\' {
				// Escaped backslash - write single backslash
				current.WriteByte('\\')
				i += 2
				continue
			} else if next == 'n' {
				current.WriteByte('\n')
				i += 2
				continue
			} else if next == 't' {
				current.WriteByte('\t')
				i += 2
				continue
			}
			// Not an escape sequence we handle - write backslash and continue
			current.WriteByte(ch)
			i++
			continue
		}

		if ch == delimiter {
			parts = append(parts, current.String())
			current.Reset()
			i++
			continue
		}

		current.WriteByte(ch)
		i++
	}

	// The last part (after final delimiter) contains flags
	// It's okay if there's content - that's the flags portion
	parts = append(parts, current.String())

	return parts, nil
}

// parseSubstitution creates a SubstitutionRule from parsed parts.
// Expected parts: [pattern, replacement, flags]
// The trailing delimiter is required, so we need at least 3 parts.
func parseSubstitution(parts []string) (*rule.SubstitutionRule, error) {
	if len(parts) < 2 {
		return nil, fmt.Errorf("substitution requires pattern and replacement with trailing delimiter")
	}

	pattern := parts[0]
	replace := parts[1]
	opts := flagsFromParts(parts, 2)

	return rule.NewSubstitutionRule(pattern, replace, opts...)
}

// parseSubstitutionLineNum creates a SubLineNumRule for line number replacement.
// Expected parts: [lineRange, replacement]
func parseSubstitutionLineNum(parts []string) (rule.LineRule, error) {
	if len(parts) < 2 {
		return nil, fmt.Errorf("substitution requires a line range and replacement")
	}

	lineRange, err := rule.ParseLineRange(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid line range: %w", err)
	}
	return rule.NewSubLineNumRule(lineRange, parts[1]), nil
}

// parsePrint creates a PrintLineRule for pattern matching.
func parsePrint(parts []string) (rule.LineRule, error) {
	if len(parts) < 1 {
		return nil, fmt.Errorf("print requires a pattern")
	}

	opts := flagsFromParts(parts, 1)
	return rule.NewPrintLineRule(parts[0], opts...)
}

// parsePrintLineNum creates a PrintLineNumRule for line number filtering.
func parsePrintLineNum(parts []string) (rule.LineRule, error) {
	if len(parts) < 1 {
		return nil, fmt.Errorf("print requires a line range")
	}

	lineRange, err := rule.ParseLineRange(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid line range: %w", err)
	}
	return rule.NewPrintLineNumRule(lineRange), nil
}

// parseDelete creates a DeleteLineRule for pattern matching.
func parseDelete(parts []string) (rule.LineRule, error) {
	if len(parts) < 1 {
		return nil, fmt.Errorf("delete requires a pattern")
	}

	opts := flagsFromParts(parts, 1)
	return rule.NewDeleteLineRule(parts[0], opts...)
}

// parseDeleteLineNum creates a DeleteLineNumRule for line number filtering.
func parseDeleteLineNum(parts []string) (rule.LineRule, error) {
	if len(parts) < 1 {
		return nil, fmt.Errorf("delete requires a line range")
	}

	lineRange, err := rule.ParseLineRange(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid line range: %w", err)
	}
	return rule.NewDeleteLineNumRule(lineRange), nil
}

// parseControl parses "name/pattern/" for control rules (on, off, after, toggle).
func parseControl(input string, name string) (rule.LineRule, error) {
	rest := input[len(name):]
	if len(rest) == 0 {
		return nil, fmt.Errorf("%s requires a pattern", name)
	}

	delimiter := rest[0]
	parts, err := splitByDelimiter(rest[1:], delimiter)
	if err != nil {
		return nil, err
	}

	if len(parts) < 1 || parts[0] == "" {
		return nil, fmt.Errorf("%s requires a pattern", name)
	}

	pattern := parts[0]
	if delimiter == '`' || delimiter == '\'' || delimiter == '"' {
		pattern = regexp2.Escape(pattern)
	}

	opts := flagsFromParts(parts, 1)

	switch name {
	case "on":
		return rule.NewOnRule(pattern, opts...)
	case "off":
		return rule.NewOffRule(pattern, opts...)
	case "after":
		return rule.NewAfterRule(pattern, opts...)
	case "toggle":
		return rule.NewToggleRule(pattern, opts...)
	default:
		return nil, fmt.Errorf("unknown control command: %s", name)
	}
}

// condition is a parser-internal type representing a parsed if/!if condition.
// It's not a rule — it gets converted into a ConditionalRule once the inner
// rules are collected from the { } block.
type condition struct {
	pattern  *regexp2.Regexp
	inverted bool
}

// parseIf parses "if/pattern/" or "!if/pattern/" and returns a condition.
func parseIf(input string) (*condition, error) {
	inverted := false
	rest := input

	if strings.HasPrefix(rest, "!if") {
		inverted = true
		rest = rest[3:]
	} else {
		rest = rest[2:]
	}

	if len(rest) == 0 {
		return nil, fmt.Errorf("missing pattern in if condition")
	}

	delimiter := rest[0]
	parts, err := splitByDelimiter(rest[1:], delimiter)
	if err != nil {
		return nil, err
	}

	if len(parts) < 1 || parts[0] == "" {
		return nil, fmt.Errorf("missing pattern in if condition")
	}

	pattern := parts[0]
	if delimiter == '`' || delimiter == '\'' || delimiter == '"' {
		pattern = regexp2.Escape(pattern)
	}

	opts := flagsFromParts(parts, 1)
	compiled, err := rule.CompilePattern(pattern, opts...)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern in if condition: %w", err)
	}

	return &condition{
		pattern:  compiled,
		inverted: inverted,
	}, nil
}

// betweenCondition is a parser-internal type representing a parsed between condition.
// Like condition, it gets assembled with inner rules from { } blocks in parseArgs.
type betweenCondition struct {
	startPattern *regexp2.Regexp
	endPattern   *regexp2.Regexp
	inverted     bool
}

// parseBetween parses "between/start/end/" or "!between/start/end/" and returns a betweenCondition.
func parseBetween(input string) (*betweenCondition, error) {
	inverted := false
	rest := input

	if strings.HasPrefix(rest, "!between") {
		inverted = true
		rest = rest[8:]
	} else {
		rest = rest[7:]
	}

	if len(rest) == 0 {
		return nil, fmt.Errorf("between requires start and end patterns")
	}

	delimiter := rest[0]
	parts, err := splitByDelimiter(rest[1:], delimiter)
	if err != nil {
		return nil, err
	}

	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("between requires start and end patterns")
	}

	startPattern := parts[0]
	endPattern := parts[1]
	if delimiter == '`' || delimiter == '\'' || delimiter == '"' {
		startPattern = regexp2.Escape(startPattern)
		endPattern = regexp2.Escape(endPattern)
	}

	opts := flagsFromParts(parts, 2)
	startCompiled, err := rule.CompilePattern(startPattern, opts...)
	if err != nil {
		return nil, fmt.Errorf("invalid start pattern in between: %w", err)
	}

	endCompiled, err := rule.CompilePattern(endPattern, opts...)
	if err != nil {
		return nil, fmt.Errorf("invalid end pattern in between: %w", err)
	}

	return &betweenCondition{
		startPattern: startCompiled,
		endPattern:   endCompiled,
		inverted:     inverted,
	}, nil
}
