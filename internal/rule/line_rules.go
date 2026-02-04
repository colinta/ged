package rule

import "regexp"

// SubstitutionRule replaces text matching a pattern.
type SubstitutionRule struct {
	pattern *regexp.Regexp
	replace string
	global  bool
}

// SubstitutionOption configures a SubstitutionRule.
type SubstitutionOption func(*SubstitutionRule)

// WithGlobal makes the substitution replace all matches, not just the first.
func WithGlobal() SubstitutionOption {
	return func(r *SubstitutionRule) {
		r.global = true
	}
}

// NewSubstitutionRule creates a rule that replaces pattern matches with replacement text.
// By default, only the first match is replaced. Use WithGlobal() to replace all matches.
func NewSubstitutionRule(pattern, replace string, opts ...SubstitutionOption) (*SubstitutionRule, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	r := &SubstitutionRule{
		pattern: re,
		replace: replace,
		global:  false,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r, nil
}

// Apply performs the substitution on the given line.
func (r *SubstitutionRule) Apply(line string) ([]string, error) {
	var result string
	if r.global {
		result = r.pattern.ReplaceAllString(line, r.replace)
	} else {
		// Actually, the above is wrong - ReplaceAllStringFunc still replaces all.
		// We need a different approach for first-only replacement.
		loc := r.pattern.FindStringIndex(line)
		if loc == nil {
			result = line
		} else {
			prefix := line[:loc[0]]
			postfix := line[loc[1]:]
			middle := r.pattern.ReplaceAllString(line[loc[0]:loc[1]], r.replace)
			result = prefix + middle + postfix
		}
	}

	return []string{result}, nil
}
