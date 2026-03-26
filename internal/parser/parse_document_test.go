package parser

import (
	"testing"

	"github.com/colinta/ged/internal/rule"
)

func TestParseRule_Sort(t *testing.T) {
	r, err := ParseRule("sort")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := r.(*rule.SortRule)
	if !ok {
		t.Fatalf("expected *SortRule, got %T", r)
	}
}

func TestParseRule_Reverse(t *testing.T) {
	r, err := ParseRule("reverse")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := r.(*rule.ReverseRule)
	if !ok {
		t.Fatalf("expected *ReverseRule, got %T", r)
	}
}

func TestParseRule_JoinBare(t *testing.T) {
	r, err := ParseRule("join")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	joinRule, ok := r.(*rule.JoinRule)
	if !ok {
		t.Fatalf("expected *JoinRule, got %T", r)
	}

	got, err := joinRule.ApplyDocument([]string{"a", "b", "c"})
	if err != nil {
		t.Fatalf("unexpected error applying join rule: %v", err)
	}
	if len(got) != 1 || got[0] != "a b c" {
		t.Fatalf("got %v, want [\"a b c\"]", got)
	}
}

func TestParseRule_JoinWithSeparator(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"comma", "join/,/"},
		{"space", "join/ /"},
		{"pipe delimiter", "join|,|"},
		{"multi-char separator", "join/, /"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := ParseRule(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			_, ok := r.(*rule.JoinRule)
			if !ok {
				t.Fatalf("expected *JoinRule, got %T", r)
			}
		})
	}
}

func TestParseRule_SortDoesNotMatchSubstitution(t *testing.T) {
	// "sort" should parse as SortRule, not as substitution with "o" delimiter
	r, err := ParseRule("sort")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := r.(*rule.SortRule)
	if !ok {
		t.Fatalf("expected *SortRule, got %T (sort should not be parsed as substitution)", r)
	}
}
