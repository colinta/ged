package parser

import (
	"testing"

	"github.com/colinta/ged/internal/rule"
)

func TestParseCols(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "whitespace default", input: "cols//1,2"},
		{name: "comma pattern", input: "cols/,/1,2"},
		{name: "with joiner", input: "cols/,/1,2/|"},
		{name: "regex pattern", input: `cols/\s*,\s*/1,3`},
		{name: "literal backtick", input: "cols`,`1,2"},
		{name: "negative index", input: "cols//1,-1"},
		{name: "range spec", input: "cols//2-4"},

		// Errors
		{name: "missing spec", input: "cols//", wantErr: true},
		{name: "missing delimiter", input: "cols", wantErr: true},
		{name: "invalid spec", input: "cols//abc", wantErr: true},
		{name: "invalid pattern", input: "cols/[invalid/1", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseRule(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Should return a LineRule (ColumnsRule)
			if _, ok := result.(rule.LineRule); !ok {
				t.Fatalf("expected LineRule, got %T", result)
			}
		})
	}
}
