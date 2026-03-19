package rule

import (
	"strings"

	"github.com/dlclark/regexp2"
)

// UniqRule removes consecutive duplicate lines, like the Unix `uniq` command.
// When a pattern is provided, the matched portion (or a specific capture group)
// is used as the unique key instead of the full line.
type UniqRule struct {
	pattern    *regexp2.Regexp // nil means compare full lines
	groupNum   int             // 0 = full match, 1-9 = capture group
	ignoreCase bool            // normalize keys to lowercase
}

// NewUniqRule creates a new UniqRule that compares full lines.
func NewUniqRule() *UniqRule {
	return &UniqRule{}
}

// NewUniqPatternRule creates a UniqRule that uses the matched portion as the key.
// groupNum 0 means use the full match; 1-9 means use that capture group.
func NewUniqPatternRule(pattern string, groupNum int, opts ...RuleOption) (*UniqRule, error) {
	cfg := buildConfig(opts)
	compiled, err := CompilePattern(pattern, opts...)
	if err != nil {
		return nil, err
	}
	return &UniqRule{pattern: compiled, groupNum: groupNum, ignoreCase: cfg.ignoreCase}, nil
}

// keyFor extracts the unique key for a line.
// If no pattern is set, returns the full line.
// If the pattern doesn't match, returns the full line (so non-matching lines
// are compared by their full content, same as plain uniq).
func (r *UniqRule) keyFor(line string) string {
	if r.pattern == nil {
		return line
	}

	m, err := r.pattern.FindStringMatch(line)
	if err != nil || m == nil {
		return line
	}

	groups := m.Groups()
	if r.groupNum > 0 && r.groupNum < len(groups) && groups[r.groupNum].Length > 0 {
		key := groups[r.groupNum].String()
		if r.ignoreCase {
			key = strings.ToLower(key)
		}
		return key
	}

	// Full match (group 0)
	key := groups[0].String()
	if r.ignoreCase {
		key = strings.ToLower(key)
	}
	return key
}

// ApplyDocument removes duplicates.
// Without a pattern: removes consecutive duplicates (like Unix `uniq`).
// With a pattern: removes all lines whose key was already seen (global dedup).
func (r *UniqRule) ApplyDocument(lines []string) ([]string, error) {
	if len(lines) == 0 {
		return lines, nil
	}

	if r.pattern == nil {
		// Consecutive dedup (original behavior)
		result := []string{lines[0]}
		for i := 1; i < len(lines); i++ {
			if lines[i] != lines[i-1] {
				result = append(result, lines[i])
			}
		}
		return result, nil
	}

	// Global dedup by pattern key
	seen := make(map[string]bool)
	var result []string
	for _, line := range lines {
		key := r.keyFor(line)
		if !seen[key] {
			seen[key] = true
			result = append(result, line)
		}
	}
	return result, nil
}
