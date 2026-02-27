package rule

import "github.com/dlclark/regexp2"

// DeleteLineRule removes lines that match a pattern, keeps non-matching lines.
type DeleteLineRule struct {
	patternStr string
	pattern    *regexp2.Regexp
}

// Pattern returns the original pattern string.
func (r *DeleteLineRule) Pattern() string { return r.patternStr }

// NewDeleteLineRule creates a rule that removes lines matching the pattern.
// Use WithIgnoreCase() for case-insensitive matching.
func NewDeleteLineRule(patternStr string, opts ...RuleOption) (*DeleteLineRule, error) {
	patternRegex, err := CompilePattern(patternStr, opts...)
	if err != nil {
		return nil, err
	}
	return &DeleteLineRule{
		patternStr: patternStr,
		pattern:    patternRegex,
	}, nil
}

// Apply returns empty slice if line matches, keeps the line if not.
func (r *DeleteLineRule) Apply(line string, ctx *LineContext) ([]string, error) {
	matched, err := r.pattern.MatchString(line)
	if err != nil {
		return nil, err
	}
	if matched {
		return []string{}, nil // Delete: line matches
	}
	return []string{line}, nil // Keep: line doesn't match
}
