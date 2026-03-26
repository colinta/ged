package rule

import "github.com/dlclark/regexp2"

// RemovePrintRule combines Print and Remove: only matching lines are kept,
// and the matching portion is removed from the output.
// Non-matching lines are removed from output.
type RemovePrintRule struct {
	patternStr string
	pattern    *regexp2.Regexp
	global     bool
}

// NewRemovePrintRule creates a rule that keeps only matching lines and removes the match.
// With WithGlobal(), all matches are removed.
func NewRemovePrintRule(patternStr string, opts ...RuleOption) (*RemovePrintRule, error) {
	cfg := buildConfig(opts)
	patternRegex, err := CompilePattern(patternStr, opts...)
	if err != nil {
		return nil, err
	}
	return &RemovePrintRule{
		patternStr: patternStr,
		pattern:    patternRegex,
		global:     cfg.global,
	}, nil
}

// Apply removes the matching portion from the line, or removes non-matching lines entirely.
func (r *RemovePrintRule) Apply(line string, ctx *LineContext) ([]string, error) {
	matched, err := r.pattern.MatchString(line)
	if err != nil {
		return nil, err
	}
	if !matched {
		return []string{}, nil // no match — remove line
	}

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
