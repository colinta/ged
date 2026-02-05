package rule

import "regexp"

// DeleteLineRule removes lines that match a pattern, keeps non-matching lines.
type DeleteLineRule struct {
	patternStr string
	pattern    *regexp.Regexp
}

// Pattern returns the original pattern string.
func (r *DeleteLineRule) Pattern() string { return r.patternStr }

// NewDeleteLineRule creates a rule that removes lines matching the pattern.
func NewDeleteLineRule(patternStr string) (*DeleteLineRule, error) {
	patternRegex, err := regexp.Compile(patternStr)
	if err != nil {
		return nil, err
	}
	return &DeleteLineRule{
		patternStr: patternStr,
		pattern:    patternRegex,
	}, nil
}

// Apply returns empty slice if line matches, keeps the line if not.
func (r *DeleteLineRule) Apply(line string) ([]string, error) {
	if r.pattern.MatchString(line) {
		return []string{}, nil // Delete: line matches
	}
	return []string{line}, nil // Keep: line doesn't match
}
