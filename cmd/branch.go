// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bmf-san/ggc/v6/git"
	"github.com/bmf-san/ggc/v6/internal/prompt"
)

// Brancher provides functionality for the branch command.
type Brancher struct {
	gitClient    git.BranchOps
	prompter     prompt.Interface
	outputWriter io.Writer
	helper       *Helper
}

// NewBrancher creates a new Brancher.
func NewBrancher(client git.BranchOps) *Brancher {
	output := os.Stdout
	helper := NewHelper()
	helper.outputWriter = output
	return &Brancher{
		gitClient:    client,
		prompter:     prompt.New(os.Stdin, output),
		outputWriter: output,
		helper:       helper,
	}
}

// Branch executes the branch command with the given arguments.
func (b *Brancher) Branch(args []string) {
	if len(args) == 0 {
		b.helper.ShowBranchHelp()
		return
	}

	b.handleBranchCommand(args[0], args[1:])
}

// handleBranchCommand processes the specific branch subcommand
func (b *Brancher) handleBranchCommand(cmd string, args []string) {
	branchCommands := map[string]func([]string){
		"current":  func([]string) { b.handleCurrentBranch() },
		"checkout": b.handleCheckoutCommand,
		"create":   func([]string) { b.branchCreate() },
		"delete":   b.handleDeleteCommand,
		"rename":   func([]string) { b.branchRename() },
		"move":     func([]string) { b.branchMove() },
		"set":      b.handleSetCommand,
		"info":     func([]string) { b.branchInfo() },
		"list":     b.handleListCommand,
		"sort":     func([]string) { b.branchSort() },
		"contains": func([]string) { b.branchContains() },
	}

	if handler, exists := branchCommands[cmd]; exists {
		handler(args)
		return
	}
	b.helper.ShowBranchHelp()
}

// handleCurrentBranch shows the current branch
func (b *Brancher) handleCurrentBranch() {
	branch, err := b.gitClient.GetCurrentBranch()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	_, _ = fmt.Fprintln(b.outputWriter, branch)
}

// handleCheckoutCommand handles checkout subcommand
func (b *Brancher) handleCheckoutCommand(args []string) {
	if len(args) > 0 && args[0] == "remote" {
		b.branchCheckoutRemote()
	} else {
		b.branchCheckout()
	}
}

// handleDeleteCommand handles delete subcommand
func (b *Brancher) handleDeleteCommand(args []string) {
	if len(args) > 0 && args[0] == "merged" {
		b.branchDeleteMerged()
	} else {
		b.branchDelete()
	}
}

// handleSetCommand handles set subcommand
func (b *Brancher) handleSetCommand(args []string) {
	if len(args) > 0 && args[0] == "upstream" {
		b.branchSetUpstream()
		return
	}
	b.helper.ShowBranchHelp()
}

// handleListCommand handles list subcommand
func (b *Brancher) handleListCommand(args []string) {
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "verbose", "--verbose", "-v":
		b.branchListVerbose()
	case "local":
		b.branchListLocal()
	case "remote":
		b.branchListRemote()
	}
}

func (b *Brancher) branchCheckout() {
	branches, err := b.gitClient.ListLocalBranches()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	if len(branches) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No local branches found.")
		return
	}
	idx, ok := b.promptSelectIndex("Local branches:", branches, "Enter the number to checkout: ")
	if !ok {
		return
	}
	branch := branches[idx]
	if err := b.gitClient.CheckoutBranch(branch); err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
	}
}

func (b *Brancher) branchCheckoutRemote() {
	branches, err := b.gitClient.ListRemoteBranches()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	if len(branches) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No remote branches found.")
		return
	}
	idx, ok := b.promptSelectIndex("Remote branches:", branches, "Enter the number to checkout: ")
	if !ok {
		return
	}
	remoteBranch := branches[idx]
	localBranch, valid := deriveLocalFromRemote(remoteBranch)
	if !valid || validateBranchName(localBranch) != nil {
		_, _ = fmt.Fprintln(b.outputWriter, "Invalid remote branch name.")
		return
	}
	if err := b.gitClient.CheckoutNewBranchFromRemote(localBranch, remoteBranch); err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
	}
}

// promptSelectIndex prints a list with title and asks for selection, returns 0-based index
func (b *Brancher) promptSelectIndex(title string, items []string, promptText string) (int, bool) {
	if b.prompter == nil {
		return 0, false
	}
	idx, canceled, err := b.prompter.Select(title, items, promptText)
	if canceled {
		return 0, false
	}
	if err != nil {
		if errors.Is(err, prompt.ErrInvalidSelection) {
			_, _ = fmt.Fprintln(b.outputWriter, "Invalid number.")
		} else {
			_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		}
		return 0, false
	}
	return idx, true
}

// deriveLocalFromRemote converts "origin/foo" -> "foo"
func deriveLocalFromRemote(remoteBranch string) (string, bool) {
	parts := strings.SplitN(remoteBranch, "/", 2)
	if len(parts) != 2 {
		return "", false
	}
	local := strings.TrimSpace(parts[1])
	if local == "" {
		return "", false
	}
	return local, true
}

func (b *Brancher) readLine(promptText string) (string, bool) {
	if b.prompter == nil {
		return "", false
	}
	line, canceled, err := b.prompter.Input(promptText)
	if canceled {
		return "", false
	}
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return "", false
	}
	return line, true
}

func (b *Brancher) branchCreate() {
	input, ok := b.readLine("Enter new branch name: ")
	if !ok {
		return
	}
	branchName := strings.TrimSpace(input)
	if branchName == "" {
		_, _ = fmt.Fprintln(b.outputWriter, "Canceled.")
		return
	}
	if err := validateBranchName(branchName); err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: invalid branch name: %v\n", err)
		return
	}

	if err := b.gitClient.CheckoutNewBranch(branchName); err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: failed to create and checkout branch: %v\n", err)
		return
	}
}

func (b *Brancher) branchDelete() {
	branches, err := b.gitClient.ListLocalBranches()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	// Exclude current branch from deletion candidates to avoid failing UX
	if curr, err := b.gitClient.GetCurrentBranch(); err == nil {
		filtered := make([]string, 0, len(branches))
		for _, br := range branches {
			if br != curr {
				filtered = append(filtered, br)
			}
		}
		branches = filtered
	}
	if len(branches) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No local branches found.")
		return
	}

	b.runBranchDeleteLoop(branches)
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
			_, _ = fmt.Fprintln(b.outputWriter, "Canceled.")
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
	_, _ = fmt.Fprintln(b.outputWriter, "\033[1;36mSelect local branches to delete by number (space separated, all: select all, none: deselect all, e.g. 1 3 5):\033[0m")
	for i, br := range branches {
		_, _ = fmt.Fprintf(b.outputWriter, "  [\033[1;33m%d\033[0m] %s\n", i+1, br)
	}
	_, _ = fmt.Fprint(b.outputWriter, "> ")
}

// handleBranchSpecialCommands processes "all" and "none" commands for branches
func (b *Brancher) handleBranchSpecialCommands(input string, branches []string) bool {
	if input == "all" {
		for _, br := range branches {
			if err := b.gitClient.DeleteBranch(br); err != nil {
				_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
			}
		}
		_, _ = fmt.Fprintln(b.outputWriter, "All branches deleted.")
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
			_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		}
	}
	_, _ = fmt.Fprintln(b.outputWriter, "Selected branches deleted.")
	return true
}

// parseBranchIndices parses user input into selected branches
func (b *Brancher) parseBranchIndices(input string, branches []string) ([]string, bool) {
	indices := strings.Fields(input)
	selectedBranches := []string{}

	for _, idx := range indices {
		n, err := strconv.Atoi(idx)
		if err != nil || n < 1 || n > len(branches) {
			_, _ = fmt.Fprintf(b.outputWriter, "\033[1;31mInvalid number: %s\033[0m\n", idx)
			return nil, false
		}
		selectedBranches = append(selectedBranches, branches[n-1])
	}
	return selectedBranches, true
}

func (b *Brancher) branchDeleteMerged() {
	branches, err := b.getMergedBranchesForDeletion()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	if len(branches) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No merged local branches.")
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
			_, _ = fmt.Fprintln(b.outputWriter, "Canceled.")
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
	_, _ = fmt.Fprintln(b.outputWriter, "\033[1;36mSelect merged local branches to delete by number (space separated, all: select all, none: deselect all, e.g. 1 3 5):\033[0m")
	for i, br := range branches {
		_, _ = fmt.Fprintf(b.outputWriter, "  [\033[1;33m%d\033[0m] %s\n", i+1, br)
	}
	_, _ = fmt.Fprint(b.outputWriter, "> ")
}

// handleMergedBranchSpecialCommands processes "all" and "none" commands for merged branches
func (b *Brancher) handleMergedBranchSpecialCommands(input string, branches []string) bool {
	if input == "all" {
		for _, br := range branches {
			if err := b.gitClient.DeleteBranch(br); err != nil {
				_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
			}
		}
		_, _ = fmt.Fprintln(b.outputWriter, "All merged branches deleted.")
		return true
	}
	if input == "none" {
		return false // Continue loop
	}
	return false
}

// handleMergedBranchSelection processes numeric merged branch selection
func (b *Brancher) handleMergedBranchSelection(input string, branches []string) bool {
	selectedBranches, valid := b.parseMergedBranchIndices(input, branches)
	if !valid {
		return false // Continue loop
	}

	for _, br := range selectedBranches {
		if err := b.gitClient.DeleteBranch(br); err != nil {
			_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		}
	}
	_, _ = fmt.Fprintln(b.outputWriter, "Selected merged branches deleted.")
	return true
}

// parseMergedBranchIndices parses user input into selected merged branches
func (b *Brancher) parseMergedBranchIndices(input string, branches []string) ([]string, bool) {
	indices := strings.Fields(input)
	selectedBranches := []string{}

	for _, idx := range indices {
		n, err := strconv.Atoi(idx)
		if err != nil || n < 1 || n > len(branches) {
			_, _ = fmt.Fprintf(b.outputWriter, "\033[1;31mInvalid number: %s\033[0m\n", idx)
			return nil, false
		}
		selectedBranches = append(selectedBranches, branches[n-1])
	}
	return selectedBranches, true
}

// validateBranchName performs basic validation aligned with git ref rules for branch names.
// It rejects empty names, control characters, disallowed characters/sequences, invalid prefixes/suffixes,
// double slashes, overly long names, and non-ASCII for safety across platforms.
func validateBranchName(name string) error {
	n := strings.TrimSpace(name)
	if n == "" {
		return fmt.Errorf("branch name cannot be empty")
	}
	// Delegate validation to git to match exact refname rules.
	// Equivalent to: git check-ref-format --branch <name>
	cmd := exec.Command("git", "check-ref-format", "--branch", n)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("invalid per git check-ref-format: %w", err)
	}
	return nil
}

func (b *Brancher) branchRename() {
	branches, err := b.gitClient.ListLocalBranches()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	if len(branches) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No local branches found.")
		return
	}
	idx, ok := b.promptSelectIndex("Local branches:", branches, "Enter the number of the branch to rename: ")
	if !ok {
		return
	}
	oldName := branches[idx]
	newInput, ok := b.readLine("Enter new branch name: ")
	if !ok {
		return
	}
	newName := strings.TrimSpace(newInput)
	if newName == "" {
		_, _ = fmt.Fprintln(b.outputWriter, "Canceled.")
		return
	}
	if err := validateBranchName(newName); err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: invalid branch name: %v\n", err)
		return
	}
	if err := b.gitClient.RenameBranch(oldName, newName); err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
}

func (b *Brancher) branchMove() {
	branches, err := b.gitClient.ListLocalBranches()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	if len(branches) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No local branches found.")
		return
	}
	idx, ok := b.promptSelectIndex("Local branches:", branches, "Enter the number of the branch to move: ")
	if !ok {
		return
	}
	branch := branches[idx]
	commitInput, ok := b.readLine("Enter commit or ref to move to: ")
	if !ok {
		return
	}
	commit := strings.TrimSpace(commitInput)
	if commit == "" {
		_, _ = fmt.Fprintln(b.outputWriter, "Canceled.")
		return
	}
	if !b.gitClient.RevParseVerify(commit) {
		_, _ = fmt.Fprintln(b.outputWriter, "Invalid commit or ref.")
		return
	}
	if err := b.gitClient.MoveBranch(branch, commit); err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
}

func (b *Brancher) branchSetUpstream() {
	branches, err := b.gitClient.ListLocalBranches()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	if len(branches) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No local branches found.")
		return
	}

	branch := b.selectLocalBranch(branches)
	if branch == "" {
		return
	}

	upstream := b.selectUpstreamBranch()
	if upstream == "" {
		return
	}

	if err := b.gitClient.SetUpstreamBranch(branch, upstream); err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
	}
}

// selectLocalBranch prompts user to select a local branch
func (b *Brancher) selectLocalBranch(branches []string) string {
	idx, ok := b.promptSelectIndex("Local branches:", branches, "Enter the number of the branch to set upstream: ")
	if !ok {
		return ""
	}
	return branches[idx]
}

// selectUpstreamBranch prompts user to select an upstream branch
func (b *Brancher) selectUpstreamBranch() string {
	remotes, _ := b.gitClient.ListRemoteBranches()
	if len(remotes) > 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "Remote branches:")
		for i, rb := range remotes {
			_, _ = fmt.Fprintf(b.outputWriter, "[%d] %s\n", i+1, rb)
		}
	}
	upIn, ok := b.readLine("Enter upstream (name or number): ")
	if !ok {
		return ""
	}
	upIn = strings.TrimSpace(upIn)
	if upIn == "" {
		_, _ = fmt.Fprintln(b.outputWriter, "Canceled.")
		return ""
	}
	// If numeric and valid index, map to remote branch
	if id, e := strconv.Atoi(upIn); e == nil && id >= 1 && id <= len(remotes) {
		upIn = remotes[id-1]
	}
	return upIn
}

func (b *Brancher) branchInfo() {
	branches, err := b.gitClient.ListLocalBranches()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	if len(branches) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No local branches found.")
		return
	}
	idx, ok := b.promptSelectIndex("Local branches:", branches, "Enter the number to show info: ")
	if !ok {
		return
	}
	br := branches[idx]
	b.printBranchInfo(br)
}

func (b *Brancher) printBranchInfo(branch string) {
	bi, err := b.gitClient.GetBranchInfo(branch)
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	_, _ = fmt.Fprintf(b.outputWriter, "Name: %s\n", bi.Name)
	_, _ = fmt.Fprintf(b.outputWriter, "Current: %t\n", bi.IsCurrentBranch)
	if bi.Upstream != "" {
		_, _ = fmt.Fprintf(b.outputWriter, "Upstream: %s\n", bi.Upstream)
	}
	if bi.AheadBehind != "" {
		_, _ = fmt.Fprintf(b.outputWriter, "Ahead/Behind: %s\n", bi.AheadBehind)
	}
	if bi.LastCommitSHA != "" {
		_, _ = fmt.Fprintf(b.outputWriter, "Last Commit: %s %s\n", bi.LastCommitSHA, bi.LastCommitMsg)
	} else if bi.LastCommitMsg != "" {
		_, _ = fmt.Fprintf(b.outputWriter, "Last Commit: %s\n", bi.LastCommitMsg)
	}
}

func (b *Brancher) branchListVerbose() {
	infos, err := b.gitClient.ListBranchesVerbose()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	if len(infos) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No local branches found.")
		return
	}
	for _, bi := range infos {
		marker := " "
		if bi.IsCurrentBranch {
			marker = "*"
		}
		extra := ""
		if bi.Upstream != "" && bi.AheadBehind != "" {
			extra = fmt.Sprintf(" [%s: %s]", bi.Upstream, bi.AheadBehind)
		} else if bi.Upstream != "" {
			extra = fmt.Sprintf(" [%s]", bi.Upstream)
		}
		_, _ = fmt.Fprintf(b.outputWriter, "%s %s %s%s %s\n", marker, bi.Name, bi.LastCommitSHA, extra, bi.LastCommitMsg)
	}
}

func (b *Brancher) branchListLocal() {
	branches, err := b.gitClient.ListLocalBranches()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	if len(branches) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No local branches found.")
		return
	}
	for _, br := range branches {
		_, _ = fmt.Fprintln(b.outputWriter, br)
	}
}

func (b *Brancher) branchListRemote() {
	branches, err := b.gitClient.ListRemoteBranches()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	if len(branches) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No remote branches found.")
		return
	}
	for _, br := range branches {
		_, _ = fmt.Fprintln(b.outputWriter, br)
	}
}

func (b *Brancher) branchSort() {
	opts := []string{"name", "date"}
	idx, ok := b.promptSelectIndex("Sort by:", opts, "Enter number: ")
	if !ok {
		return
	}
	by := opts[idx]
	names, err := b.gitClient.SortBranches(by)
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	for _, n := range names {
		_, _ = fmt.Fprintln(b.outputWriter, n)
	}
}

func (b *Brancher) branchContains() {
	input, ok := b.readLine("Enter commit or ref: ")
	if !ok {
		return
	}
	commit := strings.TrimSpace(input)
	if commit == "" {
		_, _ = fmt.Fprintln(b.outputWriter, "Canceled.")
		return
	}
	if !b.gitClient.RevParseVerify(commit) {
		_, _ = fmt.Fprintln(b.outputWriter, "Invalid commit or ref.")
		return
	}
	branches, err := b.gitClient.BranchesContaining(commit)
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	if len(branches) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No branches contain the specified commit.")
		return
	}
	for _, br := range branches {
		_, _ = fmt.Fprintln(b.outputWriter, br)
	}
}
