// Package command provides command metadata and registry functionality for ggc.
package command

// Category groups commands for help and discovery surfaces.
type Category string

// Command categories for organizing commands in help and discovery surfaces.
const (
	CategoryBasics  Category = "Basics"
	CategoryBranch  Category = "Branch"
	CategoryCommit  Category = "Commit"
	CategoryRemote  Category = "Remote"
	CategoryStatus  Category = "Status"
	CategoryCleanup Category = "Cleanup"
	CategoryDiff    Category = "Diff"
	CategoryTag     Category = "Tag"
	CategoryConfig  Category = "Config"
	CategoryHook    Category = "Hook"
	CategoryRebase  Category = "Rebase"
	CategoryStash   Category = "Stash"
	CategoryUtility Category = "Utility"
)
