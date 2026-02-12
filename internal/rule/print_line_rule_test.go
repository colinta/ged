package rule

import "testing"

func TestPrintLineRule_KeepsMatchingLines(t *testing.T) {
	rule, err := NewPrintLineRule("foo")
	if err != nil {
		t.Fatalf("failed to create rule: %v", err)
	}

	result, err := rule.Apply("foo bar", 1)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 line, got %d", len(result))
	}
	if result[0] != "foo bar" {
		t.Errorf("got %q, want %q", result[0], "foo bar")
	}
}

func TestPrintLineRule_RemovesNonMatchingLines(t *testing.T) {
	rule, err := NewPrintLineRule("foo")
	if err != nil {
		t.Fatalf("failed to create rule: %v", err)
	}

	result, err := rule.Apply("bar baz", 1)
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if len(result) != 0 {
		t.Fatalf("expected 0 lines (deleted), got %d", len(result))
	}
}

func TestPrintLineRule_RegexPattern(t *testing.T) {
	rule, err := NewPrintLineRule(`^\d+$`)
	if err != nil {
		t.Fatalf("failed to create rule: %v", err)
	}

	// Should match
	result, _ := rule.Apply("123", 1)
	if len(result) != 1 {
		t.Errorf("expected '123' to match, got deleted")
	}

	// Should not match
	result, _ = rule.Apply("abc", 1)
	if len(result) != 0 {
		t.Errorf("expected 'abc' to be deleted, got kept")
	}
}

func TestPrintLineRule_InvalidRegex(t *testing.T) {
	_, err := NewPrintLineRule("[invalid")
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}
