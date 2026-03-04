// Package diff computes unified diffs between two sets of lines.
package diff

import "fmt"

// OpKind describes what happened to a line.
type OpKind int

const (
	Equal  OpKind = iota // line is unchanged
	Insert               // line was added
	Delete               // line was removed
)

// Change represents a single line in the diff output.
type Change struct {
	Kind OpKind
	Text string
}

// Compute returns the sequence of changes to transform `a` into `b`
// using the longest common subsequence (LCS) algorithm.
func Compute(a, b []string) []Change {
	n, m := len(a), len(b)

	// Build LCS table using dynamic programming.
	// lcs[i][j] = length of LCS of a[:i] and b[:j]
	lcs := make([][]int, n+1)
	for i := range lcs {
		lcs[i] = make([]int, m+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if a[i-1] == b[j-1] {
				lcs[i][j] = lcs[i-1][j-1] + 1
			} else if lcs[i-1][j] >= lcs[i][j-1] {
				lcs[i][j] = lcs[i-1][j]
			} else {
				lcs[i][j] = lcs[i][j-1]
			}
		}
	}

	// Backtrack to produce the diff
	var changes []Change
	i, j := n, m
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && a[i-1] == b[j-1] {
			changes = append(changes, Change{Equal, a[i-1]})
			i--
			j--
		} else if j > 0 && (i == 0 || lcs[i][j-1] >= lcs[i-1][j]) {
			changes = append(changes, Change{Insert, b[j-1]})
			j--
		} else {
			changes = append(changes, Change{Delete, a[i-1]})
			i--
		}
	}

	// Reverse — we built changes from the end
	for left, right := 0, len(changes)-1; left < right; left, right = left+1, right-1 {
		changes[left], changes[right] = changes[right], changes[left]
	}

	return changes
}

// Format renders a diff as a unified-style text.
// If color is true, deletions are red and insertions are green.
func Format(changes []Change, color bool) []string {
	var lines []string

	red := ""
	green := ""
	reset := ""
	if color {
		red = "\033[31m"
		green = "\033[32m"
		reset = "\033[0m"
	}

	for _, c := range changes {
		switch c.Kind {
		case Equal:
			lines = append(lines, " "+c.Text)
		case Delete:
			lines = append(lines, fmt.Sprintf("%s-%s%s", red, c.Text, reset))
		case Insert:
			lines = append(lines, fmt.Sprintf("%s+%s%s", green, c.Text, reset))
		}
	}

	return lines
}

// HasChanges returns true if the diff contains any insertions or deletions.
func HasChanges(changes []Change) bool {
	for _, c := range changes {
		if c.Kind != Equal {
			return true
		}
	}
	return false
}
