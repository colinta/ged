package rule

import (
	"strings"
	"testing"

	"github.com/dlclark/regexp2"
)

func TestParseColumnSpec(t *testing.T) {
	tests := []struct {
		name    string
		spec    string
		numCols int
		want    []int // expected 0-based indices
		wantErr bool
	}{
		{name: "single column", spec: "1", numCols: 3, want: []int{0}},
		{name: "single last", spec: "3", numCols: 3, want: []int{2}},
		{name: "multiple columns", spec: "1,3", numCols: 3, want: []int{0, 2}},
		{name: "reorder", spec: "3,1,2", numCols: 3, want: []int{2, 0, 1}},
		{name: "negative last", spec: "-1", numCols: 3, want: []int{2}},
		{name: "negative second-to-last", spec: "-2", numCols: 3, want: []int{1}},
		{name: "range", spec: "2-4", numCols: 5, want: []int{1, 2, 3}},
		{name: "open range", spec: "3-", numCols: 5, want: []int{2, 3, 4}},
		{name: "negative open range", spec: "-2-", numCols: 5, want: []int{3, 4}},
		{name: "reverse range", spec: "3-1", numCols: 5, want: []int{2, 1, 0}},
		{name: "mixed", spec: "1,3-5,-1", numCols: 6, want: []int{0, 2, 3, 4, 5}},
		{name: "out of bounds skipped", spec: "1,10", numCols: 3, want: []int{0}},
		{name: "all out of bounds", spec: "10", numCols: 3, want: nil},
		{name: "duplicate", spec: "1,1", numCols: 3, want: []int{0, 0}},

		// Error cases
		{name: "empty spec", spec: "", wantErr: true},
		{name: "zero index", spec: "0", wantErr: true},
		{name: "invalid text", spec: "abc", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs, err := ParseColumnSpec(tt.spec)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got := cs.Resolve(tt.numCols)
			if len(got) != len(tt.want) {
				t.Fatalf("Resolve(%d) = %v, want %v", tt.numCols, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("Resolve(%d)[%d] = %d, want %d", tt.numCols, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestColumnsRule(t *testing.T) {
	whitespace := mustCompile(t, `\s+`)
	comma := mustCompile(t, `,`)
	commaSpace := mustCompile(t, `,\s*`)

	tests := []struct {
		name    string
		line    string
		pattern string // regex pattern string
		spec    string
		joiner  string
		want    string
	}{
		{
			name:    "select single column",
			line:    "alice 25 engineer",
			pattern: `\s+`,
			spec:    "1",
			joiner:  " ",
			want:    "alice",
		},
		{
			name:    "select multiple columns",
			line:    "alice 25 engineer",
			pattern: `\s+`,
			spec:    "1,3",
			joiner:  " ",
			want:    "alice engineer",
		},
		{
			name:    "reorder columns",
			line:    "alice 25 engineer",
			pattern: `\s+`,
			spec:    "3,1,2",
			joiner:  " ",
			want:    "engineer alice 25",
		},
		{
			name:    "negative index",
			line:    "alice 25 engineer",
			pattern: `\s+`,
			spec:    "-1",
			joiner:  " ",
			want:    "engineer",
		},
		{
			name:    "range of columns",
			line:    "a b c d e",
			pattern: `\s+`,
			spec:    "2-4",
			joiner:  " ",
			want:    "b c d",
		},
		{
			name:    "open range",
			line:    "a b c d e",
			pattern: `\s+`,
			spec:    "3-",
			joiner:  " ",
			want:    "c d e",
		},
		{
			name:    "comma delimiter",
			line:    "alice,25,engineer",
			pattern: `,`,
			spec:    "1,3",
			joiner:  " ",
			want:    "alice engineer",
		},
		{
			name:    "comma delimiter with custom joiner",
			line:    "alice,25,engineer",
			pattern: `,`,
			spec:    "3,1",
			joiner:  " | ",
			want:    "engineer | alice",
		},
		{
			name:    "comma-space pattern",
			line:    "alice, 25, engineer",
			pattern: `,\s*`,
			spec:    "2",
			joiner:  " ",
			want:    "25",
		},
		{
			name:    "multiple spaces collapsed",
			line:    "alice   25   engineer",
			pattern: `\s+`,
			spec:    "1,2,3",
			joiner:  " ",
			want:    "alice 25 engineer",
		},
		{
			name:    "out of bounds skipped",
			line:    "a b c",
			pattern: `\s+`,
			spec:    "1,5",
			joiner:  " ",
			want:    "a",
		},
		{
			name:    "all out of bounds gives empty",
			line:    "a b",
			pattern: `\s+`,
			spec:    "10",
			joiner:  " ",
			want:    "",
		},
		{
			name:    "empty joiner",
			line:    "a b c",
			pattern: `\s+`,
			spec:    "1,2,3",
			joiner:  "",
			want:    "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pat = whitespace
			if tt.pattern == `,` {
				pat = comma
			} else if tt.pattern == `,\s*` {
				pat = commaSpace
			}

			spec, err := ParseColumnSpec(tt.spec)
			if err != nil {
				t.Fatalf("ParseColumnSpec(%q) error: %v", tt.spec, err)
			}

			rule := NewColumnsRule(pat, spec, tt.joiner)
			ctx := &LineContext{LineNum: 1}
			got, err := rule.Apply(tt.line, ctx)
			if err != nil {
				t.Fatalf("Apply() error: %v", err)
			}
			if len(got) != 1 {
				t.Fatalf("Apply() returned %d lines, want 1", len(got))
			}
			if got[0] != tt.want {
				t.Errorf("Apply(%q) = %q, want %q", tt.line, got[0], tt.want)
			}
		})
	}
}

func TestRegexSplit(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		input   string
		want    []string
	}{
		{name: "whitespace", pattern: `\s+`, input: "a b  c", want: []string{"a", "b", "c"}},
		{name: "comma", pattern: `,`, input: "a,b,c", want: []string{"a", "b", "c"}},
		{name: "no match", pattern: `,`, input: "abc", want: []string{"abc"}},
		{name: "leading sep", pattern: `,`, input: ",a,b", want: []string{"", "a", "b"}},
		{name: "trailing sep", pattern: `,`, input: "a,b,", want: []string{"a", "b", ""}},
		{name: "multi-char", pattern: `\s*,\s*`, input: "a , b , c", want: []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pat, err := CompilePattern(tt.pattern)
			if err != nil {
				t.Fatalf("CompilePattern(%q) error: %v", tt.pattern, err)
			}
			got, err := regexSplit(pat, tt.input)
			if err != nil {
				t.Fatalf("regexSplit() error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("regexSplit() = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("regexSplit()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func mustCompile(t *testing.T, pattern string) *regexp2.Regexp {
	t.Helper()
	re, err := CompilePattern(pattern)
	if err != nil {
		t.Fatalf("CompilePattern(%q) error: %v", pattern, err)
	}
	return re
}

// helper for joining test output
func joinLines(lines []string) string {
	return strings.Join(lines, "\n")
}
