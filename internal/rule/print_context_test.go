package rule

import (
	"testing"
)

// helper: process all lines through a rule with context, including flush
func processWithContext(r LineRule, lines []string) ([]string, error) {
	ctx := &LineContext{}
	if s, ok := r.(SetupRule); ok {
		s.Setup(ctx)
	}
	var result []string
	for i, line := range lines {
		ctx.LineNum = i + 1
		out, err := r.Apply(line, ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, out...)
	}
	if f, ok := r.(FlushRule); ok {
		flushed, err := f.Flush(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, flushed...)
	}
	return result, nil
}

func TestPrintContextRule(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		before  int
		after   int
		opts    []RuleOption
		input   []string
		want    []string
	}{
		{
			name:    "context=1",
			pattern: "MATCH",
			before:  1, after: 1,
			input: []string{"a", "b", "MATCH", "c", "d"},
			want:  []string{"b", "MATCH", "c"},
		},
		{
			name:    "context=2",
			pattern: "MATCH",
			before:  2, after: 2,
			input: []string{"a", "b", "c", "MATCH", "d", "e", "f"},
			want:  []string{"b", "c", "MATCH", "d", "e"},
		},
		{
			name:    "match near start",
			pattern: "MATCH",
			before:  3, after: 1,
			input: []string{"a", "MATCH", "b", "c"},
			want:  []string{"a", "MATCH", "b"},
		},
		{
			name:    "match near end",
			pattern: "MATCH",
			before:  1, after: 3,
			input: []string{"a", "b", "MATCH", "c"},
			want:  []string{"b", "MATCH", "c"},
		},
		{
			name:    "match at first line",
			pattern: "MATCH",
			before:  2, after: 1,
			input: []string{"MATCH", "a", "b"},
			want:  []string{"MATCH", "a"},
		},
		{
			name:    "match at last line",
			pattern: "MATCH",
			before:  1, after: 2,
			input: []string{"a", "b", "MATCH"},
			want:  []string{"b", "MATCH"},
		},
		{
			name:    "overlapping contexts merge",
			pattern: "M",
			before:  1, after: 1,
			input: []string{"a", "M1", "b", "M2", "c"},
			want:  []string{"a", "M1", "b", "M2", "c"},
		},
		{
			name:    "multiple matches no overlap",
			pattern: "M",
			before:  0, after: 0,
			input: []string{"a", "M1", "b", "c", "M2", "d"},
			want:  []string{"M1", "M2"},
		},
		{
			name:    "before only",
			pattern: "MATCH",
			before:  2, after: 0,
			input: []string{"a", "b", "c", "MATCH", "d"},
			want:  []string{"b", "c", "MATCH"},
		},
		{
			name:    "after only",
			pattern: "MATCH",
			before:  0, after: 2,
			input: []string{"a", "MATCH", "b", "c", "d"},
			want:  []string{"MATCH", "b", "c"},
		},
		{
			name:    "no match with context outputs nothing",
			pattern: "NOMATCH",
			before:  2, after: 2,
			input: []string{"a", "b", "c"},
			want:  nil,
		},
		{
			name:    "case insensitive with context",
			pattern: "match",
			before:  1, after: 1,
			opts:    []RuleOption{WithIgnoreCase()},
			input:   []string{"a", "b", "MATCH", "c", "d"},
			want:    []string{"b", "MATCH", "c"},
		},
		{
			name:    "context=0 same as plain print",
			pattern: "MATCH",
			before:  0, after: 0,
			input: []string{"a", "MATCH", "b"},
			want:  []string{"MATCH"},
		},
		{
			name:    "adjacent matches share context",
			pattern: "M",
			before:  1, after: 1,
			input: []string{"a", "M1", "M2", "b"},
			want:  []string{"a", "M1", "M2", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewPrintContextRule(tt.pattern, tt.before, tt.after, tt.opts...)
			if err != nil {
				t.Fatalf("NewPrintContextRule() error: %v", err)
			}
			got, err := processWithContext(r, tt.input)
			if err != nil {
				t.Fatalf("processWithContext() error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %v (%d lines), want %v (%d lines)", got, len(got), tt.want, len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
