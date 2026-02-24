package rule

import "testing"

func TestDeleteLineRule_RemovesMatchingLines(t *testing.T) {
	rule, err := NewDeleteLineRule("foo")
	if err != nil {
		t.Fatalf("failed to create rule: %v", err)
	}

	result, err := rule.Apply("foo bar", &LineContext{LineNum: 1})
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if len(result) != 0 {
		t.Fatalf("expected 0 lines (deleted), got %d", len(result))
	}
}

func TestDeleteLineRule_KeepsNonMatchingLines(t *testing.T) {
	rule, err := NewDeleteLineRule("foo")
	if err != nil {
		t.Fatalf("failed to create rule: %v", err)
	}

	result, err := rule.Apply("bar baz", &LineContext{LineNum: 1})
	if err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 line, got %d", len(result))
	}
	if result[0] != "bar baz" {
		t.Errorf("got %q, want %q", result[0], "bar baz")
	}
}

func TestDeleteLineRule_RegexPattern(t *testing.T) {
	rule, err := NewDeleteLineRule(`^\s*#`)
	if err != nil {
		t.Fatalf("failed to create rule: %v", err)
	}

	// Should delete (matches comment)
	result, _ := rule.Apply("  # comment", &LineContext{LineNum: 1})
	if len(result) != 0 {
		t.Errorf("expected comment line to be deleted")
	}

	// Should keep (not a comment)
	result, _ = rule.Apply("code here", &LineContext{LineNum: 1})
	if len(result) != 1 {
		t.Errorf("expected non-comment line to be kept")
	}
}

func TestDeleteLineRule_InvalidRegex(t *testing.T) {
	_, err := NewDeleteLineRule("[invalid")
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}
