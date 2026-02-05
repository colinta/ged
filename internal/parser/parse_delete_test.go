package parser

import (
	"testing"

	"github.com/colinta/ged/internal/rule"
)

func TestParseRule_Delete(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantPattern string
		wantErr     bool
	}{
		// Basic delete
		{"basic", "d/foo/", "foo", false},
		{"without trailing delimiter", "d/foo", "foo", false},
		// Different delimiters
		{"pipe delimiter", "d|foo|", "foo", false},
		{"equals delimiter", "d=foo=", "foo", false},
		// Regex pattern
		{"regex pattern", `d/^\s*#/`, `^\s*#`, false},
		// Escaped delimiter
		{"escaped delimiter", `d/foo\/bar/`, "foo/bar", false},
		// Errors
		{"too short", "d", "", true},
		{"invalid regex", "d/[invalid/", "", true},
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

			dr, ok := r.(*rule.DeleteLineRule)
			if !ok {
				t.Fatalf("expected *DeleteLineRule, got %T", r)
			}

			if dr.Pattern() != tt.wantPattern {
				t.Errorf("pattern: got %q, want %q", dr.Pattern(), tt.wantPattern)
			}
		})
	}
}
