package rule

import "regexp"

// OnRule starts printing when a line matches the pattern.
// The matching line itself is printed.
type OnRule struct {
	pattern *regexp.Regexp
}

// NewOnRule creates a rule that turns printing on at the first matching line.
func NewOnRule(patternStr string) (*OnRule, error) {
	pattern, err := regexp.Compile(patternStr)
	if err != nil {
		return nil, err
	}
	return &OnRule{pattern: pattern}, nil
}

// Setup initializes the print state to off — lines are suppressed until a match.
// Only sets the initial state if no other control rule has set it first.
func (r *OnRule) Setup(ctx *LineContext) {
	if ctx.Printing == PrintDefault {
		ctx.Printing = PrintOff
	}
}

// Apply checks for a pattern match and turns printing on.
func (r *OnRule) Apply(line string, ctx *LineContext) ([]string, error) {
	if r.pattern.MatchString(line) {
		ctx.Printing = PrintOn
	}
	return []string{line}, nil
}
