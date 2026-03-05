package rule

import "strings"

// EndRule appends lines to the end of the document.
type EndRule struct {
	text string
}

// NewEndRule creates a new EndRule. The text may contain \n for multiple lines.
func NewEndRule(text string) *EndRule {
	return &EndRule{text: text}
}

// ApplyDocument appends the text to the document.
func (r *EndRule) ApplyDocument(lines []string) ([]string, error) {
	suffix := strings.Split(r.text, "\n")
	return append(lines, suffix...), nil
}
