package command

// stash returns command definitions for the Stash category.
func stash() []Info {
	return []Info{
		{
			Name:     "stash",
			Category: CategoryStash,
			Summary:  "Save and reapply work-in-progress changes",
			Usage:    []string{"ggc stash <subcommand>"},
			Examples: []string{
				"ggc stash                              # Stash current changes",
				"ggc stash list                         # List all stashes",
				"ggc stash show [stash]                 # Show changes in stash",
				"ggc stash apply [stash]                # Apply stash without removing it",
				"ggc stash pop [stash]                  # Apply and remove stash",
				"ggc stash drop [stash]                 # Remove stash",
				"ggc stash branch <branch> [stash]      # Create branch from stash",
				"ggc stash push [-m message] [files]    # Save changes to new stash",
				"ggc stash save [message]               # Save changes to new stash",
				"ggc stash clear                        # Remove all stashes",
				"ggc stash create                       # Create stash and return object name",
				"ggc stash store <object>               # Store stash object",
			},
			HandlerID: "stash",
			Subcommands: []SubcommandInfo{
				{Name: "stash", Summary: "Stash current changes", Usage: []string{"ggc stash"}},
				{Name: "stash list", Summary: "List all stashes", Usage: []string{"ggc stash list"}},
				{Name: "stash show", Summary: "Show changes in stash", Usage: []string{"ggc stash show"}},
				{Name: "stash show <stash>", Summary: "Show changes in specific stash", Usage: []string{"ggc stash show stash@{1}"}},
				{Name: "stash apply", Summary: "Apply stash without removing it", Usage: []string{"ggc stash apply"}},
				{Name: "stash apply <stash>", Summary: "Apply specific stash without removing it", Usage: []string{"ggc stash apply stash@{1}"}},
				{Name: "stash pop", Summary: "Apply and remove the latest stash", Usage: []string{"ggc stash pop"}},
				{Name: "stash pop <stash>", Summary: "Apply and remove specific stash", Usage: []string{"ggc stash pop stash@{1}"}},
				{Name: "stash drop", Summary: "Remove the latest stash", Usage: []string{"ggc stash drop"}},
				{Name: "stash drop <stash>", Summary: "Remove specific stash", Usage: []string{"ggc stash drop stash@{1}"}},
				{Name: "stash branch <branch>", Summary: "Create branch from stash", Usage: []string{"ggc stash branch feature"}},
				{Name: "stash branch <branch> <stash>", Summary: "Create branch from specific stash", Usage: []string{"ggc stash branch feature stash@{1}"}},
				{Name: "stash push", Summary: "Save changes to new stash", Usage: []string{"ggc stash push"}},
				{Name: "stash push -m <message>", Summary: "Save changes to new stash with message", Usage: []string{"ggc stash push -m \"WIP\""}},
				{Name: "stash save <message>", Summary: "Save changes to new stash with message", Usage: []string{"ggc stash save \"WIP\""}},
				{Name: "stash clear", Summary: "Remove all stashes", Usage: []string{"ggc stash clear"}},
				{Name: "stash create", Summary: "Create stash and return object name", Usage: []string{"ggc stash create"}},
				{Name: "stash store <object>", Summary: "Store stash object", Usage: []string{"ggc stash store 1234abcd"}},
			},
		},
	}
}
