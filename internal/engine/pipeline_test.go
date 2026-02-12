package engine

import (
	"testing"

	"github.com/colinta/ged/internal/rule"
)

func TestPipeline_SingleRule(t *testing.T) {
	sub, _ := rule.NewSubstitutionRule("foo", "bar")
	p := NewPipeline(sub)

	result, err := p.Process("hello foo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 || result[0] != "hello bar" {
		t.Errorf("got %v, want [\"hello bar\"]", result)
	}
}

func TestPipeline_TwoRulesChain(t *testing.T) {
	// First substitute, then substitute again
	sub1, _ := rule.NewSubstitutionRule("a", "b")
	sub2, _ := rule.NewSubstitutionRule("b", "c")
	p := NewPipeline(sub1, sub2)

	result, err := p.Process("aaa", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// "aaa" -> "baa" (first 'a' to 'b') -> "caa" (first 'b' to 'c')
	if len(result) != 1 || result[0] != "caa" {
		t.Errorf("got %v, want [\"caa\"]", result)
	}
}

func TestPipeline_FilterThenSubstitute(t *testing.T) {
	// Keep lines with "foo", then replace all "o" with "x"
	print, _ := rule.NewPrintLineRule("foo")
	sub, _ := rule.NewSubstitutionRule("o", "x", rule.WithGlobal())
	p := NewPipeline(print, sub)

	// Line matches filter, gets transformed
	result, err := p.Process("foo bar", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0] != "fxx bar" {
		t.Errorf("got %v, want [\"fxx bar\"]", result)
	}
}

func TestPipeline_FilterStopsChain(t *testing.T) {
	// Keep lines with "foo", then substitute
	print, _ := rule.NewPrintLineRule("foo")
	sub, _ := rule.NewSubstitutionRule("bar", "baz")
	p := NewPipeline(print, sub)

	// Line doesn't match filter - should produce no output
	result, err := p.Process("bar only", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("got %v, want []", result)
	}
}

func TestPipeline_SubstituteThenFilter(t *testing.T) {
	// Replace "hello" with "foo", then keep only lines with "foo"
	sub, _ := rule.NewSubstitutionRule("hello", "foo")
	print, _ := rule.NewPrintLineRule("foo")
	p := NewPipeline(sub, print)

	// "hello world" becomes "foo world", which matches filter
	result, err := p.Process("hello world", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0] != "foo world" {
		t.Errorf("got %v, want [\"foo world\"]", result)
	}

	// "goodbye world" stays "goodbye world", doesn't match filter
	result, err = p.Process("goodbye world", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("got %v, want []", result)
	}
}

func TestPipeline_EmptyPipeline(t *testing.T) {
	p := NewPipeline()

	result, err := p.Process("hello", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// No rules means line passes through unchanged
	if len(result) != 1 || result[0] != "hello" {
		t.Errorf("got %v, want [\"hello\"]", result)
	}
}

func TestPipeline_DeleteRule(t *testing.T) {
	// Delete lines containing "secret"
	del, _ := rule.NewDeleteLineRule("secret")
	sub, _ := rule.NewSubstitutionRule("public", "PUBLIC")
	p := NewPipeline(del, sub)

	// Non-matching line passes through and gets substituted
	result, _ := p.Process("public info", 1)
	if len(result) != 1 || result[0] != "PUBLIC info" {
		t.Errorf("got %v, want [\"PUBLIC info\"]", result)
	}

	// Matching line is deleted
	result, _ = p.Process("secret data", 2)
	if len(result) != 0 {
		t.Errorf("got %v, want []", result)
	}
}
