package rule

import "testing"

func TestTakePrintRule(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		joiner  string
		opts    []RuleOption
		input   string
		want    []string
	}{
		{
			name:    "extracts match from matching line",
			pattern: `\d+`,
			joiner:  " ",
			input:   "abc 123 def",
			want:    []string{"123"},
		},
		{
			name:    "removes non-matching line",
			pattern: `\d+`,
			joiner:  " ",
			input:   "no numbers here",
			want:    []string{},
		},
		{
			name:    "extracts first capture group",
			pattern: `(\w+)@(\w+)`,
			joiner:  " ",
			input:   "email: user@host ok",
			want:    []string{"user"},
		},
		{
			name:    "global collects all matches",
			pattern: `\d+`,
			joiner:  " ",
			opts:    []RuleOption{WithGlobal()},
			input:   "a1 b2 c3",
			want:    []string{"1 2 3"},
		},
		{
			name:    "global with custom joiner",
			pattern: `\d+`,
			joiner:  ",",
			opts:    []RuleOption{WithGlobal()},
			input:   "a1 b2 c3",
			want:    []string{"1,2,3"},
		},
		{
			name:    "case insensitive",
			pattern: `hello`,
			joiner:  " ",
			opts:    []RuleOption{WithIgnoreCase()},
			input:   "say HELLO world",
			want:    []string{"HELLO"},
		},
		{
			name:    "case insensitive no match removed",
			pattern: `hello`,
			joiner:  " ",
			input:   "say HELLO world",
			want:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewTakePrintRule(tt.pattern, tt.joiner, tt.opts...)
			if err != nil {
				t.Fatalf("NewTakePrintRule() error: %v", err)
			}
			got, err := r.Apply(tt.input, &LineContext{LineNum: 1})
			if err != nil {
				t.Fatalf("Apply() error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d lines, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("line %d: got %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
