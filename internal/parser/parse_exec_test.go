package parser

import (
	"testing"

	"github.com/colinta/ged/internal/rule"
)

func TestParseXargs(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"basic xargs", "xargs/echo/", false},
		{"xargs with complex command", "xargs/echo hello/", false},
		{"xargs alternate delimiter", "xargs|echo|", false},
		{"xargs missing command", "xargs//", true},
		{"xargs no delimiter", "xargs", true},
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
			if _, ok := got.(*rule.XargsRule); !ok {
				t.Errorf("expected *rule.XargsRule, got %T", got)
			}
		})
	}
}

func TestParseExec(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"basic exec", "exec/sort/", false},
		{"exec with flags", "exec/sort -n/", false},
		{"exec with pipe", "exec/grep foo | wc -l/", false},
		{"exec alternate delimiter", "exec|cat -n|", false},
		{"exec missing command", "exec//", true},
		{"exec no delimiter", "exec", true},
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
			if _, ok := got.(*rule.ExecRule); !ok {
				t.Errorf("expected *rule.ExecRule, got %T", got)
			}
		})
	}
}
