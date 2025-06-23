// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"io"
	"os"

	"github.com/bmf-san/ggc/git"
)

// Executer is an interface for executing commands.
type Executer interface {
	Help()
	Branch(args []string)
	Commit(args []string)
	Log(args []string)
	Pull(args []string)
	Push(args []string)
	Reset(args []string)
	Clean(args []string)
	PullRebasePush()
}

// Cmd represents the command-line interface.
type Cmd struct {
	gitClient        git.Clienter
	outputWriter     io.Writer
	helper           *Helper
	brancher         *Brancher
	committer        *Committer
	logger           *Logger
	puller           *Puller
	pusher           *Pusher
	resetter         *Resetter
	cleaner          *Cleaner
	pullRebasePusher *PullRebasePusher
}

// NewCmd creates a new Cmd.
func NewCmd() *Cmd {
	client := git.NewClient()
	return &Cmd{
		gitClient:        client,
		outputWriter:     os.Stdout,
		helper:           NewHelper(),
		brancher:         NewBrancher(),
		committer:        NewCommitter(),
		logger:           NewLogger(),
		puller:           NewPuller(),
		pusher:           NewPusher(),
		resetter:         NewResetter(),
		cleaner:          NewCleaner(),
		pullRebasePusher: NewPullRebasePusher(),
	}
}

// Help displays help information.
func (c *Cmd) Help() {
	c.helper.ShowHelp()
}

// Branch executes the branch command with the given arguments.
func (c *Cmd) Branch(args []string) {
	c.brancher.Branch(args)
}

// Commit executes the commit command with the given arguments.
func (c *Cmd) Commit(args []string) {
	c.committer.Commit(args)
}

// Log executes the log command with the given arguments.
func (c *Cmd) Log(args []string) {
	c.logger.Log(args)
}

// Pull executes the pull command with the given arguments.
func (c *Cmd) Pull(args []string) {
	c.puller.Pull(args)
}

// Push executes the push command with the given arguments.
func (c *Cmd) Push(args []string) {
	c.pusher.Push(args)
}

// Reset executes the reset command.
func (c *Cmd) Reset(_ []string) {
	c.resetter.Reset()
}

// Clean executes the clean command with the given arguments.
func (c *Cmd) Clean(args []string) {
	c.cleaner.Clean(args)
}

// PullRebasePush executes the pull-rebase-push command.
func (c *Cmd) PullRebasePush() {
	c.pullRebasePusher.PullRebasePush()
}
