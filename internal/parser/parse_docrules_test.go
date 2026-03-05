package parser

import (
	"fmt"
	"testing"

	"github.com/colinta/ged/internal/rule"
)

func TestParseDocRuleWords(t *testing.T) {
	tests := []struct {
		input    string
		wantType string
	}{
		{"lines", "*rule.LinesRule"},
		{"count", "*rule.CountRule"},
		{"uniq", "*rule.UniqRule"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseRule(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if typeName(got) != tt.wantType {
				t.Errorf("got %s, want %s", typeName(got), tt.wantType)
			}
		})
	}
}

func TestParseBegin(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"basic", "begin/header/", false},
		{"with newline", `begin/line1\nline2/`, false},
		{"alternate delim", "begin|===|", false},
		{"missing text", "begin//", true},
		{"no delimiter", "begin", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRule(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if _, ok := got.(*rule.BeginRule); !ok {
				t.Errorf("expected *rule.BeginRule, got %T", got)
			}
		})
	}
}

func TestParseEnd(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"basic", "end/footer/", false},
		{"alternate delim", "end|---|", false},
		{"missing text", "end//", true},
		{"no delimiter", "end", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRule(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if _, ok := got.(*rule.EndRule); !ok {
				t.Errorf("expected *rule.EndRule, got %T", got)
			}
		})
	}
}

func TestParseBorder(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"basic", "border/---/", false},
		{"alternate delim", "border|===|", false},
		{"missing text", "border//", true},
		{"no delimiter", "border", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRule(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if _, ok := got.(*rule.BorderRule); !ok {
				t.Errorf("expected *rule.BorderRule, got %T", got)
			}
		})
	}
}

func typeName(v any) string {
	return fmt.Sprintf("%T", v)
}
