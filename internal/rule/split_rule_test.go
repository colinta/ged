package rule

import (
	"testing"
)

func TestSplitRule(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		opts    []RuleOption
		input   string
		want    []string
	}{
		{
			name:    "split on comma",
			pattern: `,`,
			input:   "a,b,c",
			want:    []string{"a", "b", "c"},
		},
		{
			name:    "split on comma-space",
			pattern: `,\s*`,
			input:   "a, b, c",
			want:    []string{"a", "b", "c"},
		},
		{
			name:    "split on whitespace",
			pattern: `\s+`,
			input:   "hello   world   foo",
			want:    []string{"hello", "world", "foo"},
		},
		{
			name:    "no match returns original line",
			pattern: `,`,
			input:   "no commas here",
			want:    []string{"no commas here"},
		},
		{
			name:    "split on pipe",
			pattern: `\s*\|\s*`,
			input:   "a | b | c",
			want:    []string{"a", "b", "c"},
		},
		{
			name:    "split produces empty strings at boundaries",
			pattern: `,`,
			input:   ",a,,b,",
			want:    []string{"", "a", "", "b", ""},
		},
		{
			name:    "split single character line",
			pattern: `,`,
			input:   "x",
			want:    []string{"x"},
		},
		{
			name:    "split empty line",
			pattern: `,`,
			input:   "",
			want:    []string{""},
		},
		{
			name:    "case insensitive split",
			pattern: `and`,
			opts:    []RuleOption{WithIgnoreCase()},
			input:   "fooANDbarANDbaz",
			want:    []string{"foo", "bar", "baz"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewSplitRule(tt.pattern, tt.opts...)
			if err != nil {
				t.Fatalf("NewSplitRule() error: %v", err)
			}
			got, err := r.Apply(tt.input, &LineContext{LineNum: 1})
			if err != nil {
				t.Fatalf("Apply() error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("Apply() = %v (%d lines), want %v (%d lines)", got, len(got), tt.want, len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Apply()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
