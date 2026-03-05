package main

// Most CLI tests live in testdata/*.yaml and are run by yaml_test.go.
// This file contains tests that need file I/O, temp directories, or
// special assertions (color codes, file permissions) that don't fit
// the YAML format.

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- File I/O tests ---

func TestRun_InputFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "ged-test-*")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("hello world\n")
	path := f.Name()
	f.Close()

	out := &bytes.Buffer{}
	err = run([]string{"--input=" + path, "s/world/earth/"}, nil, out, io.Discard)
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
	err := run([]string{"--input=/nonexistent/file.txt", "s/a/b/"}, nil, out, io.Discard)
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestRun_WriteBack(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("hello world\n"), 0644)

	err := run([]string{"--input=" + path, "--write", "s/world/earth/"}, nil, io.Discard, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}
	want := "hello earth\n"
	if string(data) != want {
		t.Errorf("got %q, want %q", string(data), want)
	}
}

func TestRun_WriteBackPreservesPermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("hello\n"), 0755)

	err := run([]string{"--input=" + path, "--write", "s/hello/world/"}, nil, io.Discard, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("error stating file: %v", err)
	}
	if info.Mode().Perm() != 0755 {
		t.Errorf("permissions changed: got %o, want %o", info.Mode().Perm(), 0755)
	}
}

func TestRun_WriteTo(t *testing.T) {
	dir := t.TempDir()
	inPath := filepath.Join(dir, "input.txt")
	outPath := filepath.Join(dir, "output.txt")
	os.WriteFile(inPath, []byte("hello world\n"), 0644)

	err := run([]string{"--input=" + inPath, "--write-to=" + outPath, "s/world/earth/"}, nil, io.Discard, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Original unchanged
	orig, _ := os.ReadFile(inPath)
	if string(orig) != "hello world\n" {
		t.Errorf("original changed: %q", string(orig))
	}

	// Output written
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("error reading output: %v", err)
	}
	want := "hello earth\n"
	if string(data) != want {
		t.Errorf("got %q, want %q", string(data), want)
	}
}

func TestRun_WriteToFromStdin(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "output.txt")
	in := strings.NewReader("hello world\n")

	err := run([]string{"--write-to=" + outPath, "s/world/earth/"}, in, io.Discard, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("error reading output: %v", err)
	}
	want := "hello earth\n"
	if string(data) != want {
		t.Errorf("got %q, want %q", string(data), want)
	}
}

func TestRun_MultipleInputFiles(t *testing.T) {
	dir := t.TempDir()
	path1 := filepath.Join(dir, "a.txt")
	path2 := filepath.Join(dir, "b.txt")
	os.WriteFile(path1, []byte("aaa\n"), 0644)
	os.WriteFile(path2, []byte("bbb\n"), 0644)

	out := &bytes.Buffer{}
	err := run([]string{"--input=" + path1, "--input=" + path2, "upper"}, nil, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "AAA\nBBB\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_MultipleInputFilesWriteBack(t *testing.T) {
	dir := t.TempDir()
	path1 := filepath.Join(dir, "a.txt")
	path2 := filepath.Join(dir, "b.txt")
	os.WriteFile(path1, []byte("aaa\n"), 0644)
	os.WriteFile(path2, []byte("bbb\n"), 0644)

	err := run([]string{"--input=" + path1, "--input=" + path2, "--write", "upper"}, nil, io.Discard, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data1, _ := os.ReadFile(path1)
	data2, _ := os.ReadFile(path2)
	if string(data1) != "AAA\n" {
		t.Errorf("file1: got %q, want %q", string(data1), "AAA\n")
	}
	if string(data2) != "BBB\n" {
		t.Errorf("file2: got %q, want %q", string(data2), "BBB\n")
	}
}

func TestRun_MultipleInputFilesWriteToError(t *testing.T) {
	err := run([]string{"--input=a.txt", "--input=b.txt", "--write-to=out.txt", "s/a/b/"}, nil, io.Discard, io.Discard)
	if err == nil {
		t.Error("expected error for --write-to with multiple inputs")
	}
}

func TestRun_InputSpaceSeparated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("hello\n"), 0644)

	out := &bytes.Buffer{}
	err := run([]string{"--input", path, "upper"}, nil, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "HELLO\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_InputSpaceSeparatedMultiple(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("hello\n"), 0644)

	outPath := filepath.Join(dir, "output.txt")
	err := run([]string{"--input", path, "--write-to", outPath, "upper"}, nil, io.Discard, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("error reading output: %v", err)
	}
	want := "HELLO\n"
	if string(data) != want {
		t.Errorf("got %q, want %q", string(data), want)
	}
}

func TestRun_InputMissingFilename(t *testing.T) {
	err := run([]string{"--input"}, nil, io.Discard, io.Discard)
	if err == nil {
		t.Error("expected error for --input without filename")
	}
}

func TestRun_WriteToSpaceSeparated(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "output.txt")
	in := strings.NewReader("hello\n")

	err := run([]string{"--write-to", outPath, "upper"}, in, io.Discard, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("error reading output: %v", err)
	}
	want := "HELLO\n"
	if string(data) != want {
		t.Errorf("got %q, want %q", string(data), want)
	}
}

func TestRun_WriteToMissingFilename(t *testing.T) {
	err := run([]string{"--write-to"}, nil, io.Discard, io.Discard)
	if err == nil {
		t.Error("expected error for --write-to without filename")
	}
}

func TestRun_BareRulesAfterDash(t *testing.T) {
	in := strings.NewReader("hello\n")
	out := &bytes.Buffer{}

	// -- ensures the next arg isn't treated as a flag
	err := run([]string{"--", "s/hello/world/"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "world\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_BareRulesAfterDash_FlagLikeRule(t *testing.T) {
	in := strings.NewReader("--flag\n")
	out := &bytes.Buffer{}

	err := run([]string{"--", "s/--flag/value/"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "value\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

func TestRun_InputFileWithSort(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("c\na\nb\n"), 0644)

	out := &bytes.Buffer{}
	err := run([]string{"--input=" + path, "sort"}, nil, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "a\nb\nc\n"
	if out.String() != want {
		t.Errorf("got %q, want %q", out.String(), want)
	}
}

// --- Diff with file tests ---

func TestRun_DiffWithFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("foo\nbar\nbaz\n"), 0644)

	out := &bytes.Buffer{}
	err := run([]string{"--diff", "--no-color", "--input=" + path, "s/bar/BAR/"}, nil, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "--- "+path) {
		t.Errorf("expected header with filename, got:\n%s", got)
	}
	if !strings.Contains(got, "-bar") {
		t.Errorf("expected '-bar' in diff output, got:\n%s", got)
	}
	if !strings.Contains(got, "+BAR") {
		t.Errorf("expected '+BAR' in diff output, got:\n%s", got)
	}
}

func TestRun_DiffWithColor(t *testing.T) {
	in := strings.NewReader("hello\nworld")
	out := &bytes.Buffer{}

	err := run([]string{"--diff", "--color", "s/hello/goodbye/"}, in, out, io.Discard)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "\033[31m") {
		t.Errorf("expected red ANSI code in colored diff, got:\n%s", got)
	}
	if !strings.Contains(got, "\033[32m") {
		t.Errorf("expected green ANSI code in colored diff, got:\n%s", got)
	}
}
