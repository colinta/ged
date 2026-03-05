package rule

import (
	"strings"
	"testing"

	"github.com/dlclark/regexp2"
)

func TestConditionalLineRule_ElseApplied(t *testing.T) {
	// if/hello/ { upper } else { lower }
	cond := NewConditionalLineRule(
		regexp2.MustCompile("hello", 0),
		false,
		[]LineRule{NewUpperRule()},
		[]LineRule{NewLowerRule()},
	)
	ctx := &LineContext{LineNum: 1}

	// Matching line → upper
	got, err := cond.Apply("hello world", ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got[0] != "HELLO WORLD" {
		t.Errorf("matching: got %q, want %q", got[0], "HELLO WORLD")
	}

	// Non-matching line → lower
	got, err = cond.Apply("GOODBYE WORLD", ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got[0] != "goodbye world" {
		t.Errorf("non-matching: got %q, want %q", got[0], "goodbye world")
	}
}

func TestConditionalDocRule_ElseApplied(t *testing.T) {
	// if/hello/ { upper } else { lower }
	cond := NewConditionalDocRule(
		regexp2.MustCompile("hello", 0),
		false,
		[]DocumentRule{NewApplyAllRule([]LineRule{NewUpperRule()})},
		[]DocumentRule{NewApplyAllRule([]LineRule{NewLowerRule()})},
	)
	lines := []string{"hello World", "GOODBYE World"}
	got, err := cond.ApplyDocument(lines)
	if err != nil {
		t.Fatal(err)
	}
	// "hello World" matches → upper → "HELLO WORLD"
	// "GOODBYE World" doesn't match → lower → "goodbye world"
	want := "HELLO WORLD\ngoodbye world"
	if strings.Join(got, "\n") != want {
		t.Errorf("got %q, want %q", strings.Join(got, "\n"), want)
	}
}

func TestBetweenLineRule_ElseApplied(t *testing.T) {
	// between/START/END/ { upper } else { lower }
	start := regexp2.MustCompile("START", 0)
	end := regexp2.MustCompile("END", 0)
	r := NewBetweenLineRule(start, end, false,
		[]LineRule{NewUpperRule()},
		[]LineRule{NewLowerRule()},
	)
	lines := []string{"BEFORE", "START", "middle", "END", "AFTER"}
	got := applyBetweenLine(t, r, lines)
	want := "before\nSTART\nMIDDLE\nEND\nafter"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBetweenDocRule_ElseApplied(t *testing.T) {
	// between/START/END/ { sort } else { reverse }
	start := regexp2.MustCompile("START", 0)
	end := regexp2.MustCompile("END", 0)
	r := NewBetweenDocRule(start, end, false,
		[]DocumentRule{NewSortRule()},
		[]DocumentRule{NewReverseRule()},
	)
	lines := []string{"c", "a", "START", "z", "m", "END", "b", "d"}
	got, err := r.ApplyDocument(lines)
	if err != nil {
		t.Fatal(err)
	}
	// Outside before: ["c", "a"] → reversed → ["a", "c"]
	// Inside: ["START", "z", "m", "END"] → sorted → ["END", "START", "m", "z"]
	// Outside after: ["b", "d"] → reversed → ["d", "b"]
	want := "a\nc\nEND\nSTART\nm\nz\nd\nb"
	if strings.Join(got, "\n") != want {
		t.Errorf("got %q, want %q", strings.Join(got, "\n"), want)
	}
}
