package cmd

import (
	"fmt"
	"strings"
)

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
