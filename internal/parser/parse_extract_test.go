package parser

import (
	"testing"

	"github.com/colinta/ged/internal/rule"
)

func TestParseTake(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "basic take", input: `t/\d+/`},
		{name: "take with flags", input: `t/hello/i`},
		{name: "take with global", input: `t/\w+/g`},
		{name: "take with pipe delimiter", input: `t|\d+|`},
		{name: "take with backtick literal", input: "t`foo.bar`"},
		{name: "empty pattern errors", input: `t//`, wantErr: true},
		{name: "too short errors", input: `t`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseRule(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseRule() error: %v", err)
			}
			if _, ok := result.(*rule.TakeRule); !ok {
				t.Fatalf("expected *rule.TakeRule, got %T", result)
			}
		})
	}
}

func TestParseRemove(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "basic remove", input: `r/\d+/`},
		{name: "remove with flags", input: `r/hello/gi`},
		{name: "remove with pipe delimiter", input: `r|\s+|`},
		{name: "remove with backtick literal", input: "r`foo.bar`"},
		{name: "empty pattern errors", input: `r//`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseRule(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseRule() error: %v", err)
			}
			if _, ok := result.(*rule.RemoveRule); !ok {
				t.Fatalf("expected *rule.RemoveRule, got %T", result)
			}
		})
	}
}

func TestParseGroup(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "group 1", input: `1/(\w+) (\w+)/`},
		{name: "group 2", input: `2/(\w+) (\w+)/`},
		{name: "group 9", input: `9/(.)/`},
		{name: "group with flags", input: `1/(hello)/i`},
		{name: "group with pipe delimiter", input: `1|(\d+)|`},
		{name: "empty pattern errors", input: `1//`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseRule(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseRule() error: %v", err)
			}
			if _, ok := result.(*rule.GroupRule); !ok {
				t.Fatalf("expected *rule.GroupRule, got %T", result)
			}
		})
	}
}
