package rule

import "github.com/dlclark/regexp2"

// applyLineRules applies a pipeline of LineRules to a single line.
func applyLineRules(line string, ctx *LineContext, rules []LineRule) ([]string, error) {
	current := []string{line}
	for _, r := range rules {
		var next []string
		for _, l := range current {
			out, err := r.Apply(l, ctx)
			if err != nil {
				return nil, err
			}
			next = append(next, out...)
		}
		if len(next) == 0 {
			return nil, nil
		}
		current = next
	}
	return current, nil
}

// ConditionalLineRule implements LineRule. It applies inner LineRules only to
// lines matching (or not matching) a condition. Non-matching lines pass through
// unchanged (or have elseRules applied). Because all inner rules are LineRules,
// this can stream.
type ConditionalLineRule struct {
	condition *regexp2.Regexp
	inverted  bool
	rules     []LineRule
	elseRules []LineRule
}

// NewConditionalLineRule creates a ConditionalLineRule.
func NewConditionalLineRule(condition *regexp2.Regexp, inverted bool, rules []LineRule, elseRules []LineRule) *ConditionalLineRule {
	return &ConditionalLineRule{
		condition: condition,
		inverted:  inverted,
		rules:     rules,
		elseRules: elseRules,
	}
}

// Apply checks the condition and either runs inner rules or else rules.
func (r *ConditionalLineRule) Apply(line string, ctx *LineContext) ([]string, error) {
	matches, err := r.condition.MatchString(line)
	if err != nil {
		return nil, err
	}
	if r.inverted {
		matches = !matches
	}

	if matches {
		return applyLineRules(line, ctx, r.rules)
	}
	if len(r.elseRules) > 0 {
		return applyLineRules(line, ctx, r.elseRules)
	}
	return []string{line}, nil
}

// ConditionalDocRule implements DocumentRule. It applies inner rules to matching
// lines and elseRules to non-matching lines. Each line is processed independently
// as a mini-document.
type ConditionalDocRule struct {
	condition *regexp2.Regexp
	inverted  bool
	rules     []DocumentRule
	elseRules []DocumentRule
}

// NewConditionalDocRule creates a ConditionalDocRule.
func NewConditionalDocRule(condition *regexp2.Regexp, inverted bool, rules []DocumentRule, elseRules []DocumentRule) *ConditionalDocRule {
	return &ConditionalDocRule{
		condition: condition,
		inverted:  inverted,
		rules:     rules,
		elseRules: elseRules,
	}
}

// ApplyDocument processes each line independently: matching lines have rules
// applied, non-matching lines have elseRules applied (or pass through unchanged).
func (r *ConditionalDocRule) ApplyDocument(lines []string) ([]string, error) {
	var result []string

	for _, line := range lines {
		matches, err := r.condition.MatchString(line)
		if err != nil {
			return nil, err
		}
		if r.inverted {
			matches = !matches
		}

		var activeRules []DocumentRule
		if matches {
			activeRules = r.rules
		} else if len(r.elseRules) > 0 {
			activeRules = r.elseRules
		}

		if activeRules == nil {
			result = append(result, line)
			continue
		}

		processed := []string{line}
		for _, dr := range activeRules {
			processed, err = dr.ApplyDocument(processed)
			if err != nil {
				return nil, err
			}
		}
		result = append(result, processed...)
	}

	return result, nil
}
