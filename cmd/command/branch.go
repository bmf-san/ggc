package command

// branch returns command definitions for the Branch category.
func branch() []Info {
	return []Info{
		{
			Name:     "branch",
			Category: CategoryBranch,
			Summary:  "List, create, and manage branches",
			Usage:    []string{"ggc branch <subcommand>"},
			Examples: []string{
				"ggc branch current                # Show current branch",
				"ggc branch checkout               # Switch to an existing branch",
				"ggc branch checkout remote        # Create and checkout a local branch from the remote",
				"ggc branch create feature/login   # Create and checkout new branch",
				"ggc branch delete feature/login   # Delete local branch",
				"ggc branch delete merged          # Delete local merged branch",
				"ggc branch rename old new         # Rename a branch",
				"ggc branch move feature abc123    # Move branch to specified commit",
				"ggc branch set upstream feature origin/feature  # Set upstream branch",
				"ggc branch info feature           # Show detailed branch information",
				"ggc branch list verbose           # Show detailed branch listing",
				"ggc branch sort date              # List branches sorted by date",
				"ggc branch contains abc123        # Show branches containing a commit",
			},
			HandlerID: "branch",
			Subcommands: []SubcommandInfo{
				{Name: "branch current", Summary: "Show current branch name", Usage: []string{"ggc branch current"}},
				{Name: "branch checkout", Summary: "Switch to an existing branch", Usage: []string{"ggc branch checkout"}},
				{Name: "branch checkout remote", Summary: "Create and checkout a local branch from the remote", Usage: []string{"ggc branch checkout remote"}},
				{Name: "branch create", Summary: "Create and checkout a new branch", Usage: []string{"ggc branch create feature/login"}},
				{Name: "branch delete", Summary: "Delete local branch", Usage: []string{"ggc branch delete feature/login"}, Examples: []string{
					"ggc branch delete feature/123          # Delete a branch",
					"ggc branch delete feature/123 --force  # Force delete a branch",
				}},
				{Name: "branch delete merged", Summary: "Delete local merged branch", Usage: []string{"ggc branch delete merged"}},
				{Name: "branch rename <old> <new>", Summary: "Rename a branch", Usage: []string{"ggc branch rename old new"}},
				{Name: "branch move <branch> <commit>", Summary: "Move branch to specified commit", Usage: []string{"ggc branch move feature abc123"}},
				{Name: "branch set upstream <branch> <upstream>", Summary: "Set upstream for a branch", Usage: []string{"ggc branch set upstream feature origin/feature"}},
				{Name: "branch info <branch>", Summary: "Show detailed branch information", Usage: []string{"ggc branch info feature"}},
				{Name: "branch list verbose", Summary: "Show detailed branch listing", Usage: []string{"ggc branch list verbose"}},
				{Name: "branch list local", Summary: "List local branches", Usage: []string{"ggc branch list local"}},
				{Name: "branch list remote", Summary: "List remote branches", Usage: []string{"ggc branch list remote"}},
				{Name: "branch sort [date|name]", Summary: "List branches sorted by date or name", Usage: []string{"ggc branch sort date"}},
				{Name: "branch contains <commit>", Summary: "Show branches containing a commit", Usage: []string{"ggc branch contains abc123"}},
			},
		},
	}
}
