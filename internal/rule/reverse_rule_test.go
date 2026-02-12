package rule

import (
	"testing"
)

func TestReverseRule_ReversesOrder(t *testing.T) {
	r := NewReverseRule()
	result, err := r.ApplyDocument([]string{"a", "b", "c"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"c", "b", "a"}
	if len(result) != len(want) {
		t.Fatalf("got %v, want %v", result, want)
	}
	for i, v := range result {
		if v != want[i] {
			t.Errorf("index %d: got %q, want %q", i, v, want[i])
		}
	}
}

func TestReverseRule_EmptyInput(t *testing.T) {
	r := NewReverseRule()
	result, err := r.ApplyDocument([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("got %v, want empty", result)
	}
}

func TestReverseRule_SingleLine(t *testing.T) {
	r := NewReverseRule()
	result, err := r.ApplyDocument([]string{"only"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 || result[0] != "only" {
		t.Errorf("got %v, want [\"only\"]", result)
	}
}

func TestReverseRule_DoesNotMutateInput(t *testing.T) {
	r := NewReverseRule()
	input := []string{"a", "b", "c"}
	_, _ = r.ApplyDocument(input)

	if input[0] != "a" || input[1] != "b" || input[2] != "c" {
		t.Errorf("input was mutated: %v", input)
	}
}
