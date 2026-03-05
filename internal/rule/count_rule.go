package rule

import "fmt"

// CountRule replaces the document with a single line: the line count.
type CountRule struct{}

// NewCountRule creates a new CountRule.
func NewCountRule() *CountRule {
	return &CountRule{}
}

// ApplyDocument returns the number of lines as a single-line document.
func (r *CountRule) ApplyDocument(lines []string) ([]string, error) {
	return []string{fmt.Sprintf("%d", len(lines))}, nil
}
