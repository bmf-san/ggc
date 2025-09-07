---
name: Feature Request
about: Suggest an idea for this project
title: '[Feature]: Improve interactive mode search from prefix matching to fuzzy/substring matching'
labels: 'enhancement'
assignees: ''

---

## Related Problem
Is your feature request related to a problem? Please describe.
The current interactive mode uses prefix matching (`strings.HasPrefix`) which can be limiting for users. For example:
- To find `branch delete`, users must type `branch` first - typing `delete` won't match
- To find `commit amend`, users must type `commit` first - typing `amend` won't match
- Users often remember keywords from the middle or end of commands rather than the beginning

This makes the search less intuitive and requires users to remember the exact command structure.

## Proposed Solution
A clear and concise description of what you want to happen.
Replace the current prefix matching with one of these more flexible approaches:

### Option 1: Substring matching
```go
if strings.Contains(cmd.Command, s.input) {
    s.filtered = append(s.filtered, cmd)
}
```

### Option 2: Fuzzy matching
Implement fuzzy matching that allows for:
- Non-contiguous character matching (e.g., "bd" matches "branch delete")
- Case-insensitive matching
- Scoring based on match quality

### Option 3: Multi-word search
Allow searching for multiple keywords separated by spaces:
```go
words := strings.Fields(strings.ToLower(s.input))
cmdLower := strings.ToLower(cmd.Command)
matches := true
for _, word := range words {
    if !strings.Contains(cmdLower, word) {
        matches = false
        break
    }
}
```

## Alternative Solutions
A clear and concise description of any alternative solutions or features you've considered.
1. **Hybrid approach**: Start with substring matching, then add fuzzy matching as an enhancement
2. **Configuration option**: Allow users to choose their preferred search mode in `~/.ggcconfig.yaml`
3. **Search mode toggle**: Allow users to cycle between search modes with a hotkey (e.g., Ctrl+T)

## Additional Context
Add any other context or screenshots about the feature request here.
- Current implementation is in `cmd/interactive.go` line 194: `strings.HasPrefix(cmd.Command, s.input)`
- Many modern CLI tools use fuzzy matching (fzf, telescope.nvim, etc.) which users are familiar with
- Substring matching would be a good first step that's easy to implement and significantly improves usability
- The search should remain fast and responsive even with the improved matching algorithm
