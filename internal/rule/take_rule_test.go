package rule

import "testing"

func TestTakeRule(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		joiner  string
		opts    []RuleOption
		input   string
		want    []string
	}{
		{
			name:    "extracts match",
			pattern: `\d+`,
			joiner:  " ",
			input:   "abc 123 def",
			want:    []string{"123"},
		},
		{
			name:    "no match passes through",
			pattern: `\d+`,
			joiner:  " ",
			input:   "no numbers here",
			want:    []string{"no numbers here"},
		},
		{
			name:    "extracts first capture group",
			pattern: `(\w+)@(\w+)`,
			joiner:  " ",
			input:   "email: user@host ok",
			want:    []string{"user"},
		},
		{
			name:    "no capture group returns full match",
			pattern: `\w+@\w+`,
			joiner:  " ",
			input:   "email: user@host ok",
			want:    []string{"user@host"},
		},
		{
			name:    "global collects all matches with default joiner",
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
			name:    "global with multi-char joiner",
			pattern: `\d+`,
			joiner:  ", ",
			opts:    []RuleOption{WithGlobal()},
			input:   "a1 b2 c3",
			want:    []string{"1, 2, 3"},
		},
		{
			name:    "global with empty joiner",
			pattern: `\d+`,
			joiner:  "",
			opts:    []RuleOption{WithGlobal()},
			input:   "a1 b2 c3",
			want:    []string{"123"},
		},
		{
			name:    "global with capture groups",
			pattern: `(\d+)\w`,
			joiner:  " ",
			opts:    []RuleOption{WithGlobal()},
			input:   "1a 2b 3c",
			want:    []string{"1 2 3"},
		},
		{
			name:    "first match only without global",
			pattern: `\d+`,
			joiner:  " ",
			input:   "a1 b2 c3",
			want:    []string{"1"},
		},
		{
			name:    "case insensitive",
			pattern: `hello`,
			joiner:  " ",
			opts:    []RuleOption{WithIgnoreCase()},
			input:   "say HELLO world",
			want:    []string{"HELLO"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewTakeRule(tt.pattern, tt.joiner, tt.opts...)
			if err != nil {
				t.Fatalf("NewTakeRule() error: %v", err)
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
