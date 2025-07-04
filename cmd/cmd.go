// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

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
	Interactive()
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
	adder            *Adder
	remoteer         *Remoteer
	rebaser          *Rebaser
	stasher          *Stasher
	commitPusher     *CommitPusher
	addCommitPusher  *AddCommitPusher
	completer        *Completer
	fetcher          *Fetcher
	stashPullPopper  *StashPullPopper
	resetCleaner     *ResetCleaner
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
		adder:            NewAdder(),
		remoteer:         NewRemoteer(),
		rebaser:          NewRebaser(),
		stasher:          NewStasher(),
		commitPusher:     NewCommitPusher(),
		addCommitPusher:  NewAddCommitPusher(),
		completer:        NewCompleter(),
		fetcher:          NewFetcher(),
		stashPullPopper:  NewStashPullPopper(),
		resetCleaner:     NewResetCleaner(),
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
func (c *Cmd) Reset(args []string) {
	c.resetter.Reset(args)
}

// Clean executes the clean command with the given arguments.
func (c *Cmd) Clean(args []string) {
	c.cleaner.Clean(args)
}

// PullRebasePush executes the pull-rebase-push command.
func (c *Cmd) PullRebasePush() {
	c.pullRebasePusher.PullRebasePush()
}

// Interactive starts the interactive UI mode.
func (c *Cmd) Interactive() {
	// 既存のシグナルハンドラーをリセット
	signal.Reset(os.Interrupt, syscall.SIGTERM)

	// 全体でのCtrl+C処理を設定
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nExiting...")
		os.Exit(0)
	}()

	for {
		args := InteractiveUI()
		if args == nil {
			break
		}

		// "quit" コマンドをチェック
		if len(args) >= 2 && args[1] == "quit" {
			break
		}

		c.Route(args[1:]) // Skip "ggc" in args

		// コマンド実行後、結果を確認するための待機
		c.waitForContinue()
	}
}

// Route routes the command to the appropriate handler based on args.
func (c *Cmd) Route(args []string) {
	if len(args) == 0 {
		c.Help()
		return
	}

	switch args[0] {
	case "help":
		c.Help()
	case "add":
		c.adder.Add(args[1:])
	case "branch":
		c.Branch(args[1:])
	case "commit":
		c.Commit(args[1:])
	case "log":
		c.Log(args[1:])
	case "pull":
		c.Pull(args[1:])
	case "push":
		c.Push(args[1:])
	case "reset":
		c.Reset(args[1:])
	case "clean":
		c.Clean(args[1:])
	case "clean-interactive":
		c.cleaner.CleanInteractive()
	case "pull-rebase-push":
		c.PullRebasePush()
	case "remote":
		c.remoteer.Remote(args[1:])
	case "rebase":
		c.rebaser.Rebase(args[1:])
	case "stash":
		c.stasher.Stash(args[1:])
	case "commit-push-interactive":
		c.commitPusher.CommitPushInteractive()
	case "add-commit-push":
		c.addCommitPusher.AddCommitPush()
	case "complete":
		c.completer.Complete(args[1:])
	case "fetch":
		c.fetcher.Fetch(args[1:])
	case "stash-pull-pop":
		c.stashPullPopper.StashPullPop()
	case "reset-clean":
		c.resetCleaner.ResetClean()
	default:
		c.Help()
	}
}

func (c *Cmd) waitForContinue() {
	fmt.Println("\nPress Enter to continue...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
