package parser

import (
	"testing"

	"github.com/colinta/ged/internal/rule"
)

func TestParseBetween(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"basic between", "between/start/end/", false},
		{"different delimiter", "between|start|end|", false},
		{"inverted", "!between/start/end/", false},
		{"literal delimiters", "between`start`end`", false},
		{"missing end pattern", "between/start/", true},
		{"missing patterns", "between", true},
		{"empty start", "between//end/", true},
		{"empty end", "between/start//", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseRule(tt.input)
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestParseArgsBetweenLineRule(t *testing.T) {
	rules, err := ParseArgs([]string{"between/START/END/", "{", "s/x/X/g", "}"})
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

func TestParseArgsBetweenDocRule(t *testing.T) {
	rules, err := ParseArgs([]string{"between/START/END/", "{", "sort", "}"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if _, ok := rules[0].(*rule.BetweenDocRule); !ok {
		t.Errorf("expected *rule.BetweenDocRule, got %T", rules[0])
	}
}

func TestParseArgsBetweenInverted(t *testing.T) {
	rules, err := ParseArgs([]string{"!between/START/END/", "{", "s/x/X/", "}"})
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

func TestParseArgsBetweenMissingBrace(t *testing.T) {
	_, err := ParseArgs([]string{"between/START/END/", "s/x/X/"})
	if err == nil {
		t.Error("expected error for missing '{', got nil")
	}
}
