package rule

import "strings"

// TrimRule removes whitespace from lines.
// Mode controls which end(s) are trimmed: "both", "left", or "right".
type TrimRule struct {
	mode string // "both", "left", "right"
}

// NewTrimRule creates a rule that trims whitespace from both ends.
func NewTrimRule() *TrimRule {
	return &TrimRule{mode: "both"}
}

// NewTrimLeftRule creates a rule that trims leading whitespace.
func NewTrimLeftRule() *TrimRule {
	return &TrimRule{mode: "left"}
}

// NewTrimRightRule creates a rule that trims trailing whitespace.
func NewTrimRightRule() *TrimRule {
	return &TrimRule{mode: "right"}
}

// Apply trims whitespace from the line based on the mode.
func (r *TrimRule) Apply(line string, ctx *LineContext) ([]string, error) {
	switch r.mode {
	case "left":
		return []string{strings.TrimLeft(line, " \t")}, nil
	case "right":
		return []string{strings.TrimRight(line, " \t")}, nil
	default:
		return []string{strings.TrimSpace(line)}, nil
	}
}
