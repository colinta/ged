package rule

import (
	"strings"
	"testing"

	"github.com/dlclark/regexp2"
)

func TestIfAnyDocRule_AppliesWhenAnyMatches(t *testing.T) {
	// If any line contains "error", uppercase all lines
	r := NewIfAnyDocRule(
		regexp2.MustCompile("error", 0),
		false,
		[]DocumentRule{NewApplyAllRule([]LineRule{mustUpper(t)})},
		nil,
	)
	lines := []string{"ok", "error here", "fine"}
	got, err := r.ApplyDocument(lines)
	if err != nil {
		t.Fatal(err)
	}
	want := "OK\nERROR HERE\nFINE"
	if strings.Join(got, "\n") != want {
		t.Errorf("got %q, want %q", strings.Join(got, "\n"), want)
	}
}

func TestIfAnyDocRule_PassesThroughWhenNoMatch(t *testing.T) {
	r := NewIfAnyDocRule(
		regexp2.MustCompile("error", 0),
		false,
		[]DocumentRule{NewApplyAllRule([]LineRule{mustUpper(t)})},
		nil,
	)
	lines := []string{"ok", "fine", "good"}
	got, err := r.ApplyDocument(lines)
	if err != nil {
		t.Fatal(err)
	}
	want := "ok\nfine\ngood"
	if strings.Join(got, "\n") != want {
		t.Errorf("got %q, want %q", strings.Join(got, "\n"), want)
	}
}

func TestIfAnyDocRule_Inverted(t *testing.T) {
	// !ifany/error/ — applies rules when NO line matches "error"
	r := NewIfAnyDocRule(
		regexp2.MustCompile("error", 0),
		true,
		[]DocumentRule{NewApplyAllRule([]LineRule{mustUpper(t)})},
		nil,
	)
	lines := []string{"ok", "fine", "good"}
	got, err := r.ApplyDocument(lines)
	if err != nil {
		t.Fatal(err)
	}
	want := "OK\nFINE\nGOOD"
	if strings.Join(got, "\n") != want {
		t.Errorf("got %q, want %q", strings.Join(got, "\n"), want)
	}
}

func TestIfAnyDocRule_InvertedWithMatch(t *testing.T) {
	// !ifany/error/ with a match — should pass through
	r := NewIfAnyDocRule(
		regexp2.MustCompile("error", 0),
		true,
		[]DocumentRule{NewApplyAllRule([]LineRule{mustUpper(t)})},
		nil,
	)
	lines := []string{"ok", "error", "good"}
	got, err := r.ApplyDocument(lines)
	if err != nil {
		t.Fatal(err)
	}
	want := "ok\nerror\ngood"
	if strings.Join(got, "\n") != want {
		t.Errorf("got %q, want %q", strings.Join(got, "\n"), want)
	}
}

func TestIfAnyDocRule_ElseApplied(t *testing.T) {
	// ifany/error/ { upper } else { lower }
	r := NewIfAnyDocRule(
		regexp2.MustCompile("error", 0),
		false,
		[]DocumentRule{NewApplyAllRule([]LineRule{mustUpper(t)})},
		[]DocumentRule{NewApplyAllRule([]LineRule{NewLowerRule()})},
	)
	// No match — else applies
	lines := []string{"OK", "FINE"}
	got, err := r.ApplyDocument(lines)
	if err != nil {
		t.Fatal(err)
	}
	want := "ok\nfine"
	if strings.Join(got, "\n") != want {
		t.Errorf("got %q, want %q", strings.Join(got, "\n"), want)
	}
}

func TestIfAnyDocRule_ElseNotAppliedOnMatch(t *testing.T) {
	r := NewIfAnyDocRule(
		regexp2.MustCompile("error", 0),
		false,
		[]DocumentRule{NewApplyAllRule([]LineRule{mustUpper(t)})},
		[]DocumentRule{NewApplyAllRule([]LineRule{NewLowerRule()})},
	)
	// Has match — main rules apply
	lines := []string{"error", "Fine"}
	got, err := r.ApplyDocument(lines)
	if err != nil {
		t.Fatal(err)
	}
	want := "ERROR\nFINE"
	if strings.Join(got, "\n") != want {
		t.Errorf("got %q, want %q", strings.Join(got, "\n"), want)
	}
}

func TestIfNoneDocRule_AppliesWhenNoMatch(t *testing.T) {
	r := NewIfNoneDocRule(
		regexp2.MustCompile("error", 0),
		false,
		[]DocumentRule{NewApplyAllRule([]LineRule{mustUpper(t)})},
		nil,
	)
	lines := []string{"ok", "fine", "good"}
	got, err := r.ApplyDocument(lines)
	if err != nil {
		t.Fatal(err)
	}
	want := "OK\nFINE\nGOOD"
	if strings.Join(got, "\n") != want {
		t.Errorf("got %q, want %q", strings.Join(got, "\n"), want)
	}
}

func TestIfNoneDocRule_PassesThroughWhenMatch(t *testing.T) {
	r := NewIfNoneDocRule(
		regexp2.MustCompile("error", 0),
		false,
		[]DocumentRule{NewApplyAllRule([]LineRule{mustUpper(t)})},
		nil,
	)
	lines := []string{"ok", "error here", "fine"}
	got, err := r.ApplyDocument(lines)
	if err != nil {
		t.Fatal(err)
	}
	want := "ok\nerror here\nfine"
	if strings.Join(got, "\n") != want {
		t.Errorf("got %q, want %q", strings.Join(got, "\n"), want)
	}
}

func TestIfNoneDocRule_Inverted(t *testing.T) {
	// !ifnone/error/ — applies when SOME line matches "error"
	r := NewIfNoneDocRule(
		regexp2.MustCompile("error", 0),
		true,
		[]DocumentRule{NewApplyAllRule([]LineRule{mustUpper(t)})},
		nil,
	)
	lines := []string{"ok", "error", "good"}
	got, err := r.ApplyDocument(lines)
	if err != nil {
		t.Fatal(err)
	}
	want := "OK\nERROR\nGOOD"
	if strings.Join(got, "\n") != want {
		t.Errorf("got %q, want %q", strings.Join(got, "\n"), want)
	}
}

func TestIfNoneDocRule_ElseApplied(t *testing.T) {
	// ifnone/error/ { upper } else { lower }
	r := NewIfNoneDocRule(
		regexp2.MustCompile("error", 0),
		false,
		[]DocumentRule{NewApplyAllRule([]LineRule{mustUpper(t)})},
		[]DocumentRule{NewApplyAllRule([]LineRule{NewLowerRule()})},
	)
	// Has match — else applies (lowercases everything)
	lines := []string{"OK", "error"}
	got, err := r.ApplyDocument(lines)
	if err != nil {
		t.Fatal(err)
	}
	want := "ok\nerror"
	if strings.Join(got, "\n") != want {
		t.Errorf("got %q, want %q", strings.Join(got, "\n"), want)
	}
}

func TestIfNoneDocRule_ElseNotAppliedOnNoMatch(t *testing.T) {
	r := NewIfNoneDocRule(
		regexp2.MustCompile("error", 0),
		false,
		[]DocumentRule{NewApplyAllRule([]LineRule{mustUpper(t)})},
		[]DocumentRule{NewApplyAllRule([]LineRule{NewLowerRule()})},
	)
	// No match — main rules apply
	lines := []string{"ok", "fine"}
	got, err := r.ApplyDocument(lines)
	if err != nil {
		t.Fatal(err)
	}
	want := "OK\nFINE"
	if strings.Join(got, "\n") != want {
		t.Errorf("got %q, want %q", strings.Join(got, "\n"), want)
	}
}

// Helper to create UpperRule
func mustUpper(t *testing.T) *UpperRule {
	t.Helper()
	return NewUpperRule()
}
