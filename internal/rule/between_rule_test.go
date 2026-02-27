package rule

import (
	"strings"
	"testing"

	"github.com/dlclark/regexp2"
)

func TestBetweenLineRule_AppliesInsideRange(t *testing.T) {
	r := NewBetweenLineRule(
		regexp2.MustCompile("START", 0),
		regexp2.MustCompile("END", 0),
		false,
		[]LineRule{mustSub(t, "x", "X")},
	)
	lines := []string{"before x", "START x", "middle x", "END x", "after x"}
	got := applyBetweenLine(t, r, lines)
	want := "before x\nSTART X\nmiddle X\nEND X\nafter x"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBetweenLineRule_StartLineIncluded(t *testing.T) {
	r := NewBetweenLineRule(
		regexp2.MustCompile("START", 0),
		regexp2.MustCompile("END", 0),
		false,
		[]LineRule{mustSub(t, "^", "> ")},
	)
	lines := []string{"before", "START", "middle", "END", "after"}
	got := applyBetweenLine(t, r, lines)
	want := "before\n> START\n> middle\n> END\nafter"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBetweenLineRule_EndLineIncluded(t *testing.T) {
	r := NewBetweenLineRule(
		regexp2.MustCompile("START", 0),
		regexp2.MustCompile("END", 0),
		false,
		[]LineRule{mustSub(t, "o", "0")},
	)
	lines := []string{"foo", "START foo", "END foo", "foo"}
	got := applyBetweenLine(t, r, lines)
	want := "foo\nSTART f00\nEND f00\nfoo"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBetweenLineRule_NoMatchPassesThrough(t *testing.T) {
	r := NewBetweenLineRule(
		regexp2.MustCompile("START", 0),
		regexp2.MustCompile("END", 0),
		false,
		[]LineRule{mustSub(t, "x", "X")},
	)
	lines := []string{"a x", "b x", "c x"}
	got := applyBetweenLine(t, r, lines)
	want := "a x\nb x\nc x"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBetweenLineRule_MultipleRanges(t *testing.T) {
	r := NewBetweenLineRule(
		regexp2.MustCompile("START", 0),
		regexp2.MustCompile("END", 0),
		false,
		[]LineRule{mustSub(t, "x", "X")},
	)
	lines := []string{"x", "START x", "END x", "x", "START x", "END x", "x"}
	got := applyBetweenLine(t, r, lines)
	want := "x\nSTART X\nEND X\nx\nSTART X\nEND X\nx"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBetweenLineRule_Inverted(t *testing.T) {
	r := NewBetweenLineRule(
		regexp2.MustCompile("START", 0),
		regexp2.MustCompile("END", 0),
		true,
		[]LineRule{mustSub(t, "x", "X")},
	)
	lines := []string{"x", "START x", "middle x", "END x", "x"}
	got := applyBetweenLine(t, r, lines)
	want := "X\nSTART x\nmiddle x\nEND x\nX"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBetweenLineRule_InnerDeleteRemovesLine(t *testing.T) {
	del, _ := NewDeleteLineRule("middle")
	r := NewBetweenLineRule(
		regexp2.MustCompile("START", 0),
		regexp2.MustCompile("END", 0),
		false,
		[]LineRule{del},
	)
	lines := []string{"before", "START", "middle", "keep", "END", "after"}
	got := applyBetweenLine(t, r, lines)
	want := "before\nSTART\nkeep\nEND\nafter"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBetweenLineRule_UnclosedRangeGoesToEnd(t *testing.T) {
	r := NewBetweenLineRule(
		regexp2.MustCompile("START", 0),
		regexp2.MustCompile("END", 0),
		false,
		[]LineRule{mustSub(t, "x", "X")},
	)
	lines := []string{"x", "START x", "x", "x"}
	got := applyBetweenLine(t, r, lines)
	want := "x\nSTART X\nX\nX"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBetweenDocRule_SortsInsideRange(t *testing.T) {
	r := NewBetweenDocRule(
		regexp2.MustCompile("START", 0),
		regexp2.MustCompile("END", 0),
		false,
		[]DocumentRule{NewSortRule()},
	)
	lines := []string{"before", "START", "c", "a", "b", "END", "after"}
	got, err := r.ApplyDocument(lines)
	if err != nil {
		t.Fatal(err)
	}
	want := "before\nEND\nSTART\na\nb\nc\nafter"
	if strings.Join(got, "\n") != want {
		t.Errorf("got %q, want %q", strings.Join(got, "\n"), want)
	}
}

func TestBetweenDocRule_Inverted(t *testing.T) {
	r := NewBetweenDocRule(
		regexp2.MustCompile("START", 0),
		regexp2.MustCompile("END", 0),
		true,
		[]DocumentRule{NewSortRule()},
	)
	lines := []string{"c", "a", "START", "middle", "END", "b", "d"}
	got, err := r.ApplyDocument(lines)
	if err != nil {
		t.Fatal(err)
	}
	want := "a\nb\nSTART\nmiddle\nEND\nc\nd"
	if strings.Join(got, "\n") != want {
		t.Errorf("got %q, want %q", strings.Join(got, "\n"), want)
	}
}

// helpers

func mustSub(t *testing.T, pattern, replace string) *SubstitutionRule {
	t.Helper()
	r, err := NewSubstitutionRule(pattern, replace, WithGlobal())
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func applyBetweenLine(t *testing.T, r *BetweenLineRule, lines []string) string {
	t.Helper()
	var result []string
	ctx := &LineContext{}
	for i, line := range lines {
		ctx.LineNum = i + 1
		out, err := r.Apply(line, ctx)
		if err != nil {
			t.Fatal(err)
		}
		result = append(result, out...)
	}
	return strings.Join(result, "\n")
}
