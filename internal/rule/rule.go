// Package rule defines the Rule interfaces and implementations for text transformation.
package rule

// LineContext carries per-line state through the processing pipeline.
// It replaces the bare lineNum parameter so we can add shared state later
// (e.g., print-on/off flags, between-range tracking).
type LineContext struct {
	LineNum int
}

// LineRule is the core interface for per-line text transformation.
// Apply takes a line of text and a context (carrying line number and shared state)
// and returns:
//   - []string with transformed line(s) - could be 0, 1, or many lines
//   - error if something goes wrong
type LineRule interface {
	Apply(line string, ctx *LineContext) ([]string, error)
}

// DocumentRule operates on all lines at once.
// ApplyDocument takes the entire document as a slice of lines and returns
// the transformed document.
type DocumentRule interface {
	ApplyDocument(lines []string) ([]string, error)
}
