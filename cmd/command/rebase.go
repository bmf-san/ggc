package command

// rebase returns command definitions for the Rebase category.
func rebase() []Info {
	return []Info{
		{
			Name:     "rebase",
			Category: CategoryRebase,
			Summary:  "Reapply commits on top of another base tip",
			Usage:    []string{"ggc rebase <subcommand>"},
			Examples: []string{
				"ggc rebase interactive  # Interactive rebase",
				"ggc rebase main         # Rebase current branch onto 'main'",
				"ggc rebase continue     # Continue an in-progress rebase",
				"ggc rebase abort        # Abort an in-progress rebase",
				"ggc rebase skip         # Skip current patch and continue",
			},
			HandlerID: "rebase",
			Subcommands: []SubcommandInfo{
				{Name: "rebase interactive", Summary: "Interactive rebase", Usage: []string{"ggc rebase interactive"}},
				{Name: "rebase <upstream>", Summary: "Rebase current branch onto <upstream>", Usage: []string{"ggc rebase main"}},
				{Name: "rebase continue", Summary: "Continue an in-progress rebase", Usage: []string{"ggc rebase continue"}},
				{Name: "rebase abort", Summary: "Abort an in-progress rebase", Usage: []string{"ggc rebase abort"}},
				{Name: "rebase skip", Summary: "Skip current patch and continue", Usage: []string{"ggc rebase skip"}},
			},
		},
	}
}
