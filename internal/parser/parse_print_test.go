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
