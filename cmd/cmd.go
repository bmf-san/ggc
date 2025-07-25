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
	Diff(args []string)
	Config(args []string)
	Hook(args []string)
	Version(args []string)
	Status(args []string)
	Tag(args []string)
	Clean(args []string)
	Restore(args []string)
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

// NewCmd creates a new Cmd.
func NewCmd() *Cmd {
	client := git.NewClient()
	return &Cmd{
		gitClient:    client,
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		brancher:     NewBrancher(),
		committer:    NewCommitter(),
		logger:       NewLogger(),
		puller:       NewPuller(),
		pusher:       NewPusher(),
		resetter:     NewResetter(),
		cleaner:      NewCleaner(),
		adder:        NewAdder(),
		remoteer:     NewRemoteer(),
		rebaser:      NewRebaser(),
		stasher:      NewStasher(),
		configureer:  NewConfigureer(),
		hooker:       NewHooker(),
		tagger:       NewTagger(),
		statuseer:    NewStatuseer(),
		versioneer:   NewVersioneer(),
		completer:    NewCompleter(),
		differ:       NewDiffer(),
		restoreer:    NewRestoreer(),
		fetcher:      NewFetcher(),
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
		args := InteractiveUI()
		if args == nil {
			break
		}

		// Check for "quit" command
		if len(args) >= 2 && args[1] == "quit" {
			break
		}

		c.Route(args[1:]) // Skip "ggc" in args
		
		// Smart wait: only pause for commands that need it
		if len(args) >= 2 {
			c.smartWaitForContinue(args[1])
		}
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
	case "version":
		c.Version(args[1:])
	case "clean-interactive":
		c.cleaner.CleanInteractive()
	case "remote":
		c.remoteer.Remote(args[1:])
	case "rebase":
		c.rebaser.Rebase(args[1:])
	case "stash":
		c.stasher.Stash(args[1:])
	case "config":
		c.configureer.Config(args[1:])
	case "hook":
		c.hooker.Hook(args[1:])
	case "tag":
		c.tagger.Tag(args[1:])
	case "status":
		c.statuseer.Status(args[1:])
	case "complete":
		c.completer.Complete(args[1:])
	case "fetch":
		c.fetcher.Fetch(args[1:])
	case "diff":
		c.differ.Diff(args[1:])
	case "restore":
		c.restoreer.Restore(args[1:])
	default:
		c.Help()
	}
}

// smartWaitForContinue only pauses for commands that need user attention
func (c *Cmd) smartWaitForContinue(command string) {
	// Commands that DO need a pause (user should see results)
	waitCommands := map[string]bool{
		"status":  true,  // Status shows important info
		"log":     true,  // Log can be long
		"diff":    true,  // Diff can be long
		"branch":  true,  // Branch operations show results
		"commit":  true,  // Commit shows confirmation
		"reset":   true,  // Reset is destructive
		"clean":   true,  // Clean is destructive
		"stash":   true,  // Stash shows what was stashed
		"tag":     true,  // Tag operations show results
		"remote":  true,  // Remote operations show results
		"rebase":  true,  // Rebase can be complex
		"config":  true,  // Config shows current settings
		"hook":    true,  // Hook operations show results
		"restore": true,  // Restore is important
	}
	
	// If it's a wait command, pause
	if waitCommands[command] {
		fmt.Println("\nPress Enter to continue...")
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
	}
	// If it's not a wait command or unknown command, continue immediately
}