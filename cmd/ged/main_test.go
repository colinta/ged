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
