package rule

import "regexp"

// betweenState tracks whether we are currently inside a start/end range.
type betweenState struct {
	inside bool
}

// BetweenLineRule implements LineRule. It applies inner LineRules only to lines
// that fall between a start pattern and an end pattern (inclusive on both ends).
// The range can re-open if the start pattern is seen again after an end.
// State is stored on LineContext via GetState/SetState so the rule is reusable
// across multiple documents.
type BetweenLineRule struct {
	startPattern *regexp.Regexp
	endPattern   *regexp.Regexp
	inverted     bool
	rules        []LineRule
}

// NewBetweenLineRule creates a BetweenLineRule.
func NewBetweenLineRule(startPattern, endPattern *regexp.Regexp, inverted bool, rules []LineRule) *BetweenLineRule {
	return &BetweenLineRule{
		startPattern: startPattern,
		endPattern:   endPattern,
		inverted:     inverted,
		rules:        rules,
	}
}

// Apply checks whether the current line is inside a between range,
// applying inner rules if so (or if inverted, applying when outside).
func (r *BetweenLineRule) Apply(line string, ctx *LineContext) ([]string, error) {
	bs := GetState(ctx, r, betweenState{})

	// Check for start/end transitions
	if !bs.inside && r.startPattern.MatchString(line) {
		bs.inside = true
		SetState(ctx, r, bs)
	}

	active := bs.inside
	if r.inverted {
		active = !active
	}

	// Check for end pattern before applying rules — the end line is still "inside"
	closingThisLine := bs.inside && r.endPattern.MatchString(line)

	var result []string
	var err error
	if active {
		// Apply inner rules as a pipeline
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
		result = current
	} else {
		result = []string{line}
		err = nil
	}

	if closingThisLine {
		bs.inside = false
		SetState(ctx, r, bs)
	}

	return result, err
}

// BetweenDocRule implements DocumentRule. It collects lines inside between
// ranges into a sub-document, applies inner DocumentRules to that sub-document,
// then weaves the results back into their original positions.
type BetweenDocRule struct {
	startPattern *regexp.Regexp
	endPattern   *regexp.Regexp
	inverted     bool
	rules        []DocumentRule
}

// NewBetweenDocRule creates a BetweenDocRule.
func NewBetweenDocRule(startPattern, endPattern *regexp.Regexp, inverted bool, rules []DocumentRule) *BetweenDocRule {
	return &BetweenDocRule{
		startPattern: startPattern,
		endPattern:   endPattern,
		inverted:     inverted,
		rules:        rules,
	}
}

// ApplyDocument collects lines inside between ranges, applies inner rules,
// then reconstructs the output.
func (r *BetweenDocRule) ApplyDocument(lines []string) ([]string, error) {
	var activeLines []string
	isActive := make([]bool, len(lines))
	inside := false

	for i, line := range lines {
		if !inside && r.startPattern.MatchString(line) {
			inside = true
		}

		active := inside
		if r.inverted {
			active = !active
		}

		if active {
			activeLines = append(activeLines, line)
			isActive[i] = true
		}

		if inside && r.endPattern.MatchString(line) {
			inside = false
		}
	}

	// Apply inner document rules to the active lines
	processed := activeLines
	for _, dr := range r.rules {
		var err error
		processed, err = dr.ApplyDocument(processed)
		if err != nil {
			return nil, err
		}
	}

	// Reconstruct: inactive lines stay in place, processed lines fill active slots
	var result []string
	processedIdx := 0
	for i, line := range lines {
		if isActive[i] {
			if processedIdx < len(processed) {
				result = append(result, processed[processedIdx])
				processedIdx++
			}
		} else {
			result = append(result, line)
		}
	}
	for processedIdx < len(processed) {
		result = append(result, processed[processedIdx])
		processedIdx++
	}

	return result, nil
}
