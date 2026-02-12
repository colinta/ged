package rule

import (
	"testing"
)

func TestSortRule_SortsAlphabetically(t *testing.T) {
	r := NewSortRule()
	result, err := r.ApplyDocument([]string{"c", "a", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"a", "b", "c"}
	if len(result) != len(want) {
		t.Fatalf("got %v, want %v", result, want)
	}
	for i, v := range result {
		if v != want[i] {
			t.Errorf("index %d: got %q, want %q", i, v, want[i])
		}
	}
}

func TestSortRule_EmptyInput(t *testing.T) {
	r := NewSortRule()
	result, err := r.ApplyDocument([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("got %v, want empty", result)
	}
}

func TestSortRule_SingleLine(t *testing.T) {
	r := NewSortRule()
	result, err := r.ApplyDocument([]string{"only"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 || result[0] != "only" {
		t.Errorf("got %v, want [\"only\"]", result)
	}
}

func TestSortRule_DoesNotMutateInput(t *testing.T) {
	r := NewSortRule()
	input := []string{"c", "a", "b"}
	_, _ = r.ApplyDocument(input)

	// Original should be unchanged
	if input[0] != "c" || input[1] != "a" || input[2] != "b" {
		t.Errorf("input was mutated: %v", input)
	}
}
