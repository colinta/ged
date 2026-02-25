package parser

import (
	"testing"

	"github.com/colinta/ged/internal/rule"
)

func TestParseOn(t *testing.T) {
	parsed, err := ParseRule("on/start/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := parsed.(*rule.OnRule); !ok {
		t.Errorf("expected *OnRule, got %T", parsed)
	}
}

func TestParseOff(t *testing.T) {
	parsed, err := ParseRule("off/end/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := parsed.(*rule.OffRule); !ok {
		t.Errorf("expected *OffRule, got %T", parsed)
	}
}

func TestParseAfter(t *testing.T) {
	parsed, err := ParseRule("after/marker/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := parsed.(*rule.AfterRule); !ok {
		t.Errorf("expected *AfterRule, got %T", parsed)
	}
}

func TestParseToggle(t *testing.T) {
	parsed, err := ParseRule("toggle/---/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := parsed.(*rule.ToggleRule); !ok {
		t.Errorf("expected *ToggleRule, got %T", parsed)
	}
}

func TestParseControlLiteralDelimiter(t *testing.T) {
	parsed, err := ParseRule("on`hello`")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := parsed.(*rule.OnRule); !ok {
		t.Errorf("expected *OnRule, got %T", parsed)
	}
}

func TestParseControlMissingPattern(t *testing.T) {
	_, err := ParseRule("on")
	if err == nil {
		t.Error("expected error for missing pattern")
	}
}

func TestParseControlEmptyPattern(t *testing.T) {
	_, err := ParseRule("on//")
	if err == nil {
		t.Error("expected error for empty pattern")
	}
}
