package rule

// PrintLineNumRule keeps lines that match a line number range, deletes others.
type PrintLineNumRule struct {
	lineRange LineRange
}

// NewPrintLineNumRule creates a rule that keeps only lines matching the line range.
func NewPrintLineNumRule(lineRange LineRange) *PrintLineNumRule {
	return &PrintLineNumRule{lineRange: lineRange}
}

// Apply returns the line if its line number matches the range, empty slice if not.
func (r *PrintLineNumRule) Apply(line string, ctx *LineContext) ([]string, error) {
	if r.lineRange.Contains(ctx.LineNum) {
		return []string{line}, nil // Keep: line number matches
	}
	return []string{}, nil // Delete: line number doesn't match
}
