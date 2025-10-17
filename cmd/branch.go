// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bmf-san/ggc/v7/internal/prompt"
	"github.com/bmf-san/ggc/v7/pkg/git"
)

const errMsgBranchNameEmpty = "Error: branch name cannot be empty."

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
		"create":   b.branchCreate,
		"delete":   b.handleDeleteCommand,
		"rename":   b.branchRename,
		"move":     b.branchMove,
		"set":      b.handleSetCommand,
		"info":     b.branchInfo,
		"list":     b.handleListCommand,
		"sort":     b.branchSort,
		"contains": b.branchContains,
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
		b.branchDeleteArgs(args)
	}
}

// handleSetCommand handles set subcommand
func (b *Brancher) handleSetCommand(args []string) {
	if len(args) > 0 && args[0] == "upstream" {
		b.branchSetUpstream(args[1:])
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
	if !valid || git.ValidateBranchName(localBranch) != nil {
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

func (b *Brancher) branchCreate(args []string) {
	var branchName string
	if len(args) > 0 {
		branchName = strings.TrimSpace(args[0])
	} else {
		input, ok := b.readLine("Enter new branch name: ")
		if !ok {
			return
		}
		branchName = strings.TrimSpace(input)
		if branchName == "" {
			_, _ = fmt.Fprintln(b.outputWriter, "Canceled.")
			return
		}
	}
	if err := git.ValidateBranchName(branchName); err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: invalid branch name: %v\n", err)
		return
	}

	if err := b.gitClient.CheckoutNewBranch(branchName); err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: failed to create and checkout branch: %v\n", err)
		return
	}
}

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
		_, _ = fmt.Fprintln(b.outputWriter, "No local branches found.")
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
			_, _ = fmt.Fprintf(b.outputWriter, "Skipping current branch: %s\n", br)
			continue
		}
		if err := b.gitClient.DeleteBranch(br); err != nil {
			_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		}
	}
}

func (b *Brancher) collectDeletableBranches() ([]string, bool) {
	branches, err := b.gitClient.ListLocalBranches()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
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
	selected, ok := b.parseBranchIndices(input, branches)
	if !ok {
		return nil, false
	}
	return selected, true
}

func (b *Brancher) branchRename(args []string) {
	if len(args) >= 2 {
		oldName := strings.TrimSpace(args[0])
		newName := strings.TrimSpace(args[1])
		if oldName == "" {
			_, _ = fmt.Fprintln(b.outputWriter, errMsgBranchNameEmpty)
			return
		}
		if newName == "" {
			_, _ = fmt.Fprintln(b.outputWriter, "Error: new branch name cannot be empty.")
			return
		}
		if err := git.ValidateBranchName(newName); err != nil {
			_, _ = fmt.Fprintf(b.outputWriter, "Error: invalid branch name: %v\n", err)
			return
		}
		if err := b.gitClient.RenameBranch(oldName, newName); err != nil {
			_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		}
		return
	}

	b.branchRenameInteractive()
}

func (b *Brancher) branchRenameInteractive() {
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
	if err := git.ValidateBranchName(newName); err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: invalid branch name: %v\n", err)
		return
	}
	if err := b.gitClient.RenameBranch(oldName, newName); err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
}

func (b *Brancher) branchMove(args []string) {
	if len(args) >= 2 {
		branch := strings.TrimSpace(args[0])
		commit := strings.TrimSpace(args[1])
		if branch == "" {
			_, _ = fmt.Fprintln(b.outputWriter, errMsgBranchNameEmpty)
			return
		}
		if commit == "" {
			_, _ = fmt.Fprintln(b.outputWriter, "Error: commit or ref cannot be empty.")
			return
		}
		if !b.gitClient.RevParseVerify(commit) {
			_, _ = fmt.Fprintln(b.outputWriter, "Invalid commit or ref.")
			return
		}
		if err := b.gitClient.MoveBranch(branch, commit); err != nil {
			_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		}
		return
	}

	b.branchMoveInteractive()
}

func (b *Brancher) branchMoveInteractive() {
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

func (b *Brancher) branchSetUpstream(args []string) {
	switch len(args) {
	case 0:
		b.branchSetUpstreamInteractive()
	case 2:
		branch := strings.TrimSpace(args[0])
		if branch == "" {
			_, _ = fmt.Fprintln(b.outputWriter, errMsgBranchNameEmpty)
			return
		}
		upstream, ok := b.resolveUpstreamArgument(strings.TrimSpace(args[1]))
		if !ok {
			return
		}
		if err := b.gitClient.SetUpstreamBranch(branch, upstream); err != nil {
			_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		}
	default:
		_, _ = fmt.Fprintln(b.outputWriter, "Error: branch set upstream expects <branch> <upstream>.")
	}
}

func (b *Brancher) branchSetUpstreamInteractive() {
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

func (b *Brancher) resolveUpstreamArgument(input string) (string, bool) {
	if input == "" {
		_, _ = fmt.Fprintln(b.outputWriter, "Error: upstream cannot be empty.")
		return "", false
	}

	if idx, err := strconv.Atoi(input); err == nil {
		remotes, listErr := b.gitClient.ListRemoteBranches()
		if listErr != nil {
			_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", listErr)
			return "", false
		}
		if idx < 1 || idx > len(remotes) {
			_, _ = fmt.Fprintf(b.outputWriter, "Error: invalid remote selection: %d\n", idx)
			return "", false
		}
		return remotes[idx-1], true
	}

	return input, true
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

func (b *Brancher) branchInfo(args []string) {
	if len(args) > 1 {
		_, _ = fmt.Fprintln(b.outputWriter, "Error: branch info accepts at most one branch name.")
		return
	}

	if len(args) == 1 {
		branch := strings.TrimSpace(args[0])
		if branch == "" {
			_, _ = fmt.Fprintln(b.outputWriter, errMsgBranchNameEmpty)
			return
		}
		b.printBranchInfo(branch)
		return
	}

	b.branchInfoInteractive()
}

func (b *Brancher) branchInfoInteractive() {
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

func (b *Brancher) branchSort(args []string) {
	if len(args) > 1 {
		_, _ = fmt.Fprintln(b.outputWriter, "Error: branch sort accepts at most one option (name|date).")
		return
	}

	if len(args) == 1 {
		choice := strings.ToLower(strings.TrimSpace(args[0]))
		if choice == "" {
			_, _ = fmt.Fprintln(b.outputWriter, "Error: sort option cannot be empty.")
			return
		}
		if choice != "name" && choice != "date" {
			_, _ = fmt.Fprintf(b.outputWriter, "Error: invalid sort option %q. Use 'name' or 'date'.\n", args[0])
			return
		}
		b.printSortedBranches(choice)
		return
	}

	b.branchSortInteractive()
}

func (b *Brancher) branchSortInteractive() {
	opts := []string{"name", "date"}
	idx, ok := b.promptSelectIndex("Sort by:", opts, "Enter number: ")
	if !ok {
		return
	}
	by := opts[idx]
	b.printSortedBranches(by)
}

func (b *Brancher) printSortedBranches(by string) {
	names, err := b.gitClient.SortBranches(by)
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	if len(names) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No local branches found.")
		return
	}
	for _, n := range names {
		_, _ = fmt.Fprintln(b.outputWriter, n)
	}
}

func (b *Brancher) branchContains(args []string) {
	if len(args) > 1 {
		_, _ = fmt.Fprintln(b.outputWriter, "Error: branch contains accepts at most one commit or ref.")
		return
	}

	if len(args) == 1 {
		commit := strings.TrimSpace(args[0])
		if commit == "" {
			_, _ = fmt.Fprintln(b.outputWriter, "Error: commit or ref cannot be empty.")
			return
		}
		b.branchContainsForCommit(commit)
		return
	}

	b.branchContainsInteractive()
}

func (b *Brancher) branchContainsInteractive() {
	input, ok := b.readLine("Enter commit or ref: ")
	if !ok {
		return
	}
	commit := strings.TrimSpace(input)
	if commit == "" {
		_, _ = fmt.Fprintln(b.outputWriter, "Canceled.")
		return
	}
	b.branchContainsForCommit(commit)
}

func (b *Brancher) branchContainsForCommit(commit string) {
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
