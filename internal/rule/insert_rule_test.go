package rule

import (
	"testing"
)

func TestInsertRule(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		text    string
		opts    []RuleOption
		input   string
		want    []string
	}{
		{
			name:    "insert after matching line",
			pattern: `^#`,
			text:    "---",
			input:   "# heading",
			want:    []string{"# heading", "---"},
		},
		{
			name:    "no match passes through",
			pattern: `^#`,
			text:    "---",
			input:   "plain text",
			want:    []string{"plain text"},
		},
		{
			name:    "multi-line insertion",
			pattern: `start`,
			text:    "line1\nline2\nline3",
			input:   "start here",
			want:    []string{"start here", "line1", "line2", "line3"},
		},
		{
			name:    "insert empty string adds one empty line",
			pattern: `match`,
			text:    "",
			input:   "match",
			want:    []string{"match", ""},
		},
		{
			name:    "regex pattern match",
			pattern: `\d+`,
			text:    "  ^ number found",
			input:   "item 42",
			want:    []string{"item 42", "  ^ number found"},
		},
		{
			name:    "case insensitive",
			pattern: `todo`,
			text:    "  FIXME!",
			opts:    []RuleOption{WithIgnoreCase()},
			input:   "TODO: fix this",
			want:    []string{"TODO: fix this", "  FIXME!"},
		},
		{
			name:    "original line is unchanged",
			pattern: `foo`,
			text:    "bar",
			input:   "foo",
			want:    []string{"foo", "bar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewInsertRule(tt.pattern, tt.text, tt.opts...)
			if err != nil {
				t.Fatalf("NewInsertRule() error: %v", err)
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

func TestInsertRuleInvalidPattern(t *testing.T) {
	_, err := NewInsertRule("[invalid", "text")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}
