package rule

// UniqRule removes consecutive duplicate lines, like the Unix `uniq` command.
type UniqRule struct{}

// NewUniqRule creates a new UniqRule.
func NewUniqRule() *UniqRule {
	return &UniqRule{}
}

// ApplyDocument removes consecutive duplicates.
func (r *UniqRule) ApplyDocument(lines []string) ([]string, error) {
	if len(lines) == 0 {
		return lines, nil
	}

	result := []string{lines[0]}
	for i := 1; i < len(lines); i++ {
		if lines[i] != lines[i-1] {
			result = append(result, lines[i])
		}
	}
	return result, nil
}
