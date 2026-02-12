package rule

import (
	"testing"
)

func TestJoinRule_JoinsWithComma(t *testing.T) {
	r := NewJoinRule(",")
	result, err := r.ApplyDocument([]string{"a", "b", "c"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 || result[0] != "a,b,c" {
		t.Errorf("got %v, want [\"a,b,c\"]", result)
	}
}

func TestJoinRule_JoinsWithSpace(t *testing.T) {
	r := NewJoinRule(" ")
	result, err := r.ApplyDocument([]string{"hello", "world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 || result[0] != "hello world" {
		t.Errorf("got %v, want [\"hello world\"]", result)
	}
}

func TestJoinRule_JoinsWithEmptySeparator(t *testing.T) {
	r := NewJoinRule("")
	result, err := r.ApplyDocument([]string{"a", "b", "c"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 || result[0] != "abc" {
		t.Errorf("got %v, want [\"abc\"]", result)
	}
}

func TestJoinRule_EmptyInput(t *testing.T) {
	r := NewJoinRule(",")
	result, err := r.ApplyDocument([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("got %v, want empty", result)
	}
}

func TestJoinRule_SingleLine(t *testing.T) {
	r := NewJoinRule(",")
	result, err := r.ApplyDocument([]string{"only"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 || result[0] != "only" {
		t.Errorf("got %v, want [\"only\"]", result)
	}
}
