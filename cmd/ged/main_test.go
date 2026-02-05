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
