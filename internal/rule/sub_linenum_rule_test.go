package rule

import "testing"

func TestSubLineNumRule(t *testing.T) {
	tests := []struct {
		name        string
		lineRange   LineRange
		replacement string
		line        string
		lineNum     int
		want        string
	}{
		{"matching line replaced", SingleLine(2), "replaced", "original", 2, "replaced"},
		{"non-matching line kept", SingleLine(2), "replaced", "original", 3, "original"},
		{"range replaces all in range", &Range{Start: 2, End: 4}, "new", "old", 3, "new"},
		{"range keeps lines outside", &Range{Start: 2, End: 4}, "new", "old", 5, "old"},
		{"empty replacement", SingleLine(1), "", "hello", 1, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSubLineNumRule(tt.lineRange, tt.replacement)
			result, err := r.Apply(tt.line, tt.lineNum)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != 1 {
				t.Fatalf("expected 1 line, got %d", len(result))
			}
			if result[0] != tt.want {
				t.Errorf("got %q, want %q", result[0], tt.want)
			}
		})
	}
}

func TestSubLineNumRule_NewlineInReplacement(t *testing.T) {
	r := NewSubLineNumRule(SingleLine(2), "hello\nworld")
	result, err := r.Apply("original", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 lines, got %d: %q", len(result), result)
	}
	if result[0] != "hello" {
		t.Errorf("line 0: got %q, want %q", result[0], "hello")
	}
	if result[1] != "world" {
		t.Errorf("line 1: got %q, want %q", result[1], "world")
	}
}
