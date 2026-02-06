package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseSelectionInput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		maxIndex    int
		wantResult  SelectionResult
		wantIndices []int
		wantInvalid string
	}{
		{
			name:       "empty input",
			input:      "",
			maxIndex:   5,
			wantResult: SelectionCanceled,
		},
		{
			name:       "whitespace only",
			input:      "   ",
			maxIndex:   5,
			wantResult: SelectionCanceled,
		},
		{
			name:       "all command",
			input:      "all",
			maxIndex:   5,
			wantResult: SelectionAll,
		},
		{
			name:       "none command",
			input:      "none",
			maxIndex:   5,
			wantResult: SelectionNone,
		},
		{
			name:        "single valid number",
			input:       "3",
			maxIndex:    5,
			wantResult:  SelectionItems,
			wantIndices: []int{2}, // 0-based
		},
		{
			name:        "multiple valid numbers",
			input:       "1 3 5",
			maxIndex:    5,
			wantResult:  SelectionItems,
			wantIndices: []int{0, 2, 4}, // 0-based
		},
		{
			name:        "numbers with extra spaces",
			input:       "  1   3  ",
			maxIndex:    5,
			wantResult:  SelectionItems,
			wantIndices: []int{0, 2},
		},
		{
			name:        "invalid number zero",
			input:       "0",
			maxIndex:    5,
			wantInvalid: "0",
		},
		{
			name:        "invalid number too high",
			input:       "6",
			maxIndex:    5,
			wantInvalid: "6",
		},
		{
			name:        "non-numeric input",
			input:       "abc",
			maxIndex:    5,
			wantInvalid: "abc",
		},
		{
			name:        "mixed valid and invalid",
			input:       "1 abc 3",
			maxIndex:    5,
			wantInvalid: "abc",
		},
		{
			name:        "negative number",
			input:       "-1",
			maxIndex:    5,
			wantInvalid: "-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, invalid := ParseSelectionInput(tt.input, tt.maxIndex)

			if tt.wantInvalid != "" {
				if invalid != tt.wantInvalid {
					t.Errorf("ParseSelectionInput() invalid = %q, want %q", invalid, tt.wantInvalid)
				}
				return
			}

			if result.Result != tt.wantResult {
				t.Errorf("ParseSelectionInput() result = %v, want %v", result.Result, tt.wantResult)
			}

			if tt.wantIndices != nil {
				if len(result.Indices) != len(tt.wantIndices) {
					t.Errorf("ParseSelectionInput() indices len = %d, want %d", len(result.Indices), len(tt.wantIndices))
				} else {
					for i, idx := range result.Indices {
						if idx != tt.wantIndices[i] {
							t.Errorf("ParseSelectionInput() indices[%d] = %d, want %d", i, idx, tt.wantIndices[i])
						}
					}
				}
			}
		})
	}
}

func TestSelectionLoop_Display(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)
	items := []string{"item1", "item2", "item3"}
	sl := NewSelectionLoop(f, "Select items:", items)

	sl.Display()

	output := buf.String()
	if !strings.Contains(output, "Select items:") {
		t.Error("Display() should contain header")
	}
	if !strings.Contains(output, "item1") {
		t.Error("Display() should contain item1")
	}
	if !strings.Contains(output, "item2") {
		t.Error("Display() should contain item2")
	}
	if !strings.Contains(output, "item3") {
		t.Error("Display() should contain item3")
	}
	if !strings.Contains(output, "> ") {
		t.Error("Display() should contain prompt")
	}
}

func TestSelectionLoop_ParseInput(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)
	items := []string{"a", "b", "c"}
	sl := NewSelectionLoop(f, "Select:", items)

	result, invalid := sl.ParseInput("1 3")
	if invalid != "" {
		t.Errorf("ParseInput() returned invalid = %q", invalid)
	}
	if result.Result != SelectionItems {
		t.Errorf("ParseInput() result = %v, want SelectionItems", result.Result)
	}
	if len(result.Indices) != 2 || result.Indices[0] != 0 || result.Indices[1] != 2 {
		t.Errorf("ParseInput() indices = %v, want [0 2]", result.Indices)
	}
}

func TestSelectionLoop_GetSelectedItems(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)
	items := []string{"apple", "banana", "cherry"}
	sl := NewSelectionLoop(f, "Select:", items)

	selected := sl.GetSelectedItems([]int{0, 2})

	if len(selected) != 2 {
		t.Fatalf("GetSelectedItems() len = %d, want 2", len(selected))
	}
	if selected[0] != "apple" {
		t.Errorf("GetSelectedItems()[0] = %q, want %q", selected[0], "apple")
	}
	if selected[1] != "cherry" {
		t.Errorf("GetSelectedItems()[1] = %q, want %q", selected[1], "cherry")
	}
}

func TestSelectionLoop_GetSelectedItems_OutOfRange(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)
	items := []string{"a", "b"}
	sl := NewSelectionLoop(f, "Select:", items)

	// Should skip out-of-range indices
	selected := sl.GetSelectedItems([]int{0, 5, 1, -1})

	if len(selected) != 2 {
		t.Fatalf("GetSelectedItems() len = %d, want 2", len(selected))
	}
	if selected[0] != "a" || selected[1] != "b" {
		t.Errorf("GetSelectedItems() = %v, want [a b]", selected)
	}
}

func TestSelectionLoop_Items(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)
	items := []string{"x", "y", "z"}
	sl := NewSelectionLoop(f, "Select:", items)

	got := sl.Items()

	if len(got) != 3 {
		t.Fatalf("Items() len = %d, want 3", len(got))
	}
	for i, item := range items {
		if got[i] != item {
			t.Errorf("Items()[%d] = %q, want %q", i, got[i], item)
		}
	}
}
