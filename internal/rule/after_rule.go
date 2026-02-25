package rule

import "regexp"

// AfterRule starts printing after a line matches the pattern.
// The matching line itself is not printed — printing starts on the next line.
type AfterRule struct {
	pattern *regexp.Regexp
	matched bool
}

// NewAfterRule creates a rule that turns printing on after the first matching line.
func NewAfterRule(patternStr string) (*AfterRule, error) {
	pattern, err := regexp.Compile(patternStr)
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
	if r.matched {
		ctx.Printing = PrintOn
	}
	if r.pattern.MatchString(line) {
		r.matched = true
	}
	return []string{line}, nil
}
