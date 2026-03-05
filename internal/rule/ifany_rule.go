package rule

import "github.com/dlclark/regexp2"

// IfAnyDocRule implements DocumentRule. It scans all lines — if ANY line
// matches the pattern, inner rules are applied to ALL lines. If no line
// matches, the document passes through unchanged (or elseRules are applied).
type IfAnyDocRule struct {
	condition *regexp2.Regexp
	inverted  bool
	rules     []DocumentRule
	elseRules []DocumentRule
}

// NewIfAnyDocRule creates an IfAnyDocRule.
func NewIfAnyDocRule(condition *regexp2.Regexp, inverted bool, rules []DocumentRule, elseRules []DocumentRule) *IfAnyDocRule {
	return &IfAnyDocRule{
		condition: condition,
		inverted:  inverted,
		rules:     rules,
		elseRules: elseRules,
	}
}

// ApplyDocument scans all lines for the condition, then applies rules to
// the entire document if any match was found.
func (r *IfAnyDocRule) ApplyDocument(lines []string) ([]string, error) {
	anyMatch, err := r.anyMatches(lines)
	if err != nil {
		return nil, err
	}
	if r.inverted {
		anyMatch = !anyMatch
	}

	if anyMatch {
		return applyDocRules(lines, r.rules)
	}
	if len(r.elseRules) > 0 {
		return applyDocRules(lines, r.elseRules)
	}
	return lines, nil
}

func (r *IfAnyDocRule) anyMatches(lines []string) (bool, error) {
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

// applyDocRules applies a sequence of DocumentRules to a document.
func applyDocRules(lines []string, rules []DocumentRule) ([]string, error) {
	var err error
	for _, dr := range rules {
		lines, err = dr.ApplyDocument(lines)
		if err != nil {
			return nil, err
		}
	}
	return lines, nil
}
