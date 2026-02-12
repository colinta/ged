// Package rule defines the Rule interfaces and implementations for text transformation.
package rule

// LineRule is the core interface for per-line text transformation.
// Apply takes a line of text and its line number (1-indexed) and returns:
//   - []string with transformed line(s) - could be 0, 1, or many lines
//   - error if something goes wrong
type LineRule interface {
	Apply(line string, lineNum int) ([]string, error)
}

// DocumentRule operates on all lines at once.
// ApplyDocument takes the entire document as a slice of lines and returns
// the transformed document.
type DocumentRule interface {
	ApplyDocument(lines []string) ([]string, error)
}
