package ui

import (
	"strconv"
	"strings"
)

// SelectionResult represents the result of a selection operation.
type SelectionResult int

const (
	// SelectionCanceled indicates the user canceled the selection.
	SelectionCanceled SelectionResult = iota
	// SelectionAll indicates the user selected all items.
	SelectionAll
	// SelectionNone indicates the user wants to deselect/continue.
	SelectionNone
	// SelectionItems indicates specific items were selected.
	SelectionItems
)

// SelectionInput represents parsed user input for selection.
type SelectionInput struct {
	Result  SelectionResult
	Indices []int // 0-based indices of selected items
}

// ParseSelectionInput parses user input for a selection prompt.
// It recognizes:
//   - Empty input or whitespace: returns SelectionCanceled
//   - "all": returns SelectionAll
//   - "none": returns SelectionNone
//   - Space-separated numbers: returns SelectionItems with 0-based indices
//
// The maxIndex parameter specifies the maximum valid 1-based index (typically len(items)).
// Invalid indices result in an error string being returned.
func ParseSelectionInput(input string, maxIndex int) (SelectionInput, string) {
	input = strings.TrimSpace(input)

	if input == "" {
		return SelectionInput{Result: SelectionCanceled}, ""
	}

	if input == "all" {
		return SelectionInput{Result: SelectionAll}, ""
	}

	if input == "none" {
		return SelectionInput{Result: SelectionNone}, ""
	}

	// Parse space-separated numbers
	fields := strings.Fields(input)
	indices := make([]int, 0, len(fields))

	for _, field := range fields {
		n, err := strconv.Atoi(field)
		if err != nil || n < 1 || n > maxIndex {
			return SelectionInput{}, field // Return invalid field
		}
		indices = append(indices, n-1) // Convert to 0-based
	}

	return SelectionInput{Result: SelectionItems, Indices: indices}, ""
}

// SelectionLoop provides a reusable pattern for interactive selection.
type SelectionLoop struct {
	formatter *Formatter
	header    string
	items     []string
}

// NewSelectionLoop creates a new selection loop with the given formatter, header, and items.
func NewSelectionLoop(formatter *Formatter, header string, items []string) *SelectionLoop {
	return &SelectionLoop{
		formatter: formatter,
		header:    header,
		items:     items,
	}
}

// Display shows the selection interface.
func (s *SelectionLoop) Display() {
	s.formatter.Header(s.header)
	for i, item := range s.items {
		s.formatter.NumberedItem(i+1, item)
	}
	s.formatter.Prompt()
}

// ParseInput parses user input and returns the selection result.
// Returns the SelectionInput and any invalid field that was encountered.
func (s *SelectionLoop) ParseInput(input string) (SelectionInput, string) {
	return ParseSelectionInput(input, len(s.items))
}

// GetSelectedItems returns the items corresponding to the given 0-based indices.
func (s *SelectionLoop) GetSelectedItems(indices []int) []string {
	selected := make([]string, 0, len(indices))
	for _, idx := range indices {
		if idx >= 0 && idx < len(s.items) {
			selected = append(selected, s.items[idx])
		}
	}
	return selected
}

// Items returns all items in the selection.
func (s *SelectionLoop) Items() []string {
	return s.items
}
