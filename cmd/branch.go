// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bmf-san/ggc/v4/git"
)

// Brancher provides functionality for the branch command.
type Brancher struct {
	gitClient    git.Clienter
	inputReader  *bufio.Reader
	outputWriter io.Writer
	helper       *Helper
}

// NewBrancher creates a new Brancher.
func NewBrancher(client git.Clienter) *Brancher {
	return &Brancher{
		gitClient:    client,
		inputReader:  bufio.NewReader(os.Stdin),
		outputWriter: os.Stdout,
		helper:       NewHelper(),
	}
}

// Branch executes the branch command with the given arguments.
func (b *Brancher) Branch(args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "current":
			branch, err := b.gitClient.GetCurrentBranch()
			if err != nil {
				_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
				return
			}
			_, _ = fmt.Fprintln(b.outputWriter, branch)
			return
		case "checkout":
			b.branchCheckout()
			return
		case "checkout-remote":
			b.branchCheckoutRemote()
			return
		case "create":
			b.branchCreate()
			return
		case "delete":
			b.branchDelete()
			return
		case "delete-merged":
			b.branchDeleteMerged()
			return
		case "rename":
			b.branchRename()
			return
		case "move":
			b.branchMove()
			return
		case "set-upstream":
			b.branchSetUpstream()
			return
		case "info":
			b.branchInfo()
			return
		case "list":
			// Support list --verbose
			if len(args) > 1 && (args[1] == "--verbose" || args[1] == "-v") {
				b.branchListVerbose()
				return
			}
		case "sort":
			b.branchSort()
			return
		case "contains":
			b.branchContains()
			return
		}
	}
	b.helper.ShowBranchHelp()
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
	_, _ = fmt.Fprintln(b.outputWriter, "Local branches:")
	for i, br := range branches {
		_, _ = fmt.Fprintf(b.outputWriter, "[%d] %s\n", i+1, br)
	}
	_, _ = fmt.Fprint(b.outputWriter, "Enter the number to checkout: ")
	input, _ := b.inputReader.ReadString('\n')
	input = strings.TrimSpace(input)
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(branches) {
		_, _ = fmt.Fprintln(b.outputWriter, "Invalid number.")
		return
	}
	branch := branches[idx-1]
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
	_, _ = fmt.Fprintln(b.outputWriter, "Remote branches:")
	for i, br := range branches {
		_, _ = fmt.Fprintf(b.outputWriter, "[%d] %s\n", i+1, br)
	}
	_, _ = fmt.Fprint(b.outputWriter, "Enter the number to checkout: ")
	input, _ := b.inputReader.ReadString('\n')
	input = strings.TrimSpace(input)
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(branches) {
		_, _ = fmt.Fprintln(b.outputWriter, "Invalid number.")
		return
	}
	remoteBranch := branches[idx-1]
	// origin/feature/foo → feature/foo
	parts := strings.SplitN(remoteBranch, "/", 2)
	if len(parts) != 2 {
		_, _ = fmt.Fprintln(b.outputWriter, "Invalid remote branch name.")
		return
	}
	localBranch := parts[1]
	if strings.TrimSpace(localBranch) == "" {
		_, _ = fmt.Fprintln(b.outputWriter, "Invalid remote branch name.")
		return
	}
	if err := validateBranchName(localBranch); err != nil {
		// Keep message consistent with existing expectations
		_, _ = fmt.Fprintln(b.outputWriter, "Invalid remote branch name.")
		return
	}
	if err := b.gitClient.CheckoutNewBranchFromRemote(localBranch, remoteBranch); err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
	}
}

func (b *Brancher) branchCreate() {
	_, _ = fmt.Fprint(b.outputWriter, "Enter new branch name: ")
	input, _ := b.inputReader.ReadString('\n')
	branchName := strings.TrimSpace(input)
	if branchName == "" {
		_, _ = fmt.Fprintln(b.outputWriter, "Cancelled.")
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
	if len(branches) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No local branches found.")
		return
	}
	for {
		_, _ = fmt.Fprintln(b.outputWriter, "\033[1;36mSelect local branches to delete by number (space separated, all: select all, none: deselect all, e.g. 1 3 5):\033[0m")
		for i, br := range branches {
			_, _ = fmt.Fprintf(b.outputWriter, "  [\033[1;33m%d\033[0m] %s\n", i+1, br)
		}
		_, _ = fmt.Fprint(b.outputWriter, "> ")
		input, _ := b.inputReader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			_, _ = fmt.Fprintln(b.outputWriter, "Cancelled.")
			return
		}
		if input == "all" {
			for _, br := range branches {
				if err := b.gitClient.DeleteBranch(br); err != nil {
					_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
				}
			}
			_, _ = fmt.Fprintln(b.outputWriter, "All branches deleted.")
			break
		}
		if input == "none" {
			continue
		}
		indices := strings.Fields(input)
		tmp := []string{}
		valid := true
		for _, idx := range indices {
			n, err := strconv.Atoi(idx)
			if err != nil || n < 1 || n > len(branches) {
				_, _ = fmt.Fprintf(b.outputWriter, "\033[1;31mInvalid number: %s\033[0m\n", idx)
				valid = false
				break
			}
			tmp = append(tmp, branches[n-1])
		}
		if !valid {
			continue
		}
		for _, br := range tmp {
			if err := b.gitClient.DeleteBranch(br); err != nil {
				_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
			}
		}
		_, _ = fmt.Fprintln(b.outputWriter, "Selected branches deleted.")
		break
	}
}

func (b *Brancher) branchDeleteMerged() {
	current, err := b.gitClient.GetCurrentBranch()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: failed to get current branch: %v\n", err)
		return
	}
	mergedBranches, err := b.gitClient.ListMergedBranches()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}

	branches := []string{}
	for _, br := range mergedBranches {
		if br != current {
			branches = append(branches, br)
		}
	}
	if len(branches) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No merged local branches.")
		return
	}
	for {
		_, _ = fmt.Fprintln(b.outputWriter, "\033[1;36mSelect merged local branches to delete by number (space separated, all: select all, none: deselect all, e.g. 1 3 5):\033[0m")
		for i, br := range branches {
			_, _ = fmt.Fprintf(b.outputWriter, "  [\033[1;33m%d\033[0m] %s\n", i+1, br)
		}
		_, _ = fmt.Fprint(b.outputWriter, "> ")
		input, _ := b.inputReader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			_, _ = fmt.Fprintln(b.outputWriter, "Cancelled.")
			return
		}
		if input == "all" {
			for _, br := range branches {
				if err := b.gitClient.DeleteBranch(br); err != nil {
					_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
				}
			}
			_, _ = fmt.Fprintln(b.outputWriter, "All merged branches deleted.")
			break
		}
		if input == "none" {
			continue
		}
		indices := strings.Fields(input)
		tmp := []string{}
		valid := true
		for _, idx := range indices {
			n, err := strconv.Atoi(idx)
			if err != nil || n < 1 || n > len(branches) {
				_, _ = fmt.Fprintf(b.outputWriter, "\033[1;31mInvalid number: %s\033[0m\n", idx)
				valid = false
				break
			}
			tmp = append(tmp, branches[n-1])
		}
		if !valid {
			continue
		}
		for _, br := range tmp {
			if err := b.gitClient.DeleteBranch(br); err != nil {
				_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
			}
		}
		_, _ = fmt.Fprintln(b.outputWriter, "Selected merged branches deleted.")
		break
	}
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
	_, _ = fmt.Fprintln(b.outputWriter, "Local branches:")
	for i, br := range branches {
		_, _ = fmt.Fprintf(b.outputWriter, "[%d] %s\n", i+1, br)
	}
	_, _ = fmt.Fprint(b.outputWriter, "Enter the number of the branch to rename: ")
	input, _ := b.inputReader.ReadString('\n')
	input = strings.TrimSpace(input)
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(branches) {
		_, _ = fmt.Fprintln(b.outputWriter, "Invalid number.")
		return
	}
	oldName := branches[idx-1]
	_, _ = fmt.Fprint(b.outputWriter, "Enter new branch name: ")
	newInput, _ := b.inputReader.ReadString('\n')
	newName := strings.TrimSpace(newInput)
	if newName == "" {
		_, _ = fmt.Fprintln(b.outputWriter, "Cancelled.")
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
	_, _ = fmt.Fprintln(b.outputWriter, "Local branches:")
	for i, br := range branches {
		_, _ = fmt.Fprintf(b.outputWriter, "[%d] %s\n", i+1, br)
	}
	_, _ = fmt.Fprint(b.outputWriter, "Enter the number of the branch to move: ")
	input, _ := b.inputReader.ReadString('\n')
	input = strings.TrimSpace(input)
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(branches) {
		_, _ = fmt.Fprintln(b.outputWriter, "Invalid number.")
		return
	}
	branch := branches[idx-1]
	_, _ = fmt.Fprint(b.outputWriter, "Enter commit or ref to move to: ")
	commit, _ := b.inputReader.ReadString('\n')
	commit = strings.TrimSpace(commit)
	if commit == "" {
		_, _ = fmt.Fprintln(b.outputWriter, "Cancelled.")
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
	_, _ = fmt.Fprintln(b.outputWriter, "Local branches:")
	for i, br := range branches {
		_, _ = fmt.Fprintf(b.outputWriter, "[%d] %s\n", i+1, br)
	}
	_, _ = fmt.Fprint(b.outputWriter, "Enter the number of the branch to set upstream: ")
	input, _ := b.inputReader.ReadString('\n')
	input = strings.TrimSpace(input)
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(branches) {
		_, _ = fmt.Fprintln(b.outputWriter, "Invalid number.")
		return
	}
	branch := branches[idx-1]

	remotes, _ := b.gitClient.ListRemoteBranches()
	if len(remotes) > 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "Remote branches:")
		for i, rb := range remotes {
			_, _ = fmt.Fprintf(b.outputWriter, "[%d] %s\n", i+1, rb)
		}
	}
	_, _ = fmt.Fprint(b.outputWriter, "Enter upstream (name or number): ")
	upIn, _ := b.inputReader.ReadString('\n')
	upIn = strings.TrimSpace(upIn)
	if upIn == "" {
		_, _ = fmt.Fprintln(b.outputWriter, "Cancelled.")
		return
	}
	// If numeric and valid index, map to remote branch
	if id, e := strconv.Atoi(upIn); e == nil && id >= 1 && id <= len(remotes) {
		upIn = remotes[id-1]
	}
	if err := b.gitClient.SetUpstreamBranch(branch, upIn); err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
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
	_, _ = fmt.Fprintln(b.outputWriter, "Local branches:")
	for i, br := range branches {
		_, _ = fmt.Fprintf(b.outputWriter, "[%d] %s\n", i+1, br)
	}
	_, _ = fmt.Fprint(b.outputWriter, "Enter the number to show info: ")
	input, _ := b.inputReader.ReadString('\n')
	input = strings.TrimSpace(input)
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(branches) {
		_, _ = fmt.Fprintln(b.outputWriter, "Invalid number.")
		return
	}
	branch := branches[idx-1]
	bi, err := b.gitClient.GetBranchInfo(branch)
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	// Display detailed info
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

func (b *Brancher) branchSort() {
	_, _ = fmt.Fprintln(b.outputWriter, "Sort by:")
	_, _ = fmt.Fprintln(b.outputWriter, "[1] name")
	_, _ = fmt.Fprintln(b.outputWriter, "[2] date")
	_, _ = fmt.Fprint(b.outputWriter, "Enter number: ")
	input, _ := b.inputReader.ReadString('\n')
	input = strings.TrimSpace(input)
	by := "name"
	if input == "2" {
		by = "date"
	}
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
	_, _ = fmt.Fprint(b.outputWriter, "Enter commit or ref: ")
	input, _ := b.inputReader.ReadString('\n')
	commit := strings.TrimSpace(input)
	if commit == "" {
		_, _ = fmt.Fprintln(b.outputWriter, "Cancelled.")
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
