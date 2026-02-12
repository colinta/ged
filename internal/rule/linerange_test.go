package rule

import "testing"

func TestSingleLine_Contains(t *testing.T) {
	r := SingleLine(5)

	if !r.Contains(5) {
		t.Error("SingleLine(5) should contain 5")
	}
	if r.Contains(4) {
		t.Error("SingleLine(5) should not contain 4")
	}
	if r.Contains(6) {
		t.Error("SingleLine(5) should not contain 6")
	}
}

func TestLineRange_Contains(t *testing.T) {
	r := &Range{Start: 2, End: 4}

	if r.Contains(1) {
		t.Error("Range{2,4} should not contain 1")
	}
	if !r.Contains(2) {
		t.Error("Range{2,4} should contain 2")
	}
	if !r.Contains(3) {
		t.Error("Range{2,4} should contain 3")
	}
	if !r.Contains(4) {
		t.Error("Range{2,4} should contain 4")
	}
	if r.Contains(5) {
		t.Error("Range{2,4} should not contain 5")
	}
}

func TestOpenRangeFrom_Contains(t *testing.T) {
	// 5- means "from line 5 onwards"
	r := &OpenRange{From: 5, ToEnd: true}

	if r.Contains(4) {
		t.Error("OpenRange{From:5} should not contain 4")
	}
	if !r.Contains(5) {
		t.Error("OpenRange{From:5} should contain 5")
	}
	if !r.Contains(100) {
		t.Error("OpenRange{From:5} should contain 100")
	}
}

func TestOpenRangeTo_Contains(t *testing.T) {
	// -5 means "up to line 5"
	r := &OpenRange{To: 5, ToEnd: false}

	if !r.Contains(1) {
		t.Error("OpenRange{To:5} should contain 1")
	}
	if !r.Contains(5) {
		t.Error("OpenRange{To:5} should contain 5")
	}
	if r.Contains(6) {
		t.Error("OpenRange{To:5} should not contain 6")
	}
}

func TestParseLineRange_Single(t *testing.T) {
	r, err := ParseLineRange("5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !r.Contains(5) {
		t.Error("parsed '5' should contain 5")
	}
	if r.Contains(4) {
		t.Error("parsed '5' should not contain 4")
	}
}

func TestParseLineRange_Range(t *testing.T) {
	r, err := ParseLineRange("2-4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if r.Contains(1) {
		t.Error("parsed '2-4' should not contain 1")
	}
	if !r.Contains(2) {
		t.Error("parsed '2-4' should contain 2")
	}
	if !r.Contains(3) {
		t.Error("parsed '2-4' should contain 3")
	}
	if !r.Contains(4) {
		t.Error("parsed '2-4' should contain 4")
	}
	if r.Contains(5) {
		t.Error("parsed '2-4' should not contain 5")
	}
}

func TestParseLineRange_OpenFrom(t *testing.T) {
	r, err := ParseLineRange("5-")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if r.Contains(4) {
		t.Error("parsed '5-' should not contain 4")
	}
	if !r.Contains(5) {
		t.Error("parsed '5-' should contain 5")
	}
	if !r.Contains(100) {
		t.Error("parsed '5-' should contain 100")
	}
}

func TestParseLineRange_OpenTo(t *testing.T) {
	r, err := ParseLineRange("-5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !r.Contains(1) {
		t.Error("parsed '-5' should contain 1")
	}
	if !r.Contains(5) {
		t.Error("parsed '-5' should contain 5")
	}
	if r.Contains(6) {
		t.Error("parsed '-5' should not contain 6")
	}
}

func TestParseLineRange_Composite(t *testing.T) {
	// "1,3,5-7" means lines 1, 3, 5, 6, 7
	r, err := ParseLineRange("1,3,5-7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !r.Contains(1) {
		t.Error("parsed '1,3,5-7' should contain 1")
	}
	if r.Contains(2) {
		t.Error("parsed '1,3,5-7' should not contain 2")
	}
	if !r.Contains(3) {
		t.Error("parsed '1,3,5-7' should contain 3")
	}
	if r.Contains(4) {
		t.Error("parsed '1,3,5-7' should not contain 4")
	}
	if !r.Contains(5) {
		t.Error("parsed '1,3,5-7' should contain 5")
	}
	if !r.Contains(6) {
		t.Error("parsed '1,3,5-7' should contain 6")
	}
	if !r.Contains(7) {
		t.Error("parsed '1,3,5-7' should contain 7")
	}
	if r.Contains(8) {
		t.Error("parsed '1,3,5-7' should not contain 8")
	}
}

func TestParseLineRange_Invalid(t *testing.T) {
	_, err := ParseLineRange("abc")
	if err == nil {
		t.Error("expected error for 'abc'")
	}

	_, err = ParseLineRange("")
	if err == nil {
		t.Error("expected error for empty string")
	}
}
