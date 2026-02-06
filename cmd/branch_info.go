package cmd

import (
	"fmt"
	"strings"
)

func (b *Brancher) branchInfo(args []string) {
	if len(args) > 1 {
		WriteLine(b.outputWriter, "Error: branch info accepts at most one branch name.")
		return
	}

	if len(args) == 1 {
		branch := strings.TrimSpace(args[0])
		if branch == "" {
			WriteLine(b.outputWriter, errMsgBranchNameEmpty)
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
		WriteError(b.outputWriter, err)
		return
	}
	if len(branches) == 0 {
		WriteLine(b.outputWriter, "No local branches found.")
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
		WriteError(b.outputWriter, err)
		return
	}
	WriteLinef(b.outputWriter, "Name: %s", bi.Name)
	WriteLinef(b.outputWriter, "Current: %t", bi.IsCurrentBranch)
	if bi.Upstream != "" {
		WriteLinef(b.outputWriter, "Upstream: %s", bi.Upstream)
	}
	if bi.AheadBehind != "" {
		WriteLinef(b.outputWriter, "Ahead/Behind: %s", bi.AheadBehind)
	}
	if bi.LastCommitSHA != "" {
		WriteLinef(b.outputWriter, "Last Commit: %s %s", bi.LastCommitSHA, bi.LastCommitMsg)
	} else if bi.LastCommitMsg != "" {
		WriteLinef(b.outputWriter, "Last Commit: %s", bi.LastCommitMsg)
	}
}

func (b *Brancher) branchListVerbose() {
	infos, err := b.gitClient.ListBranchesVerbose()
	if err != nil {
		WriteError(b.outputWriter, err)
		return
	}
	if len(infos) == 0 {
		WriteLine(b.outputWriter, "No local branches found.")
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
		WriteLinef(b.outputWriter, "%s %s %s%s %s", marker, bi.Name, bi.LastCommitSHA, extra, bi.LastCommitMsg)
	}
}

func (b *Brancher) branchListLocal() {
	branches, err := b.gitClient.ListLocalBranches()
	if err != nil {
		WriteError(b.outputWriter, err)
		return
	}
	if len(branches) == 0 {
		WriteLine(b.outputWriter, "No local branches found.")
		return
	}
	for _, br := range branches {
		WriteLine(b.outputWriter, br)
	}
}

func (b *Brancher) branchListRemote() {
	branches, err := b.gitClient.ListRemoteBranches()
	if err != nil {
		WriteError(b.outputWriter, err)
		return
	}
	if len(branches) == 0 {
		WriteLine(b.outputWriter, "No remote branches found.")
		return
	}
	for _, br := range branches {
		WriteLine(b.outputWriter, br)
	}
}

func (b *Brancher) branchSort(args []string) {
	if len(args) > 1 {
		WriteLine(b.outputWriter, "Error: branch sort accepts at most one option (name|date).")
		return
	}

	if len(args) == 1 {
		choice := strings.ToLower(strings.TrimSpace(args[0]))
		if choice == "" {
			WriteLine(b.outputWriter, "Error: sort option cannot be empty.")
			return
		}
		if choice != "name" && choice != "date" {
			WriteErrorf(b.outputWriter, "invalid sort option %q. Use 'name' or 'date'.", args[0])
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
		WriteError(b.outputWriter, err)
		return
	}
	if len(names) == 0 {
		WriteLine(b.outputWriter, "No local branches found.")
		return
	}
	for _, n := range names {
		WriteLine(b.outputWriter, n)
	}
}

func (b *Brancher) branchContains(args []string) {
	if len(args) > 1 {
		WriteLine(b.outputWriter, "Error: branch contains accepts at most one commit or ref.")
		return
	}

	if len(args) == 1 {
		commit := strings.TrimSpace(args[0])
		if commit == "" {
			WriteLine(b.outputWriter, "Error: commit or ref cannot be empty.")
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
		WriteLine(b.outputWriter, "Canceled.")
		return
	}
	b.branchContainsForCommit(commit)
}

func (b *Brancher) branchContainsForCommit(commit string) {
	if !b.gitClient.RevParseVerify(commit) {
		WriteLine(b.outputWriter, "Invalid commit or ref.")
		return
	}
	branches, err := b.gitClient.BranchesContaining(commit)
	if err != nil {
		WriteError(b.outputWriter, err)
		return
	}
	if len(branches) == 0 {
		WriteLine(b.outputWriter, "No branches contain the specified commit.")
		return
	}
	for _, br := range branches {
		WriteLine(b.outputWriter, br)
	}
}
