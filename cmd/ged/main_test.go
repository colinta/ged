package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestRun_BasicSubstitution(t *testing.T) {
	in := strings.NewReader("hello world")
	out := &bytes.Buffer{}

	err := run([]string{"s/world/earth"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "hello earth\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_GlobalSubstitution(t *testing.T) {
	in := strings.NewReader("hello world world")
	out := &bytes.Buffer{}

	err := run([]string{"s/world/earth/g"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "hello earth earth\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_MultipleLines(t *testing.T) {
	in := strings.NewReader("line1\nline2\nline3")
	out := &bytes.Buffer{}

	err := run([]string{"s/line/row"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "row1\nrow2\nrow3\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_RegexPattern(t *testing.T) {
	in := strings.NewReader("foo 123 bar 456")
	out := &bytes.Buffer{}

	err := run([]string{`s/\d+/NUM/g`}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "foo NUM bar NUM\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_NoArgs(t *testing.T) {
	in := strings.NewReader("hello")
	out := &bytes.Buffer{}

	err := run([]string{}, in, out, io.Discard)
	if err == nil {
		t.Error("expected error for no args, got nil")
	}
}

func TestRun_InvalidRule(t *testing.T) {
	in := strings.NewReader("hello")
	out := &bytes.Buffer{}

	err := run([]string{"x/invalid"}, in, out, io.Discard)
	if err == nil {
		t.Error("expected error for invalid rule, got nil")
	}
}

func TestRun_InvalidRegex(t *testing.T) {
	in := strings.NewReader("hello")
	out := &bytes.Buffer{}

	err := run([]string{"s/[invalid/replacement"}, in, out, io.Discard)
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}

func TestRun_EmptyInput(t *testing.T) {
	in := strings.NewReader("")
	out := &bytes.Buffer{}

	err := run([]string{"s/foo/bar"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.String() != "" {
		t.Errorf("expected empty output, got %q", out.String())
	}
}

func TestRun_PrintKeepsMatchingLines(t *testing.T) {
	in := strings.NewReader("foo\nbar\nfoo baz")
	out := &bytes.Buffer{}

	err := run([]string{"p/foo"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "foo\nfoo baz\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_PrintWithRegex(t *testing.T) {
	in := strings.NewReader("123\nabc\n456")
	out := &bytes.Buffer{}

	err := run([]string{`p/^\d+$`}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "123\n456\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_DeleteRemovesMatchingLines(t *testing.T) {
	in := strings.NewReader("foo\nbar\nfoo baz")
	out := &bytes.Buffer{}

	err := run([]string{"d/foo"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "bar\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_DeleteComments(t *testing.T) {
	in := strings.NewReader("code\n# comment\nmore code\n  # indented comment")
	out := &bytes.Buffer{}

	err := run([]string{`d/^\s*#`}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "code\nmore code\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_MultipleRules(t *testing.T) {
	in := strings.NewReader("hello\nworld\nhello world")
	out := &bytes.Buffer{}

	// Keep lines with "hello", then replace "o" with "0"
	err := run([]string{"p/hello", "s/o/0/g"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "hell0\nhell0 w0rld\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_ChainedSubstitutions(t *testing.T) {
	in := strings.NewReader("abc")
	out := &bytes.Buffer{}

	// a->b, then b->c (first match only each time)
	err := run([]string{"s/a/b", "s/b/c"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// "abc" -> "bbc" -> "cbc"
	want := "cbc\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_FilterDeletesBeforeSubstitute(t *testing.T) {
	in := strings.NewReader("keep this\ndelete this\nkeep that")
	out := &bytes.Buffer{}

	// Delete lines with "delete", then substitute "keep" with "KEEP"
	err := run([]string{"d/delete", "s/keep/KEEP"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "KEEP this\nKEEP that\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_Sort(t *testing.T) {
	in := strings.NewReader("c\na\nb")
	out := &bytes.Buffer{}

	err := run([]string{"sort"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "a\nb\nc\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_Reverse(t *testing.T) {
	in := strings.NewReader("a\nb\nc")
	out := &bytes.Buffer{}

	err := run([]string{"reverse"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "c\nb\na\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_JoinWithComma(t *testing.T) {
	in := strings.NewReader("a\nb\nc")
	out := &bytes.Buffer{}

	err := run([]string{"join/,/"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "a,b,c\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_LineRulesThenSort(t *testing.T) {
	in := strings.NewReader("c3\na1\nb2")
	out := &bytes.Buffer{}

	// Remove digits, then sort
	err := run([]string{`s/[0-9]//g`, "sort"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "a\nb\nc\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_SortThenLineRules(t *testing.T) {
	in := strings.NewReader("cherry\napple\nbanana")
	out := &bytes.Buffer{}

	// Sort, then uppercase the first letter
	err := run([]string{"sort", "s/a/A"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Apple\nbAnana\ncherry\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_JoinBare(t *testing.T) {
	in := strings.NewReader("a\nb\nc")
	out := &bytes.Buffer{}

	err := run([]string{"join"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "abc\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_IfCondition(t *testing.T) {
	in := strings.NewReader("hello\nworld\nhello world")
	out := &bytes.Buffer{}

	// Only substitute on lines containing "hello"
	err := run([]string{"if/hello/", "{", "s/o/x/", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "hellx\nworld\nhellx world\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_IfConditionInverted(t *testing.T) {
	in := strings.NewReader("hello\nworld\nhello world")
	out := &bytes.Buffer{}

	// Substitute on lines NOT containing "hello"
	err := run([]string{"!if/hello/", "{", "s/o/x/", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "hello\nwxrld\nhello world\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_IfWithMultipleInnerRules(t *testing.T) {
	in := strings.NewReader("hello\nworld\nhello world")
	out := &bytes.Buffer{}

	// On "hello" lines: replace "h" then "e"
	err := run([]string{"if/hello/", "{", "s/h/H/", "s/e/E/", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "HEllo\nworld\nHEllo world\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_IfThenSort(t *testing.T) {
	in := strings.NewReader("b_hello\na_hello\nc_world")
	out := &bytes.Buffer{}

	// Conditional then sort — conditional is a LineRule, sort is DocumentRule
	err := run([]string{"if/hello/", "{", "s/_hello//", "}", "sort"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "a\nb\nc_world\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_IfWithDocumentRule(t *testing.T) {
	in := strings.NewReader("b_item\na_item\nc_other\nd_item")
	out := &bytes.Buffer{}

	// Sort only lines matching "item"
	err := run([]string{"if/item/", "{", "sort", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// items sorted: a_item, b_item, d_item woven back into positions 0,1,3
	// c_other stays at position 2
	want := "a_item\nb_item\nc_other\nd_item\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_NestedIf(t *testing.T) {
	in := strings.NewReader("ab\nac\nbd\nbc")
	out := &bytes.Buffer{}

	// Nested: only apply to lines with "a" AND "b"
	err := run([]string{"if/a/", "{", "if/b/", "{", "s/ab/AB/", "}", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "AB\nac\nbd\nbc\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_On(t *testing.T) {
	in := strings.NewReader("a\nstart\nb\nc")
	out := &bytes.Buffer{}

	err := run([]string{"on/start/"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "start\nb\nc\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_Off(t *testing.T) {
	in := strings.NewReader("a\nb\nstop\nc")
	out := &bytes.Buffer{}

	err := run([]string{"off/stop/"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "a\nb\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_After(t *testing.T) {
	in := strings.NewReader("a\nmarker\nb\nc")
	out := &bytes.Buffer{}

	err := run([]string{"after/marker/"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "b\nc\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_Toggle(t *testing.T) {
	in := strings.NewReader("off1\n---\non1\non2\n---\noff2")
	out := &bytes.Buffer{}

	err := run([]string{"toggle/---/"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "---\non1\non2\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_OnWithSubstitution(t *testing.T) {
	in := strings.NewReader("a\nstart\nb\nc")
	out := &bytes.Buffer{}

	err := run([]string{"on/start/", "s/b/B/"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "start\nB\nc\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_OnOffCombined(t *testing.T) {
	in := strings.NewReader("before\nstart\nmiddle\nend\nafter")
	out := &bytes.Buffer{}

	err := run([]string{"on/start/", "off/end/"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "start\nmiddle\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}
