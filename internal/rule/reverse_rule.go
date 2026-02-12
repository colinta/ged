package rule

import "slices"

// ReverseRule reverses the order of all lines.
type ReverseRule struct{}

// NewReverseRule creates a new ReverseRule.
func NewReverseRule() *ReverseRule {
	return &ReverseRule{}
}

// ApplyDocument reverses the line order and returns a new slice.
func (r *ReverseRule) ApplyDocument(lines []string) ([]string, error) {
	reversed := make([]string, len(lines))
	copy(reversed, lines)
	slices.Reverse(reversed)
	return reversed, nil
}
