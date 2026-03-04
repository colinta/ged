package rule

import "testing"

func TestXargsRule(t *testing.T) {
	tests := []struct {
		name    string
		command string
		input   string
		want    []string
	}{
		{
			name:    "echo with argument",
			command: "echo hi",
			input:   "world",
			want:    []string{"hi world"},
		},
		{
			name:    "transforms each line",
			command: "echo prefix:",
			input:   "hello",
			want:    []string{"prefix: hello"},
		},
		{
			name:    "command that ignores argument",
			command: "echo constant",
			input:   "anything",
			want:    []string{"constant anything"},
		},
		{
			name:    "command producing multiple lines",
			command: "echo first; echo second",
			input:   "hello",
			want:    []string{"first", "second hello"},
		},
		{
			name:    "command producing empty output",
			command: "true",
			input:   "hello",
			want:    []string{},
		},
		{
			name:    "failing command passes line through",
			command: "false",
			input:   "hello",
			want:    []string{"hello"},
		},
		{
			name:    "handles special characters in line",
			command: "echo",
			input:   "hello 'world' $(date)",
			want:    []string{"hello 'world' $(date)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewXargsRule(tt.command)
			ctx := &LineContext{LineNum: 1}
			got, err := r.Apply(tt.input, ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d lines, want %d: %v", len(got), len(tt.want), got)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("line %d: got %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestShellQuote(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "'hello'"},
		{"hello world", "'hello world'"},
		{"it's", "'it'\\''s'"},
		{"$(date)", "'$(date)'"},
		{"`ls`", "'`ls`'"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := shellQuote(tt.input)
			if got != tt.want {
				t.Errorf("shellQuote(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
