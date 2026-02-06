package command

// diff returns command definitions for the Diff category.
func diff() []Info {
	return []Info{
		{
			Name:     "diff",
			Category: CategoryDiff,
			Summary:  "Inspect changes between commits, the index, and the working tree",
			Usage: []string{
				"ggc diff [staged|unstaged|head] [--stat|--name-only|--name-status] [<commit>|<commit1> <commit2>] [--] [<path>...]",
			},
			Examples: []string{
				"ggc diff --stat                     # Show staged + unstaged changes with summary",
				"ggc diff staged cmd/diff.go         # Diff staged changes for a file",
				"ggc diff abc123 def456              # Compare two commits",
				"ggc diff abc123 cmd/diff.go         # Compare commit to working tree for a path",
				"ggc diff -- cmd/deleted_file.go     # Diff a path using -- for disambiguation",
			},
			HandlerID: "diff",
			Subcommands: []SubcommandInfo{
				{Name: "diff", Summary: "Show changes (git diff HEAD)", Usage: []string{"ggc diff"}},
				{Name: "diff unstaged", Summary: "Show unstaged changes", Usage: []string{"ggc diff unstaged"}},
				{Name: "diff staged", Summary: "Show staged changes", Usage: []string{"ggc diff staged"}},
				{Name: "diff head", Summary: "Alias for default diff against HEAD", Usage: []string{"ggc diff head"}},
			},
		},
	}
}
