package rule

import "regexp"

// ToggleRule flips the print state each time a line matches the pattern.
type ToggleRule struct {
	pattern *regexp.Regexp
}

// NewToggleRule creates a rule that toggles printing on each matching line.
func NewToggleRule(patternStr string) (*ToggleRule, error) {
	pattern, err := regexp.Compile(patternStr)
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
	if r.pattern.MatchString(line) {
		if ctx.Printing == PrintOn {
			ctx.Printing = PrintOff
		} else {
			ctx.Printing = PrintOn
		}
	}
	return []string{line}, nil
}
