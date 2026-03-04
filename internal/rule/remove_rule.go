package rule

import "github.com/dlclark/regexp2"

// RemoveRule removes the matching portion of each line, keeping the rest.
// If no match, the line passes through unchanged.
type RemoveRule struct {
	patternStr string
	pattern    *regexp2.Regexp
	global     bool
}

// NewRemoveRule creates a rule that removes matching text from each line.
// With WithGlobal(), all matches are removed.
func NewRemoveRule(patternStr string, opts ...RuleOption) (*RemoveRule, error) {
	cfg := buildConfig(opts)
	patternRegex, err := CompilePattern(patternStr, opts...)
	if err != nil {
		return nil, err
	}
	return &RemoveRule{
		patternStr: patternStr,
		pattern:    patternRegex,
		global:     cfg.global,
	}, nil
}

// Apply removes the matching portion from the line.
func (r *RemoveRule) Apply(line string, ctx *LineContext) ([]string, error) {
	count := 1
	if r.global {
		count = -1
	}
	result, err := r.pattern.Replace(line, "", 0, count)
	if err != nil {
		return nil, err
	}
	return []string{result}, nil
}
