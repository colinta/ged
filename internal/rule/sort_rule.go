package rule

import "sort"

// SortRule sorts all lines alphabetically.
type SortRule struct{}

// NewSortRule creates a new SortRule.
func NewSortRule() *SortRule {
	return &SortRule{}
}

// ApplyDocument sorts the lines alphabetically and returns a new slice.
func (r *SortRule) ApplyDocument(lines []string) ([]string, error) {
	sorted := make([]string, len(lines))
	copy(sorted, lines)
	sort.Strings(sorted)
	return sorted, nil
}
