package rule

import "testing"

func TestSubstitutionRule_ReplacesFirstMatch(t *testing.T) {
	rule, err := NewSubstitutionRule("world", "earth")
	if err != nil {
		t.Fatalf("failed to create rule: %v", err)
	}

	result, err := rule.Apply("hello world world")
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 line, got %d", len(result))
	}

	want := "hello earth world" // only first match replaced
	if result[0] != want {
		t.Errorf("got %q, want %q", result[0], want)
	}
}

func TestSubstitutionRule_GlobalReplacesAll(t *testing.T) {
	rule, err := NewSubstitutionRule("o", "0", WithGlobal())
	if err != nil {
		t.Fatalf("failed to create rule: %v", err)
	}

	result, err := rule.Apply("hello world")
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	want := "hell0 w0rld"
	if result[0] != want {
		t.Errorf("got %q, want %q", result[0], want)
	}
}

func TestSubstitutionRule_NoMatch(t *testing.T) {
	rule, err := NewSubstitutionRule("xyz", "abc")
	if err != nil {
		t.Fatalf("failed to create rule: %v", err)
	}

	result, err := rule.Apply("hello world")
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	want := "hello world" // unchanged
	if result[0] != want {
		t.Errorf("got %q, want %q", result[0], want)
	}
}

func TestSubstitutionRule_RegexPattern(t *testing.T) {
	rule, err := NewSubstitutionRule(`\d+`, "NUM", WithGlobal())
	if err != nil {
		t.Fatalf("failed to create rule: %v", err)
	}

	result, err := rule.Apply("foo 123 bar 456")
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	want := "foo NUM bar NUM"
	if result[0] != want {
		t.Errorf("got %q, want %q", result[0], want)
	}
}

func TestSubstitutionRule_RegexReplacePattern_WithGlobal(t *testing.T) {
	rule, err := NewSubstitutionRule(`(\d)(\d+)`, "$1$2$1", WithGlobal())
	if err != nil {
		t.Fatalf("failed to create rule: %v", err)
	}

	result, err := rule.Apply("foo 123 bar 456")
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	want := "foo 1231 bar 4564"
	if result[0] != want {
		t.Errorf("got %q, want %q", result[0], want)
	}
}

func TestSubstitutionRule_RegexReplacePattern_OnlyFirst(t *testing.T) {
	rule, err := NewSubstitutionRule(`(\d)(\d+)`, "$1$2$1")
	if err != nil {
		t.Fatalf("failed to create rule: %v", err)
	}

	result, err := rule.Apply("foo 123 bar 456")
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	want := "foo 1231 bar 456"
	if result[0] != want {
		t.Errorf("got %q, want %q", result[0], want)
	}
}

func TestSubstitutionRule_InvalidRegex(t *testing.T) {
	_, err := NewSubstitutionRule("[invalid", "replacement")
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}
