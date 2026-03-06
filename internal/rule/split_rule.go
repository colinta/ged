package rule

import "github.com/dlclark/regexp2"

// SplitRule splits each line on matches of a regex pattern,
// producing one output line per segment.
type SplitRule struct {
	pattern *regexp2.Regexp
}

// NewSplitRule creates a rule that splits lines on pattern matches.
func NewSplitRule(patternStr string, opts ...RuleOption) (*SplitRule, error) {
	compiled, err := CompilePattern(patternStr, opts...)
	if err != nil {
		return nil, err
	}
	return &SplitRule{pattern: compiled}, nil
}

// Apply splits the line on pattern matches and returns multiple lines.
func (r *SplitRule) Apply(line string, ctx *LineContext) ([]string, error) {
	parts, err := regexSplit(r.pattern, line)
	if err != nil {
		return nil, err
	}
	// If no match, regexSplit returns []string{line} (one element), which is correct.
	return parts, nil
}
