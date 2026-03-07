package cmd

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "plain words",
			input: "commit -m message",
			want:  []string{"commit", "-m", "message"},
		},
		{
			name:  "double-quoted argument",
			input: `commit -m "fix bug"`,
			want:  []string{"commit", "-m", "fix bug"},
		},
		{
			name:  "single-quoted argument",
			input: "commit -m 'fix the bug'",
			want:  []string{"commit", "-m", "fix the bug"},
		},
		{
			name:  "backslash escape inside double quotes",
			input: `echo "it's a \"test\""`,
			want:  []string{"echo", `it's a "test"`},
		},
		{
			name:  "multiple spaces between tokens",
			input: "git   log  --oneline",
			want:  []string{"git", "log", "--oneline"},
		},
		{
			name:  "tab separator",
			input: "git\tlog",
			want:  []string{"git", "log"},
		},
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "only whitespace",
			input: "   ",
			want:  nil,
		},
		{
			name:  "quoted empty string",
			input: `cmd ""`,
			want:  []string{"cmd", ""},
		},
		{
			name:  "single word",
			input: "help",
			want:  []string{"help"},
		},
		{
			name:  "mixed quotes",
			input: `git commit -m "fix: it's done"`,
			want:  []string{"git", "commit", "-m", "fix: it's done"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tokenize(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokenize(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
