// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/bmf-san/ggc/v7/internal/prompt"
	"github.com/bmf-san/ggc/v7/pkg/git"
)

const errMsgBranchNameEmpty = "Error: branch name cannot be empty."

// Brancher provides functionality for the branch command.
type Brancher struct {
	gitClient    git.BranchOps
	prompter     prompt.Prompter
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
		WriteError(b.outputWriter, err)
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
