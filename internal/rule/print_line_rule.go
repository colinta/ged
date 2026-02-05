package rule

import "regexp"

// PrintLineRule keeps lines that match a pattern, deletes non-matching lines.
type PrintLineRule struct {
	patternStr string
	pattern    *regexp.Regexp
}

// Pattern returns the original pattern string.
func (r *PrintLineRule) Pattern() string { return r.patternStr }

// NewPrintLineRule creates a rule that keeps only lines matching the pattern.
func NewPrintLineRule(patternStr string) (*PrintLineRule, error) {
	patternRegex, err := regexp.Compile(patternStr)
	if err != nil {
		return nil, err
	}
	return &PrintLineRule{
		patternStr: patternStr,
		pattern:    patternRegex,
	}, nil
}

// Apply returns the line if it matches, empty slice if not.
func (r *PrintLineRule) Apply(line string) ([]string, error) {
	if r.pattern.MatchString(line) {
		return []string{line}, nil // Keep: line matches
	}
	return []string{}, nil // Delete: line doesn't match
}
