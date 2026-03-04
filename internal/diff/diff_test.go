package diff

import (
	"strings"
	"testing"
)

func TestCompute(t *testing.T) {
	tests := []struct {
		name string
		a    []string
		b    []string
		want []Change
	}{
		{
			name: "identical",
			a:    []string{"a", "b", "c"},
			b:    []string{"a", "b", "c"},
			want: []Change{
				{Equal, "a"}, {Equal, "b"}, {Equal, "c"},
			},
		},
		{
			name: "insertion",
			a:    []string{"a", "c"},
			b:    []string{"a", "b", "c"},
			want: []Change{
				{Equal, "a"}, {Insert, "b"}, {Equal, "c"},
			},
		},
		{
			name: "deletion",
			a:    []string{"a", "b", "c"},
			b:    []string{"a", "c"},
			want: []Change{
				{Equal, "a"}, {Delete, "b"}, {Equal, "c"},
			},
		},
		{
			name: "substitution",
			a:    []string{"a", "b", "c"},
			b:    []string{"a", "x", "c"},
			want: []Change{
				{Equal, "a"}, {Delete, "b"}, {Insert, "x"}, {Equal, "c"},
			},
		},
		{
			name: "all different",
			a:    []string{"a", "b"},
			b:    []string{"x", "y"},
			want: []Change{
				{Delete, "a"}, {Delete, "b"}, {Insert, "x"}, {Insert, "y"},
			},
		},
		{
			name: "empty original",
			a:    []string{},
			b:    []string{"a", "b"},
			want: []Change{
				{Insert, "a"}, {Insert, "b"},
			},
		},
		{
			name: "empty result",
			a:    []string{"a", "b"},
			b:    []string{},
			want: []Change{
				{Delete, "a"}, {Delete, "b"},
			},
		},
		{
			name: "both empty",
			a:    []string{},
			b:    []string{},
			want: nil,
		},
		{
			name: "multiple changes",
			a:    []string{"a", "b", "c", "d", "e"},
			b:    []string{"a", "x", "c", "y", "e"},
			want: []Change{
				{Equal, "a"}, {Delete, "b"}, {Insert, "x"},
				{Equal, "c"}, {Delete, "d"}, {Insert, "y"},
				{Equal, "e"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Compute(tt.a, tt.b)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d changes, want %d:\ngot:  %v\nwant: %v", len(got), len(tt.want), got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("change[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestFormat_NoColor(t *testing.T) {
	changes := []Change{
		{Equal, "a"},
		{Delete, "b"},
		{Insert, "x"},
		{Equal, "c"},
	}

	got := Format(changes, false)
	want := []string{" a", "-b", "+x", " c"}

	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("line[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestFormat_WithColor(t *testing.T) {
	changes := []Change{
		{Delete, "old"},
		{Insert, "new"},
	}

	got := Format(changes, true)

	if !strings.Contains(got[0], "\033[31m") {
		t.Errorf("expected red for delete, got %q", got[0])
	}
	if !strings.Contains(got[0], "-old") {
		t.Errorf("expected '-old' in delete, got %q", got[0])
	}
	if !strings.Contains(got[1], "\033[32m") {
		t.Errorf("expected green for insert, got %q", got[1])
	}
	if !strings.Contains(got[1], "+new") {
		t.Errorf("expected '+new' in insert, got %q", got[1])
	}
}

func TestHasChanges(t *testing.T) {
	tests := []struct {
		name    string
		changes []Change
		want    bool
	}{
		{"no changes", []Change{{Equal, "a"}, {Equal, "b"}}, false},
		{"has insert", []Change{{Equal, "a"}, {Insert, "b"}}, true},
		{"has delete", []Change{{Delete, "a"}, {Equal, "b"}}, true},
		{"empty", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasChanges(tt.changes)
			if got != tt.want {
				t.Errorf("HasChanges() = %v, want %v", got, tt.want)
			}
		})
	}
}
