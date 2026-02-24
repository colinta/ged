package rule

import "regexp"

// ConditionalLineRule implements LineRule. It applies inner LineRules only to
// lines matching (or not matching) a condition. Non-matching lines pass through
// unchanged. Because all inner rules are LineRules, this can stream.
type ConditionalLineRule struct {
	condition *regexp.Regexp
	inverted  bool
	rules     []LineRule
}

// NewConditionalLineRule creates a ConditionalLineRule.
func NewConditionalLineRule(condition *regexp.Regexp, inverted bool, rules []LineRule) *ConditionalLineRule {
	return &ConditionalLineRule{
		condition: condition,
		inverted:  inverted,
		rules:     rules,
	}
}

// Apply checks the condition and either runs inner rules or passes the line through.
func (r *ConditionalLineRule) Apply(line string, ctx *LineContext) ([]string, error) {
	matches := r.condition.MatchString(line)
	if r.inverted {
		matches = !matches
	}

	if !matches {
		return []string{line}, nil
	}

	// Apply inner rules as a pipeline — same pattern as ApplyAllRule
	current := []string{line}
	for _, innerRule := range r.rules {
		var next []string
		for _, l := range current {
			out, err := innerRule.Apply(l, ctx)
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

// ConditionalDocRule implements DocumentRule. It collects lines matching the
// condition into a sub-document, applies inner DocumentRules to that sub-document,
// then weaves the results back into the original positions. Non-matching lines
// stay in place.
type ConditionalDocRule struct {
	condition *regexp.Regexp
	inverted  bool
	rules     []DocumentRule
}

// NewConditionalDocRule creates a ConditionalDocRule.
func NewConditionalDocRule(condition *regexp.Regexp, inverted bool, rules []DocumentRule) *ConditionalDocRule {
	return &ConditionalDocRule{
		condition: condition,
		inverted:  inverted,
		rules:     rules,
	}
}

// ApplyDocument collects matching lines, applies inner rules, then reconstructs
// the output with processed lines replacing their original positions.
func (r *ConditionalDocRule) ApplyDocument(lines []string) ([]string, error) {
	var matchingLines []string
	isMatch := make([]bool, len(lines))

	for i, line := range lines {
		matches := r.condition.MatchString(line)
		if r.inverted {
			matches = !matches
		}
		if matches {
			matchingLines = append(matchingLines, line)
			isMatch[i] = true
		}
	}

	// Apply inner document rules to the matching lines
	processed := matchingLines
	for _, dr := range r.rules {
		var err error
		processed, err = dr.ApplyDocument(processed)
		if err != nil {
			return nil, err
		}
	}

	// Reconstruct: non-matching lines stay in place,
	// processed lines fill in where matching lines were.
	var result []string
	processedIdx := 0
	for i, line := range lines {
		if isMatch[i] {
			if processedIdx < len(processed) {
				result = append(result, processed[processedIdx])
				processedIdx++
			}
			// else: inner rules consumed this line (e.g. join reduced line count)
		} else {
			result = append(result, line)
		}
	}
	// Append any extra lines produced by inner rules
	for processedIdx < len(processed) {
		result = append(result, processed[processedIdx])
		processedIdx++
	}

	return result, nil
}
