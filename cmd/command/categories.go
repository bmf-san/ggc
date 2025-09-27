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

// CategoryOrder returns the display order for command categories.
func CategoryOrder(cat Category) int {
	order := map[Category]int{
		CategoryBasics:  1,
		CategoryBranch:  2,
		CategoryCommit:  3,
		CategoryRemote:  4,
		CategoryStatus:  5,
		CategoryCleanup: 6,
		CategoryDiff:    7,
		CategoryTag:     8,
		CategoryConfig:  9,
		CategoryHook:    10,
		CategoryRebase:  11,
		CategoryStash:   12,
		CategoryUtility: 13,
	}
	if o, exists := order[cat]; exists {
		return o
	}
	return 999 // Unknown categories at the end
}

// OrderedCategories returns all categories in their display order.
func OrderedCategories() []Category {
	return []Category{
		CategoryBasics,
		CategoryBranch,
		CategoryCommit,
		CategoryRemote,
		CategoryStatus,
		CategoryCleanup,
		CategoryDiff,
		CategoryTag,
		CategoryConfig,
		CategoryHook,
		CategoryRebase,
		CategoryStash,
		CategoryUtility,
	}
}
