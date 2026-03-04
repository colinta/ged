package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
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

func TestRun_Between(t *testing.T) {
	in := strings.NewReader("before\nSTART\n1\n2\nEND\nafter")
	out := &bytes.Buffer{}

	err := run([]string{"between/START/END/", "{", "s/\\d/x/", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "before\nSTART\nx\nx\nEND\nafter\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_BetweenInverted(t *testing.T) {
	in := strings.NewReader("x\nSTART\nx\nEND\nx")
	out := &bytes.Buffer{}

	err := run([]string{"!between/START/END/", "{", "s/x/X/", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "X\nSTART\nx\nEND\nX\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_BetweenMultipleRanges(t *testing.T) {
	in := strings.NewReader("x\nA\nx\nB\nx\nA\nx\nB\nx")
	out := &bytes.Buffer{}

	err := run([]string{"between/A/B/", "{", "s/x/X/g", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "x\nA\nX\nB\nx\nA\nX\nB\nx\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_BetweenWithDocRule(t *testing.T) {
	in := strings.NewReader("before\nSTART\nc\na\nb\nEND\nafter")
	out := &bytes.Buffer{}

	err := run([]string{"between/START/END/", "{", "sort", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "before\nEND\nSTART\na\nb\nc\nafter\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

// --- IgnoreCase flag tests ---

func TestRun_SubstitutionIgnoreCase(t *testing.T) {
	in := strings.NewReader("Hello World")
	out := &bytes.Buffer{}

	err := run([]string{"s/hello/HI/i"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "HI World\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_SubstitutionIgnoreCaseGlobal(t *testing.T) {
	in := strings.NewReader("Hello hello HELLO")
	out := &bytes.Buffer{}

	err := run([]string{"s/hello/x/gi"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "x x x\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_PrintIgnoreCase(t *testing.T) {
	in := strings.NewReader("Hello\nworld\nHELLO")
	out := &bytes.Buffer{}

	err := run([]string{"p/hello/i"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Hello\nHELLO\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_DeleteIgnoreCase(t *testing.T) {
	in := strings.NewReader("Hello\nworld\nHELLO")
	out := &bytes.Buffer{}

	err := run([]string{"d/hello/i"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "world\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_IfIgnoreCase(t *testing.T) {
	in := strings.NewReader("Hello\nworld\nHELLO")
	out := &bytes.Buffer{}

	err := run([]string{"if/hello/i", "{", "s/l/L/g", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "HeLLo\nworld\nHELLO\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_OnIgnoreCase(t *testing.T) {
	in := strings.NewReader("a\nSTART\nb\nc")
	out := &bytes.Buffer{}

	err := run([]string{"on/start/i"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "START\nb\nc\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_BetweenIgnoreCase(t *testing.T) {
	in := strings.NewReader("before\nStart\n1\n2\nEnd\nafter")
	out := &bytes.Buffer{}

	err := run([]string{"between/start/end/i", "{", "s/\\d/x/", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "before\nStart\nx\nx\nEnd\nafter\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

// --- File I/O tests ---

// writeTestFile creates a temporary file with the given content and returns its path.
// The file is automatically cleaned up when the test finishes.
func writeTestFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "ged-test-*")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestRun_InputFile(t *testing.T) {
	path := writeTestFile(t, "hello world\n")
	out := &bytes.Buffer{}

	err := run([]string{"s/world/earth", "--input=" + path}, nil, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "hello earth\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_InputFileNotFound(t *testing.T) {
	out := &bytes.Buffer{}

	err := run([]string{"s/a/b", "--input=/nonexistent/file.txt"}, nil, out, io.Discard)
	if err == nil {
		t.Error("expected error for missing input file, got nil")
	}
}

func TestRun_WriteBack(t *testing.T) {
	path := writeTestFile(t, "hello world\n")
	out := &bytes.Buffer{}

	err := run([]string{"s/world/earth", "--input=" + path, "--write"}, nil, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// stdout should be empty (output went to file)
	if out.String() != "" {
		t.Errorf("expected no stdout, got %q", out.String())
	}

	// File should be updated
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := "hello earth\n"
	if string(data) != want {
		t.Errorf("file content: got %q, want %q", string(data), want)
	}
}

func TestRun_WriteBackPreservesPermissions(t *testing.T) {
	path := writeTestFile(t, "hello\n")
	os.Chmod(path, 0755)
	out := &bytes.Buffer{}

	err := run([]string{"s/hello/world", "--input=" + path, "--write"}, nil, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0755 {
		t.Errorf("permissions: got %o, want 0755", info.Mode().Perm())
	}
}

func TestRun_WriteTo(t *testing.T) {
	inPath := writeTestFile(t, "hello world\n")
	outPath := filepath.Join(t.TempDir(), "output.txt")
	out := &bytes.Buffer{}

	err := run([]string{"s/world/earth", "--input=" + inPath, "--write-to=" + outPath}, nil, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Original should be unchanged
	orig, _ := os.ReadFile(inPath)
	if string(orig) != "hello world\n" {
		t.Errorf("original changed: %q", string(orig))
	}

	// Output file should have transformed content
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	want := "hello earth\n"
	if string(data) != want {
		t.Errorf("output file: got %q, want %q", string(data), want)
	}
}

func TestRun_WriteToFromStdin(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "output.txt")
	in := strings.NewReader("hello world\n")
	out := &bytes.Buffer{}

	err := run([]string{"s/world/earth", "--write-to=" + outPath}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	want := "hello earth\n"
	if string(data) != want {
		t.Errorf("output file: got %q, want %q", string(data), want)
	}
}

func TestRun_WriteRequiresInput(t *testing.T) {
	err := run([]string{"s/a/b", "--write"}, nil, &bytes.Buffer{}, io.Discard)
	if err == nil {
		t.Error("expected error for --write without --input, got nil")
	}
}

func TestRun_WriteAndWriteToMutuallyExclusive(t *testing.T) {
	path := writeTestFile(t, "hello\n")
	err := run([]string{"s/a/b", "--input=" + path, "--write", "--write-to=out.txt"}, nil, &bytes.Buffer{}, io.Discard)
	if err == nil {
		t.Error("expected error for --write + --write-to, got nil")
	}
}

func TestRun_MultipleInputFiles(t *testing.T) {
	path1 := writeTestFile(t, "aaa\n")
	path2 := writeTestFile(t, "bbb\n")
	out := &bytes.Buffer{}

	err := run([]string{"s/a/x/g", "s/b/y/g", "--input=" + path1, "--input=" + path2}, nil, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "xxx\nyyy\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_MultipleInputFilesWriteBack(t *testing.T) {
	path1 := writeTestFile(t, "aaa\n")
	path2 := writeTestFile(t, "bbb\n")
	out := &bytes.Buffer{}

	err := run([]string{"s/a/x/g", "s/b/y/g", "--input=" + path1, "--input=" + path2, "--write"}, nil, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data1, _ := os.ReadFile(path1)
	data2, _ := os.ReadFile(path2)
	if string(data1) != "xxx\n" {
		t.Errorf("file1: got %q, want %q", string(data1), "xxx\n")
	}
	if string(data2) != "yyy\n" {
		t.Errorf("file2: got %q, want %q", string(data2), "yyy\n")
	}
}

func TestRun_MultipleInputFilesWriteToError(t *testing.T) {
	path1 := writeTestFile(t, "a\n")
	path2 := writeTestFile(t, "b\n")
	err := run([]string{"s/a/b", "--input=" + path1, "--input=" + path2, "--write-to=out.txt"}, nil, &bytes.Buffer{}, io.Discard)
	if err == nil {
		t.Error("expected error for --write-to with multiple inputs, got nil")
	}
}

func TestRun_InputSpaceSeparated(t *testing.T) {
	path := writeTestFile(t, "hello world\n")
	out := &bytes.Buffer{}

	err := run([]string{"s/world/earth", "--input", path}, nil, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "hello earth\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_InputSpaceSeparatedMultiple(t *testing.T) {
	path1 := writeTestFile(t, "aaa\n")
	path2 := writeTestFile(t, "bbb\n")
	out := &bytes.Buffer{}

	err := run([]string{"s/a/x/g", "s/b/y/g", "--input", path1, "--input", path2}, nil, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "xxx\nyyy\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_InputMissingFilename(t *testing.T) {
	err := run([]string{"s/a/b", "--input"}, nil, &bytes.Buffer{}, io.Discard)
	if err == nil {
		t.Error("expected error for --input without filename, got nil")
	}
}

func TestRun_WriteToSpaceSeparated(t *testing.T) {
	inPath := writeTestFile(t, "hello world\n")
	outPath := filepath.Join(t.TempDir(), "output.txt")
	out := &bytes.Buffer{}

	err := run([]string{"s/world/earth", "--input", inPath, "--write-to", outPath}, nil, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	want := "hello earth\n"
	if string(data) != want {
		t.Errorf("output file: got %q, want %q", string(data), want)
	}
}

func TestRun_WriteToMissingFilename(t *testing.T) {
	err := run([]string{"s/a/b", "--write-to"}, nil, &bytes.Buffer{}, io.Discard)
	if err == nil {
		t.Error("expected error for --write-to without filename, got nil")
	}
}

func TestRun_BareRulesAfterDash(t *testing.T) {
	input := strings.NewReader("hello world\n")
	out := &bytes.Buffer{}

	// "--" prevents "--write" from being treated as a flag
	err := run([]string{"--", "s/world/earth"}, input, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "hello earth\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_BareRulesAfterDash_FlagLikeRule(t *testing.T) {
	input := strings.NewReader("--write\n--input=foo\n")
	out := &bytes.Buffer{}

	// Rule args after -- that look like flags should be treated as rules
	err := run([]string{"--", "s/--/++/"}, input, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "++write\n++input=foo\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_InputFileWithSort(t *testing.T) {
	path := writeTestFile(t, "c\na\nb\n")
	out := &bytes.Buffer{}

	err := run([]string{"sort", "--input=" + path}, nil, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "a\nb\nc\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

// --- Phase 10: Text modification rules ---

func TestRun_Trim(t *testing.T) {
	in := strings.NewReader("  hello  \n\tworld\t\n  foo bar  ")
	out := &bytes.Buffer{}

	err := run([]string{"trim"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "hello\nworld\nfoo bar\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_TrimLeft(t *testing.T) {
	in := strings.NewReader("  hello  ")
	out := &bytes.Buffer{}

	err := run([]string{"triml"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "hello  \n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_TrimRight(t *testing.T) {
	in := strings.NewReader("  hello  ")
	out := &bytes.Buffer{}

	err := run([]string{"trimr"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "  hello\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_Upper(t *testing.T) {
	in := strings.NewReader("hello\nWorld")
	out := &bytes.Buffer{}

	err := run([]string{"upper"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "HELLO\nWORLD\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_Lower(t *testing.T) {
	in := strings.NewReader("HELLO\nWorld")
	out := &bytes.Buffer{}

	err := run([]string{"lower"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "hello\nworld\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_Prepend(t *testing.T) {
	in := strings.NewReader("hello\nworld")
	out := &bytes.Buffer{}

	err := run([]string{"prepend/>> /"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := ">> hello\n>> world\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_Append(t *testing.T) {
	in := strings.NewReader("hello\nworld")
	out := &bytes.Buffer{}

	err := run([]string{"append/;/"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "hello;\nworld;\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_Surround(t *testing.T) {
	in := strings.NewReader("hello\nworld")
	out := &bytes.Buffer{}

	err := run([]string{"surround/(/)/)"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "(hello)\n(world)\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_TrimThenUpper(t *testing.T) {
	in := strings.NewReader("  hello  \n  world  ")
	out := &bytes.Buffer{}

	err := run([]string{"trim", "upper"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "HELLO\nWORLD\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_IfWithPrepend(t *testing.T) {
	in := strings.NewReader("error: bad\ninfo: ok\nerror: worse")
	out := &bytes.Buffer{}

	err := run([]string{"if/error/", "{", "prepend/!! /", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "!! error: bad\ninfo: ok\n!! error: worse\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

// --- Phase 10: Between + text modification ---

func TestRun_BetweenWithTrim(t *testing.T) {
	in := strings.NewReader("  before  \nSTART\n  hello  \n  world  \nEND\n  after  ")
	out := &bytes.Buffer{}

	err := run([]string{"between/START/END/", "{", "trim", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "  before  \nSTART\nhello\nworld\nEND\n  after  \n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_BetweenWithUpper(t *testing.T) {
	in := strings.NewReader("before\nSTART\nhello\nworld\nEND\nafter")
	out := &bytes.Buffer{}

	err := run([]string{"between/START/END/", "{", "upper", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "before\nSTART\nHELLO\nWORLD\nEND\nafter\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_BetweenWithLower(t *testing.T) {
	in := strings.NewReader("BEFORE\nSTART\nHELLO\nWORLD\nEND\nAFTER")
	out := &bytes.Buffer{}

	err := run([]string{"between/START/END/", "{", "lower", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "BEFORE\nstart\nhello\nworld\nend\nAFTER\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_BetweenWithPrepend(t *testing.T) {
	in := strings.NewReader("before\nSTART\nhello\nworld\nEND\nafter")
	out := &bytes.Buffer{}

	err := run([]string{"between/START/END/", "{", "prepend/  /", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "before\n  START\n  hello\n  world\n  END\nafter\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_BetweenWithAppend(t *testing.T) {
	in := strings.NewReader("before\nSTART\nhello\nworld\nEND\nafter")
	out := &bytes.Buffer{}

	err := run([]string{"between/START/END/", "{", "append/!/", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "before\nSTART!\nhello!\nworld!\nEND!\nafter\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_BetweenWithSurround(t *testing.T) {
	in := strings.NewReader("before\nSTART\nhello\nworld\nEND\nafter")
	out := &bytes.Buffer{}

	err := run([]string{"between/START/END/", "{", "surround/[/]/", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "before\n[START]\n[hello]\n[world]\n[END]\nafter\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_BetweenInvertedWithUpper(t *testing.T) {
	in := strings.NewReader("hello\nSTART\ninside\nEND\nworld")
	out := &bytes.Buffer{}

	err := run([]string{"!between/START/END/", "{", "upper", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "HELLO\nSTART\ninside\nEND\nWORLD\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_BetweenWithChainedTextRules(t *testing.T) {
	in := strings.NewReader("before\nSTART\n  hello  \n  world  \nEND\nafter")
	out := &bytes.Buffer{}

	err := run([]string{"between/START/END/", "{", "trim", "upper", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "before\nSTART\nHELLO\nWORLD\nEND\nafter\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_BetweenMultipleRangesWithPrepend(t *testing.T) {
	in := strings.NewReader("x\nA\ny\nB\nx\nA\nz\nB\nx")
	out := &bytes.Buffer{}

	err := run([]string{"between/A/B/", "{", "prepend/>> /", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "x\n>> A\n>> y\n>> B\nx\n>> A\n>> z\n>> B\nx\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_ColsWhitespace(t *testing.T) {
	in := strings.NewReader("alice 25 engineer\nbob 30 designer")
	out := &bytes.Buffer{}

	err := run([]string{"cols//1,3"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "alice engineer\nbob designer\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_ColsComma(t *testing.T) {
	in := strings.NewReader("alice,25,engineer")
	out := &bytes.Buffer{}

	err := run([]string{"cols/,/3,1"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "engineer alice\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_ColsWithJoiner(t *testing.T) {
	in := strings.NewReader("a,b,c")
	out := &bytes.Buffer{}

	err := run([]string{"cols/,/1,3/ | "}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "a | c\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_ColsNegativeIndex(t *testing.T) {
	in := strings.NewReader("a b c d e")
	out := &bytes.Buffer{}

	err := run([]string{"cols//1,-1"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "a e\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_ColsThenSubstitute(t *testing.T) {
	in := strings.NewReader("alice 25 engineer\nbob 30 designer")
	out := &bytes.Buffer{}

	err := run([]string{"cols//1,3", "s/e/E/g"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "alicE EnginEEr\nbob dEsignEr\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_ColsWithIfCondition(t *testing.T) {
	in := strings.NewReader("alice,25,engineer\nbob,30,designer\ncharlie,35,engineer")
	out := &bytes.Buffer{}

	err := run([]string{"if/engineer/", "{", "cols/,/1,3", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "alice engineer\nbob,30,designer\ncharlie engineer\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_Take(t *testing.T) {
	in := strings.NewReader("abc 123 def\nno numbers\n456 xyz")
	out := &bytes.Buffer{}

	err := run([]string{`t/\d+/`}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "123\nno numbers\n456\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_TakeGlobal(t *testing.T) {
	in := strings.NewReader("a1 b2 c3")
	out := &bytes.Buffer{}

	err := run([]string{`t/\d+/g`}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "1 2 3\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_Remove(t *testing.T) {
	in := strings.NewReader("hello 123 world\nabc def")
	out := &bytes.Buffer{}

	err := run([]string{`r/\d+/`}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "hello  world\nabc def\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_RemoveGlobal(t *testing.T) {
	in := strings.NewReader("a1b2c3")
	out := &bytes.Buffer{}

	err := run([]string{`r/\d/g`}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "abc\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_GroupCapture(t *testing.T) {
	in := strings.NewReader("name: alice, age: 30\nname: bob, age: 25")
	out := &bytes.Buffer{}

	err := run([]string{`2/name: (\w+), age: (\d+)/`}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "30\n25\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_TakeThenSubstitute(t *testing.T) {
	in := strings.NewReader("email: user@host.com ok")
	out := &bytes.Buffer{}

	err := run([]string{`t/\S+@\S+/`, `s/@/ at /`}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "user at host.com\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_RemoveWithTrim(t *testing.T) {
	in := strings.NewReader("hello  # comment")
	out := &bytes.Buffer{}

	err := run([]string{`r/\s*#.*/`, `trimr`}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "hello\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_TakeGlobalWithJoiner(t *testing.T) {
	in := strings.NewReader("abc 123 def 456")
	out := &bytes.Buffer{}

	err := run([]string{`t/\d+/g/,`}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "123,456\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_XargsEcho(t *testing.T) {
	in := strings.NewReader("hello\nworld")
	out := &bytes.Buffer{}

	err := run([]string{"xargs/echo hi/"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "hi hello\nhi world\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_XargsWithSubstitution(t *testing.T) {
	in := strings.NewReader("hello\nworld")
	out := &bytes.Buffer{}

	err := run([]string{"s/o/0/", "xargs/echo/"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "hell0\nw0rld\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_ExecSort(t *testing.T) {
	in := strings.NewReader("c\na\nb")
	out := &bytes.Buffer{}

	err := run([]string{"exec/sort/"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "a\nb\nc\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_ExecWithPipe(t *testing.T) {
	in := strings.NewReader("hello\nworld\nhello again")
	out := &bytes.Buffer{}

	err := run([]string{"exec/grep hello/"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "hello\nhello again\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_ExecAfterSubstitution(t *testing.T) {
	in := strings.NewReader("3\n1\n2")
	out := &bytes.Buffer{}

	err := run([]string{"s/$/x/", "exec/sort/"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "1x\n2x\n3x\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_XargsInsideIf(t *testing.T) {
	in := strings.NewReader("hello\nworld\nhello")
	out := &bytes.Buffer{}

	err := run([]string{"if/hello/", "{", "xargs/echo matched:/", "}"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "matched: hello\nworld\nmatched: hello\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}
