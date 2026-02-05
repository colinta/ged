// Package parser handles parsing of ged command syntax.
package parser

import (
	"fmt"
	"strings"

	"github.com/colinta/ged/internal/rule"
)

// ParseRule parses a rule string and returns the appropriate Rule.
// It handles delimiter detection and dispatches to command-specific parsers.
func ParseRule(input string) (rule.Rule, error) {
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

	switch command {
	case 's':
		return parseSubstitution(parts)
	case 'p':
		return parsePrint(parts)
	case 'd':
		return parseDelete(parts)
	default:
		return nil, fmt.Errorf("unknown command: %c", command)
	}
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

// parsePrint creates a PrintLineRule from parsed parts.
// Expected parts: [pattern] or [pattern, ""]
func parsePrint(parts []string) (*rule.PrintLineRule, error) {
	if len(parts) < 1 {
		return nil, fmt.Errorf("print requires a pattern")
	}

	pattern := parts[0]
	return rule.NewPrintLineRule(pattern)
}

// parseDelete creates a DeleteLineRule from parsed parts.
// Expected parts: [pattern] or [pattern, ""]
func parseDelete(parts []string) (*rule.DeleteLineRule, error) {
	if len(parts) < 1 {
		return nil, fmt.Errorf("delete requires a pattern")
	}

	pattern := parts[0]
	return rule.NewDeleteLineRule(pattern)
}
