package command

// status returns command definitions for the Status category.
func status() []Info {
	return []Info{
		{
			Name:     "status",
			Category: CategoryStatus,
			Summary:  "Show working tree status",
			Usage:    []string{"ggc status", "ggc status short"},
			Examples: []string{
				"ggc status        # Full detailed status output",
				"ggc status short  # Short, concise output (porcelain format)",
			},
			HandlerID: "status",
			Subcommands: []SubcommandInfo{
				{Name: "status", Summary: "Show working tree status", Usage: []string{"ggc status"}},
				{Name: "status short", Summary: "Show concise status (porcelain format)", Usage: []string{"ggc status short"}},
			},
		},
	}
}
