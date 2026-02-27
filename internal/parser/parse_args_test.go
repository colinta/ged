package parser

import (
	"testing"

	"github.com/colinta/ged/internal/rule"
)

func TestParseIf_Basic(t *testing.T) {
	result, err := ParseRule("if/hello/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cond, ok := result.(*condition)
	if !ok {
		t.Fatalf("expected *condition, got %T", result)
	}
	if cond.inverted {
		t.Error("expected inverted=false")
	}
	matched, matchErr := cond.pattern.MatchString("hello world")
	if matchErr != nil {
		t.Fatalf("match error: %v", matchErr)
	}
	if !matched {
		t.Error("pattern should match 'hello world'")
	}
}

func TestParseIf_Inverted(t *testing.T) {
	result, err := ParseRule("!if/hello/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cond, ok := result.(*condition)
	if !ok {
		t.Fatalf("expected *condition, got %T", result)
	}
	if !cond.inverted {
		t.Error("expected inverted=true")
	}
}

func TestParseIf_LiteralDelimiter(t *testing.T) {
	result, err := ParseRule("if`foo.bar`")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cond, ok := result.(*condition)
	if !ok {
		t.Fatalf("expected *condition, got %T", result)
	}
	// Literal: "foo.bar" should NOT match "fooXbar"
	matched1, err1 := cond.pattern.MatchString("fooXbar")
	if err1 != nil {
		t.Fatalf("match error: %v", err1)
	}
	if matched1 {
		t.Error("literal pattern should not match 'fooXbar'")
	}
	matched2, err2 := cond.pattern.MatchString("foo.bar")
	if err2 != nil {
		t.Fatalf("match error: %v", err2)
	}
	if !matched2 {
		t.Error("literal pattern should match 'foo.bar'")
	}
}

func TestParseIf_MissingPattern(t *testing.T) {
	_, err := ParseRule("if//")
	if err == nil {
		t.Error("expected error for empty pattern")
	}
}

func TestParseIf_NoDelimiter(t *testing.T) {
	_, err := ParseRule("if")
	if err == nil {
		t.Error("expected error for bare 'if'")
	}
}

func TestParseArgs_SimpleRules(t *testing.T) {
	// Without conditionals, ParseArgs works like calling ParseRule on each arg
	results, err := ParseArgs([]string{"s/a/b/", "d/x/"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(results))
	}
}

func TestParseArgs_ConditionalBlock(t *testing.T) {
	results, err := ParseArgs([]string{"if/hello/", "{", "s/o/x/", "}"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(results))
	}
	_, ok := results[0].(rule.LineRule)
	if !ok {
		t.Errorf("expected LineRule, got %T", results[0])
	}
}

func TestParseArgs_ConditionalWithMultipleInnerRules(t *testing.T) {
	results, err := ParseArgs([]string{"if/x/", "{", "s/a/b/", "s/c/d/", "}"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(results))
	}
}

func TestParseArgs_NestedConditionals(t *testing.T) {
	results, err := ParseArgs([]string{
		"if/foo/", "{", "if/bar/", "{", "s/x/y/", "}", "}",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(results))
	}
}

func TestParseArgs_MixedRules(t *testing.T) {
	// Line rules before and after a conditional
	results, err := ParseArgs([]string{"s/a/b/", "if/x/", "{", "s/c/d/", "}", "s/e/f/"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(results))
	}
}

func TestParseArgs_MissingOpenBrace(t *testing.T) {
	_, err := ParseArgs([]string{"if/hello/", "s/a/b/"})
	if err == nil {
		t.Error("expected error for missing '{'")
	}
}

func TestParseArgs_MissingCloseBrace(t *testing.T) {
	_, err := ParseArgs([]string{"if/hello/", "{", "s/a/b/"})
	if err == nil {
		t.Error("expected error for missing '}'")
	}
}

func TestParseArgs_UnexpectedCloseBrace(t *testing.T) {
	_, err := ParseArgs([]string{"s/a/b/", "}"})
	if err == nil {
		t.Error("expected error for unexpected '}'")
	}
}

func TestParseArgs_UnexpectedOpenBrace(t *testing.T) {
	_, err := ParseArgs([]string{"{", "s/a/b/", "}"})
	if err == nil {
		t.Error("expected error for unexpected '{'")
	}
}

func TestParseArgs_DocumentRuleInsideBlock_CreatesDocRule(t *testing.T) {
	// A document rule inside a block creates a ConditionalDocRule
	results, err := ParseArgs([]string{"if/x/", "{", "sort", "}"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(results))
	}
	_, ok := results[0].(rule.DocumentRule)
	if !ok {
		t.Errorf("expected DocumentRule, got %T", results[0])
	}
}

func TestParseArgs_AllLineRulesInsideBlock_CreatesLineRule(t *testing.T) {
	results, err := ParseArgs([]string{"if/x/", "{", "s/a/b/", "}"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(results))
	}
	_, ok := results[0].(rule.LineRule)
	if !ok {
		t.Errorf("expected LineRule, got %T", results[0])
	}
}

func TestParseArgs_MixedRulesInsideBlock_CreatesDocRule(t *testing.T) {
	// Line rule + document rule inside block → ConditionalDocRule
	results, err := ParseArgs([]string{"if/x/", "{", "s/a/b/", "sort", "}"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(results))
	}
	_, ok := results[0].(rule.DocumentRule)
	if !ok {
		t.Errorf("expected DocumentRule, got %T", results[0])
	}
}
