package rule

// PrependRule adds text before each line.
type PrependRule struct {
	text string
}

// NewPrependRule creates a rule that prepends text to each line.
func NewPrependRule(text string) *PrependRule {
	return &PrependRule{text: text}
}

// Apply returns the line with text prepended.
func (r *PrependRule) Apply(line string, ctx *LineContext) ([]string, error) {
	return []string{r.text + line}, nil
}

// AppendRule adds text after each line.
type AppendRule struct {
	text string
}

// NewAppendRule creates a rule that appends text to each line.
func NewAppendRule(text string) *AppendRule {
	return &AppendRule{text: text}
}

// Apply returns the line with text appended.
func (r *AppendRule) Apply(line string, ctx *LineContext) ([]string, error) {
	return []string{line + r.text}, nil
}

// SurroundRule wraps each line with before and after text.
type SurroundRule struct {
	before string
	after  string
}

// NewSurroundRule creates a rule that wraps lines with before and after text.
func NewSurroundRule(before, after string) *SurroundRule {
	return &SurroundRule{before: before, after: after}
}

// Apply returns the line wrapped with before and after text.
func (r *SurroundRule) Apply(line string, ctx *LineContext) ([]string, error) {
	return []string{r.before + line + r.after}, nil
}
