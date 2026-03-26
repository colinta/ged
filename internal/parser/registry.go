package parser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/colinta/ged/internal/rule"
	"github.com/dlclark/regexp2"
)

// Command represents a parsed command with its delimiter-separated parts.
// Matchers populate this; handlers consume it.
type Command struct {
	Name      string
	Inverted  bool
	Delimiter byte
	Parts     []string
}

// IsQuoted returns true if the delimiter is a quote character (`, ', ").
func (c Command) IsQuoted() bool {
	return c.Delimiter == '`' || c.Delimiter == '\'' || c.Delimiter == '"'
}

// EscapePattern returns s with regex metacharacters escaped if the delimiter
// is a quote character. Otherwise returns s unchanged.
func (c Command) EscapePattern(s string) string {
	if c.IsQuoted() {
		return regexp2.Escape(s)
	}
	return s
}

// Pattern returns parts[0] with regex escaping applied for quoted delimiters.
// Returns "" if there are no parts.
func (c Command) Pattern() string {
	if len(c.Parts) == 0 {
		return ""
	}
	return c.EscapePattern(c.Parts[0])
}

// RequireParts returns an error if fewer than n non-empty parts are present.
func (c Command) RequireParts(n int) error {
	if len(c.Parts) < n {
		return fmt.Errorf("%s requires %d argument(s)", c.Name, n)
	}
	return nil
}

// RequirePattern returns an error if parts[0] is missing or empty.
func (c Command) RequirePattern() error {
	if len(c.Parts) == 0 || c.Parts[0] == "" {
		return fmt.Errorf("%s requires a pattern", c.Name)
	}
	return nil
}

// RequireDelimiter returns an error if no delimiter was parsed (bare word command).
func (c Command) RequireDelimiter() error {
	if c.Delimiter == 0 {
		return fmt.Errorf("%s requires a delimiter and argument(s)", c.Name)
	}
	return nil
}

// Flags returns rule options parsed from the given part index.
func (c Command) Flags(index int) []rule.RuleOption {
	return flagsFromParts(c.Parts, index)
}

// FlagsAndContext returns both flags and context options from trailing parts.
func (c Command) FlagsAndContext(startIndex int) ([]rule.RuleOption, contextOptions, error) {
	return flagsAndOptionsFromParts(c.Parts, startIndex)
}

// CompilePattern compiles parts[0] as a regex pattern, applying quote-escaping
// and the given flags. Convenience for the common handler pattern.
func (c Command) CompilePattern(opts ...rule.RuleOption) (*regexp2.Regexp, error) {
	return rule.CompilePattern(c.Pattern(), opts...)
}

// ---------------------------------------------------------------------------
// Matcher and Handler types
// ---------------------------------------------------------------------------

// Matcher tests whether input matches a command pattern. On match, it returns
// a populated Command and true. The Command's Parts are already split by
// delimiter if applicable.
type Matcher func(input string) (Command, bool)

// Handler processes a matched Command and returns a rule (or condition).
type Handler func(cmd Command) (any, error)

// ---------------------------------------------------------------------------
// Registry
// ---------------------------------------------------------------------------

type registration struct {
	matcher Matcher
	handler Handler
}

type registry struct {
	entries []registration
}

func (r *registry) Add(m Matcher, h Handler) {
	r.entries = append(r.entries, registration{m, h})
}

func (r *registry) Parse(input string) (any, error) {
	for _, e := range r.entries {
		if cmd, ok := e.matcher(input); ok {
			return e.handler(cmd)
		}
	}
	return nil, fmt.Errorf("unknown command: %s", input)
}

// ---------------------------------------------------------------------------
// Matchers
// ---------------------------------------------------------------------------

// Exact matches when the full input equals one of the given names.
// No delimiter parsing is performed.
func Exact(names ...string) Matcher {
	return func(input string) (Command, bool) {
		for _, n := range names {
			if input == n {
				return Command{Name: n}, true
			}
		}
		return Command{}, false
	}
}

// Prefix matches when the input starts with one of the given names followed
// by a valid delimiter character. The first name is the canonical name stored
// in Command.Name; the rest are aliases. The text after the name is split
// by the delimiter into Command.Parts.
func Prefix(names ...string) Matcher {
	canonical := names[0]
	// Sort by length descending so longer names match first
	// (e.g., "substitute" before "sub")
	sorted := make([]string, len(names))
	copy(sorted, names)
	sort.Slice(sorted, func(i, j int) bool {
		return len(sorted[i]) > len(sorted[j])
	})

	return func(input string) (Command, bool) {
		for _, n := range sorted {
			if !strings.HasPrefix(input, n) {
				continue
			}
			rest := input[len(n):]
			if len(rest) == 0 {
				continue // bare word — let Exact handle it
			}
			delim := rest[0]
			if !isValidDelimiter(delim) {
				continue
			}
			parts, err := splitByDelimiter(rest[1:], delim)
			if err != nil {
				continue
			}
			return Command{
				Name:      canonical,
				Delimiter: delim,
				Parts:     parts,
			}, true
		}
		return Command{}, false
	}
}

// NegPrefix is like Prefix but also matches a leading "!" which sets
// Command.Inverted to true. The canonical name is names[0] (without "!").
func NegPrefix(names ...string) Matcher {
	positive := Prefix(names...)

	negNames := make([]string, len(names))
	for i, n := range names {
		negNames[i] = "!" + n
	}
	negative := Prefix(negNames...)

	return func(input string) (Command, bool) {
		if cmd, ok := negative(input); ok {
			cmd.Name = names[0] // canonical, without "!"
			cmd.Inverted = true
			return cmd, true
		}
		return positive(input)
	}
}

// CharRange matches when the first byte of input is in [lo, hi] and the
// second byte is a valid delimiter. The matched character is stored as
// Command.Name.
func CharRange(lo, hi byte) Matcher {
	return func(input string) (Command, bool) {
		if len(input) < 2 {
			return Command{}, false
		}
		ch := input[0]
		if ch < lo || ch > hi {
			return Command{}, false
		}
		delim := input[1]
		if !isValidDelimiter(delim) {
			return Command{}, false
		}
		parts, err := splitByDelimiter(input[2:], delim)
		if err != nil {
			return Command{}, false
		}
		return Command{
			Name:      string(ch),
			Delimiter: delim,
			Parts:     parts,
		}, true
	}
}

// ---------------------------------------------------------------------------
// Handler helpers
// ---------------------------------------------------------------------------

// Returns creates a Handler that ignores the Command and returns fn().
// Useful for zero-argument word commands like "sort", "trim", etc.
func Returns[T any](fn func() T) Handler {
	return func(_ Command) (any, error) {
		return fn(), nil
	}
}
