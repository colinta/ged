package rule

import (
	"strings"

	"github.com/dlclark/regexp2"
)

// InsertRule inserts text after each line that matches a pattern.
// Non-matching lines pass through unchanged.
type InsertRule struct {
	pattern *regexp2.Regexp
	text    string
}

// NewInsertRule creates a rule that inserts text after matching lines.
// The text may contain \n for multi-line insertions.
func NewInsertRule(patternStr string, text string, opts ...RuleOption) (*InsertRule, error) {
	compiled, err := CompilePattern(patternStr, opts...)
	if err != nil {
		return nil, err
	}
	return &InsertRule{pattern: compiled, text: text}, nil
}

// Apply returns the original line, followed by the insertion text if the line matches.
func (r *InsertRule) Apply(line string, ctx *LineContext) ([]string, error) {
	matched, err := r.pattern.MatchString(line)
	if err != nil {
		return nil, err
	}
	if !matched {
		return []string{line}, nil
	}
	// Split text on \n to support multi-line insertions
	inserted := strings.Split(r.text, "\n")
	result := make([]string, 0, 1+len(inserted))
	result = append(result, line)
	result = append(result, inserted...)
	return result, nil
}
