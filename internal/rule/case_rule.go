package rule

import "strings"

// UpperRule converts each line to uppercase.
type UpperRule struct{}

// NewUpperRule creates a rule that converts lines to uppercase.
func NewUpperRule() *UpperRule {
	return &UpperRule{}
}

// Apply returns the line converted to uppercase.
func (r *UpperRule) Apply(line string, ctx *LineContext) ([]string, error) {
	return []string{strings.ToUpper(line)}, nil
}

// LowerRule converts each line to lowercase.
type LowerRule struct{}

// NewLowerRule creates a rule that converts lines to lowercase.
func NewLowerRule() *LowerRule {
	return &LowerRule{}
}

// Apply returns the line converted to lowercase.
func (r *LowerRule) Apply(line string, ctx *LineContext) ([]string, error) {
	return []string{strings.ToLower(line)}, nil
}
