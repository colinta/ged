package parser

import (
	"testing"

	"github.com/colinta/ged/internal/rule"
)

func TestParseRule_Print(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantPattern string
		wantErr     bool
	}{
		// Basic print
		{"basic", "p/foo/", "foo", false},
		{"without trailing delimiter", "p/foo", "foo", false},
		// Different delimiters
		{"pipe delimiter", "p|foo|", "foo", false},
		{"equals delimiter", "p=foo=", "foo", false},
		{"hash delimiter", "p#foo#", "foo", false},
		// Regex pattern
		{"regex pattern", `p/\d+/`, `\d+`, false},
		// Escaped delimiter
		{"escaped delimiter", `p/foo\/bar/`, "foo/bar", false},
		// Empty pattern (matches all lines)
		{"empty pattern", "p//", "", false},
		// Literal matching (quote delimiters)
		{"backtick literal dot", "p`foo.bar`", `foo\.bar`, false},
		{"backtick literal brackets", "p`[test]`", `\[test\]`, false},
		{"single quote literal", "p'foo.bar'", `foo\.bar`, false},
		{"double quote literal", `p"foo.bar"`, `foo\.bar`, false},
		// Errors
		{"too short", "p", "", true},
		{"invalid regex", "p/[invalid/", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := ParseRule(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			pr, ok := r.(*rule.PrintLineRule)
			if !ok {
				t.Fatalf("expected *PrintLineRule, got %T", r)
			}

			if pr.Pattern() != tt.wantPattern {
				t.Errorf("pattern: got %q, want %q", pr.Pattern(), tt.wantPattern)
			}
		})
	}
}

func TestParseRule_PrintLineNum(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"single line", "p:5", false},
		{"range", "p:2-4", false},
		{"open from", "p:5-", false},
		{"open to", "p:-5", false},
		{"composite", "p:1,3,5-7", false},
		{"invalid", "p:abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := ParseRule(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			_, ok := r.(*rule.PrintLineNumRule)
			if !ok {
				t.Fatalf("expected *PrintLineNumRule, got %T", r)
			}
		})
	}
}
