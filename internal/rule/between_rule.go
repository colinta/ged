package rule

import "github.com/dlclark/regexp2"

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
	startPattern *regexp2.Regexp
	endPattern   *regexp2.Regexp
	inverted     bool
	rules        []LineRule
	elseRules    []LineRule
}

// NewBetweenLineRule creates a BetweenLineRule.
func NewBetweenLineRule(startPattern, endPattern *regexp2.Regexp, inverted bool, rules []LineRule, elseRules []LineRule) *BetweenLineRule {
	return &BetweenLineRule{
		startPattern: startPattern,
		endPattern:   endPattern,
		inverted:     inverted,
		rules:        rules,
		elseRules:    elseRules,
	}
}

// Apply checks whether the current line is inside a between range,
// applying inner rules if so (or if inverted, applying when outside).
func (r *BetweenLineRule) Apply(line string, ctx *LineContext) ([]string, error) {
	bs := GetState(ctx, r, betweenState{})

	// Check for start/end transitions
	if !bs.inside {
		matched, err := r.startPattern.MatchString(line)
		if err != nil {
			return nil, err
		}
		if matched {
			bs.inside = true
			SetState(ctx, r, bs)
		}
	}

	active := bs.inside
	if r.inverted {
		active = !active
	}

	// Check for end pattern before applying rules — the end line is still "inside"
	closingThisLine := false
	if bs.inside {
		matched, err := r.endPattern.MatchString(line)
		if err != nil {
			return nil, err
		}
		closingThisLine = matched
	}

	var result []string
	var err error
	if active {
		result, err = applyLineRules(line, ctx, r.rules)
	} else if len(r.elseRules) > 0 {
		result, err = applyLineRules(line, ctx, r.elseRules)
	} else {
		result = []string{line}
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
	startPattern *regexp2.Regexp
	endPattern   *regexp2.Regexp
	inverted     bool
	rules        []DocumentRule
	elseRules    []DocumentRule
}

// NewBetweenDocRule creates a BetweenDocRule.
func NewBetweenDocRule(startPattern, endPattern *regexp2.Regexp, inverted bool, rules []DocumentRule, elseRules []DocumentRule) *BetweenDocRule {
	return &BetweenDocRule{
		startPattern: startPattern,
		endPattern:   endPattern,
		inverted:     inverted,
		rules:        rules,
		elseRules:    elseRules,
	}
}

// ApplyDocument processes each contiguous between range as its own sub-document.
// Inner rules are applied independently to each range, and the results replace
// the range in-place. Non-active lines pass through unchanged.
// When inverted, the "outside" segments are each treated as their own sub-document.
func (r *BetweenDocRule) ApplyDocument(lines []string) ([]string, error) {
	// First pass: identify contiguous segments of active/inactive lines.
	type segment struct {
		lines  []string
		active bool
	}
	var segments []segment
	inside := false

	for _, line := range lines {
		if !inside {
			matched, err := r.startPattern.MatchString(line)
			if err != nil {
				return nil, err
			}
			if matched {
				inside = true
			}
		}

		active := inside
		if r.inverted {
			active = !active
		}

		// Append to current segment or start a new one
		if len(segments) == 0 || segments[len(segments)-1].active != active {
			segments = append(segments, segment{active: active})
		}
		segments[len(segments)-1].lines = append(segments[len(segments)-1].lines, line)

		if inside {
			matched, err := r.endPattern.MatchString(line)
			if err != nil {
				return nil, err
			}
			if matched {
				inside = false
			}
		}
	}

	// Second pass: apply inner rules to active segments, elseRules to inactive
	var result []string
	for _, seg := range segments {
		if seg.active {
			processed, err := applyDocRules(seg.lines, r.rules)
			if err != nil {
				return nil, err
			}
			result = append(result, processed...)
		} else if len(r.elseRules) > 0 {
			processed, err := applyDocRules(seg.lines, r.elseRules)
			if err != nil {
				return nil, err
			}
			result = append(result, processed...)
		} else {
			result = append(result, seg.lines...)
		}
	}

	return result, nil
}
