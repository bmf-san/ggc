// Package interactive houses interactive UI types and helpers shared across the application.
package interactive

import "unicode"

// fuzzyMatch performs fuzzy matching between text and pattern
// Returns true if all characters in pattern appear in text in order (but not necessarily consecutive)
func fuzzyMatch(text, pattern string) bool {
	matched, _ := fuzzyMatchScore(text, pattern)
	return matched
}

// fuzzyMatchScore returns whether the pattern matches the text and a relevance score for sorting results.
// Lower scores indicate a tighter, earlier match.
func fuzzyMatchScore(text, pattern string) (bool, matchScore) {
	if pattern == "" {
		return true, matchScore{length: len([]rune(text))}
	}

	textRunes := []rune(text)
	patternRunes := []rune(pattern)

	matched, meta := matchPattern(textRunes, patternRunes)
	if !matched {
		return false, matchScore{}
	}

	trailing := len(textRunes) - meta.lastIndex - 1
	continuation := continuationPenalty(textRunes, meta.lastIndex)
	score := matchScore{
		first:        meta.firstIndex,
		gap:          meta.gapScore,
		trailing:     trailing,
		continuation: continuation,
		length:       len(textRunes),
	}

	return true, score
}

type matchMetadata struct {
	firstIndex int
	lastIndex  int
	gapScore   int
}

func matchPattern(textRunes, patternRunes []rune) (bool, matchMetadata) {
	meta := matchMetadata{
		firstIndex: -1,
		lastIndex:  -1,
	}

	textIdx := 0
	patternIdx := 0

	for textIdx < len(textRunes) && patternIdx < len(patternRunes) {
		if textRunes[textIdx] == patternRunes[patternIdx] {
			if meta.firstIndex == -1 {
				meta.firstIndex = textIdx
			}
			if meta.lastIndex != -1 {
				meta.gapScore += textIdx - meta.lastIndex - 1
			}
			meta.lastIndex = textIdx
			patternIdx++
		}
		textIdx++
	}

	if patternIdx != len(patternRunes) {
		return false, meta
	}

	return true, meta
}

func continuationPenalty(textRunes []rune, lastMatchIdx int) int {
	if lastMatchIdx < 0 || lastMatchIdx+1 >= len(textRunes) {
		return 0
	}

	nextIdx := lastMatchIdx + 1
	spaceSkipped := false
	for nextIdx < len(textRunes) && textRunes[nextIdx] == ' ' {
		spaceSkipped = true
		nextIdx++
	}

	if spaceSkipped && nextIdx < len(textRunes) && (unicode.IsLetter(textRunes[nextIdx]) || unicode.IsDigit(textRunes[nextIdx])) {
		return 1
	}

	return 0
}

type matchScore struct {
	first        int
	gap          int
	trailing     int
	continuation int
	length       int
}

func (m matchScore) less(other matchScore) bool {
	if m.first != other.first {
		return m.first < other.first
	}
	if m.gap != other.gap {
		return m.gap < other.gap
	}
	if m.continuation != other.continuation {
		return m.continuation < other.continuation
	}
	if m.trailing != other.trailing {
		return m.trailing < other.trailing
	}
	if m.length != other.length {
		return m.length < other.length
	}
	return false
}
