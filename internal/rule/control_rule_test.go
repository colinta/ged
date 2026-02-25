package rule

import "testing"

// helper to process multiple lines through a single rule with Setup
func processLines(t *testing.T, r LineRule, lines []string) []string {
	t.Helper()
	ctx := &LineContext{}
	if s, ok := r.(SetupRule); ok {
		s.Setup(ctx)
	}

	var result []string
	for i, line := range lines {
		ctx.LineNum = i + 1
		out, err := r.Apply(line, ctx)
		if err != nil {
			t.Fatalf("Apply error: %v", err)
		}
		if ctx.Printing != PrintOff {
			result = append(result, out...)
		}
	}
	return result
}

func TestOnRule_StartsAtMatch(t *testing.T) {
	r, _ := NewOnRule("start")
	result := processLines(t, r, []string{"before", "start", "after"})
	if len(result) != 2 || result[0] != "start" || result[1] != "after" {
		t.Errorf("expected [start after], got %v", result)
	}
}

func TestOnRule_NoMatch(t *testing.T) {
	r, _ := NewOnRule("start")
	result := processLines(t, r, []string{"a", "b", "c"})
	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

func TestOffRule_StopsAtMatch(t *testing.T) {
	r, _ := NewOffRule("stop")
	result := processLines(t, r, []string{"before", "stop", "after"})
	if len(result) != 1 || result[0] != "before" {
		t.Errorf("expected [before], got %v", result)
	}
}

func TestOffRule_NoMatch(t *testing.T) {
	r, _ := NewOffRule("stop")
	result := processLines(t, r, []string{"a", "b", "c"})
	if len(result) != 3 {
		t.Errorf("expected all 3 lines, got %v", result)
	}
}

func TestAfterRule_StartsAfterMatch(t *testing.T) {
	r, _ := NewAfterRule("marker")
	result := processLines(t, r, []string{"before", "marker", "after1", "after2"})
	if len(result) != 2 || result[0] != "after1" || result[1] != "after2" {
		t.Errorf("expected [after1 after2], got %v", result)
	}
}

func TestAfterRule_NoMatch(t *testing.T) {
	r, _ := NewAfterRule("marker")
	result := processLines(t, r, []string{"a", "b", "c"})
	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

func TestToggleRule_FlipsOnMatch(t *testing.T) {
	r, _ := NewToggleRule("---")
	result := processLines(t, r, []string{"off", "---", "on1", "on2", "---", "off2", "---", "on3"})
	// "off" → Off. "---" → On (printed). "on1","on2" → On. "---" → Off (not printed).
	// "off2" → Off. "---" → On (printed). "on3" → On.
	expected := []string{"---", "on1", "on2", "---", "on3"}
	if len(result) != len(expected) {
		t.Errorf("expected %v, got %v", expected, result)
		return
	}
	for i, want := range expected {
		if result[i] != want {
			t.Errorf("expected %v, got %v", expected, result)
			return
		}
	}
}

func TestToggleRule_MatchLineItself(t *testing.T) {
	// The toggle line itself: first toggle turns ON, so the toggle line is printed
	r, _ := NewToggleRule("---")
	result := processLines(t, r, []string{"---", "a", "---", "b"})
	// "---" toggles ON → printed. "a" → printed. "---" toggles OFF → not printed. "b" → not printed.
	if len(result) != 2 || result[0] != "---" || result[1] != "a" {
		t.Errorf("expected [--- a], got %v", result)
	}
}

func TestOnOff_Combined(t *testing.T) {
	on, _ := NewOnRule("start")
	off, _ := NewOffRule("end")

	ctx := &LineContext{}
	on.Setup(ctx)
	// off.Setup would set PrintOn, but on.Setup already set PrintOff.
	// on runs first so its initial state wins.

	lines := []string{"before", "start", "middle", "end", "after"}
	var result []string
	for i, line := range lines {
		ctx.LineNum = i + 1
		on.Apply(line, ctx)
		off.Apply(line, ctx)
		if ctx.Printing != PrintOff {
			result = append(result, line)
		}
	}

	// "before" → off. "start" → on turns on, off doesn't match → printed.
	// "middle" → on. "end" → off turns off → not printed. "after" → off.
	if len(result) != 2 || result[0] != "start" || result[1] != "middle" {
		t.Errorf("expected [start middle], got %v", result)
	}
}
