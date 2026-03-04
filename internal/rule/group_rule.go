package rule

import (
	"fmt"

	"github.com/dlclark/regexp2"
)

// GroupRule extracts a specific capture group from each line.
// Group numbers are 1-based (1 = first capture group, 2 = second, etc.).
// If no match or the group didn't participate, the line passes through unchanged.
type GroupRule struct {
	patternStr string
	pattern    *regexp2.Regexp
	group      int // 1-based group number
}

// NewGroupRule creates a rule that extracts the nth capture group.
// groupNum must be >= 1.
func NewGroupRule(patternStr string, groupNum int, opts ...RuleOption) (*GroupRule, error) {
	if groupNum < 1 {
		return nil, fmt.Errorf("group number must be >= 1, got %d", groupNum)
	}
	patternRegex, err := CompilePattern(patternStr, opts...)
	if err != nil {
		return nil, err
	}
	return &GroupRule{
		patternStr: patternStr,
		pattern:    patternRegex,
		group:      groupNum,
	}, nil
}

// Apply extracts the specified capture group from the line.
func (r *GroupRule) Apply(line string, ctx *LineContext) ([]string, error) {
	match, err := r.pattern.FindStringMatch(line)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return []string{line}, nil // no match — pass through
	}

	groups := match.Groups()
	if r.group >= len(groups) {
		return []string{line}, nil // group doesn't exist — pass through
	}

	g := groups[r.group]
	if g.Length == 0 {
		return []string{line}, nil // group didn't participate — pass through
	}

	return []string{g.String()}, nil
}
