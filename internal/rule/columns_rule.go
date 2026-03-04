package rule

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
)

// ColumnSpec represents an ordered list of column selections.
// Each entry is either a single index or a range.
type ColumnSpec struct {
	entries []colEntry
}

type colEntry struct {
	start int // 1-based, negative means from end
	end   int // 0 means single column, otherwise range end (1-based, negative ok)
}

// ParseColumnSpec parses a column specification string like "1,3-5,-1".
// Indices are 1-based. Negative indices count from the end (-1 = last).
// Ranges are inclusive: "2-4" means columns 2, 3, 4.
// Open ranges: "3-" means column 3 through last, "-3-" means third-from-end through last.
func ParseColumnSpec(spec string) (*ColumnSpec, error) {
	if spec == "" {
		return nil, fmt.Errorf("empty column spec")
	}

	parts := strings.Split(spec, ",")
	var entries []colEntry

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		entry, err := parseColEntry(part)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("empty column spec")
	}

	return &ColumnSpec{entries: entries}, nil
}

// parseColEntry parses a single entry like "3", "-1", "2-4", "3-", "-2-".
func parseColEntry(s string) (colEntry, error) {
	// Find the range dash. We need to distinguish negative signs from range dashes.
	// Strategy: find the first dash that isn't at position 0 (which would be a negative sign).
	dashIdx := findRangeDash(s)

	if dashIdx == -1 {
		// Single column
		n, err := strconv.Atoi(s)
		if err != nil || n == 0 {
			return colEntry{}, fmt.Errorf("invalid column index: %q", s)
		}
		return colEntry{start: n}, nil
	}

	// Range: start-end
	startStr := s[:dashIdx]
	endStr := s[dashIdx+1:]

	if startStr == "" {
		return colEntry{}, fmt.Errorf("invalid column range: %q", s)
	}

	start, err := strconv.Atoi(startStr)
	if err != nil || start == 0 {
		return colEntry{}, fmt.Errorf("invalid column index: %q", startStr)
	}

	if endStr == "" {
		// Open range like "3-" — end is 0 to signal "through last"
		return colEntry{start: start, end: -1}, nil // -1 sentinel: resolve to last
	}

	end, err := strconv.Atoi(endStr)
	if err != nil || end == 0 {
		return colEntry{}, fmt.Errorf("invalid column index: %q", endStr)
	}

	return colEntry{start: start, end: end}, nil
}

// findRangeDash finds the index of the dash that separates a range (not a negative sign).
// Returns -1 if there's no range dash.
func findRangeDash(s string) int {
	// Skip leading negative sign
	i := 0
	if i < len(s) && s[i] == '-' {
		i++
	}
	// Skip digits
	for i < len(s) && (s[i] >= '0' && s[i] <= '9') {
		i++
	}
	// If we're at a dash, that's the range separator
	if i < len(s) && s[i] == '-' {
		return i
	}
	return -1
}

// Resolve returns the 0-based indices for a given number of columns.
// Out-of-bounds indices are silently skipped.
func (cs *ColumnSpec) Resolve(numCols int) []int {
	var indices []int
	for _, e := range cs.entries {
		if e.end == 0 {
			// Single column
			idx, ok := resolveColIndex(e.start, numCols)
			if ok {
				indices = append(indices, idx)
			}
		} else {
			// Range
			startIdx, ok1 := resolveColIndex(e.start, numCols)
			endEntry := e.end
			if endEntry == -1 && e.start > 0 {
				// Open range "3-" → through last
				endEntry = numCols
			} else if endEntry == -1 && e.start < 0 {
				// Open range "-3-" → through last
				endEntry = numCols
			}
			endIdx, ok2 := resolveColIndex(endEntry, numCols)
			if !ok1 || !ok2 {
				continue
			}
			if startIdx <= endIdx {
				for i := startIdx; i <= endIdx; i++ {
					indices = append(indices, i)
				}
			} else {
				// Reverse range
				for i := startIdx; i >= endIdx; i-- {
					indices = append(indices, i)
				}
			}
		}
	}
	return indices
}

// resolveColIndex converts a 1-based (possibly negative) column index to 0-based.
// Returns (index, true) on success, (0, false) if out of bounds.
func resolveColIndex(idx, numCols int) (int, bool) {
	if idx < 0 {
		idx = numCols + idx + 1 // -1 → numCols (last, 1-based)
	}
	if idx < 1 || idx > numCols {
		return 0, false
	}
	return idx - 1, true // 1-based → 0-based
}

// ColumnsRule splits each line by a pattern and selects/reorders columns.
type ColumnsRule struct {
	pattern *regexp2.Regexp
	spec    *ColumnSpec
	joiner  string
}

// NewColumnsRule creates a rule that splits lines by pattern and selects columns.
// pattern is the split regex, spec selects which columns, joiner reassembles them.
func NewColumnsRule(pattern *regexp2.Regexp, spec *ColumnSpec, joiner string) *ColumnsRule {
	return &ColumnsRule{
		pattern: pattern,
		spec:    spec,
		joiner:  joiner,
	}
}

// Apply splits the line by the pattern, selects columns per the spec, and joins them.
func (r *ColumnsRule) Apply(line string, ctx *LineContext) ([]string, error) {
	fields, err := regexSplit(r.pattern, line)
	if err != nil {
		return nil, err
	}

	indices := r.spec.Resolve(len(fields))
	if len(indices) == 0 {
		return []string{""}, nil
	}

	selected := make([]string, len(indices))
	for i, idx := range indices {
		selected[i] = fields[idx]
	}

	return []string{strings.Join(selected, r.joiner)}, nil
}

// regexSplit splits a string by a regexp2 pattern, similar to strings.Split but with regex.
func regexSplit(pattern *regexp2.Regexp, s string) ([]string, error) {
	var parts []string
	lastEnd := 0

	m, err := pattern.FindStringMatch(s)
	if err != nil {
		return nil, err
	}

	for m != nil {
		start := m.Index
		end := start + m.Length
		parts = append(parts, s[lastEnd:start])
		lastEnd = end

		m, err = pattern.FindNextMatch(m)
		if err != nil {
			return nil, err
		}
	}

	// Append the remainder after the last match
	parts = append(parts, s[lastEnd:])
	return parts, nil
}
