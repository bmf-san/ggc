package ui

// Ellipsis truncates the provided string to maxLen characters, appending an ellipsis when
// truncation occurs. For zero or negative lengths it returns an empty string. The function
// is intentionally ASCII-focused to match existing interactive behavior.
func Ellipsis(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 1 {
		return "…"
	}
	return s[:maxLen-1] + "…"
}
