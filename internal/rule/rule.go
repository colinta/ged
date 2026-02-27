// Package rule defines the Rule interfaces and implementations for text transformation.
package rule

import "github.com/dlclark/regexp2"

// PrintState controls whether lines are included in output.
type PrintState int

const (
	PrintDefault PrintState = iota // no control rule active, print everything
	PrintOn                        // printing is enabled
	PrintOff                       // printing is suppressed
)

// LineContext carries per-line state through the processing pipeline.
// Rules that need per-document mutable state store it here via GetState/SetState
// rather than on the rule struct, so a single rule pipeline can be shared across
// multiple documents processed in parallel.
type LineContext struct {
	LineNum  int
	Printing PrintState
	state    map[any]any // rule-local state, lazily initialized
}

// GetState retrieves rule-local state from the context.
// The key should be the rule's own pointer (r) to ensure uniqueness.
// Returns defaultVal if no state has been set for this key.
func GetState[T any](ctx *LineContext, key any, defaultVal T) T {
	if ctx.state == nil {
		return defaultVal
	}
	v, ok := ctx.state[key]
	if !ok {
		return defaultVal
	}
	return v.(T)
}

// SetState stores rule-local state on the context.
func SetState(ctx *LineContext, key any, val any) {
	if ctx.state == nil {
		ctx.state = make(map[any]any)
	}
	ctx.state[key] = val
}

// LineRule is the core interface for per-line text transformation.
// Apply takes a line of text and a context (carrying line number and shared state)
// and returns:
//   - []string with transformed line(s) - could be 0, 1, or many lines
//   - error if something goes wrong
type LineRule interface {
	Apply(line string, ctx *LineContext) ([]string, error)
}

// SetupRule is an optional interface for rules that need to initialize
// shared state on the LineContext before processing begins.
// The caller checks for this with a type assertion and calls Setup once
// before the processing loop.
type SetupRule interface {
	Setup(ctx *LineContext)
}

// DocumentRule operates on all lines at once.
// ApplyDocument takes the entire document as a slice of lines and returns
// the transformed document.
type DocumentRule interface {
	ApplyDocument(lines []string) ([]string, error)
}

// --- Shared rule options and pattern compilation ---

// ruleConfig holds parsed option state used during rule construction.
type ruleConfig struct {
	ignoreCase bool
	global     bool
}

// RuleOption configures rule behavior. Shared across all regex-based rules.
// Options that don't apply to a particular rule are silently ignored.
type RuleOption func(*ruleConfig)

// WithIgnoreCase makes pattern matching case-insensitive.
func WithIgnoreCase() RuleOption {
	return func(c *ruleConfig) {
		c.ignoreCase = true
	}
}

// WithGlobal makes substitution replace all matches, not just the first.
// Only meaningful for SubstitutionRule.
func WithGlobal() RuleOption {
	return func(c *ruleConfig) {
		c.global = true
	}
}

// buildConfig applies options and returns the resolved config.
func buildConfig(opts []RuleOption) ruleConfig {
	var cfg ruleConfig
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

// CompilePattern compiles a regex pattern with ECMAScript mode enabled by default.
// If the config includes ignoreCase, IgnoreCase is ORed into the options.
func CompilePattern(pattern string, opts ...RuleOption) (*regexp2.Regexp, error) {
	cfg := buildConfig(opts)
	options := regexp2.RegexOptions(regexp2.ECMAScript)
	if cfg.ignoreCase {
		options |= regexp2.IgnoreCase
	}
	return regexp2.Compile(pattern, options)
}
