//nolint:unparam // Benchmark functions intentionally have unused return values
package cmd

import (
	"testing"
)

// Benchmark AddRune performance with copy approach
func BenchmarkUIState_AddRune_Copy(b *testing.B) {
	state := &UIState{
		input:     "hello world test string",
		cursorPos: 12, // Middle position
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset state for each iteration
		state.input = "hello world test string"
		state.cursorPos = 12

		// Add a rune in the middle
		state.AddRune('X')
	}
}

// Alternative implementation using nested append for comparison
func addRuneNested(input string, cursorPos int, r rune) string {
	inputRunes := []rune(input)
	if cursorPos <= len(inputRunes) {
		inputRunes = append(inputRunes[:cursorPos], append([]rune{r}, inputRunes[cursorPos:]...)...)
		return string(inputRunes)
	}
	return input
}

// Benchmark nested append approach
func BenchmarkAddRune_NestedAppend(b *testing.B) {
	input := "hello world test string"
	cursorPos := 12

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = addRuneNested(input, cursorPos, 'X')
	}
}

// Alternative implementation using copy for comparison
func addRuneCopy(input string, cursorPos int, r rune) string {
	inputRunes := []rune(input)
	if cursorPos <= len(inputRunes) {
		newRunes := make([]rune, len(inputRunes)+1)
		copy(newRunes, inputRunes[:cursorPos])
		newRunes[cursorPos] = r
		copy(newRunes[cursorPos+1:], inputRunes[cursorPos:])
		return string(newRunes)
	}
	return input
}

// Benchmark copy approach
func BenchmarkAddRune_Copy(b *testing.B) {
	input := "hello world test string"
	cursorPos := 12

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = addRuneCopy(input, cursorPos, 'X')
	}
}

// Benchmark with different string lengths
func BenchmarkAddRune_ShortString_Copy(b *testing.B) {
	input := "short"
	cursorPos := 2

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = addRuneCopy(input, cursorPos, 'X')
	}
}

func BenchmarkAddRune_ShortString_NestedAppend(b *testing.B) {
	input := "short"
	cursorPos := 2

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = addRuneNested(input, cursorPos, 'X')
	}
}

func BenchmarkAddRune_LongString_Copy(b *testing.B) {
	input := "this is a very long string that contains many characters for testing performance with longer inputs"
	cursorPos := 50

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = addRuneCopy(input, cursorPos, 'X')
	}
}

func BenchmarkAddRune_LongString_NestedAppend(b *testing.B) {
	input := "this is a very long string that contains many characters for testing performance with longer inputs"
	cursorPos := 50
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = addRuneNested(input, cursorPos, 'X')
	}
}

// Benchmark RemoveChar performance with append approach
func BenchmarkUIState_RemoveChar_Append(b *testing.B) {
	state := &UIState{
		input:     "hello world test string",
		cursorPos: 12, // Middle position
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset state for each iteration
		state.input = "hello world test string"
		state.cursorPos = 12
		
		// Remove a character
		state.RemoveChar()
	}
}

// Alternative implementation using make + copy for comparison
func removeCharCopy(input string, cursorPos int) string {
	inputRunes := []rune(input)
	if cursorPos > 0 && cursorPos <= len(inputRunes) {
		newRunes := make([]rune, len(inputRunes)-1)
		copy(newRunes[:cursorPos-1], inputRunes[:cursorPos-1])
		copy(newRunes[cursorPos-1:], inputRunes[cursorPos:])
		return string(newRunes)
	}
	return input
}

// Alternative implementation using append for comparison
func removeCharAppend(input string, cursorPos int) string {
	inputRunes := []rune(input)
	if cursorPos > 0 && cursorPos <= len(inputRunes) {
		inputRunes = append(inputRunes[:cursorPos-1], inputRunes[cursorPos:]...)
		return string(inputRunes)
	}
	return input
}

// Benchmark append vs copy approaches for removal
func BenchmarkRemoveChar_Append(b *testing.B) {
	input := "hello world test string"
	cursorPos := 12
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = removeCharAppend(input, cursorPos)
	}
}

func BenchmarkRemoveChar_Copy(b *testing.B) {
	input := "hello world test string"
	cursorPos := 12
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = removeCharCopy(input, cursorPos)
	}
}

func BenchmarkRemoveChar_ShortString_Append(b *testing.B) {
	input := "short"
	cursorPos := 3
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = removeCharAppend(input, cursorPos)
	}
}

func BenchmarkRemoveChar_ShortString_Copy(b *testing.B) {
	input := "short"
	cursorPos := 3
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = removeCharCopy(input, cursorPos)
	}
}
