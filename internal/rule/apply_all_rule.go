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
// Rules implementing SetupRule have Setup called once before processing.
// After processing each line, ctx.Printing is checked to decide inclusion.
func (r *ApplyAllRule) ApplyDocument(lines []string) ([]string, error) {
	var result []string
	ctx := &LineContext{}

	// Call Setup on any rules that need it
	for _, lr := range r.rules {
		if s, ok := lr.(SetupRule); ok {
			s.Setup(ctx)
		}
	}

	for i, line := range lines {
		ctx.LineNum = i + 1
		// Process this line through all rules
		current := []string{line}

		for _, lr := range r.rules {
			var next []string
			for _, l := range current {
				out, err := lr.Apply(l, ctx)
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

		// Check print state after processing
		if ctx.Printing == PrintOff {
			continue
		}

		result = append(result, current...)
	}

	return result, nil
}
