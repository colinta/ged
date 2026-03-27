package rule

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

// captureStderr redirects os.Stderr to a buffer for the duration of fn.
func captureStderr(fn func()) string {
	r, w, _ := os.Pipe()
	old := os.Stderr
	os.Stderr = w

	fn()

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestTapRule(t *testing.T) {
	t.Run("passes lines through unchanged", func(t *testing.T) {
		r := NewTapRule("")
		ctx := &LineContext{LineNum: 1, }
		got, err := r.Apply("hello", ctx)
		if err != nil {
			t.Fatalf("Apply() error: %v", err)
		}
		if len(got) != 1 || got[0] != "hello" {
			t.Fatalf("got %v, want [hello]", got)
		}
	})

	t.Run("writes lines to stderr", func(t *testing.T) {
		r := NewTapRule("")
		ctx := &LineContext{LineNum: 1, }
		output := captureStderr(func() {
			r.Apply("hello", ctx)
			r.Apply("world", ctx)
		})
		if !strings.Contains(output, "hello\n") || !strings.Contains(output, "world\n") {
			t.Errorf("expected lines on stderr, got: %q", output)
		}
	})

	t.Run("prints header once before first line", func(t *testing.T) {
		r := NewTapRule("Debug")
		ctx := &LineContext{LineNum: 1, }
		output := captureStderr(func() {
			r.Apply("line1", ctx)
			r.Apply("line2", ctx)
		})
		if count := strings.Count(output, "###### Debug ######"); count != 1 {
			t.Errorf("expected header once, got %d times in: %q", count, output)
		}
		// Header should come before line1
		headerIdx := strings.Index(output, "###### Debug ######")
		line1Idx := strings.Index(output, "line1")
		if headerIdx > line1Idx {
			t.Error("header should appear before first line")
		}
	})

	t.Run("no header when empty", func(t *testing.T) {
		r := NewTapRule("")
		ctx := &LineContext{LineNum: 1, }
		output := captureStderr(func() {
			r.Apply("line1", ctx)
		})
		if strings.Contains(output, "######") {
			t.Errorf("expected no header, got: %q", output)
		}
	})

	t.Run("flush prints closing line with header", func(t *testing.T) {
		r := NewTapRule("Test")
		ctx := &LineContext{LineNum: 1, }
		output := captureStderr(func() {
			r.Apply("line1", ctx)
			r.Flush(ctx)
		})
		// closing line is #'s of same length as header line
		expected := fmt.Sprintf("%s\n", strings.Repeat("#", 6+len("Test")+8))
		if !strings.HasSuffix(output, expected) {
			t.Errorf("expected closing line %q, got: %q", expected, output)
		}
	})

	t.Run("flush does nothing without header", func(t *testing.T) {
		r := NewTapRule("")
		ctx := &LineContext{LineNum: 1, }
		output := captureStderr(func() {
			r.Apply("line1", ctx)
			r.Flush(ctx)
		})
		if strings.Contains(output, "######") {
			t.Errorf("expected no decoration, got: %q", output)
		}
	})
}
