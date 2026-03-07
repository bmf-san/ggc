package interactive

// CommandInfo contains the name and description of a command available in
// interactive mode. The list is injected at construction time via NewUI so
// that this package does not depend on the cmd layer.
type CommandInfo struct {
	Command     string
	Description string
}

// extractPlaceholders extracts <...> placeholders from a string
func extractPlaceholders(s string) []string {
	var res []string
	start := -1
	for i, c := range s {
		if c == '<' {
			start = i + 1
		} else if c == '>' && start != -1 {
			res = append(res, s[start:i])
			start = -1
		}
	}
	return res
}
