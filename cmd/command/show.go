package command

// show returns command definitions for the show command (Basics category).
func show() []Info {
	return []Info{
		{
			Name:     "show",
			Category: CategoryBasics,
			Summary:  "Show various types of objects (commits, tags, trees, blobs)",
			Usage: []string{
				"ggc show [<options>] [<object>...]",
			},
			Examples: []string{
				"ggc show                              # Show HEAD commit",
				"ggc show HEAD~1                       # Show previous commit",
				"ggc show abc1234                      # Show a specific commit",
				"ggc show --stat HEAD                  # Show commit with diffstat",
				"ggc show --name-only HEAD             # Show only changed file names",
				"ggc show v1.0.0                       # Show a tag",
				"ggc show HEAD:path/to/file.go         # Show file contents at HEAD",
			},
			Subcommands: []SubcommandInfo{
				{Name: "show", Summary: "Show HEAD commit", Usage: []string{"ggc show"}},
				{Name: "show <object>", Summary: "Show a specific commit, tag, tree, or blob", Usage: []string{"ggc show HEAD~1"}},
				{Name: "show --stat <object>", Summary: "Show object with diffstat", Usage: []string{"ggc show --stat HEAD"}},
				{Name: "show --name-only <object>", Summary: "Show object with names only", Usage: []string{"ggc show --name-only HEAD"}},
			},
		},
	}
}
