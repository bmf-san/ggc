package command

// commit returns command definitions for the Commit category.
func commit() []Info {
	return []Info{
		{
			Name:     "log",
			Category: CategoryCommit,
			Summary:  "Inspect commit history",
			Usage:    []string{"ggc log simple", "ggc log graph"},
			Examples: []string{
				"ggc log simple  # Show commit logs in a simple format",
				"ggc log graph   # Show commit logs with a graph",
			},
			HandlerID: "log",
			Subcommands: []SubcommandInfo{
				{Name: "log simple", Summary: "Show simple historical log", Usage: []string{"ggc log simple"}},
				{Name: "log graph", Summary: "Show log with graph", Usage: []string{"ggc log graph"}},
			},
		},
		{
			Name:     "commit",
			Category: CategoryCommit,
			Summary:  "Create commits from staged changes",
			Usage:    []string{"ggc commit <message>", "ggc commit amend", "ggc commit allow empty"},
			Examples: []string{
				"ggc commit \"Update docs\"        # Create commit with a message",
				"ggc commit allow empty            # Create an empty commit",
				"ggc commit amend                  # Amend previous commit (editor)",
				"ggc commit amend no-edit          # Amend without editing commit message",
			},
			HandlerID: "commit",
			Subcommands: []SubcommandInfo{
				{Name: "commit <message>", Summary: "Create commit with a message", Usage: []string{"ggc commit \"Add feature\""}},
				{Name: "commit allow empty", Summary: "Create an empty commit", Usage: []string{"ggc commit allow empty"}},
				{Name: "commit amend", Summary: "Amend previous commit (editor)", Usage: []string{"ggc commit amend"}},
				{Name: "commit amend no-edit", Summary: "Amend without editing commit message", Usage: []string{"ggc commit amend no-edit"}},
			},
		},
	}
}
