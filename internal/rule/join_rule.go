package rule

import "strings"

// JoinRule joins all lines into a single line with a separator.
type JoinRule struct {
	separator string
}

// NewJoinRule creates a new JoinRule with the given separator.
func NewJoinRule(separator string) *JoinRule {
	return &JoinRule{separator: separator}
}

// ApplyDocument joins all lines with the separator and returns a single-element slice.
func (r *JoinRule) ApplyDocument(lines []string) ([]string, error) {
	if len(lines) == 0 {
		return []string{}, nil
	}
	return []string{strings.Join(lines, r.separator)}, nil
}
