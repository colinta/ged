package rule

import "testing"

func TestPrintLineNumRule_SingleLine(t *testing.T) {
	lineRange, _ := ParseLineRange("2")
	rule := NewPrintLineNumRule(lineRange)

	// Line 1 should be deleted
	result, _ := rule.Apply("first", &LineContext{LineNum: 1})
	if len(result) != 0 {
		t.Errorf("line 1 should be deleted, got %v", result)
	}

	// Line 2 should be kept
	result, _ = rule.Apply("second", &LineContext{LineNum: 2})
	if len(result) != 1 || result[0] != "second" {
		t.Errorf("line 2 should be kept, got %v", result)
	}

	// Line 3 should be deleted
	result, _ = rule.Apply("third", &LineContext{LineNum: 3})
	if len(result) != 0 {
		t.Errorf("line 3 should be deleted, got %v", result)
	}
}

func TestPrintLineNumRule_Range(t *testing.T) {
	lineRange, _ := ParseLineRange("2-4")
	rule := NewPrintLineNumRule(lineRange)

	// Line 1 should be deleted
	result, _ := rule.Apply("one", &LineContext{LineNum: 1})
	if len(result) != 0 {
		t.Errorf("line 1 should be deleted")
	}

	// Lines 2, 3, 4 should be kept
	for i := 2; i <= 4; i++ {
		result, _ = rule.Apply("content", &LineContext{LineNum: i})
		if len(result) != 1 {
			t.Errorf("line %d should be kept", &LineContext{LineNum: i})
		}
	}

	// Line 5 should be deleted
	result, _ = rule.Apply("five", &LineContext{LineNum: 5})
	if len(result) != 0 {
		t.Errorf("line 5 should be deleted")
	}
}

func TestPrintLineNumRule_OpenRange(t *testing.T) {
	// "3-" means from line 3 onwards
	lineRange, _ := ParseLineRange("3-")
	rule := NewPrintLineNumRule(lineRange)

	result, _ := rule.Apply("one", &LineContext{LineNum: 1})
	if len(result) != 0 {
		t.Errorf("line 1 should be deleted")
	}

	result, _ = rule.Apply("two", &LineContext{LineNum: 2})
	if len(result) != 0 {
		t.Errorf("line 2 should be deleted")
	}

	result, _ = rule.Apply("three", &LineContext{LineNum: 3})
	if len(result) != 1 {
		t.Errorf("line 3 should be kept")
	}

	result, _ = rule.Apply("hundred", &LineContext{LineNum: 100})
	if len(result) != 1 {
		t.Errorf("line 100 should be kept")
	}
}
