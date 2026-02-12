// Package parser handles parsing of ged command syntax.
package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/colinta/ged/internal/rule"
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
		parts[0] = regexp.QuoteMeta(parts[0])
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

	// Parse flags (part 2 if present)
	var opts []rule.SubstitutionOption
	if len(parts) >= 3 {
		flags := parts[2]
		if strings.Contains(flags, "g") {
			opts = append(opts, rule.WithGlobal())
		}
	}

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

	return rule.NewPrintLineRule(parts[0])
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

	return rule.NewDeleteLineRule(parts[0])
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
