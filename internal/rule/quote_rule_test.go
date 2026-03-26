package rule

import "testing"

func TestQuoteRule(t *testing.T) {
	tests := []struct {
		name string
		char string
		input string
		want  string
	}{
		{name: "default double quote", char: "", input: "hello", want: `"hello"`},
		{name: "single quote", char: "'", input: "hello", want: "'hello'"},
		{name: "backtick", char: "`", input: "hello", want: "`hello`"},
		{name: "square brackets", char: "[", input: "hello", want: "[hello]"},
		{name: "parens", char: "(", input: "hello", want: "(hello)"},
		{name: "curly braces", char: "{", input: "hello", want: "{hello}"},
		{name: "angle brackets", char: "<", input: "hello", want: "<hello>"},
		{name: "empty line", char: "", input: "", want: `""`},
		{name: "line with spaces", char: "'", input: "hello world", want: "'hello world'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewQuoteRule(tt.char)
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
