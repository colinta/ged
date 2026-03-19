package rule

import "testing"

func TestUniqPatternRule(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		groupNum int
		opts     []RuleOption
		input    []string
		want     []string
	}{
		{
			name:    "full match as key",
			pattern: `^\w+`,
			input:   []string{"abc this is a line", "123 this is something", "abc this is nice", "456 hello", "123 goodbye"},
			want:    []string{"abc this is a line", "123 this is something", "456 hello"},
		},
		{
			name:    "consecutive same key",
			pattern: `^\w+`,
			input:   []string{"abc first", "abc second", "def third"},
			want:    []string{"abc first", "def third"},
		},
		{
			name:    "non-consecutive same key removed (global dedup)",
			pattern: `^\w+`,
			input:   []string{"abc first", "def second", "abc third"},
			want:    []string{"abc first", "def second"},
		},
		{
			name:    "no match falls back to full line",
			pattern: `\d+`,
			input:   []string{"abc", "abc", "def"},
			want:    []string{"abc", "def"},
		},
		{
			name:     "capture group as key",
			pattern:  `^\w+ (\w+) `,
			groupNum: 1,
			input:    []string{"111 abc this is a line", "111 123 this is something", "111 abc this is nice", "111 456 hello", "111 123 goodbye"},
			want:     []string{"111 abc this is a line", "111 123 this is something", "111 456 hello"},
		},
		{
			name:     "group number out of range uses full match",
			pattern:  `^\w+`,
			groupNum: 5,
			input:    []string{"abc first", "abc second", "def third"},
			want:     []string{"abc first", "def third"},
		},
		{
			name:    "case insensitive",
			pattern: `^\w+`,
			opts:    []RuleOption{WithIgnoreCase()},
			input:   []string{"Abc first", "abc second", "ABC third", "def fourth"},
			want:    []string{"Abc first", "def fourth"},
		},
		{
			name:    "empty input",
			pattern: `\w+`,
			input:   []string{},
			want:    []string{},
		},
		{
			name:    "single line",
			pattern: `\w+`,
			input:   []string{"hello"},
			want:    []string{"hello"},
		},
		{
			name:     "mixed match and no-match lines",
			pattern:  `#(\w+)`,
			groupNum: 1,
			input:    []string{"#foo bar", "no tag here", "#foo baz", "#bar qux"},
			want:     []string{"#foo bar", "no tag here", "#bar qux"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewUniqPatternRule(tt.pattern, tt.groupNum, tt.opts...)
			if err != nil {
				t.Fatalf("unexpected error creating rule: %v", err)
			}
			got, err := r.ApplyDocument(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d lines, want %d:\ngot:  %v\nwant: %v", len(got), len(tt.want), got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("line[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestUniqPatternRuleInvalidPattern(t *testing.T) {
	_, err := NewUniqPatternRule("[invalid", 0)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}
