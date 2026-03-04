package rule

import (
	"testing"
)

func TestTrimRule(t *testing.T) {
	tests := []struct {
		name  string
		rule  *TrimRule
		input string
		want  string
	}{
		{"both ends", NewTrimRule(), "  hello  ", "hello"},
		{"tabs", NewTrimRule(), "\thello\t", "hello"},
		{"mixed whitespace", NewTrimRule(), "  \thello\t  ", "hello"},
		{"no whitespace", NewTrimRule(), "hello", "hello"},
		{"empty string", NewTrimRule(), "", ""},
		{"only whitespace", NewTrimRule(), "   ", ""},
		{"internal whitespace preserved", NewTrimRule(), "  hello world  ", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.rule.Apply(tt.input, &LineContext{LineNum: 1})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 || result[0] != tt.want {
				t.Errorf("got %q, want %q", result, tt.want)
			}
		})
	}
}

func TestTrimLeftRule(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"removes leading", "  hello  ", "hello  "},
		{"keeps trailing", "hello  ", "hello  "},
		{"tabs", "\thello", "hello"},
	}

	rule := NewTrimLeftRule()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := rule.Apply(tt.input, &LineContext{LineNum: 1})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 || result[0] != tt.want {
				t.Errorf("got %q, want %q", result, tt.want)
			}
		})
	}
}

func TestTrimRightRule(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"removes trailing", "  hello  ", "  hello"},
		{"keeps leading", "  hello", "  hello"},
		{"tabs", "hello\t", "hello"},
	}

	rule := NewTrimRightRule()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := rule.Apply(tt.input, &LineContext{LineNum: 1})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 || result[0] != tt.want {
				t.Errorf("got %q, want %q", result, tt.want)
			}
		})
	}
}

func TestUpperRule(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "HELLO"},
		{"Hello World", "HELLO WORLD"},
		{"ALREADY", "ALREADY"},
		{"123", "123"},
		{"", ""},
	}

	rule := NewUpperRule()
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := rule.Apply(tt.input, &LineContext{LineNum: 1})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 || result[0] != tt.want {
				t.Errorf("got %q, want %q", result, tt.want)
			}
		})
	}
}

func TestLowerRule(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"HELLO", "hello"},
		{"Hello World", "hello world"},
		{"already", "already"},
		{"123", "123"},
		{"", ""},
	}

	rule := NewLowerRule()
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := rule.Apply(tt.input, &LineContext{LineNum: 1})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 || result[0] != tt.want {
				t.Errorf("got %q, want %q", result, tt.want)
			}
		})
	}
}

func TestPrependRule(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		input string
		want  string
	}{
		{"prefix", ">> ", "hello", ">> hello"},
		{"number", "1. ", "item", "1. item"},
		{"empty prefix", "", "hello", "hello"},
		{"empty line", ">> ", "", ">> "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewPrependRule(tt.text)
			result, err := r.Apply(tt.input, &LineContext{LineNum: 1})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 || result[0] != tt.want {
				t.Errorf("got %q, want %q", result, tt.want)
			}
		})
	}
}

func TestAppendRule(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		input string
		want  string
	}{
		{"suffix", ";", "hello", "hello;"},
		{"comma", ",", "item", "item,"},
		{"empty suffix", "", "hello", "hello"},
		{"empty line", ";", "", ";"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewAppendRule(tt.text)
			result, err := r.Apply(tt.input, &LineContext{LineNum: 1})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 || result[0] != tt.want {
				t.Errorf("got %q, want %q", result, tt.want)
			}
		})
	}
}

func TestSurroundRule(t *testing.T) {
	tests := []struct {
		name   string
		before string
		after  string
		input  string
		want   string
	}{
		{"parens", "(", ")", "hello", "(hello)"},
		{"quotes", "\"", "\"", "hello", "\"hello\""},
		{"tags", "<b>", "</b>", "text", "<b>text</b>"},
		{"empty line", "[", "]", "", "[]"},
		{"markdown bold", "**", "**", "word", "**word**"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSurroundRule(tt.before, tt.after)
			result, err := r.Apply(tt.input, &LineContext{LineNum: 1})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 || result[0] != tt.want {
				t.Errorf("got %q, want %q", result, tt.want)
			}
		})
	}
}
