package rule

import (
	"strings"

	"github.com/dlclark/regexp2"
)

// SubstitutionRule replaces text matching a pattern.
type SubstitutionRule struct {
	patternStr string          // original pattern string
	pattern    *regexp2.Regexp // compiled regex
	replace    string
	global     bool
}

// Pattern returns the original pattern string.
func (r *SubstitutionRule) Pattern() string { return r.patternStr }

// Replace returns the replacement string.
func (r *SubstitutionRule) Replace() string { return r.replace }

// Global returns whether all matches are replaced.
func (r *SubstitutionRule) Global() bool { return r.global }

// NewSubstitutionRule creates a rule that replaces pattern matches with replacement text.
// By default, only the first match is replaced. Use WithGlobal() to replace all matches.
// Use WithIgnoreCase() for case-insensitive matching.
func NewSubstitutionRule(patternStr, replace string, opts ...RuleOption) (*SubstitutionRule, error) {
	cfg := buildConfig(opts)
	patternRegex, err := CompilePattern(patternStr, opts...)
	if err != nil {
		return nil, err
	}

	return &SubstitutionRule{
		patternStr: patternStr,
		pattern:    patternRegex,
		replace:    replace,
		global:     cfg.global,
	}, nil
}

// Apply performs the substitution on the given line.
func (r *SubstitutionRule) Apply(line string, ctx *LineContext) ([]string, error) {
	count := 1
	if r.global {
		count = -1 // -1 means replace all
	}
	result, err := r.pattern.Replace(line, r.replace, 0, count)
	if err != nil {
		return nil, err
	}

	return strings.Split(result, "\n"), nil
}
