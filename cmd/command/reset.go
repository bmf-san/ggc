package command

// reset returns command definitions for the Reset category.
func reset() []Info {
	return []Info{
		{
			Name:     "reset",
			Category: CategoryBasics,
			Summary:  "Reset current HEAD to the specified state",
			Usage:    []string{"ggc reset", "ggc reset hard <commit>", "ggc reset soft <commit>"},
			Examples: []string{
				"ggc reset               # Hard reset to origin/<current-branch> and clean",
				"ggc reset hard HEAD~1   # Hard reset to previous commit",
				"ggc reset soft HEAD~1   # Soft reset: keep changes staged",
				"ggc reset soft HEAD~3   # Soft reset 3 commits, keeping changes staged",
			},
			Subcommands: []SubcommandInfo{
				{Name: "reset", Summary: "Hard reset to origin/<branch> and clean working directory", Usage: []string{"ggc reset"}},
				{Name: "reset hard <commit>", Summary: "Hard reset to specified commit", Usage: []string{"ggc reset hard HEAD~1"}},
				{Name: "reset soft <commit>", Summary: "Soft reset: move HEAD but keep changes staged", Usage: []string{"ggc reset soft HEAD~1"}},
			},
		},
	}
}
