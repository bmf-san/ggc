package cmd

import (
	"strconv"
	"strings"

	"github.com/bmf-san/ggc/v7/pkg/git"
)

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
			WriteLine(b.outputWriter, "Canceled.")
			return
		}
	}
	if err := git.ValidateBranchName(branchName); err != nil {
		WriteErrorf(b.outputWriter, "invalid branch name: %v", err)
		return
	}

	if err := b.gitClient.CheckoutNewBranch(branchName); err != nil {
		WriteErrorf(b.outputWriter, "failed to create and checkout branch: %v", err)
		return
	}
}

func (b *Brancher) branchRename(args []string) {
	if len(args) >= 2 {
		oldName := strings.TrimSpace(args[0])
		newName := strings.TrimSpace(args[1])
		if oldName == "" {
			WriteLine(b.outputWriter, errMsgBranchNameEmpty)
			return
		}
		if newName == "" {
			WriteLine(b.outputWriter, "Error: new branch name cannot be empty.")
			return
		}
		if err := git.ValidateBranchName(newName); err != nil {
			WriteErrorf(b.outputWriter, "invalid branch name: %v", err)
			return
		}
		if err := b.gitClient.RenameBranch(oldName, newName); err != nil {
			WriteError(b.outputWriter, err)
		}
		return
	}

	b.branchRenameInteractive()
}

func (b *Brancher) branchRenameInteractive() {
	branches, err := b.gitClient.ListLocalBranches()
	if err != nil {
		WriteError(b.outputWriter, err)
		return
	}
	if len(branches) == 0 {
		WriteLine(b.outputWriter, "No local branches found.")
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
		WriteLine(b.outputWriter, "Canceled.")
		return
	}
	if err := git.ValidateBranchName(newName); err != nil {
		WriteErrorf(b.outputWriter, "invalid branch name: %v", err)
		return
	}
	if err := b.gitClient.RenameBranch(oldName, newName); err != nil {
		WriteError(b.outputWriter, err)
		return
	}
}

func (b *Brancher) branchMove(args []string) {
	if len(args) >= 2 {
		branch := strings.TrimSpace(args[0])
		commit := strings.TrimSpace(args[1])
		if branch == "" {
			WriteLine(b.outputWriter, errMsgBranchNameEmpty)
			return
		}
		if commit == "" {
			WriteLine(b.outputWriter, "Error: commit or ref cannot be empty.")
			return
		}
		if !b.gitClient.RevParseVerify(commit) {
			WriteLine(b.outputWriter, "Invalid commit or ref.")
			return
		}
		if err := b.gitClient.MoveBranch(branch, commit); err != nil {
			WriteError(b.outputWriter, err)
		}
		return
	}

	b.branchMoveInteractive()
}

func (b *Brancher) branchMoveInteractive() {
	branches, err := b.gitClient.ListLocalBranches()
	if err != nil {
		WriteError(b.outputWriter, err)
		return
	}
	if len(branches) == 0 {
		WriteLine(b.outputWriter, "No local branches found.")
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
		WriteLine(b.outputWriter, "Canceled.")
		return
	}
	if !b.gitClient.RevParseVerify(commit) {
		WriteLine(b.outputWriter, "Invalid commit or ref.")
		return
	}
	if err := b.gitClient.MoveBranch(branch, commit); err != nil {
		WriteError(b.outputWriter, err)
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
			WriteLine(b.outputWriter, errMsgBranchNameEmpty)
			return
		}
		upstream, ok := b.resolveUpstreamArgument(strings.TrimSpace(args[1]))
		if !ok {
			return
		}
		if err := b.gitClient.SetUpstreamBranch(branch, upstream); err != nil {
			WriteError(b.outputWriter, err)
		}
	default:
		WriteLine(b.outputWriter, "Error: branch set upstream expects <branch> <upstream>.")
	}
}

func (b *Brancher) branchSetUpstreamInteractive() {
	branches, err := b.gitClient.ListLocalBranches()
	if err != nil {
		WriteError(b.outputWriter, err)
		return
	}
	if len(branches) == 0 {
		WriteLine(b.outputWriter, "No local branches found.")
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
		WriteError(b.outputWriter, err)
	}
}

func (b *Brancher) resolveUpstreamArgument(input string) (string, bool) {
	if input == "" {
		WriteLine(b.outputWriter, "Error: upstream cannot be empty.")
		return "", false
	}

	if idx, err := strconv.Atoi(input); err == nil {
		remotes, listErr := b.gitClient.ListRemoteBranches()
		if listErr != nil {
			WriteError(b.outputWriter, listErr)
			return "", false
		}
		if idx < 1 || idx > len(remotes) {
			WriteErrorf(b.outputWriter, "invalid remote selection: %d", idx)
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
	remotes, err := b.getValidRemoteBranches()
	if err != nil {
		WriteLinef(b.outputWriter, "Error listing remote branches: %v", err)
		return ""
	}

	if len(remotes) == 0 {
		WriteLine(b.outputWriter, "No remote branches found.")
	}
	b.displayRemoteBranches(remotes)

	upIn, ok := b.readLine("Enter upstream (name or number): ")
	if !ok {
		return ""
	}
	upIn = strings.TrimSpace(upIn)
	if upIn == "" {
		WriteLine(b.outputWriter, "Canceled.")
		return ""
	}
	return b.resolveUpstreamInput(upIn, remotes)
}

// getValidRemoteBranches retrieves and filters remote branches
func (b *Brancher) getValidRemoteBranches() ([]string, error) {
	remotes, err := b.gitClient.ListRemoteBranches()
	if err != nil {
		return nil, err
	}
	// Filter out empty strings from the remote branches list
	validRemotes := make([]string, 0, len(remotes))
	for _, rb := range remotes {
		if strings.TrimSpace(rb) != "" {
			validRemotes = append(validRemotes, rb)
		}
	}
	return validRemotes, nil
}

// displayRemoteBranches shows the list of remote branches
func (b *Brancher) displayRemoteBranches(remotes []string) {
	if len(remotes) > 0 {
		WriteLine(b.outputWriter, "Remote branches:")
		for i, rb := range remotes {
			WriteLinef(b.outputWriter, "[%d] %s", i+1, rb)
		}
	}
}

// resolveUpstreamInput converts user input to upstream branch name
func (b *Brancher) resolveUpstreamInput(input string, remotes []string) string {
	// If numeric and valid index, map to remote branch
	if id, e := strconv.Atoi(input); e == nil && id >= 1 && id <= len(remotes) {
		return remotes[id-1]
	}
	return input
}
