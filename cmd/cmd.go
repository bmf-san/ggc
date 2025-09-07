// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/bmf-san/ggc/v5/git"
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
	Diff(args []string)
	Config(args []string)
	Hook(args []string)
	Rebase(args []string)
	Remote(args []string)
	Version(args []string)
	Stash(args []string)
	Fetch(args []string)
	Status(args []string)
	Tag(args []string)
	Clean(args []string)
	Restore(args []string)
	Add(args []string)
	Interactive()
}

// Cmd represents the command-line interface.
type Cmd struct {
	gitClient    git.Clienter
	outputWriter io.Writer
	helper       *Helper
	brancher     *Brancher
	committer    *Committer
	logger       *Logger
	puller       *Puller
	pusher       *Pusher
	resetter     *Resetter
	cleaner      *Cleaner
	adder        *Adder
	remoteer     *Remoteer
	rebaser      *Rebaser
	stasher      *Stasher
	configureer  *Configureer
	hooker       *Hooker
	tagger       *Tagger
	statuseer    *Statuseer
	versioneer   *Versioneer
	completer    *Completer
	differ       *Differ
	restoreer    *Restoreer
	fetcher      *Fetcher
}

// NewCmd creates a new Cmd with the provided git client.
func NewCmd(client git.Clienter) *Cmd {
	return &Cmd{
		gitClient:    client,
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		brancher:     NewBrancher(client),
		committer:    NewCommitter(client),
		logger:       NewLogger(client),
		puller:       NewPuller(client),
		pusher:       NewPusher(client),
		resetter:     NewResetter(client),
		cleaner:      NewCleaner(client),
		adder:        NewAdder(client),
		remoteer:     NewRemoteer(client),
		rebaser:      NewRebaser(client),
		stasher:      NewStasher(client),
		configureer:  NewConfigureer(client),
		hooker:       NewHooker(client),
		tagger:       NewTagger(client),
		statuseer:    NewStatuseer(client),
		versioneer:   NewVersioneer(client),
		completer:    NewCompleter(client),
		differ:       NewDiffer(client),
		restoreer:    NewRestoreer(client),
		fetcher:      NewFetcher(client),
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

// Remote executes the remote command with the given arguments.
func (c *Cmd) Remote(args []string) {
	c.remoteer.Remote(args)
}

// Rebase executes the rebase command with the given arguments.
func (c *Cmd) Rebase(args []string) {
	c.rebaser.Rebase(args)
}

// Stash executes the stash command with the given arguments.
func (c *Cmd) Stash(args []string) {
	c.stasher.Stash(args)
}

// Fetch executes the fetch command with the given arguments.
func (c *Cmd) Fetch(args []string) {
	c.fetcher.Fetch(args)
}

// Commit executes the commit command with the given arguments.
func (c *Cmd) Commit(args []string) {
	c.committer.Commit(args)
}

// Log executes the log command with the given arguments.
func (c *Cmd) Log(args []string) {
	c.logger.Log(args)
}

// Add executes the add command with the given arguments.
func (c *Cmd) Add(args []string) {
	c.adder.Add(args)
}

// Status executes the status command with the given arguments.
func (c *Cmd) Status(args []string) {
	c.statuseer.Status(args)
}

// Config executes the status command with the given arguments.
func (c *Cmd) Config(args []string) {
	c.configureer.Config(args)
}

// Hook executes the hook command with the given arguments.
func (c *Cmd) Hook(args []string) {
	c.hooker.Hook(args)
}

// Tag executes the tag command with the given arguments.
func (c *Cmd) Tag(args []string) {
	c.tagger.Tag(args)
}

// Diff executes the diff command with the given arguments.
func (c *Cmd) Diff(args []string) {
	c.differ.Diff(args)
}

// Restore executes the restore command with the given arguments.
func (c *Cmd) Restore(args []string) {
	c.restoreer.Restore(args)
}

// Version executes the version command with the given arguments.
func (c *Cmd) Version(args []string) {
	c.versioneer.Version(args)
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

// Interactive starts the interactive UI mode.
func (c *Cmd) Interactive() {
	// Reset existing signal handlers
	signal.Reset(os.Interrupt, syscall.SIGTERM)

	// Set up global Ctrl+C handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nExiting...")
		os.Exit(0)
	}()

	for {
		args := InteractiveUI(c.gitClient)
		if args == nil {
			break
		}

		// Check for "quit" command
		if len(args) >= 2 && args[1] == "quit" {
			break
		}

		c.Route(args[1:]) // Skip "ggc" in args

		// Wait for user to continue
		c.waitForContinue()
	}
}

// commandHandler represents a function that handles a command
type commandHandler func([]string)

// getCommandHandlers returns a map of command names to their handlers
func (c *Cmd) getCommandHandlers() map[string]commandHandler {
	return map[string]commandHandler{
		"help":              func(_ []string) { c.Help() },
		"add":               c.adder.Add,
		"branch":            c.Branch,
		"commit":            c.Commit,
		"log":               c.Log,
		"pull":              c.Pull,
		"push":              c.Push,
		"reset":             c.Reset,
		"clean":             c.Clean,
		"version":           c.Version,
		"clean-interactive": func(_ []string) { c.cleaner.CleanInteractive() },
		"remote":            c.remoteer.Remote,
		"rebase":            c.rebaser.Rebase,
		"stash":             c.stasher.Stash,
		"config":            c.configureer.Config,
		"hook":              c.hooker.Hook,
		"tag":               c.tagger.Tag,
		"status":            c.statuseer.Status,
		"complete":          c.completer.Complete,
		"fetch":             c.fetcher.Fetch,
		"diff":              c.differ.Diff,
		"restore":           c.restoreer.Restore,
	}
}

// Route routes the command to the appropriate handler based on args.
func (c *Cmd) Route(args []string) {
	if len(args) == 0 {
		c.Help()
		return
	}

	handlers := c.getCommandHandlers()
	if handler, exists := handlers[args[0]]; exists {
		handler(args[1:])
	} else {
		c.Help()
	}
}

// waitForContinue waits for user input to continue
func (c *Cmd) waitForContinue() {
	fmt.Print("\nPress Enter to continue...")
	_, _ = fmt.Scanln()
}
