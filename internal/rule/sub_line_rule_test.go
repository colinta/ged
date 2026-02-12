package rule

import "testing"

func TestSubstitutionRule_ReplacesFirstMatch(t *testing.T) {
	rule, err := NewSubstitutionRule("world", "earth")
	if err != nil {
		t.Fatalf("failed to create rule: %v", err)
	}

	result, err := rule.Apply("hello world world", 1)
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

	result, err := rule.Apply("hello world", 1)
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

	result, err := rule.Apply("hello world", 1)
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

	result, err := rule.Apply("foo 123 bar 456", 1)
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

	result, err := rule.Apply("foo 123 bar 456", 1)
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

	result, err := rule.Apply("foo 123 bar 456", 1)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	want := "foo 1231 bar 456"
	if result[0] != want {
		t.Errorf("got %q, want %q", result[0], want)
	}
}

func TestSubstitutionRule_LiteralPattern(t *testing.T) {
	// QuoteMeta'd pattern — dots are literal, not regex wildcards
	rule, err := NewSubstitutionRule(`foo\.bar`, "baz")
	if err != nil {
		t.Fatalf("failed to create rule: %v", err)
	}

	// Should match literal "foo.bar"
	result, _ := rule.Apply("foo.bar", 1)
	if result[0] != "baz" {
		t.Errorf("got %q, want %q", result[0], "baz")
	}

	// Should NOT match "fooXbar"
	result, _ = rule.Apply("fooXbar", 1)
	if result[0] != "fooXbar" {
		t.Errorf("got %q, want %q", result[0], "fooXbar")
	}
}

func TestSubstitutionRule_NewlineInReplacement(t *testing.T) {
	rule, err := NewSubstitutionRule("foo", "bar\nbaz")
	if err != nil {
		t.Fatalf("failed to create rule: %v", err)
	}

	result, _ := rule.Apply("foo", 1)
	if len(result) != 2 {
		t.Fatalf("expected 2 lines, got %d: %q", len(result), result)
	}
	if result[0] != "bar" {
		t.Errorf("line 0: got %q, want %q", result[0], "bar")
	}
	if result[1] != "baz" {
		t.Errorf("line 1: got %q, want %q", result[1], "baz")
	}
}

func TestSubstitutionRule_InvalidRegex(t *testing.T) {
	_, err := NewSubstitutionRule("[invalid", "replacement")
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}
