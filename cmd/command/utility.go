package command

// utility returns command definitions for the Utility category.
func utility() []Info {
	return []Info{
		{
			Name:     "version",
			Category: CategoryUtility,
			Summary:  "Display current ggc version",
			Usage:    []string{"ggc version"},
			Examples: []string{
				"ggc version   # Shows build time, latest commit and version number",
			},
			HandlerID: "version",
		},
		{
			Name:     "debug-keys",
			Category: CategoryUtility,
			Summary:  "Debug keybinding issues and capture raw key sequences",
			Usage: []string{
				"ggc debug-keys",
				"ggc debug-keys raw",
				"ggc debug-keys raw <file>",
			},
			Examples: []string{
				"ggc debug-keys                 # Show active keybindings",
				"ggc debug-keys raw             # Capture key sequences interactively",
				"ggc debug-keys raw keys.txt    # Capture and save to keys.txt",
			},
			HandlerID: "debug-keys",
			Subcommands: []SubcommandInfo{
				{
					Name:    "debug-keys",
					Summary: "Show current keybindings",
					Usage:   []string{"ggc debug-keys"},
				},
				{
					Name:    "debug-keys raw",
					Summary: "Capture key sequences interactively",
					Usage:   []string{"ggc debug-keys raw"},
				},
				{
					Name:    "debug-keys raw <file>",
					Summary: "Capture key sequences and save them to a file",
					Usage:   []string{"ggc debug-keys raw keys.txt"},
				},
			},
		},
		{
			Name:      "quit",
			Category:  CategoryUtility,
			Summary:   "Exit interactive mode",
			Usage:     []string{"quit"},
			Examples:  []string{"quit"},
			HandlerID: "quit",
		},
	}
}
