package rule

import "testing"

func TestUnquoteRule(t *testing.T) {
	tests := []struct {
		name  string
		chars string
		input string
		want  string
	}{
		// Default chars (' and ")
		{name: "double quotes", chars: "", input: `"hello"`, want: "hello"},
		{name: "single quotes", chars: "", input: "'hello'", want: "hello"},
		{name: "no quotes passthrough", chars: "", input: "hello", want: "hello"},
		{name: "mismatched quotes", chars: "", input: `"hello'`, want: `"hello'`},
		{name: "backtick not default", chars: "", input: "`hello`", want: "`hello`"},
		{name: "empty string", chars: "", input: "", want: ""},
		{name: "single char", chars: "", input: `"`, want: `"`},
		{name: "only quotes", chars: "", input: `""`, want: ""},

		// Custom chars
		{name: "backtick custom", chars: "`", input: "`hello`", want: "hello"},
		{name: "backtick custom no match", chars: "`", input: `"hello"`, want: `"hello"`},
		{name: "multiple custom", chars: "'\"`", input: "`hello`", want: "hello"},
		{name: "multiple custom double", chars: "'\"`", input: `"hello"`, want: "hello"},

		// Bracket pairs
		{name: "square brackets", chars: "[", input: "[hello]", want: "hello"},
		{name: "parens", chars: "(", input: "(hello)", want: "hello"},
		{name: "curly braces", chars: "{", input: "{hello}", want: "hello"},
		{name: "angle brackets", chars: "<", input: "<hello>", want: "hello"},
		{name: "bracket mismatch", chars: "[", input: "[hello)", want: "[hello)"},
		{name: "bracket not in chars", chars: "'\"", input: "[hello]", want: "[hello]"},

		// Mixed brackets and quotes
		{name: "brackets and quotes", chars: `['"`, input: "'hello'", want: "hello"},
		{name: "brackets and quotes 2", chars: `['"`, input: "[hello]", want: "hello"},

		// Interior content preserved
		{name: "inner quotes preserved", chars: "", input: `"it's a "test""`, want: `it's a "test"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewUnquoteRule(tt.chars)
			got, err := r.Apply(tt.input, &LineContext{LineNum: 1})
			if err != nil {
				t.Fatalf("Apply() error: %v", err)
			}
			if len(got) != 1 {
				t.Fatalf("got %d lines, want 1", len(got))
			}
			if got[0] != tt.want {
				t.Errorf("got %q, want %q", got[0], tt.want)
			}
		})
	}
}
