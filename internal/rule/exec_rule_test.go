package rule

import "testing"

func TestExecRule(t *testing.T) {
	tests := []struct {
		name    string
		command string
		input   []string
		want    []string
		wantErr bool
	}{
		{
			name:    "sort command",
			command: "sort",
			input:   []string{"c", "a", "b"},
			want:    []string{"a", "b", "c"},
		},
		{
			name:    "cat passes through",
			command: "cat",
			input:   []string{"hello", "world"},
			want:    []string{"hello", "world"},
		},
		{
			name:    "grep filters lines",
			command: "grep hello",
			input:   []string{"hello", "world", "hello again"},
			want:    []string{"hello", "hello again"},
		},
		{
			name:    "awk transforms",
			command: "awk '{print NR\": \"$0}'",
			input:   []string{"a", "b", "c"},
			want:    []string{"1: a", "2: b", "3: c"},
		},
		{
			name:    "wc -l counts lines",
			command: "wc -l | tr -d ' '",
			input:   []string{"a", "b", "c"},
			want:    []string{"3"},
		},
		{
			name:    "empty output",
			command: "grep nomatch",
			input:   []string{"hello", "world"},
			want:    []string{},
			wantErr: true, // grep returns exit 1 on no match
		},
		{
			name:    "failing command returns error",
			command: "nonexistent_command_xyz",
			input:   []string{"hello"},
			wantErr: true,
		},
		{
			name:    "head limits output",
			command: "head -2",
			input:   []string{"a", "b", "c", "d"},
			want:    []string{"a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewExecRule(tt.command)
			got, err := r.ApplyDocument(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d lines, want %d: %v", len(got), len(tt.want), got)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("line %d: got %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
