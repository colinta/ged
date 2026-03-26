package rule

// UnquoteRule removes matching quote characters from the start and end of a line.
// Each character in chars is an opening quote; its closing pair is determined
// automatically (e.g. '[' → ']', '(' → ')', '{' → '}', '<' → '>').
// For non-bracket characters the closing quote equals the opening quote.
type UnquoteRule struct {
	pairs []quotePair
}

type quotePair struct {
	open  byte
	close byte
}

// closingFor returns the matching closing character for an opener.
func closingFor(open byte) byte {
	switch open {
	case '[':
		return ']'
	case '(':
		return ')'
	case '{':
		return '}'
	case '<':
		return '>'
	default:
		return open
	}
}

// NewUnquoteRule creates a rule that strips matching quotes.
// Default chars (when chars is empty) are ' and ".
func NewUnquoteRule(chars string) *UnquoteRule {
	if chars == "" {
		chars = `'"`
	}
	pairs := make([]quotePair, len(chars))
	for i := 0; i < len(chars); i++ {
		pairs[i] = quotePair{open: chars[i], close: closingFor(chars[i])}
	}
	return &UnquoteRule{pairs: pairs}
}

// Apply removes a matching open/close quote pair from the line ends.
func (r *UnquoteRule) Apply(line string, ctx *LineContext) ([]string, error) {
	if len(line) < 2 {
		return []string{line}, nil
	}
	first := line[0]
	last := line[len(line)-1]
	for _, p := range r.pairs {
		if first == p.open && last == p.close {
			return []string{line[1 : len(line)-1]}, nil
		}
	}
	return []string{line}, nil
}
