package parser

import (
	"testing"

	"github.com/colinta/ged/internal/rule"
)

func TestParseSplit(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"basic split", "split/,/", false},
		{"split on whitespace", `split/\s+/`, false},
		{"split alternate delimiter", "split|,|", false},
		{"split with flags", `split/and/i`, false},
		{"split literal backtick", "split`|`", false},
		{"split missing pattern", "split//", true},
		{"split no delimiter", "split", true},
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
			if _, ok := got.(*rule.SplitRule); !ok {
				t.Errorf("expected *rule.SplitRule, got %T", got)
			}
		})
	}
}

func TestParseInsert(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"basic insert", "insert/^#/---/", false},
		{"insert alternate delimiter", "insert|pattern|text|", false},
		{"insert with flags", `insert/todo/FIXME/i`, false},
		{"insert literal backtick", "insert`#`---`", false},
		{"insert empty text inserts blank line", "insert/pattern//", false},
		{"insert missing text", "insert/pattern", true},
		{"insert missing pattern", "insert//text/", true},
		{"insert no delimiter", "insert", true},
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
			if _, ok := got.(*rule.InsertRule); !ok {
				t.Errorf("expected *rule.InsertRule, got %T", got)
			}
		})
	}
}
