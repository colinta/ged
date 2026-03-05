package rule

import "strings"

// BorderRule prepends and appends lines to the document.
type BorderRule struct {
	text string
}

// NewBorderRule creates a new BorderRule. The text is added to both beginning and end.
func NewBorderRule(text string) *BorderRule {
	return &BorderRule{text: text}
}

// ApplyDocument adds the border text to both ends of the document.
func (r *BorderRule) ApplyDocument(lines []string) ([]string, error) {
	border := strings.Split(r.text, "\n")
	result := make([]string, 0, len(border)+len(lines)+len(border))
	result = append(result, border...)
	result = append(result, lines...)
	result = append(result, border...)
	return result, nil
}
