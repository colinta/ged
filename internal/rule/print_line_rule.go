package rule

import "github.com/dlclark/regexp2"

// printState holds the buffered state for context-aware printing.
type printState struct {
	buffer     []string // rolling buffer of recent non-matching lines (for "before")
	afterCount int      // remaining lines to include after a match
}

// PrintLineRule keeps lines that match a pattern, deletes non-matching lines.
// When before/after context is configured, it buffers surrounding lines
// and emits them alongside matches.
type PrintLineRule struct {
	patternStr string
	pattern    *regexp2.Regexp
	before     int // lines of context before match
	after      int // lines of context after match
}

// Pattern returns the original pattern string.
func (r *PrintLineRule) Pattern() string { return r.patternStr }

// NewPrintLineRule creates a rule that keeps only lines matching the pattern.
// Use WithIgnoreCase() for case-insensitive matching.
func NewPrintLineRule(patternStr string, opts ...RuleOption) (*PrintLineRule, error) {
	patternRegex, err := CompilePattern(patternStr, opts...)
	if err != nil {
		return nil, err
	}
	return &PrintLineRule{
		patternStr: patternStr,
		pattern:    patternRegex,
	}, nil
}

// NewPrintContextRule creates a print rule with before/after context lines.
func NewPrintContextRule(patternStr string, before, after int, opts ...RuleOption) (*PrintLineRule, error) {
	patternRegex, err := CompilePattern(patternStr, opts...)
	if err != nil {
		return nil, err
	}
	return &PrintLineRule{
		patternStr: patternStr,
		pattern:    patternRegex,
		before:     before,
		after:      after,
	}, nil
}

// HasContext returns true if this rule has before/after context configured.
func (r *PrintLineRule) HasContext() bool {
	return r.before > 0 || r.after > 0
}

// Apply returns the line if it matches, empty slice if not.
// When context is configured, uses LineContext state to buffer surrounding lines.
func (r *PrintLineRule) Apply(line string, ctx *LineContext) ([]string, error) {
	matched, err := r.pattern.MatchString(line)
	if err != nil {
		return nil, err
	}

	// No context — simple mode
	if !r.HasContext() {
		if matched {
			return []string{line}, nil
		}
		return []string{}, nil
	}

	// Context mode — use state for buffering
	st := GetState(ctx, r, printState{})

	if matched {
		// Flush before-buffer + match line
		result := make([]string, 0, len(st.buffer)+1)
		result = append(result, st.buffer...)
		result = append(result, line)
		st.buffer = nil
		st.afterCount = r.after
		SetState(ctx, r, st)
		return result, nil
	}

	// Non-matching line
	if st.afterCount > 0 {
		// Within after-context range
		st.afterCount--
		SetState(ctx, r, st)
		return []string{line}, nil
	}

	// Outside any context — buffer for potential future match
	if r.before > 0 {
		st.buffer = append(st.buffer, line)
		if len(st.buffer) > r.before {
			st.buffer = st.buffer[len(st.buffer)-r.before:]
		}
	}
	SetState(ctx, r, st)
	return []string{}, nil
}

// Flush discards the unflushed buffer — no match followed, so before-lines are dropped.
func (r *PrintLineRule) Flush(ctx *LineContext) ([]string, error) {
	if !r.HasContext() {
		return nil, nil
	}
	// Clear state, return nothing (unflushed = no match came)
	SetState(ctx, r, printState{})
	return nil, nil
}
