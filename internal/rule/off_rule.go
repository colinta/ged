package rule

import "github.com/dlclark/regexp2"

// OffRule stops printing when a line matches the pattern.
// The matching line itself is not printed.
type OffRule struct {
	pattern *regexp2.Regexp
}

// NewOffRule creates a rule that turns printing off at the first matching line.
// Use WithIgnoreCase() for case-insensitive matching.
func NewOffRule(patternStr string, opts ...RuleOption) (*OffRule, error) {
	pattern, err := CompilePattern(patternStr, opts...)
	if err != nil {
		return nil, err
	}
	return &OffRule{pattern: pattern}, nil
}

// Setup initializes the print state to on — lines are printed until a match.
// Only sets the initial state if no other control rule has set it first.
func (r *OffRule) Setup(ctx *LineContext) {
	if ctx.Printing == PrintDefault {
		ctx.Printing = PrintOn
	}
}

// Apply checks for a pattern match and turns printing off.
func (r *OffRule) Apply(line string, ctx *LineContext) ([]string, error) {
	matched, err := r.pattern.MatchString(line)
	if err != nil {
		return nil, err
	}
	if matched {
		ctx.Printing = PrintOff
	}
	return []string{line}, nil
}
