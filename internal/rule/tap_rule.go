package rule

import (
	"fmt"
	"os"
	"strings"
)

// TapRule writes every line to stderr for debugging, then passes it through unchanged.
// An optional header is printed once before the first line.
type TapRule struct {
	header string
}

// NewTapRule creates a rule that echoes lines to stderr.
// If header is non-empty, it is printed (decorated) before the first line.
func NewTapRule(header string) *TapRule {
	return &TapRule{header: header}
}

// Apply writes the line to stderr and returns it unchanged.
func (r *TapRule) Apply(line string, ctx *LineContext) ([]string, error) {
	if r.header != "" {
		printed := GetState(ctx, r, false)
		if !printed {
			fmt.Fprintf(os.Stderr, "###### %s ######\n", r.header)
			SetState(ctx, r, true)
		}
	}
	fmt.Fprintln(os.Stderr, line)
	return []string{line}, nil
}

// Flush prints a trailing blank line to stderr when a header was used,
// to visually separate tap sections.
func (r *TapRule) Flush(ctx *LineContext) ([]string, error) {
	if r.header != "" {
		printed := GetState(ctx, r, false)
		if printed {
			fmt.Fprintln(os.Stderr, strings.Repeat("#", 6+len(r.header)+8))
		}
	}
	return nil, nil
}
