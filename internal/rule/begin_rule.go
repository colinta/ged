package rule

import "strings"

// BeginRule prepends lines to the beginning of the document.
type BeginRule struct {
	text string
}

// NewBeginRule creates a new BeginRule. The text may contain \n for multiple lines.
func NewBeginRule(text string) *BeginRule {
	return &BeginRule{text: text}
}

// ApplyDocument prepends the text to the document.
func (r *BeginRule) ApplyDocument(lines []string) ([]string, error) {
	prefix := strings.Split(r.text, "\n")
	return append(prefix, lines...), nil
}
