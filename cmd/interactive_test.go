package cmd

import (
	"reflect"
	"testing"
)

func TestExtractPlaceholders(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{"add <file>", []string{"file"}},
		{"remote add <name> <url>", []string{"name", "url"}},
		{"no placeholder", nil},
		{"<onlyone>", []string{"onlyone"}},
		{"<first> then text", []string{"first"}},
		{"text then <last>", []string{"last"}},
		{"<multiple> <placeholders>", []string{"multiple", "placeholders"}},
		{"<incomplete", nil},
		{"incomplete>", nil},
		{"<>", []string{""}},
		{"< >", []string{" "}},
	}

	for _, tc := range cases {
		result := extractPlaceholders(tc.input)
		if tc.expected == nil {
			if result != nil {
				t.Errorf("input: %s, expected nil, but got: %v", tc.input, result)
			}
		} else if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("input: %s, expected: %v, but got: %v", tc.input, tc.expected, result)
		}
	}
}
