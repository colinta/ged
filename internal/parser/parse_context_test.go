package parser

import (
	"testing"

	"github.com/colinta/ged/internal/rule"
)

func TestParsePrintContext(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantErr    bool
		hasContext bool
	}{
		{"plain print", "p/error/", false, false},
		{"print with context", "p/error/context=2", false, true},
		{"print flags then context", "p/error/i/context=2", false, true},
		{"print before only", "p/error/before=1", false, true},
		{"print after only", "p/error/after=3", false, true},
		{"print before and after", "p/error/before=1/after=3", false, true},
		{"print context=0", "p/error/context=0", false, false},
		{"print flags and context", "p/error/i/before=2", false, true},
		{"print invalid value", "p/error/context=abc", true, false},
		{"print negative value", "p/error/context=-1", true, false},
		{"print unknown option", "p/error/foo=1", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRule(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.input, err)
			}
			pr, ok := got.(*rule.PrintLineRule)
			if !ok {
				t.Fatalf("expected *rule.PrintLineRule, got %T", got)
			}
			if pr.HasContext() != tt.hasContext {
				t.Errorf("HasContext() = %v, want %v", pr.HasContext(), tt.hasContext)
			}
		})
	}
}

func TestParseDeleteContext(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantErr    bool
		hasContext bool
	}{
		{"plain delete", "d/debug/", false, false},
		{"delete with context", "d/debug/context=1", false, true},
		{"delete after only", "d/debug/after=2", false, true},
		{"delete before only", "d/debug/before=1", false, true},
		{"delete flags and context", "d/debug/i/context=1", false, true},
		{"delete invalid value", "d/debug/context=x", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRule(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.input, err)
			}
			dr, ok := got.(*rule.DeleteLineRule)
			if !ok {
				t.Fatalf("expected *rule.DeleteLineRule, got %T", got)
			}
			if dr.HasContext() != tt.hasContext {
				t.Errorf("HasContext() = %v, want %v", dr.HasContext(), tt.hasContext)
			}
		})
	}
}
