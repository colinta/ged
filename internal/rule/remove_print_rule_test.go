package rule

import "testing"

func TestRemovePrintRule(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		opts    []RuleOption
		input   string
		want    []string
	}{
		{
			name:    "removes match from matching line",
			pattern: `\d+`,
			input:   "abc 123 def 456",
			want:    []string{"abc  def 456"},
		},
		{
			name:    "removes non-matching line",
			pattern: `\d+`,
			input:   "no numbers here",
			want:    []string{},
		},
		{
			name:    "global removes all matches",
			pattern: `\d+`,
			opts:    []RuleOption{WithGlobal()},
			input:   "a1 b2 c3",
			want:    []string{"a b c"},
		},
		{
			name:    "removes regex match",
			pattern: `\s*#.*$`,
			input:   "code  # comment",
			want:    []string{"code"},
		},
		{
			name:    "case insensitive",
			pattern: `hello`,
			opts:    []RuleOption{WithIgnoreCase()},
			input:   "say HELLO world",
			want:    []string{"say  world"},
		},
		{
			name:    "case insensitive no match removed",
			pattern: `hello`,
			input:   "say HELLO world",
			want:    []string{},
		},
		{
			name:    "global removes all from matching line",
			pattern: `x`,
			opts:    []RuleOption{WithGlobal()},
			input:   "xaxbxc",
			want:    []string{"abc"},
		},
		{
			name:    "non-global removes only first",
			pattern: `x`,
			input:   "xaxbxc",
			want:    []string{"axbxc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewRemovePrintRule(tt.pattern, tt.opts...)
			if err != nil {
				t.Fatalf("NewRemovePrintRule() error: %v", err)
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
