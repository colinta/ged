package rule

import "testing"

func TestGroupRule(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		group   int
		opts    []RuleOption
		input   string
		want    []string
	}{
		{
			name:    "extracts group 1",
			pattern: `(\w+) (\w+)`,
			group:   1,
			input:   "hello world",
			want:    []string{"hello"},
		},
		{
			name:    "extracts group 2",
			pattern: `(\w+) (\w+)`,
			group:   2,
			input:   "hello world",
			want:    []string{"world"},
		},
		{
			name:    "no match passes through",
			pattern: `(\d+)`,
			group:   1,
			input:   "no numbers",
			want:    []string{"no numbers"},
		},
		{
			name:    "group out of range passes through",
			pattern: `(\w+)`,
			group:   5,
			input:   "hello",
			want:    []string{"hello"},
		},
		{
			name:    "optional group that didnt participate",
			pattern: `(\w+)(?:\s+(\d+))?`,
			group:   2,
			input:   "hello",
			want:    []string{"hello"},
		},
		{
			name:    "case insensitive",
			pattern: `(hello) (world)`,
			group:   1,
			opts:    []RuleOption{WithIgnoreCase()},
			input:   "HELLO WORLD",
			want:    []string{"HELLO"},
		},
		{
			name:    "extracts from middle of line",
			pattern: `name: (\w+), age: (\d+)`,
			group:   2,
			input:   "name: alice, age: 30",
			want:    []string{"30"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewGroupRule(tt.pattern, tt.group, tt.opts...)
			if err != nil {
				t.Fatalf("NewGroupRule() error: %v", err)
			}
			got, err := r.Apply(tt.input, &LineContext{LineNum: 1})
			if err != nil {
				t.Fatalf("Apply() error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d lines, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("line %d: got %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestGroupRule_InvalidGroupNum(t *testing.T) {
	_, err := NewGroupRule(`(\w+)`, 0)
	if err == nil {
		t.Fatal("expected error for group 0")
	}

	_, err = NewGroupRule(`(\w+)`, -1)
	if err == nil {
		t.Fatal("expected error for negative group")
	}
}
