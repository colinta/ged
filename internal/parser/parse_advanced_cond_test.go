package parser

import (
	"testing"

	"github.com/colinta/ged/internal/rule"
)

func TestParseRule_IfAny(t *testing.T) {
	parsed, err := ParseRule("ifany/error/")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := parsed.(*ifAnyCondition); !ok {
		t.Errorf("expected *ifAnyCondition, got %T", parsed)
	}
}

func TestParseRule_IfAnyInverted(t *testing.T) {
	parsed, err := ParseRule("!ifany/error/")
	if err != nil {
		t.Fatal(err)
	}
	cond, ok := parsed.(*ifAnyCondition)
	if !ok {
		t.Fatalf("expected *ifAnyCondition, got %T", parsed)
	}
	if !cond.inverted {
		t.Error("expected inverted to be true")
	}
}

func TestParseRule_IfNone(t *testing.T) {
	parsed, err := ParseRule("ifnone/error/")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := parsed.(*ifNoneCondition); !ok {
		t.Errorf("expected *ifNoneCondition, got %T", parsed)
	}
}

func TestParseRule_IfNoneInverted(t *testing.T) {
	parsed, err := ParseRule("!ifnone/error/")
	if err != nil {
		t.Fatal(err)
	}
	cond, ok := parsed.(*ifNoneCondition)
	if !ok {
		t.Fatalf("expected *ifNoneCondition, got %T", parsed)
	}
	if !cond.inverted {
		t.Error("expected inverted to be true")
	}
}

func TestParseRule_IfAnyLiteral(t *testing.T) {
	parsed, err := ParseRule("ifany`foo.bar`")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := parsed.(*ifAnyCondition); !ok {
		t.Errorf("expected *ifAnyCondition, got %T", parsed)
	}
}

func TestParseRule_IfAnyMissingPattern(t *testing.T) {
	_, err := ParseRule("ifany//")
	if err == nil {
		t.Error("expected error for empty pattern")
	}
}

func TestParseRule_IfNoneMissingPattern(t *testing.T) {
	_, err := ParseRule("ifnone//")
	if err == nil {
		t.Error("expected error for empty pattern")
	}
}

func TestParseArgs_IfAnyBlock(t *testing.T) {
	rules, err := ParseArgs([]string{"ifany/error/", "{", "upper", "}"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if _, ok := rules[0].(*rule.IfAnyDocRule); !ok {
		t.Errorf("expected *rule.IfAnyDocRule, got %T", rules[0])
	}
}

func TestParseArgs_IfNoneBlock(t *testing.T) {
	rules, err := ParseArgs([]string{"ifnone/error/", "{", "upper", "}"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if _, ok := rules[0].(*rule.IfNoneDocRule); !ok {
		t.Errorf("expected *rule.IfNoneDocRule, got %T", rules[0])
	}
}

func TestParseArgs_IfElseLineRules(t *testing.T) {
	rules, err := ParseArgs([]string{"if/hello/", "{", "upper", "}", "else", "{", "lower", "}"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	// Both branches are line rules, so it should be ConditionalLineRule
	if _, ok := rules[0].(*rule.ConditionalLineRule); !ok {
		t.Errorf("expected *rule.ConditionalLineRule, got %T", rules[0])
	}
}

func TestParseArgs_IfElseDocRules(t *testing.T) {
	rules, err := ParseArgs([]string{"if/hello/", "{", "sort", "}", "else", "{", "reverse", "}"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if _, ok := rules[0].(*rule.ConditionalDocRule); !ok {
		t.Errorf("expected *rule.ConditionalDocRule, got %T", rules[0])
	}
}

func TestParseArgs_IfAnyElse(t *testing.T) {
	rules, err := ParseArgs([]string{"ifany/error/", "{", "upper", "}", "else", "{", "lower", "}"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if _, ok := rules[0].(*rule.IfAnyDocRule); !ok {
		t.Errorf("expected *rule.IfAnyDocRule, got %T", rules[0])
	}
}

func TestParseArgs_BetweenElse(t *testing.T) {
	rules, err := ParseArgs([]string{"between/START/END/", "{", "upper", "}", "else", "{", "lower", "}"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if _, ok := rules[0].(*rule.BetweenLineRule); !ok {
		t.Errorf("expected *rule.BetweenLineRule, got %T", rules[0])
	}
}

func TestParseArgs_ElseWithoutIf(t *testing.T) {
	_, err := ParseArgs([]string{"upper", "else", "{", "lower", "}"})
	if err == nil {
		t.Error("expected error for else without condition")
	}
}

func TestParseArgs_ElseMissingBlock(t *testing.T) {
	_, err := ParseArgs([]string{"if/hello/", "{", "upper", "}", "else"})
	if err == nil {
		t.Error("expected error for else without block")
	}
}
