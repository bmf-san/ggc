// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bmf-san/ggc/v4/git"
	"golang.org/x/term"
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
		args := InteractiveUI()
		if args == nil {
			break
		}

		// Check for "quit" command
		if len(args) >= 2 && args[1] == "quit" {
			break
		}

		c.Route(args[1:]) // Skip "ggc" in args

		// Smart continue logic based on command type
		c.smartWaitForContinue(args[1:])
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

// smartWaitForContinue decides whether to wait based on command type
func (c *Cmd) smartWaitForContinue(args []string) {
	if len(args) == 0 {
		return
	}

	command := args[0]

	// Commands with quick execution and minimal output - auto return
	quickCommands := map[string]bool{
		"add":     true, // File staging
		"commit":  true, // Create commit
		"push":    true, // Push to remote
		"pull":    true, // Pull from remote
		"fetch":   true, // Fetch from remote
		"reset":   true, // Reset changes
		"tag":     true, // Create/manage tags
		"restore": true, // Restore files
	}

	// Interactive commands that handle their own user input - no additional wait
	interactiveCommands := map[string]bool{
		"clean":             true, // clean-interactive file selection
		"clean-interactive": true, // explicit clean-interactive
		"rebase":            true, // interactive rebase with commit selection
		"config":            true, // interactive configuration
	}

	// Commands with substantial output that users need to review - wait for key
	reviewCommands := map[string]bool{
		"log":    true, // Git history
		"diff":   true, // File differences
		"status": true, // Working tree status
		"stash":  true, // Stash list/operations
		"remote": true, // Remote repository info
		"branch": true, // Branch listing/info
	}

	// Informational commands - brief pause then auto-return
	infoCommands := map[string]bool{
		"help":     true, // Help information
		"version":  true, // Version information
		"complete": true, // Completion information
	}

	// Special commands - handle appropriately
	specialCommands := map[string]bool{
		"hook": true, // Hook management (varies by operation)
	}

	switch {
	case quickCommands[command]:
		fmt.Print("\n\033[90m✓ Command completed\033[0m")
		time.Sleep(800 * time.Millisecond)
		fmt.Print("\r\033[K")
		return

	case interactiveCommands[command]:
		fmt.Print("\n\033[90m✓ Interactive session completed\033[0m")
		time.Sleep(1200 * time.Millisecond)
		fmt.Print("\r\033[K")
		return

	case reviewCommands[command]:
		c.waitForContinue()
		return

	case infoCommands[command]:
		fmt.Print("\n\033[90m✓ Information displayed\033[0m")
		time.Sleep(1000 * time.Millisecond)
		fmt.Print("\r\033[K")
		return

	case specialCommands[command]:
		// Hook operations vary - some are quick, some need review
		// Default to brief pause for safety
		fmt.Print("\n\033[90m✓ Operation completed\033[0m")
		time.Sleep(1000 * time.Millisecond)
		fmt.Print("\r\033[K")
		return

	default:
		// Unknown commands - safe default with key wait
		fmt.Print("\n\033[90m⚠ Unknown command - press any key to continue\033[0m")
		c.waitForContinue()
	}
}

func (c *Cmd) waitForContinue() {
	fmt.Print("\n\033[90mPress any key to continue...\033[0m")

	// Set terminal to raw mode for single key press
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		// Fallback to Enter key if raw mode fails
		fmt.Println("\nPress Enter to continue...")
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
		return
	}
	defer func() {
		_ = term.Restore(fd, oldState)
	}()

	// Read single key
	buf := make([]byte, 1)
	_, _ = os.Stdin.Read(buf)
	fmt.Println() // Add newline after key press
}
