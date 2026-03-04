package rule

import (
	"strings"

	"github.com/dlclark/regexp2"
)

// TakeRule extracts the matching portion of each line.
// If the pattern has capture groups, the first group's text is returned.
// If no match, the line passes through unchanged.
type TakeRule struct {
	patternStr string
	pattern    *regexp2.Regexp
	global     bool
	joiner     string // separator for global mode (default " ")
}

// NewTakeRule creates a rule that extracts matching text from each line.
// With WithGlobal(), all matches are joined by joiner (default space).
func NewTakeRule(patternStr string, joiner string, opts ...RuleOption) (*TakeRule, error) {
	cfg := buildConfig(opts)
	patternRegex, err := CompilePattern(patternStr, opts...)
	if err != nil {
		return nil, err
	}
	return &TakeRule{
		patternStr: patternStr,
		pattern:    patternRegex,
		global:     cfg.global,
		joiner:     joiner,
	}, nil
}

// Apply extracts the matching portion from the line.
func (r *TakeRule) Apply(line string, ctx *LineContext) ([]string, error) {
	match, err := r.pattern.FindStringMatch(line)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return []string{line}, nil // no match — pass through
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

// extractMatch returns the first capture group if present, otherwise the full match.
func extractMatch(match *regexp2.Match) string {
	groups := match.Groups()
	// Groups()[0] is the full match; Groups()[1+] are capture groups
	if len(groups) > 1 && groups[1].Length > 0 {
		return groups[1].String()
	}
	return match.String()
}
