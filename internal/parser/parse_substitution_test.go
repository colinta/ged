package parser

import (
	"testing"

	"github.com/colinta/ged/internal/rule"
)

func TestParseRule_Substitution(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantPattern string
		wantReplace string
		wantGlobal  bool
		wantErr     bool
	}{
		// Basic substitution
		{"basic", "s/foo/bar", "foo", "bar", false, false},
		{"basic with trailing delimiter", "s/foo/bar/", "foo", "bar", false, false},
		{"empty replacement", "s/foo/", "foo", "", false, false},
		// Global flag
		{"global", "s/foo/bar/g", "foo", "bar", true, false},
		// Different delimiters
		{"pipe delimiter", "s|foo|bar", "foo", "bar", false, false},
		{"pipe with global", "s|foo|bar|g", "foo", "bar", true, false},
		{"equals delimiter", "s=foo=bar", "foo", "bar", false, false},
		{"hash delimiter", "s#foo#bar", "foo", "bar", false, false},
		// Empty replacement
		{"empty replacement", "s/foo//", "foo", "", false, false},
		// Empty pattern (valid for parser)
		{"empty pattern", "s//bar/", "", "bar", false, false},
		// Regex pattern
		{"regex pattern", `s/\d+/NUM/g`, `\d+`, "NUM", true, false},
		// Escaped delimiter in pattern
		{"escaped delimiter in pattern", `s/foo\/bar/baz/`, "foo/bar", "baz", false, false},
		// Escaped delimiter in replacement
		{"escaped delimiter in replacement", `s/foo/bar\/baz/`, "foo", "bar/baz", false, false},
		// Multiple escaped delimiters
		{"multiple escapes", `s/a\/b\/c/d\/e/`, "a/b/c", "d/e", false, false},
		// Delimiter in replacement using different delimiter
		{"alt delimiter avoids escape", "s|foo/bar|baz|", "foo/bar", "baz", false, false},
		// Escaped backslash (not a delimiter escape)
		{"escaped backslash", `s/foo\\bar/baz/`, `foo\bar`, "baz", false, false},
		// Whitespace preserved
		{"whitespace preserved", "s/ foo / bar /", " foo ", " bar ", false, false},
		// Literal matching (quote delimiters)
		{"backtick literal dot", "s`foo.bar`baz`", `foo\.bar`, "baz", false, false},
		{"backtick literal star", "s`a*b`c`", `a\*b`, "c", false, false},
		{"backtick no metacharacters", "s`foo`bar`", "foo", "bar", false, false},
		{"single quote literal dot", "s'foo.bar'baz'", `foo\.bar`, "baz", false, false},
		{"double quote literal dot", `s"foo.bar"baz"`, `foo\.bar`, "baz", false, false},
		// Escape sequences
		{"newline in replacement", `s/foo/bar\nbaz/`, "foo", "bar\nbaz", false, false},
		{"tab in replacement", `s/foo/bar\tbaz/`, "foo", "bar\tbaz", false, false},
		{"newline in pattern", `s/foo\nbar/baz/`, "foo\nbar", "baz", false, false},
		// Errors
		{"too short", "s/", "", "", false, true},
		{"missing replacement", "s/foo", "", "", false, true},
		{"unknown command", "x/foo/bar/", "", "", false, true},
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

			// Type assert to SubstitutionRule to check internals
			sub, ok := r.(*rule.SubstitutionRule)
			if !ok {
				t.Fatalf("expected *SubstitutionRule, got %T", r)
			}

			if sub.Pattern() != tt.wantPattern {
				t.Errorf("pattern: got %q, want %q", sub.Pattern(), tt.wantPattern)
			}
			if sub.Replace() != tt.wantReplace {
				t.Errorf("replace: got %q, want %q", sub.Replace(), tt.wantReplace)
			}
			if sub.Global() != tt.wantGlobal {
				t.Errorf("global: got %v, want %v", sub.Global(), tt.wantGlobal)
			}
		})
	}
}

func TestParseRule_SubstitutionLineNum(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"single line", "s:5:replaced", false},
		{"range", "s:2-4:new content", false},
		{"open from", "s:5-:replaced", false},
		{"open to", "s:-5:replaced", false},
		{"composite", "s:1,3,5-7:replaced", false},
		{"empty replacement", "s:2:", false},
		{"invalid range", "s:abc:replaced", true},
		{"missing replacement", "s:5", true},
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

			_, ok := r.(*rule.SubLineNumRule)
			if !ok {
				t.Fatalf("expected *SubLineNumRule, got %T", r)
			}
		})
	}
}
