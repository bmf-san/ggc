// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
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
	remoter      *Remoter
	rebaser      *Rebaser
	stasher      *Stasher
	configurer   *Configurer
	hooker       *Hooker
	tagger       *Tagger
	statuser     *Statuser
	versioner    *Versioner
	completer    *Completer
	differ       *Differ
	restorer     *Restorer
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
		remoter:      NewRemoter(client),
		rebaser:      NewRebaser(client),
		stasher:      NewStasher(client),
		configurer:   NewConfigurer(client),
		hooker:       NewHooker(client),
		tagger:       NewTagger(client),
		statuser:     NewStatuser(client),
		versioner:    NewVersioner(client),
		completer:    NewCompleter(client),
		differ:       NewDiffer(client),
		restorer:     NewRestorer(client),
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
	c.remoter.Remote(args)
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
	c.statuser.Status(args)
}

// Config executes the status command with the given arguments.
func (c *Cmd) Config(args []string) {
	c.configurer.Config(args)
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
	c.restorer.Restore(args)
}

// Version executes the version command with the given arguments.
func (c *Cmd) Version(args []string) {
	c.versioner.Version(args)
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

// Route routes the command to the appropriate handler based on args.
func (c *Cmd) Route(args []string) {
	if len(args) == 0 {
		c.Help()
		return
	}

	// Friendly hint for legacy syntax (flags/hyphens) without executing
	if suggestion, ok := suggestNewSyntax(args); ok {
		_, _ = fmt.Fprintf(c.outputWriter, "Legacy syntax detected.\nUse: ggc %s\n", suggestion)
		return
	}

	c.routeCommand(args[0], args[1:])
}

// routeCommand routes to the appropriate command handler
func (c *Cmd) routeCommand(cmd string, args []string) {
	// Handle core commands first
	if c.handleCoreCommand(cmd, args) {
		return
	}

	// Handle extended commands
	c.routeExtendedCommand(cmd, args)
}

// handleCoreCommand handles core git commands
func (c *Cmd) handleCoreCommand(cmd string, args []string) bool {
	coreCommands := map[string]func([]string){
		"help":    func([]string) { c.Help() },
		"add":     c.adder.Add,
		"branch":  c.Branch,
		"commit":  c.Commit,
		"log":     c.Log,
		"pull":    c.Pull,
		"push":    c.Push,
		"reset":   c.Reset,
		"clean":   c.Clean,
		"version": c.Version,
	}

	if handler, exists := coreCommands[cmd]; exists {
		handler(args)
		return true
	}
	return false
}

// routeExtendedCommand routes to extended command handlers
func (c *Cmd) routeExtendedCommand(cmd string, args []string) {
	extendedCommands := map[string]func([]string){
		"remote":   c.remoter.Remote,
		"rebase":   c.rebaser.Rebase,
		"stash":    c.stasher.Stash,
		"config":   c.configurer.Config,
		"hook":     c.hooker.Hook,
		"tag":      c.tagger.Tag,
		"status":   c.statuser.Status,
		"complete": c.completer.Complete,
		"fetch":    c.fetcher.Fetch,
		"diff":     c.differ.Diff,
		"restore":  c.restorer.Restore,
	}

	if handler, exists := extendedCommands[cmd]; exists {
		handler(args)
		return
	}
	c.Help()
}

// suggestNewSyntax returns a friendly suggestion for unified syntax if a legacy
// pattern is detected. It does not execute the command.
func suggestNewSyntax(args []string) (string, bool) {
	if len(args) == 0 {
		return "", false
	}

	// Top-level hyphenated commands
	if args[0] == "clean-interactive" {
		return "clean interactive", true
	}

	// Check command-specific legacy syntax
	switch args[0] {
	case "rebase":
		return checkRebaseSyntax(args[1:])
	case "add":
		return checkAddSyntax(args[1:])
	case "restore":
		return checkRestoreSyntax(args[1:])
	case "fetch":
		return checkFetchSyntax(args[1:])
	case "commit":
		return checkCommitSyntax(args[1:])
	case "branch":
		return checkBranchSyntax(args[1:])
	}

	return "", false
}

// checkRebaseSyntax checks for legacy rebase syntax
func checkRebaseSyntax(args []string) (string, bool) {
	for _, a := range args {
		if a == "-i" || a == "--interactive" {
			return "rebase interactive", true
		}
	}
	return "", false
}

// checkAddSyntax checks for legacy add syntax
func checkAddSyntax(args []string) (string, bool) {
	for _, a := range args {
		switch a {
		case "-p":
			return "add patch", true
		case "-i", "--interactive":
			return "add interactive", true
		}
	}
	return "", false
}

// checkRestoreSyntax checks for legacy restore syntax
func checkRestoreSyntax(args []string) (string, bool) {
	out := []string{"restore"}
	staged := false
	rest := []string{}
	for _, a := range args {
		if a == "--staged" {
			staged = true
			continue
		}
		rest = append(rest, a)
	}
	if staged {
		out = append(out, "staged")
		out = append(out, rest...)
		return strings.Join(out, " "), true
	}
	return "", false
}

// checkFetchSyntax checks for legacy fetch syntax
func checkFetchSyntax(args []string) (string, bool) {
	for _, a := range args {
		if a == "--prune" {
			return "fetch prune", true
		}
	}
	return "", false
}

// checkCommitSyntax checks for legacy commit syntax
func checkCommitSyntax(args []string) (string, bool) {
	// Check for direct patterns first
	if suggestion := checkDirectCommitPatterns(args); suggestion != "" {
		return suggestion, true
	}

	// Check for amend patterns
	if suggestion := checkAmendPatterns(args); suggestion != "" {
		return suggestion, true
	}

	return "", false
}

// checkDirectCommitPatterns checks for direct commit patterns
func checkDirectCommitPatterns(args []string) string {
	for i, a := range args {
		if a == "--allow-empty" || a == "allow-empty" {
			return "commit allow empty"
		}
		// Check for amend --no-edit pattern
		if a == "amend" && i+1 < len(args) {
			for _, b := range args[i+1:] {
				if b == "--no-edit" {
					return "commit amend no-edit"
				}
			}
		}
	}
	return ""
}

// checkAmendPatterns checks for amend-related patterns
func checkAmendPatterns(args []string) string {
	hasAmend := false
	hasNoEdit := false
	for _, a := range args {
		if a == "amend" || a == "--amend" {
			hasAmend = true
		}
		if a == "--no-edit" {
			hasNoEdit = true
		}
	}
	if hasAmend {
		if hasNoEdit {
			return "commit amend no-edit"
		}
		return "commit amend"
	}
	return ""
}

// checkBranchSyntax checks for legacy branch syntax
func checkBranchSyntax(args []string) (string, bool) {
	if len(args) >= 1 {
		switch args[0] {
		case "checkout-remote":
			return "branch checkout remote", true
		case "delete-merged":
			return "branch delete merged", true
		case "set-upstream":
			return "branch set upstream", true
		case "list-local":
			return "branch list local", true
		case "list-remote":
			return "branch list remote", true
		}
	}
	return "", false
}

// waitForContinue waits for user input to continue
func (c *Cmd) waitForContinue() {
	fmt.Print("\nPress Enter to continue...")
	_, _ = fmt.Scanln()
}
