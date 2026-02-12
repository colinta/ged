package rule

// DeleteLineNumRule removes lines that match a line number range, keeps others.
type DeleteLineNumRule struct {
	lineRange LineRange
}

// NewDeleteLineNumRule creates a rule that removes lines matching the line range.
func NewDeleteLineNumRule(lineRange LineRange) *DeleteLineNumRule {
	return &DeleteLineNumRule{lineRange: lineRange}
}

// Apply returns empty slice if line number matches the range, keeps the line if not.
func (r *DeleteLineNumRule) Apply(line string, lineNum int) ([]string, error) {
	if r.lineRange.Contains(lineNum) {
		return []string{}, nil // Delete: line number matches
	}
	return []string{line}, nil // Keep: line number doesn't match
}
