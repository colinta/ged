package rule

import (
	"regexp"
	"testing"
)

// --- ConditionalLineRule tests ---

func TestConditionalLineRule_MatchingLine(t *testing.T) {
	sub, _ := NewSubstitutionRule("o", "x")
	cond := NewConditionalLineRule(regexp.MustCompile("hello"), false, []LineRule{sub})

	result, err := cond.Apply("hello world", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0] != "hellx world" {
		t.Errorf("got %v, want [hellx world]", result)
	}
}

func TestConditionalLineRule_NonMatchingLine(t *testing.T) {
	sub, _ := NewSubstitutionRule("o", "x")
	cond := NewConditionalLineRule(regexp.MustCompile("hello"), false, []LineRule{sub})

	result, err := cond.Apply("goodbye world", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0] != "goodbye world" {
		t.Errorf("got %v, want [goodbye world]", result)
	}
}

func TestConditionalLineRule_Inverted(t *testing.T) {
	sub, _ := NewSubstitutionRule("o", "x")
	cond := NewConditionalLineRule(regexp.MustCompile("hello"), true, []LineRule{sub})

	// "hello" matches the pattern, so inverted means rules DON'T apply
	result, err := cond.Apply("hello world", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0] != "hello world" {
		t.Errorf("got %v, want [hello world]", result)
	}

	// "goodbye" doesn't match, so inverted means rules DO apply
	// First match only: first "o" in "goodbye" → "x"
	result, err = cond.Apply("goodbye world", 1)
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
	cond := NewConditionalLineRule(regexp.MustCompile("x"), false, []LineRule{sub1, sub2})

	result, err := cond.Apply("xab", 1)
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
	cond := NewConditionalLineRule(regexp.MustCompile("hello"), false, []LineRule{del})

	result, err := cond.Apply("hello world", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("got %v, want empty", result)
	}
}

func TestConditionalLineRule_PassesLineNum(t *testing.T) {
	lineNumRule := NewPrintLineNumRule(SingleLine(3))
	cond := NewConditionalLineRule(regexp.MustCompile(".*"), false, []LineRule{lineNumRule})

	// Line 3 should be kept
	result, err := cond.Apply("hello", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0] != "hello" {
		t.Errorf("line 3: got %v, want [hello]", result)
	}

	// Line 2 should be filtered out by PrintLineNumRule
	result, err = cond.Apply("hello", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("line 2: got %v, want empty", result)
	}
}

// --- ConditionalDocRule tests ---

func TestConditionalDocRule_SortMatchingLines(t *testing.T) {
	// Sort only lines matching "item"
	cond := NewConditionalDocRule(
		regexp.MustCompile("item"),
		false,
		[]DocumentRule{NewSortRule()},
	)

	input := []string{"header", "item c", "item a", "footer", "item b"}
	result, err := cond.ApplyDocument(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"header", "item a", "item b", "footer", "item c"}
	if len(result) != len(want) {
		t.Fatalf("got %v, want %v", result, want)
	}
	for i := range want {
		if result[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, result[i], want[i])
		}
	}
}

func TestConditionalDocRule_ReverseMatchingLines(t *testing.T) {
	cond := NewConditionalDocRule(
		regexp.MustCompile("x"),
		false,
		[]DocumentRule{NewReverseRule()},
	)

	input := []string{"a", "x1", "b", "x2", "x3"}
	result, err := cond.ApplyDocument(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// x1, x2, x3 reversed → x3, x2, x1, placed at original matching positions
	want := []string{"a", "x3", "b", "x2", "x1"}
	if len(result) != len(want) {
		t.Fatalf("got %v, want %v", result, want)
	}
	for i := range want {
		if result[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, result[i], want[i])
		}
	}
}

func TestConditionalDocRule_JoinMatchingLines(t *testing.T) {
	// Join reduces matching lines to one — extra matching positions are consumed
	cond := NewConditionalDocRule(
		regexp.MustCompile("item"),
		false,
		[]DocumentRule{NewJoinRule(",")},
	)

	input := []string{"header", "item a", "item b", "item c", "footer"}
	result, err := cond.ApplyDocument(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Joined result fills first matching position, others consumed
	want := []string{"header", "item a,item b,item c", "footer"}
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
	// Sort non-matching lines, leave matching ones in place
	cond := NewConditionalDocRule(
		regexp.MustCompile("KEEP"),
		true,
		[]DocumentRule{NewSortRule()},
	)

	input := []string{"c", "KEEP", "a", "b"}
	result, err := cond.ApplyDocument(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Non-matching: c, a, b → sorted: a, b, c placed at positions 0, 2, 3
	want := []string{"a", "KEEP", "b", "c"}
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
		regexp.MustCompile("NOMATCH"),
		false,
		[]DocumentRule{NewSortRule()},
	)

	input := []string{"c", "a", "b"}
	result, err := cond.ApplyDocument(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Nothing matches, everything passes through
	if len(result) != 3 || result[0] != "c" || result[1] != "a" || result[2] != "b" {
		t.Errorf("got %v, want [c a b]", result)
	}
}

func TestConditionalDocRule_SubThenSort(t *testing.T) {
	// Mix of line rules (wrapped in ApplyAllRule) and document rules
	sub, _ := NewSubstitutionRule("item ", "")
	cond := NewConditionalDocRule(
		regexp.MustCompile("item"),
		false,
		[]DocumentRule{NewApplyAllRule([]LineRule{sub}), NewSortRule()},
	)

	input := []string{"header", "item c", "item a", "footer", "item b"}
	result, err := cond.ApplyDocument(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Strip "item ", then sort: a, b, c
	want := []string{"header", "a", "b", "footer", "c"}
	if len(result) != len(want) {
		t.Fatalf("got %v, want %v", result, want)
	}
	for i := range want {
		if result[i] != want[i] {
			t.Errorf("index %d: got %q, want %q", i, result[i], want[i])
		}
	}
}
