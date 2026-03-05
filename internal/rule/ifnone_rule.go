package rule

import "github.com/dlclark/regexp2"

// IfNoneDocRule implements DocumentRule. It scans all lines — if NO lines
// match the pattern, inner rules are applied to ALL lines. If any line
// matches, the document passes through unchanged (or elseRules are applied).
//
// This is semantically the inverse of IfAnyDocRule, but exists as a separate
// type for clarity in the syntax: ifnone/pattern/ vs !ifany/pattern/.
type IfNoneDocRule struct {
	condition *regexp2.Regexp
	inverted  bool
	rules     []DocumentRule
	elseRules []DocumentRule
}

// NewIfNoneDocRule creates an IfNoneDocRule.
func NewIfNoneDocRule(condition *regexp2.Regexp, inverted bool, rules []DocumentRule, elseRules []DocumentRule) *IfNoneDocRule {
	return &IfNoneDocRule{
		condition: condition,
		inverted:  inverted,
		rules:     rules,
		elseRules: elseRules,
	}
}

// ApplyDocument scans all lines — if none match, apply rules to all lines.
func (r *IfNoneDocRule) ApplyDocument(lines []string) ([]string, error) {
	anyMatch, err := r.anyMatches(lines)
	if err != nil {
		return nil, err
	}

	noneMatch := !anyMatch
	if r.inverted {
		noneMatch = !noneMatch
	}

	if noneMatch {
		return applyDocRules(lines, r.rules)
	}
	if len(r.elseRules) > 0 {
		return applyDocRules(lines, r.elseRules)
	}
	return lines, nil
}

func (r *IfNoneDocRule) anyMatches(lines []string) (bool, error) {
	for _, line := range lines {
		matches, err := r.condition.MatchString(line)
		if err != nil {
			return false, err
		}
		if matches {
			return true, nil
		}
	}
	return false, nil
}
