package rule

import "github.com/dlclark/regexp2"

// deleteState holds the buffered state for context-aware deletion.
type deleteState struct {
	delayBuffer []string // lines waiting to be emitted (for "before" delay)
	afterCount  int      // remaining lines to suppress after a match
}

// DeleteLineRule removes lines that match a pattern, keeps non-matching lines.
// When before/after context is configured, it also removes surrounding lines.
// Lines are delayed by `before` positions so they can be retroactively suppressed.
type DeleteLineRule struct {
	patternStr string
	pattern    *regexp2.Regexp
	before     int // lines to also delete before match
	after      int // lines to also delete after match
}

// Pattern returns the original pattern string.
func (r *DeleteLineRule) Pattern() string { return r.patternStr }

// NewDeleteLineRule creates a rule that removes lines matching the pattern.
// Use WithIgnoreCase() for case-insensitive matching.
func NewDeleteLineRule(patternStr string, opts ...RuleOption) (*DeleteLineRule, error) {
	patternRegex, err := CompilePattern(patternStr, opts...)
	if err != nil {
		return nil, err
	}
	return &DeleteLineRule{
		patternStr: patternStr,
		pattern:    patternRegex,
	}, nil
}

// NewDeleteContextRule creates a delete rule with before/after context lines.
func NewDeleteContextRule(patternStr string, before, after int, opts ...RuleOption) (*DeleteLineRule, error) {
	patternRegex, err := CompilePattern(patternStr, opts...)
	if err != nil {
		return nil, err
	}
	return &DeleteLineRule{
		patternStr: patternStr,
		pattern:    patternRegex,
		before:     before,
		after:      after,
	}, nil
}

// HasContext returns true if this rule has before/after context configured.
func (r *DeleteLineRule) HasContext() bool {
	return r.before > 0 || r.after > 0
}

// Apply returns empty slice if line matches, keeps the line if not.
// When context is configured, uses LineContext state to delay output
// and suppress surrounding lines.
func (r *DeleteLineRule) Apply(line string, ctx *LineContext) ([]string, error) {
	matched, err := r.pattern.MatchString(line)
	if err != nil {
		return nil, err
	}

	// No context — simple mode
	if !r.HasContext() {
		if matched {
			return []string{}, nil
		}
		return []string{line}, nil
	}

	// Context mode — use state for delayed output
	st := GetState(ctx, r, deleteState{})

	if st.afterCount > 0 && !matched {
		// Within after-suppression range, suppress this line
		st.afterCount--
		SetState(ctx, r, st)
		return []string{}, nil
	}

	if matched {
		// Discard delay buffer (those are the "before" lines) + suppress match
		st.delayBuffer = nil
		st.afterCount = r.after
		SetState(ctx, r, st)
		return []string{}, nil
	}

	// Non-matching line outside any suppression range
	// If before > 0, delay output: buffer this line, emit oldest if buffer full
	if r.before > 0 {
		st.delayBuffer = append(st.delayBuffer, line)
		if len(st.delayBuffer) > r.before {
			// Buffer exceeded — oldest line is safe to emit
			emit := st.delayBuffer[0]
			st.delayBuffer = st.delayBuffer[1:]
			SetState(ctx, r, st)
			return []string{emit}, nil
		}
		// Buffer not full yet — hold everything
		SetState(ctx, r, st)
		return []string{}, nil
	}

	// No before context, just pass through
	return []string{line}, nil
}

// Flush emits any remaining delayed lines — they survived (no match consumed them).
func (r *DeleteLineRule) Flush(ctx *LineContext) ([]string, error) {
	if !r.HasContext() {
		return nil, nil
	}
	st := GetState(ctx, r, deleteState{})
	result := st.delayBuffer
	st.delayBuffer = nil
	SetState(ctx, r, st)
	return result, nil
}
