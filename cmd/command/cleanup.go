package command

// cleanup returns command definitions for the Cleanup category.
func cleanup() []Info {
	return []Info{
		{
			Name:     "clean",
			Category: CategoryCleanup,
			Summary:  "Remove untracked files and directories",
			Usage:    []string{"ggc clean files", "ggc clean dirs", "ggc clean interactive"},
			Examples: []string{
				"ggc clean files       # Clean untracked files",
				"ggc clean dirs        # Clean untracked directories",
				"ggc clean interactive # Clean files interactively",
			},
			HandlerID: "clean",
			Subcommands: []SubcommandInfo{
				{Name: "clean files", Summary: "Clean untracked files", Usage: []string{"ggc clean files"}},
				{Name: "clean dirs", Summary: "Clean untracked directories", Usage: []string{"ggc clean dirs"}},
				{Name: "clean interactive", Summary: "Clean files interactively", Usage: []string{"ggc clean interactive"}},
			},
		},
		{
			Name:      "restore",
			Category:  CategoryCleanup,
			Summary:   "Restore files in working tree or staging area",
			Usage:     []string{"ggc restore <file>", "ggc restore .", "ggc restore staged <file>", "ggc restore staged .", "ggc restore <commit> <file>"},
			Examples:  []string{"ggc restore staged .", "ggc restore main README.md"},
			HandlerID: "restore",
			Subcommands: []SubcommandInfo{
				{Name: "restore <file>", Summary: "Restore file in working directory from index", Usage: []string{"ggc restore README.md"}},
				{Name: "restore .", Summary: "Restore all files in working directory from index", Usage: []string{"ggc restore ."}},
				{Name: "restore staged <file>", Summary: "Unstage file (restore from HEAD to index)", Usage: []string{"ggc restore staged README.md"}},
				{Name: "restore staged .", Summary: "Unstage all files", Usage: []string{"ggc restore staged ."}},
				{Name: "restore <commit> <file>", Summary: "Restore file from specific commit", Usage: []string{"ggc restore HEAD~1 README.md"}},
			},
		},
	}
}
