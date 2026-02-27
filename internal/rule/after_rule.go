package rule

import "github.com/dlclark/regexp2"

// AfterRule starts printing after a line matches the pattern.
// The matching line itself is not printed — printing starts on the next line.
type AfterRule struct {
	pattern *regexp2.Regexp
}

// NewAfterRule creates a rule that turns printing on after the first matching line.
// Use WithIgnoreCase() for case-insensitive matching.
func NewAfterRule(patternStr string, opts ...RuleOption) (*AfterRule, error) {
	pattern, err := CompilePattern(patternStr, opts...)
	if err != nil {
		return nil, err
	}
	return &AfterRule{pattern: pattern}, nil
}

// Setup initializes the print state to off — lines are suppressed until after a match.
// Only sets the initial state if no other control rule has set it first.
func (r *AfterRule) Setup(ctx *LineContext) {
	if ctx.Printing == PrintDefault {
		ctx.Printing = PrintOff
	}
}

// Apply checks for a pattern match. On the line after a match, printing turns on.
// The matching line itself stays off.
func (r *AfterRule) Apply(line string, ctx *LineContext) ([]string, error) {
	if GetState(ctx, r, false) {
		ctx.Printing = PrintOn
	}
	matched, err := r.pattern.MatchString(line)
	if err != nil {
		return nil, err
	}
	if matched {
		SetState(ctx, r, true)
	}
	return []string{line}, nil
}
