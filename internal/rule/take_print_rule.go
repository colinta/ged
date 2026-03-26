package rule

import (
	"strings"

	"github.com/dlclark/regexp2"
)

// TakePrintRule combines Print and Take: only matching lines are kept,
// and only the matching portion is printed.
// If the pattern has capture groups, the first group's text is returned.
// Non-matching lines are removed from output.
type TakePrintRule struct {
	patternStr string
	pattern    *regexp2.Regexp
	global     bool
	joiner     string // separator for global mode (default " ")
}

// NewTakePrintRule creates a rule that keeps only matching lines and extracts the match.
// With WithGlobal(), all matches are joined by joiner (default space).
func NewTakePrintRule(patternStr string, joiner string, opts ...RuleOption) (*TakePrintRule, error) {
	cfg := buildConfig(opts)
	patternRegex, err := CompilePattern(patternStr, opts...)
	if err != nil {
		return nil, err
	}
	return &TakePrintRule{
		patternStr: patternStr,
		pattern:    patternRegex,
		global:     cfg.global,
		joiner:     joiner,
	}, nil
}

// Apply extracts the matching portion from the line, or removes non-matching lines.
func (r *TakePrintRule) Apply(line string, ctx *LineContext) ([]string, error) {
	match, err := r.pattern.FindStringMatch(line)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return []string{}, nil // no match — remove line
	}

	if !r.global {
		return []string{extractMatch(match)}, nil
	}

	// Global: collect all matches
	var results []string
	for match != nil {
		results = append(results, extractMatch(match))
		match, err = r.pattern.FindNextMatch(match)
		if err != nil {
			return nil, err
		}
	}
	return []string{strings.Join(results, r.joiner)}, nil
}
