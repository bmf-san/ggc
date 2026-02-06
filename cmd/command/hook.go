package command

// hook returns command definitions for the Hook category.
func hook() []Info {
	return []Info{
		{
			Name:     "hook",
			Category: CategoryHook,
			Summary:  "Manage Git hooks",
			Usage:    []string{"ggc hook <subcommand>"},
			Examples: []string{
				"ggc hook list                    # List all hooks and their status",
				"ggc hook install <hook>          # Install a hook",
				"ggc hook enable <hook>           # Make a hook executable",
				"ggc hook disable <hook>          # Make a hook non-executable",
				"ggc hook uninstall <hook>        # Remove a hook",
				"ggc hook edit <hook>             # Edit a hook",
			},
			HandlerID: "hook",
			Subcommands: []SubcommandInfo{
				{Name: "hook list", Summary: "List all hooks", Usage: []string{"ggc hook list"}},
				{Name: "hook install <hook>", Summary: "Install a hook", Usage: []string{"ggc hook install pre-commit"}},
				{Name: "hook enable <hook>", Summary: "Enable a hook", Usage: []string{"ggc hook enable pre-commit"}},
				{Name: "hook disable <hook>", Summary: "Disable a hook", Usage: []string{"ggc hook disable pre-commit"}},
				{Name: "hook uninstall <hook>", Summary: "Uninstall an existing hook", Usage: []string{"ggc hook uninstall pre-commit"}},
				{Name: "hook edit <hook>", Summary: "Edit a hook's contents", Usage: []string{"ggc hook edit pre-commit"}},
			},
		},
	}
}
