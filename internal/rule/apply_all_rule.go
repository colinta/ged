package rule

// ApplyAllRule wraps a slice of LineRules into a DocumentRule.
// It applies the line rules as a pipeline to each line of the document.
// This avoids a circular import with the engine package by inlining
// the pipeline chaining logic.
type ApplyAllRule struct {
	rules []LineRule
}

// NewApplyAllRule creates a DocumentRule that applies line rules to each line.
func NewApplyAllRule(rules []LineRule) *ApplyAllRule {
	return &ApplyAllRule{rules: rules}
}

// ApplyDocument applies the line rules pipeline to each line of the document.
// Each line is processed through all rules in order, with each rule's output
// feeding into the next rule. Line numbers are 1-indexed.
func (r *ApplyAllRule) ApplyDocument(lines []string) ([]string, error) {
	var result []string

	for lineNum, line := range lines {
		// Process this line through all rules
		current := []string{line}

		for _, lr := range r.rules {
			var next []string
			for _, l := range current {
				out, err := lr.Apply(l, lineNum+1)
				if err != nil {
					return nil, err
				}
				next = append(next, out...)
			}

			if len(next) == 0 {
				current = nil
				break
			}
			current = next
		}

		result = append(result, current...)
	}

	return result, nil
}
