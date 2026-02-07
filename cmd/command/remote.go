package command

// remote returns command definitions for the Remote category.
func remote() []Info {
	return []Info{
		{
			Name:     "push",
			Category: CategoryRemote,
			Summary:  "Update remote branches",
			Usage:    []string{"ggc push current", "ggc push force"},
			Examples: []string{
				"ggc push current  # Push current branch to remote",
				"ggc push force    # Force push current branch",
			},
			HandlerID: "push",
			Subcommands: []SubcommandInfo{
				{Name: "push current", Summary: "Push current branch to remote repository", Usage: []string{"ggc push current"}},
				{Name: "push force", Summary: "Force push current branch", Usage: []string{"ggc push force"}},
			},
		},
		{
			Name:     "pull",
			Category: CategoryRemote,
			Summary:  "Fetch and integrate from the remote",
			Usage:    []string{"ggc pull current", "ggc pull rebase"},
			Examples: []string{
				"ggc pull current  # Pull current branch from remote",
				"ggc pull rebase   # Pull with rebase",
			},
			HandlerID: "pull",
			Subcommands: []SubcommandInfo{
				{Name: "pull current", Summary: "Pull current branch from remote repository", Usage: []string{"ggc pull current"}},
				{Name: "pull rebase", Summary: "Pull and rebase", Usage: []string{"ggc pull rebase"}},
			},
		},
		{
			Name:     "fetch",
			Category: CategoryRemote,
			Summary:  "Download objects and refs from remotes",
			Usage:    []string{"ggc fetch", "ggc fetch prune"},
			Examples: []string{
				"ggc fetch prune   # Fetch and remove stale remote-tracking references",
			},
			HandlerID: "fetch",
			Subcommands: []SubcommandInfo{
				{Name: "fetch", Summary: "Fetch from the remote", Usage: []string{"ggc fetch"}},
				{Name: "fetch prune", Summary: "Fetch and clean stale references", Usage: []string{"ggc fetch prune"}},
			},
		},
		{
			Name:      "remote",
			Category:  CategoryRemote,
			Summary:   "Manage remotes",
			Usage:     []string{"ggc remote list", "ggc remote add <name> <url>", "ggc remote remove <name>", "ggc remote set-url <name> <url>"},
			Examples:  []string{"ggc remote list", "ggc remote add origin git@github.com:user/repo.git"},
			HandlerID: "remote",
			Subcommands: []SubcommandInfo{
				{Name: "remote list", Summary: "List all remote repositories", Usage: []string{"ggc remote list"}},
				{Name: "remote add <name> <url>", Summary: "Add remote repository", Usage: []string{"ggc remote add upstream git@github.com:user/repo.git"}},
				{Name: "remote remove <name>", Summary: "Remove remote repository", Usage: []string{"ggc remote remove upstream"}},
				{Name: "remote set-url <name> <url>", Summary: "Change remote URL", Usage: []string{"ggc remote set-url origin git@github.com:user/new.git"}},
			},
		},
	}
}
