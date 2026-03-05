package rule

import (
	"testing"

	"github.com/dlclark/regexp2"
)

// --- ConditionalLineRule tests ---

func TestConditionalLineRule_MatchingLine(t *testing.T) {
	sub, _ := NewSubstitutionRule("o", "x")
	cond := NewConditionalLineRule(regexp2.MustCompile("hello", 0), false, []LineRule{sub}, nil)

	result, err := cond.Apply("hello world", &LineContext{LineNum: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0] != "hellx world" {
		t.Errorf("got %v, want [hellx world]", result)
	}
}

func TestConditionalLineRule_NonMatchingLine(t *testing.T) {
	sub, _ := NewSubstitutionRule("o", "x")
	cond := NewConditionalLineRule(regexp2.MustCompile("hello", 0), false, []LineRule{sub}, nil)

	result, err := cond.Apply("goodbye world", &LineContext{LineNum: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0] != "goodbye world" {
		t.Errorf("got %v, want [goodbye world]", result)
	}
}

func TestConditionalLineRule_Inverted(t *testing.T) {
	sub, _ := NewSubstitutionRule("o", "x")
	cond := NewConditionalLineRule(regexp2.MustCompile("hello", 0), true, []LineRule{sub}, nil)

	// "hello" matches the pattern, so inverted means rules DON'T apply
	result, err := cond.Apply("hello world", &LineContext{LineNum: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0] != "hello world" {
		t.Errorf("got %v, want [hello world]", result)
	}

	// "goodbye" doesn't match, so inverted means rules DO apply
	// First match only: first "o" in "goodbye" → "x"
	result, err = cond.Apply("goodbye world", &LineContext{LineNum: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0] != "gxodbye world" {
		t.Errorf("got %v, want [gxodbye world]", result)
	}
}

func TestConditionalLineRule_MultipleInnerRules(t *testing.T) {
	sub1, _ := NewSubstitutionRule("a", "b")
	sub2, _ := NewSubstitutionRule("b", "c")
	cond := NewConditionalLineRule(regexp2.MustCompile("x", 0), false, []LineRule{sub1, sub2}, nil)

	result, err := cond.Apply("xab", &LineContext{LineNum: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// sub1: first "a" → "b" → "xbb"
	// sub2: first "b" → "c" → "xcb"
	if len(result) != 1 || result[0] != "xcb" {
		t.Errorf("got %v, want [xcb]", result)
	}
}

func TestConditionalLineRule_InnerDeleteRemovesLine(t *testing.T) {
	del, _ := NewDeleteLineRule("hello")
	cond := NewConditionalLineRule(regexp2.MustCompile("hello", 0), false, []LineRule{del}, nil)

	result, err := cond.Apply("hello world", &LineContext{LineNum: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("got %v, want empty", result)
	}
}

func TestConditionalLineRule_PassesLineNum(t *testing.T) {
	lineNumRule := NewPrintLineNumRule(SingleLine(3))
	cond := NewConditionalLineRule(regexp2.MustCompile(".*", 0), false, []LineRule{lineNumRule}, nil)

	// Line 3 should be kept
	result, err := cond.Apply("hello", &LineContext{LineNum: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0] != "hello" {
		t.Errorf("line 3: got %v, want [hello]", result)
	}

	// Line 2 should be filtered out by PrintLineNumRule
	result, err = cond.Apply("hello", &LineContext{LineNum: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("line 2: got %v, want empty", result)
	}
}

// --- ConditionalDocRule tests ---

func TestConditionalDocRule_PerLineProcessing(t *testing.T) {
	// Each matching line is processed independently as its own mini-document.
	// Sort on a single line is a no-op — lines stay in original order.
	cond := NewConditionalDocRule(
		regexp2.MustCompile("item", 0),
		false,
		[]DocumentRule{NewSortRule()},
		nil,
	)

	input := []string{"header", "item c", "item a", "footer", "item b"}
	result, err := cond.ApplyDocument(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Sort on single-line docs is a no-op, order unchanged
	want := []string{"header", "item c", "item a", "footer", "item b"}
	if len(result) != len(want) {
		t.Fatalf("got %v, want %v", result, want)
	}
	for i := range want {
		if result[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, result[i], want[i])
		}
	}
}

func TestConditionalDocRule_BeginPerLine(t *testing.T) {
	// begin/end/border apply per matching line, expanding each independently
	cond := NewConditionalDocRule(
		regexp2.MustCompile("item", 0),
		false,
		[]DocumentRule{NewBeginRule(">>")},
		nil,
	)

	input := []string{"header", "item a", "footer", "item b"}
	result, err := cond.ApplyDocument(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"header", ">>", "item a", "footer", ">>", "item b"}
	if len(result) != len(want) {
		t.Fatalf("got %v, want %v", result, want)
	}
	for i := range want {
		if result[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, result[i], want[i])
		}
	}
}

func TestConditionalDocRule_BorderPerLine(t *testing.T) {
	cond := NewConditionalDocRule(
		regexp2.MustCompile("x", 0),
		false,
		[]DocumentRule{NewBorderRule("---")},
		nil,
	)

	input := []string{"a", "x1", "b", "x2"}
	result, err := cond.ApplyDocument(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"a", "---", "x1", "---", "b", "---", "x2", "---"}
	if len(result) != len(want) {
		t.Fatalf("got %v, want %v", result, want)
	}
	for i := range want {
		if result[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, result[i], want[i])
		}
	}
}

func TestConditionalDocRule_Inverted(t *testing.T) {
	// Inverted: apply rules to each non-matching line independently
	cond := NewConditionalDocRule(
		regexp2.MustCompile("KEEP", 0),
		true,
		[]DocumentRule{NewApplyAllRule([]LineRule{mustSubRule("^", ">> ")})},
		nil,
	)

	input := []string{"c", "KEEP", "a", "b"}
	result, err := cond.ApplyDocument(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{">> c", "KEEP", ">> a", ">> b"}
	if len(result) != len(want) {
		t.Fatalf("got %v, want %v", result, want)
	}
	for i := range want {
		if result[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, result[i], want[i])
		}
	}
}

func TestConditionalDocRule_NoMatches(t *testing.T) {
	cond := NewConditionalDocRule(
		regexp2.MustCompile("NOMATCH", 0),
		false,
		[]DocumentRule{NewBeginRule("header")},
		nil,
	)

	input := []string{"c", "a", "b"}
	result, err := cond.ApplyDocument(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Nothing matches, everything passes through unchanged
	if len(result) != 3 || result[0] != "c" || result[1] != "a" || result[2] != "b" {
		t.Errorf("got %v, want [c a b]", result)
	}
}

func TestConditionalDocRule_SubPerLine(t *testing.T) {
	// Line rules wrapped in ApplyAllRule still work per-line
	sub, _ := NewSubstitutionRule("item ", "")
	cond := NewConditionalDocRule(
		regexp2.MustCompile("item", 0),
		false,
		[]DocumentRule{NewApplyAllRule([]LineRule{sub})},
		nil,
	)

	input := []string{"header", "item c", "item a", "footer", "item b"}
	result, err := cond.ApplyDocument(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"header", "c", "a", "footer", "b"}
	if len(result) != len(want) {
		t.Fatalf("got %v, want %v", result, want)
	}
	for i := range want {
		if result[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, result[i], want[i])
		}
	}
}

// mustSubRule is a test helper for creating substitution rules.
func mustSubRule(pattern, replace string) *SubstitutionRule {
	r, _ := NewSubstitutionRule(pattern, replace)
	return r
}
