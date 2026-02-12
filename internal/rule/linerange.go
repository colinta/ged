package rule

import (
	"fmt"
	"strconv"
	"strings"
)

// LineRange represents a set of line numbers.
type LineRange interface {
	Contains(lineNum int) bool
}

// SingleLine matches exactly one line number.
type SingleLine int

func (s SingleLine) Contains(lineNum int) bool {
	return int(s) == lineNum
}

// Range matches a range of line numbers (inclusive).
type Range struct {
	Start int
	End   int
}

func (r *Range) Contains(lineNum int) bool {
	return lineNum >= r.Start && lineNum <= r.End
}

// OpenRange matches from a start line to the end, or from the beginning to an end line.
type OpenRange struct {
	From  int  // starting line (if ToEnd is true)
	To    int  // ending line (if ToEnd is false)
	ToEnd bool // true means "From onwards", false means "up to To"
}

func (r *OpenRange) Contains(lineNum int) bool {
	if r.ToEnd {
		return lineNum >= r.From
	}
	return lineNum <= r.To
}

// CompositeRange combines multiple ranges with OR logic.
type CompositeRange struct {
	Ranges []LineRange
}

func (c *CompositeRange) Contains(lineNum int) bool {
	for _, r := range c.Ranges {
		if r.Contains(lineNum) {
			return true
		}
	}
	return false
}

// ParseLineRange parses a line range specification.
// Supported formats:
//   - "5"     - single line
//   - "2-4"   - range (inclusive)
//   - "5-"    - from line 5 to end
//   - "-5"    - from beginning to line 5
//   - "1,3,5" - comma-separated (combines with OR)
func ParseLineRange(spec string) (LineRange, error) {
	if spec == "" {
		return nil, fmt.Errorf("empty line range")
	}

	// Handle comma-separated ranges
	if strings.Contains(spec, ",") {
		parts := strings.Split(spec, ",")
		ranges := make([]LineRange, 0, len(parts))
		for _, part := range parts {
			r, err := ParseLineRange(strings.TrimSpace(part))
			if err != nil {
				return nil, err
			}
			ranges = append(ranges, r)
		}
		return &CompositeRange{Ranges: ranges}, nil
	}

	// Handle range with dash
	if strings.Contains(spec, "-") {
		// Check for open range "5-" or "-5"
		if strings.HasSuffix(spec, "-") {
			// "5-" means from line 5 onwards
			numStr := strings.TrimSuffix(spec, "-")
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, fmt.Errorf("invalid line number %q: %w", numStr, err)
			}
			return &OpenRange{From: num, ToEnd: true}, nil
		}
		if strings.HasPrefix(spec, "-") {
			// "-5" means up to line 5
			numStr := strings.TrimPrefix(spec, "-")
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, fmt.Errorf("invalid line number %q: %w", numStr, err)
			}
			return &OpenRange{To: num, ToEnd: false}, nil
		}

		// "2-4" means range from 2 to 4
		parts := strings.SplitN(spec, "-", 2)
		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid start line %q: %w", parts[0], err)
		}
		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid end line %q: %w", parts[1], err)
		}
		return &Range{Start: start, End: end}, nil
	}

	// Single line number
	num, err := strconv.Atoi(spec)
	if err != nil {
		return nil, fmt.Errorf("invalid line number %q: %w", spec, err)
	}
	return SingleLine(num), nil
}
