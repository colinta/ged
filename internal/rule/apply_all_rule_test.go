package rule

import (
	"testing"
)

func TestApplyAllRule_AppliesSubstitution(t *testing.T) {
	sub, _ := NewSubstitutionRule("foo", "bar")
	r := NewApplyAllRule([]LineRule{sub})

	result, err := r.ApplyDocument([]string{"foo baz", "hello foo"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"bar baz", "hello bar"}
	if len(result) != len(want) {
		t.Fatalf("got %v, want %v", result, want)
	}
	for i, v := range result {
		if v != want[i] {
			t.Errorf("index %d: got %q, want %q", i, v, want[i])
		}
	}
}

func TestApplyAllRule_FiltersLines(t *testing.T) {
	print, _ := NewPrintLineRule("keep")
	r := NewApplyAllRule([]LineRule{print})

	result, err := r.ApplyDocument([]string{"keep this", "drop this", "keep that"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"keep this", "keep that"}
	if len(result) != len(want) {
		t.Fatalf("got %v, want %v", result, want)
	}
	for i, v := range result {
		if v != want[i] {
			t.Errorf("index %d: got %q, want %q", i, v, want[i])
		}
	}
}

func TestApplyAllRule_ChainsRules(t *testing.T) {
	print, _ := NewPrintLineRule("hello")
	sub, _ := NewSubstitutionRule("hello", "HI", WithGlobal())
	r := NewApplyAllRule([]LineRule{print, sub})

	result, err := r.ApplyDocument([]string{"hello world", "goodbye", "hello hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"HI world", "HI HI"}
	if len(result) != len(want) {
		t.Fatalf("got %v, want %v", result, want)
	}
	for i, v := range result {
		if v != want[i] {
			t.Errorf("index %d: got %q, want %q", i, v, want[i])
		}
	}
}

func TestApplyAllRule_PreservesLineNumbering(t *testing.T) {
	// Line number rule should use original document line numbers
	lineNumRule := NewPrintLineNumRule(SingleLine(2))
	r := NewApplyAllRule([]LineRule{lineNumRule})

	result, err := r.ApplyDocument([]string{"line1", "line2", "line3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"line2"}
	if len(result) != len(want) {
		t.Fatalf("got %v, want %v", result, want)
	}
	if result[0] != "line2" {
		t.Errorf("got %q, want %q", result[0], "line2")
	}
}

func TestApplyAllRule_EmptyInput(t *testing.T) {
	sub, _ := NewSubstitutionRule("foo", "bar")
	r := NewApplyAllRule([]LineRule{sub})

	result, err := r.ApplyDocument([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("got %v, want empty", result)
	}
}
