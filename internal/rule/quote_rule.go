package rule

// QuoteRule wraps each line in a matching open/close quote pair.
// The char is the opening quote; its closing pair is determined by closingFor
// (e.g. '[' → ']', '(' → ')', '{' → '}', '<' → '>').
// For non-bracket characters the closing quote equals the opening quote.
type QuoteRule struct {
	open  string
	close string
}

// NewQuoteRule creates a rule that wraps lines in a quote pair.
// Default char (when char is empty) is ".
func NewQuoteRule(char string) *QuoteRule {
	if char == "" {
		char = `"`
	}
	return &QuoteRule{
		open:  char,
		close: string(closingFor(char[0])),
	}
}

// Apply wraps the line in the open/close quote pair.
func (r *QuoteRule) Apply(line string, ctx *LineContext) ([]string, error) {
	return []string{r.open + line + r.close}, nil
}
