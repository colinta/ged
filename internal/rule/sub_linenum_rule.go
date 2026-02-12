package rule

import "strings"

// SubLineNumRule replaces the entire content of lines matching a line number range.
type SubLineNumRule struct {
	lineRange   LineRange
	replacement string
}

// NewSubLineNumRule creates a rule that replaces matching lines with the replacement string.
func NewSubLineNumRule(lineRange LineRange, replacement string) *SubLineNumRule {
	return &SubLineNumRule{lineRange: lineRange, replacement: replacement}
}

// Apply returns the replacement if line number matches, keeps the original line if not.
func (r *SubLineNumRule) Apply(line string, lineNum int) ([]string, error) {
	if r.lineRange.Contains(lineNum) {
		return strings.Split(r.replacement, "\n"), nil
	}
	return []string{line}, nil
}
