package cmd

import (
	"fmt"
	"strconv"
	"strings"
)

func (b *Brancher) branchDeleteArgs(args []string) {
	if len(args) > 0 {
		b.deleteBranchesFromArgs(args)
		return
	}

	branches, ok := b.collectDeletableBranches()
	if !ok {
		return
	}

	if len(branches) == 0 {
		WriteLine(b.outputWriter, "No local branches found.")
		return
	}

	b.runBranchDeleteLoop(branches)
}

func (b *Brancher) deleteBranchesFromArgs(args []string) {
	current, _ := b.gitClient.GetCurrentBranch()
	for _, a := range args {
		br := strings.TrimSpace(a)
		if br == "" {
			continue
		}
		if current != "" && br == current {
			WriteLinef(b.outputWriter, "Skipping current branch: %s", br)
			continue
		}
		if err := b.gitClient.DeleteBranch(br); err != nil {
			WriteError(b.outputWriter, err)
		}
	}
}

func (b *Brancher) collectDeletableBranches() ([]string, bool) {
	branches, err := b.gitClient.ListLocalBranches()
	if err != nil {
		WriteError(b.outputWriter, err)
		return nil, false
	}

	if curr, err := b.gitClient.GetCurrentBranch(); err == nil && curr != "" {
		filtered := make([]string, 0, len(branches))
		for _, br := range branches {
			if br != curr {
				filtered = append(filtered, br)
			}
		}
		branches = filtered
	}

	return branches, true
}

// runBranchDeleteLoop runs the interactive branch deletion loop
func (b *Brancher) runBranchDeleteLoop(branches []string) {
	for {
		b.displayBranchSelection(branches)
		input, ok := b.readLine("")
		if !ok {
			return
		}
		input = strings.TrimSpace(input)

		if input == "" {
			WriteLine(b.outputWriter, "Canceled.")
			return
		}
		if b.handleBranchSpecialCommands(input, branches) {
			return
		}
		if b.handleBranchSelection(input, branches) {
			return
		}
	}
}

// displayBranchSelection shows the branch selection interface
func (b *Brancher) displayBranchSelection(branches []string) {
	WriteLine(b.outputWriter, "\033[1;36mSelect local branches to delete by number (space separated, all: select all, none: deselect all, e.g. 1 3 5):\033[0m")
	for i, br := range branches {
		WriteLinef(b.outputWriter, "  [\033[1;33m%d\033[0m] %s", i+1, br)
	}
	_, _ = fmt.Fprint(b.outputWriter, "> ")
}

// handleBranchSpecialCommands processes "all" and "none" commands for branches
func (b *Brancher) handleBranchSpecialCommands(input string, branches []string) bool {
	if input == "all" {
		for _, br := range branches {
			if err := b.gitClient.DeleteBranch(br); err != nil {
				WriteError(b.outputWriter, err)
			}
		}
		WriteLine(b.outputWriter, "All branches deleted.")
		return true
	}
	if input == "none" {
		return false // Continue loop
	}
	return false
}

// handleBranchSelection processes numeric branch selection
func (b *Brancher) handleBranchSelection(input string, branches []string) bool {
	selectedBranches, valid := b.parseBranchIndices(input, branches)
	if !valid {
		return false // Continue loop
	}

	for _, br := range selectedBranches {
		if err := b.gitClient.DeleteBranch(br); err != nil {
			WriteError(b.outputWriter, err)
		}
	}
	WriteLine(b.outputWriter, "Selected branches deleted.")
	return true
}

// parseBranchIndices parses user input into selected branches
func (b *Brancher) parseBranchIndices(input string, branches []string) ([]string, bool) {
	indices := strings.Fields(input)
	selectedBranches := []string{}

	for _, idx := range indices {
		n, err := strconv.Atoi(idx)
		if err != nil || n < 1 || n > len(branches) {
			WriteLinef(b.outputWriter, "\033[1;31mInvalid number: %s\033[0m", idx)
			return nil, false
		}
		selectedBranches = append(selectedBranches, branches[n-1])
	}
	return selectedBranches, true
}

func (b *Brancher) branchDeleteMerged() {
	branches, err := b.getMergedBranchesForDeletion()
	if err != nil {
		WriteError(b.outputWriter, err)
		return
	}
	if len(branches) == 0 {
		WriteLine(b.outputWriter, "No merged local branches.")
		return
	}

	b.runMergedBranchDeleteLoop(branches)
}

// getMergedBranchesForDeletion gets the list of merged branches that can be deleted
func (b *Brancher) getMergedBranchesForDeletion() ([]string, error) {
	current, err := b.gitClient.GetCurrentBranch()
	if err != nil {
		return nil, fmt.Errorf("failed to get current branch: %w", err)
	}

	mergedBranches, err := b.gitClient.ListMergedBranches()
	if err != nil {
		return nil, err
	}

	branches := []string{}
	for _, br := range mergedBranches {
		if br != current {
			branches = append(branches, br)
		}
	}
	return branches, nil
}

// runMergedBranchDeleteLoop runs the interactive branch deletion loop
func (b *Brancher) runMergedBranchDeleteLoop(branches []string) {
	for {
		b.displayMergedBranchSelection(branches)
		input, ok := b.readLine("")
		if !ok {
			return
		}
		input = strings.TrimSpace(input)

		if input == "" {
			WriteLine(b.outputWriter, "Canceled.")
			return
		}
		if b.handleMergedBranchSpecialCommands(input, branches) {
			return
		}
		if b.handleMergedBranchSelection(input, branches) {
			return
		}
	}
}

// displayMergedBranchSelection shows the merged branch selection interface
func (b *Brancher) displayMergedBranchSelection(branches []string) {
	WriteLine(b.outputWriter, "\033[1;36mSelect merged local branches to delete by number (space separated, all: select all, none: deselect all, e.g. 1 3 5):\033[0m")
	for i, br := range branches {
		WriteLinef(b.outputWriter, "  [\033[1;33m%d\033[0m] %s", i+1, br)
	}
	_, _ = fmt.Fprint(b.outputWriter, "> ")
}

// handleMergedBranchSpecialCommands processes "all" and "none" commands for merged branches
func (b *Brancher) handleMergedBranchSpecialCommands(input string, branches []string) bool {
	if input == "all" {
		for _, br := range branches {
			if err := b.gitClient.DeleteBranch(br); err != nil {
				WriteError(b.outputWriter, err)
			}
		}
		WriteLine(b.outputWriter, "All merged branches deleted.")
		return true
	}
	if input == "none" {
		return false // Continue loop
	}
	return false
}

// handleMergedBranchSelection processes numeric merged branch selection
func (b *Brancher) handleMergedBranchSelection(input string, branches []string) bool {
	selectedBranches, valid := b.parseBranchIndices(input, branches)
	if !valid {
		return false // Continue loop
	}

	for _, br := range selectedBranches {
		if err := b.gitClient.DeleteBranch(br); err != nil {
			WriteError(b.outputWriter, err)
		}
	}
	WriteLine(b.outputWriter, "Selected merged branches deleted.")
	return true
}
