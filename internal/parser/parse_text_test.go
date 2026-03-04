package parser

import (
	"testing"

	"github.com/colinta/ged/internal/rule"
)

func TestParseTextModificationRules(t *testing.T) {
	t.Run("trim", func(t *testing.T) {
		r, err := ParseRule("trim")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := r.(*rule.TrimRule); !ok {
			t.Fatalf("expected *rule.TrimRule, got %T", r)
		}
	})

	t.Run("triml", func(t *testing.T) {
		r, err := ParseRule("triml")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := r.(*rule.TrimRule); !ok {
			t.Fatalf("expected *rule.TrimRule, got %T", r)
		}
	})

	t.Run("trimr", func(t *testing.T) {
		r, err := ParseRule("trimr")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := r.(*rule.TrimRule); !ok {
			t.Fatalf("expected *rule.TrimRule, got %T", r)
		}
	})

	t.Run("upper", func(t *testing.T) {
		r, err := ParseRule("upper")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := r.(*rule.UpperRule); !ok {
			t.Fatalf("expected *rule.UpperRule, got %T", r)
		}
	})

	t.Run("lower", func(t *testing.T) {
		r, err := ParseRule("lower")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := r.(*rule.LowerRule); !ok {
			t.Fatalf("expected *rule.LowerRule, got %T", r)
		}
	})

	t.Run("prepend", func(t *testing.T) {
		r, err := ParseRule("prepend/>> /")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := r.(*rule.PrependRule); !ok {
			t.Fatalf("expected *rule.PrependRule, got %T", r)
		}
	})

	t.Run("prepend with different delimiter", func(t *testing.T) {
		r, err := ParseRule("prepend|>> |")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := r.(*rule.PrependRule); !ok {
			t.Fatalf("expected *rule.PrependRule, got %T", r)
		}
	})

	t.Run("prepend missing arg errors", func(t *testing.T) {
		_, err := ParseRule("prepend")
		if err == nil {
			t.Fatal("expected error for bare prepend")
		}
	})

	t.Run("append", func(t *testing.T) {
		r, err := ParseRule("append/;/")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := r.(*rule.AppendRule); !ok {
			t.Fatalf("expected *rule.AppendRule, got %T", r)
		}
	})

	t.Run("append missing arg errors", func(t *testing.T) {
		_, err := ParseRule("append")
		if err == nil {
			t.Fatal("expected error for bare append")
		}
	})

	t.Run("surround", func(t *testing.T) {
		r, err := ParseRule("surround/(/)/")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := r.(*rule.SurroundRule); !ok {
			t.Fatalf("expected *rule.SurroundRule, got %T", r)
		}
	})

	t.Run("surround missing second arg errors", func(t *testing.T) {
		_, err := ParseRule("surround/(//")
		if err == nil {
			t.Fatal("expected error for surround with missing second arg")
		}
	})

	t.Run("surround missing args errors", func(t *testing.T) {
		_, err := ParseRule("surround")
		if err == nil {
			t.Fatal("expected error for bare surround")
		}
	})
}
