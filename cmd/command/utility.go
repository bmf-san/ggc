package command

// utility returns command definitions for the Utility category.
func utility() []Info {
	return []Info{
		{
			Name:     "version",
			Category: CategoryUtility,
			Summary:  "Display current ggc version",
			Usage: []string{
				"ggc version",
				"ggc version json",
			},
			Examples: []string{
				"ggc version        # Human-readable version, commit, build time, os/arch",
				"ggc version json   # Same info as a JSON document for scripting",
			},
			Subcommands: []SubcommandInfo{
				{
					Name:    "version json",
					Summary: "Emit the version information as a JSON document",
					Usage:   []string{"ggc version json"},
				},
			},
		},
		{
			Name:     "doctor",
			Category: CategoryUtility,
			Summary:  "Diagnose the local ggc installation",
			Usage:    []string{"ggc doctor"},
			Examples: []string{
				"ggc doctor   # Check git binary, config, shell completions, TTY, etc.",
			},
		},
		{
			Name:     "completion",
			Category: CategoryUtility,
			Summary:  "Print or install shell completion scripts",
			Usage: []string{
				"ggc completion <bash|zsh|fish>",
				"ggc completion install <bash|zsh|fish>",
			},
			Examples: []string{
				"ggc completion bash                   # Print the bash completion to stdout",
				"ggc completion install zsh            # Install zsh completion under ~/.zsh/completions/",
				"ggc completion fish > ~/.config/fish/completions/ggc.fish",
			},
			Subcommands: []SubcommandInfo{
				{
					Name:    "completion bash",
					Summary: "Print bash completion script",
					Usage:   []string{"ggc completion bash"},
				},
				{
					Name:    "completion zsh",
					Summary: "Print zsh completion script",
					Usage:   []string{"ggc completion zsh"},
				},
				{
					Name:    "completion fish",
					Summary: "Print fish completion script",
					Usage:   []string{"ggc completion fish"},
				},
				{
					Name:    "completion install <shell>",
					Summary: "Install the completion script for <bash|zsh|fish>",
					Usage:   []string{"ggc completion install <bash|zsh|fish>"},
				},
			},
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
			Name:     "quit",
			Category: CategoryUtility,
			Summary:  "Exit interactive mode",
			Usage:    []string{"quit"},
			Examples: []string{"quit"},
		},
	}
}
