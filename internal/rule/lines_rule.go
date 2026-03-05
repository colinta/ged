package rule

import "fmt"

// LinesRule prepends line numbers to each line.
// Implements DocumentRule because it needs to know total line count for padding.
type LinesRule struct{}

// NewLinesRule creates a new LinesRule.
func NewLinesRule() *LinesRule {
	return &LinesRule{}
}

// ApplyDocument prepends line numbers with consistent padding.
func (r *LinesRule) ApplyDocument(lines []string) ([]string, error) {
	if len(lines) == 0 {
		return lines, nil
	}

	// Calculate width needed for the largest line number
	width := len(fmt.Sprintf("%d", len(lines)))
	format := fmt.Sprintf("%%%dd: %%s", width)

	result := make([]string, len(lines))
	for i, line := range lines {
		result[i] = fmt.Sprintf(format, i+1, line)
	}
	return result, nil
}
