package rule

import (
	"testing"
)

func TestDeleteContextRule(t *testing.T) {
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
			name:    "context=1 removes match and surrounding",
			pattern: "MATCH",
			before:  1, after: 1,
			input: []string{"a", "b", "MATCH", "c", "d"},
			want:  []string{"a", "d"},
		},
		{
			name:    "context=2",
			pattern: "MATCH",
			before:  2, after: 2,
			input: []string{"a", "b", "c", "MATCH", "d", "e", "f"},
			want:  []string{"a", "f"},
		},
		{
			name:    "after only",
			pattern: "MATCH",
			before:  0, after: 2,
			input: []string{"a", "MATCH", "b", "c", "d"},
			want:  []string{"a", "d"},
		},
		{
			name:    "before only",
			pattern: "MATCH",
			before:  2, after: 0,
			input: []string{"a", "b", "c", "MATCH", "d"},
			want:  []string{"a", "d"},
		},
		{
			name:    "no match returns all lines",
			pattern: "NOMATCH",
			before:  2, after: 2,
			input: []string{"a", "b", "c"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:    "match at start",
			pattern: "MATCH",
			before:  2, after: 1,
			input: []string{"MATCH", "a", "b", "c"},
			want:  []string{"b", "c"},
		},
		{
			name:    "match at end",
			pattern: "MATCH",
			before:  1, after: 2,
			input: []string{"a", "b", "MATCH"},
			want:  []string{"a"},
		},
		{
			name:    "overlapping delete ranges",
			pattern: "M",
			before:  1, after: 1,
			input: []string{"a", "M1", "b", "M2", "c"},
			want:  []string{},
		},
		{
			name:    "multiple deletes no overlap",
			pattern: "M",
			before:  0, after: 0,
			input: []string{"a", "M1", "b", "c", "M2", "d"},
			want:  []string{"a", "b", "c", "d"},
		},
		{
			name:    "context=0 same as plain delete",
			pattern: "MATCH",
			before:  0, after: 0,
			input: []string{"a", "MATCH", "b"},
			want:  []string{"a", "b"},
		},
		{
			name:    "case insensitive",
			pattern: "match",
			before:  1, after: 0,
			opts:    []RuleOption{WithIgnoreCase()},
			input:   []string{"a", "b", "MATCH", "c"},
			want:    []string{"a", "c"},
		},
		{
			name:    "delete all with large context",
			pattern: "MATCH",
			before:  10, after: 10,
			input: []string{"a", "b", "MATCH", "c", "d"},
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewDeleteContextRule(tt.pattern, tt.before, tt.after, tt.opts...)
			if err != nil {
				t.Fatalf("NewDeleteContextRule() error: %v", err)
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
