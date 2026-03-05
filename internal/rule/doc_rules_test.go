package rule

import "testing"

func TestLinesRule(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{"numbers lines", []string{"a", "b", "c"}, []string{"1: a", "2: b", "3: c"}},
		{"pads to width", []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
			[]string{" 1: a", " 2: b", " 3: c", " 4: d", " 5: e", " 6: f", " 7: g", " 8: h", " 9: i", "10: j"}},
		{"single line", []string{"hello"}, []string{"1: hello"}},
		{"empty", []string{}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewLinesRule()
			got, err := r.ApplyDocument(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertLines(t, got, tt.want)
		})
	}
}

func TestBeginRule(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		input []string
		want  []string
	}{
		{"prepend single line", "header", []string{"a", "b"}, []string{"header", "a", "b"}},
		{"prepend multi-line", "h1\nh2", []string{"a"}, []string{"h1", "h2", "a"}},
		{"prepend to empty", "header", []string{}, []string{"header"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewBeginRule(tt.text)
			got, err := r.ApplyDocument(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertLines(t, got, tt.want)
		})
	}
}

func TestEndRule(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		input []string
		want  []string
	}{
		{"append single line", "footer", []string{"a", "b"}, []string{"a", "b", "footer"}},
		{"append multi-line", "f1\nf2", []string{"a"}, []string{"a", "f1", "f2"}},
		{"append to empty", "footer", []string{}, []string{"footer"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewEndRule(tt.text)
			got, err := r.ApplyDocument(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertLines(t, got, tt.want)
		})
	}
}

func TestBorderRule(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		input []string
		want  []string
	}{
		{"single border", "---", []string{"a", "b"}, []string{"---", "a", "b", "---"}},
		{"multi-line border", "=\n=", []string{"a"}, []string{"=", "=", "a", "=", "="}},
		{"empty input", "---", []string{}, []string{"---", "---"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewBorderRule(tt.text)
			got, err := r.ApplyDocument(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertLines(t, got, tt.want)
		})
	}
}

func TestCountRule(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{"three lines", []string{"a", "b", "c"}, []string{"3"}},
		{"one line", []string{"hello"}, []string{"1"}},
		{"empty", []string{}, []string{"0"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewCountRule()
			got, err := r.ApplyDocument(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertLines(t, got, tt.want)
		})
	}
}

func TestUniqRule(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{"removes consecutive dupes", []string{"a", "a", "b", "b", "a"}, []string{"a", "b", "a"}},
		{"no dupes", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"all same", []string{"x", "x", "x"}, []string{"x"}},
		{"single line", []string{"a"}, []string{"a"}},
		{"empty", []string{}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewUniqRule()
			got, err := r.ApplyDocument(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertLines(t, got, tt.want)
		})
	}
}

// assertLines is a test helper for comparing string slices.
func assertLines(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d:\ngot:  %v\nwant: %v", len(got), len(want), got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("line[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
