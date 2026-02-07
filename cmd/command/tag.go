package command

// tag returns command definitions for the Tag category.
func tag() []Info {
	return []Info{
		{
			Name:     "tag",
			Category: CategoryTag,
			Summary:  "Create, list, and manage tags",
			Usage:    []string{"ggc tag list", "ggc tag annotated <tag> <message>", "ggc tag delete <tag>", "ggc tag show <tag>", "ggc tag push [<remote> <tag>]", "ggc tag create <tag>"},
			Examples: []string{
				"ggc tag                                   # List all tags",
				"ggc tag list                              # List all tags (sorted)",
				"ggc tag list v1.*                         # List tags matching pattern",
				"ggc tag create v1.0.0                     # Create tag",
				"ggc tag create v1.0.0 abc123              # Tag specific commit",
				"ggc tag annotated v1.0.0 'Release notes'  # Create annotated tag",
				"ggc tag delete v1.0.0                     # Delete tag",
				"ggc tag push                              # Push all tags to origin",
				"ggc tag push origin v1.0.0                # Push specific tag (remote first)",
				"ggc tag show v1.0.0                       # Show tag information",
			},
			HandlerID: "tag",
			Subcommands: []SubcommandInfo{
				{Name: "tag list", Summary: "List all tags", Usage: []string{"ggc tag list"}},
				{Name: "tag annotated <tag> <message>", Summary: "Create annotated tag", Usage: []string{"ggc tag annotated v1.0.0 \"Release\""}},
				{Name: "tag delete <tag>", Summary: "Delete tag", Usage: []string{"ggc tag delete v1.0.0"}},
				{Name: "tag show <tag>", Summary: "Show tag information", Usage: []string{"ggc tag show v1.0.0"}},
				{Name: "tag push", Summary: "Push tags to remote", Usage: []string{"ggc tag push", "ggc tag push <remote> <tag>"}},
				{Name: "tag create <tag>", Summary: "Create tag", Usage: []string{"ggc tag create v1.0.1"}},
			},
		},
	}
}
