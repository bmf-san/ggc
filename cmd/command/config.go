package command

// config returns command definitions for the Config category.
func config() []Info {
	return []Info{
		{
			Name:     "config",
			Category: CategoryConfig,
			Summary:  "Get and set ggc configuration",
			Usage:    []string{"ggc config list", "ggc config get <key>", "ggc config set <key> <value>"},
			Examples: []string{
				"ggc config list                  # List all configuration values",
				"ggc config get <key>             # Get a config value by key path (e.g., 'ui.color')",
				"ggc config set <key> <value>     # Set a config value by key path",
			},
			HandlerID: "config",
			Subcommands: []SubcommandInfo{
				{Name: "config list", Summary: "List all configuration", Usage: []string{"ggc config list"}},
				{Name: "config get <key>", Summary: "Get a specific config value", Usage: []string{"ggc config get core.editor"}},
				{Name: "config set <key> <value>", Summary: "Set a configuration value", Usage: []string{"ggc config set core.editor vim"}},
			},
		},
	}
}
