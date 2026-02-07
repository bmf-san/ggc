package cmd

import (
	"errors"
	"strings"

	"github.com/bmf-san/ggc/v7/internal/prompt"
	"github.com/bmf-san/ggc/v7/pkg/git"
)

func (b *Brancher) branchCheckout() {
	branches, err := b.gitClient.ListLocalBranches()
	if err != nil {
		WriteError(b.outputWriter, err)
		return
	}
	if len(branches) == 0 {
		WriteLine(b.outputWriter, "No local branches found.")
		return
	}
	idx, ok := b.promptSelectIndex("Local branches:", branches, "Enter the number to checkout: ")
	if !ok {
		return
	}
	branch := branches[idx]
	if err := b.gitClient.CheckoutBranch(branch); err != nil {
		WriteError(b.outputWriter, err)
	}
}

func (b *Brancher) branchCheckoutRemote() {
	branches, err := b.gitClient.ListRemoteBranches()
	if err != nil {
		WriteError(b.outputWriter, err)
		return
	}
	if len(branches) == 0 {
		WriteLine(b.outputWriter, "No remote branches found.")
		return
	}
	idx, ok := b.promptSelectIndex("Remote branches:", branches, "Enter the number to checkout: ")
	if !ok {
		return
	}
	remoteBranch := branches[idx]
	localBranch, valid := deriveLocalFromRemote(remoteBranch)
	if !valid || git.ValidateBranchName(localBranch) != nil {
		WriteLine(b.outputWriter, "Invalid remote branch name.")
		return
	}
	if err := b.gitClient.CheckoutNewBranchFromRemote(localBranch, remoteBranch); err != nil {
		WriteError(b.outputWriter, err)
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
			WriteLine(b.outputWriter, "Invalid number.")
		} else {
			WriteError(b.outputWriter, err)
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
		WriteError(b.outputWriter, err)
		return "", false
	}
	return line, true
}
