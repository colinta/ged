package parser

import (
	"testing"

	"github.com/colinta/ged/internal/rule"
)

func TestParseUniq(t *testing.T) {
	tests := []struct {
		name  string
		input string
		ok    bool
	}{
		{"bare uniq", "uniq", true},
		{"pattern", `uniq/^\w+/`, true},
		{"pattern with flags", `uniq/^\w+/i`, true},
		{"pattern with group", `uniq/^\w+ (\w+)/1`, true},
		{"group then flags", `uniq/^\w+ (\w+)/1/i`, true},
		{"flags then group", `uniq/^\w+ (\w+)/i/1`, true},
		{"alternate delimiter", `uniq|^\w+|`, true},
		{"literal backtick", "uniq`foo.bar`", true},
		{"empty pattern errors", `uniq//`, false},
		{"invalid pattern errors", `uniq/[invalid/`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseRule(tt.input)
			if tt.ok {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if _, ok := result.(*rule.UniqRule); !ok {
					t.Fatalf("expected *rule.UniqRule, got %T", result)
				}
			} else {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			}
		})
	}
}
