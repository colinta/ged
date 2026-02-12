// Package engine provides the rule processing pipeline.
package engine

import "github.com/colinta/ged/internal/rule"

// Pipeline chains multiple rules together.
// Each rule's output feeds into the next rule.
type Pipeline struct {
	rules []rule.LineRule
}

// NewPipeline creates a pipeline with the given rules.
// Rules are applied in order, with each rule's output feeding into the next.
func NewPipeline(rules ...rule.LineRule) *Pipeline {
	return &Pipeline{rules: rules}
}

// Process applies all rules to a line and returns the results.
// lineNum is the 1-indexed line number of the input line.
// If any rule returns an empty slice, processing stops and empty is returned.
// Each output line from a rule feeds into the next rule.
func (p *Pipeline) Process(line string, lineNum int) ([]string, error) {
	// Start with the input line
	lines := []string{line}

	for _, r := range p.rules {
		var nextLines []string

		// Apply rule to each line from previous stage
		for _, l := range lines {
			result, err := r.Apply(l, lineNum)
			if err != nil {
				return nil, err
			}
			// Collect all outputs
			nextLines = append(nextLines, result...)
		}

		// If no output, stop processing
		if len(nextLines) == 0 {
			return []string{}, nil
		}

		lines = nextLines
	}

	return lines, nil
}
