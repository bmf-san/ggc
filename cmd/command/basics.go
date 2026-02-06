package command

// basics returns command definitions for the Basics category.
func basics() []Info {
	return []Info{
		{
			Name:      "help",
			Category:  CategoryBasics,
			Summary:   "Show help information for commands",
			Usage:     []string{"ggc help", "ggc help <command>"},
			Examples:  []string{"ggc help", "ggc help branch"},
			HandlerID: "help",
			Subcommands: []SubcommandInfo{
				{
					Name:    "help",
					Summary: "Show main help message",
					Usage:   []string{"ggc help"},
				},
				{
					Name:    "help <command>",
					Summary: "Show help for a specific command",
					Usage:   []string{"ggc help branch"},
				},
			},
		},
		{
			Name:     "add",
			Category: CategoryBasics,
			Summary:  "Stage changes for the next commit",
			Usage:    []string{"ggc add <file>", "ggc add .", "ggc add interactive", "ggc add patch"},
			Examples: []string{
				"ggc add file.txt   # Add a specific file",
				"ggc add .          # Add all changes to index",
				"ggc add interactive  # Add changes interactively",
				"ggc add patch        # Add changes interactively (patch mode)",
			},
			HandlerID: "add",
			Subcommands: []SubcommandInfo{
				{
					Name:    "add <file>",
					Summary: "Add a specific file to the index",
					Usage:   []string{"ggc add README.md"},
				},
				{
					Name:    "add .",
					Summary: "Add all changes to the index",
					Usage:   []string{"ggc add ."},
				},
				{
					Name:    "add interactive",
					Summary: "Add changes interactively",
					Usage:   []string{"ggc add interactive"},
				},
				{
					Name:    "add patch",
					Summary: "Add changes interactively (patch mode)",
					Usage:   []string{"ggc add patch"},
				},
			},
		},
	}
}
