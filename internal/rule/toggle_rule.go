package rule

import "github.com/dlclark/regexp2"

// ToggleRule flips the print state each time a line matches the pattern.
type ToggleRule struct {
	pattern *regexp2.Regexp
}

// NewToggleRule creates a rule that toggles printing on each matching line.
// Use WithIgnoreCase() for case-insensitive matching.
func NewToggleRule(patternStr string, opts ...RuleOption) (*ToggleRule, error) {
	pattern, err := CompilePattern(patternStr, opts...)
	if err != nil {
		return nil, err
	}
	return &ToggleRule{pattern: pattern}, nil
}

// Setup initializes the print state to off.
// Only sets the initial state if no other control rule has set it first.
func (r *ToggleRule) Setup(ctx *LineContext) {
	if ctx.Printing == PrintDefault {
		ctx.Printing = PrintOff
	}
}

// Apply flips the print state when the line matches.
func (r *ToggleRule) Apply(line string, ctx *LineContext) ([]string, error) {
	matched, err := r.pattern.MatchString(line)
	if err != nil {
		return nil, err
	}
	if matched {
		if ctx.Printing == PrintOn {
			ctx.Printing = PrintOff
		} else {
			ctx.Printing = PrintOn
		}
	}
	return []string{line}, nil
}
