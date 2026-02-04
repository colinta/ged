// Package rule defines the Rule interface and implementations for text transformation.
package rule

// Rule is the core interface for text transformation.
// Apply takes a line of text and returns:
//   - []string with transformed line(s) - could be 0, 1, or many lines
//   - error if something goes wrong
type Rule interface {
	Apply(line string) ([]string, error)
}
