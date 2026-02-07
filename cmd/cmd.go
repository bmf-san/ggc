// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"

	commandregistry "github.com/bmf-san/ggc/v7/cmd/command"
	"github.com/bmf-san/ggc/v7/internal/interactive"
	"github.com/bmf-san/ggc/v7/pkg/git"
)

// Interactive mode special return values
// These constants are used to signal special states when returning from interactive mode
const (
	// InteractiveQuitCommand signals that the user wants to quit the interactive mode
	InteractiveQuitCommand = interactive.InteractiveQuitCommand
	// InteractiveWorkflowCommand signals that a workflow has been executed and the interactive mode should restart
	InteractiveWorkflowCommand = interactive.InteractiveWorkflowCommand
)

// Executer is an interface for executing commands.
type Executer interface {
	Help(args []string)
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
	DebugKeys(args []string)
	Interactive()
	Route(args []string)
}

// Cmd represents the command-line interface.
type Cmd struct {
	gitClient    git.StatusInfoReader
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
	differ       *Differ
	restorer     *Restorer
	fetcher      *Fetcher
	cmdRouter    *commandRouter
	debugger     *Debugger
}

// GitDeps is a composite for wiring commands that depend on git operations.
// This narrows each command to the smallest required surface while keeping
// construction ergonomic.
type GitDeps interface {
	git.BranchOps
	git.CommitWriter
	git.LogReader
	git.Puller
	git.Pusher
	git.ResetOps
	git.CleanOps
	git.Stager
	git.RemoteManager
	git.RebaseOps
	git.StashOps
	git.ConfigOps
	git.TagOps
	git.StatusInfoReader
	git.DiffReader
	git.RestoreOps
	git.FetchOps
	git.LocalBranchLister
	git.FileLister
}

// NewCmd creates a new Cmd with the provided git client.
func NewCmd(client GitDeps) *Cmd {
	cmd := &Cmd{
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
		differ:       NewDiffer(client),
		restorer:     NewRestorer(client),
		fetcher:      NewFetcher(client),
		debugger:     NewDebugger(),
	}
	cmd.cmdRouter = mustNewCommandRouter(cmd)
	return cmd
}

// SetDefaultRemote sets the default remote used by tag operations.
// Passes through to Tagger to avoid repeated config loading.
func (c *Cmd) SetDefaultRemote(remote string) {
	if c.tagger != nil && strings.TrimSpace(remote) != "" {
		c.tagger.defaultRemote = remote
	}
}

// Help displays help information.
func (c *Cmd) Help(args []string) {
	var name string
	if len(args) > 0 {
		name = strings.TrimSpace(args[0])
	}
	if name == "" {
		c.helper.ShowHelp()
		return
	}
	c.helper.renderCommandFromRegistry(name, nil, "")
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

// DebugKeys executes the debug-keys command with the given arguments.
func (c *Cmd) DebugKeys(args []string) {
	c.debugger.DebugKeys(args)
}

// Interactive starts the interactive UI mode.
func (c *Cmd) Interactive() {
	// Set up global Ctrl+C handling without introducing a reset window
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)
	go func() {
		<-sigChan
		fmt.Println("\nExiting...")
		signal.Stop(sigChan)
		os.Exit(0)
	}()

	// Create persistent UI instance to preserve state
	ui := interactive.NewUI(c.gitClient, c)

	for {
		args := ui.Run()
		if args == nil {
			break
		}

		// Check for "quit" command
		if len(args) >= 2 && args[1] == InteractiveQuitCommand {
			break
		}

		// Check for workflow execution
		if len(args) >= 2 && args[1] == InteractiveWorkflowCommand {
			// Workflow was executed, wait for user to continue
			c.waitForContinue()
			ui.ResetToSearchMode()
			continue
		}

		c.Route(args[1:]) // Skip "ggc" in args

		// Wait for user to continue
		c.waitForContinue()
		ui.ResetToSearchMode()
	}
}

// Route routes the command to the appropriate handler based on args.
func (c *Cmd) Route(args []string) {
	if len(args) == 0 {
		c.Help(nil)
		return
	}

	// Treat legacy-like syntax as a hard error (no heuristics/suggestions)
	if isLegacyLike(args) {
		_, _ = fmt.Fprintln(c.outputWriter, "Error: legacy-like syntax is not supported.")
		_, _ = fmt.Fprintln(c.outputWriter, "Use unified subcommands. See: ggc help <command>")
		return
	}

	c.routeCommand(args[0], args[1:])
}

// routeCommand routes to the appropriate command handler
func (c *Cmd) routeCommand(cmd string, args []string) {

	if c.cmdRouter.route(cmd, args) {
		return
	}

	c.Help(nil)
}

type commandRouter struct {
	handlers map[string]func([]string)
}

func mustNewCommandRouter(cmd *Cmd) *commandRouter {
	router, err := newCommandRouter(cmd)
	if err != nil {
		panic(err)
	}
	return router
}

func newCommandRouter(cmd *Cmd) (*commandRouter, error) {
	if err := commandregistry.ValidateAll(); err != nil {
		return nil, fmt.Errorf("command registry validation failed: %w", err)
	}

	handlers := map[string]func([]string){
		"help":       func(args []string) { cmd.Help(args) },
		"add":        func(args []string) { cmd.Add(args) },
		"branch":     func(args []string) { cmd.Branch(args) },
		"commit":     func(args []string) { cmd.Commit(args) },
		"log":        func(args []string) { cmd.Log(args) },
		"pull":       func(args []string) { cmd.Pull(args) },
		"push":       func(args []string) { cmd.Push(args) },
		"reset":      func(args []string) { cmd.Reset(args) },
		"clean":      func(args []string) { cmd.Clean(args) },
		"version":    func(args []string) { cmd.Version(args) },
		"remote":     func(args []string) { cmd.Remote(args) },
		"rebase":     func(args []string) { cmd.Rebase(args) },
		"stash":      func(args []string) { cmd.Stash(args) },
		"config":     func(args []string) { cmd.Config(args) },
		"hook":       func(args []string) { cmd.Hook(args) },
		"tag":        func(args []string) { cmd.Tag(args) },
		"status":     func(args []string) { cmd.Status(args) },
		"fetch":      func(args []string) { cmd.Fetch(args) },
		"diff":       func(args []string) { cmd.Diff(args) },
		"restore":    func(args []string) { cmd.Restore(args) },
		"debug-keys": func(args []string) { cmd.DebugKeys(args) },
		InteractiveQuitCommand: func([]string) {
			_, _ = fmt.Fprintln(cmd.outputWriter, "The 'quit' command is only available in interactive mode.")
		},
	}

	available := make(map[string]struct{}, len(handlers))
	for key := range handlers {
		available[key] = struct{}{}
	}

	missing := missingHandlers(available)
	if len(missing) > 0 {
		sort.Strings(missing)
		return nil, fmt.Errorf("no handler registered for commands: %s", strings.Join(missing, ", "))
	}

	return &commandRouter{handlers: handlers}, nil
}

func (r *commandRouter) route(cmd string, args []string) bool {
	info, ok := commandregistry.Find(cmd)
	if !ok {
		return false
	}
	identifier := strings.TrimSpace(info.HandlerID)
	if identifier == "" {
		return false
	}
	handler, ok := r.handlers[identifier]
	if !ok {
		return false
	}
	handler(args)
	return true
}

func missingHandlers(available map[string]struct{}) []string {
	var missing []string
	allCommands := commandregistry.All()
	for i := range allCommands {
		info := &allCommands[i]
		id := strings.TrimSpace(info.HandlerID)
		if id == "" || info.Hidden {
			continue
		}
		if _, ok := available[id]; !ok {
			missing = append(missing, id)
		}
	}
	return missing
}

// isLegacyLike reports whether args resemble the pre-v6, flag-driven CLI surface
// instead of the unified subcommand style introduced in v6.
// Definition of "legacy-like":
//   - Hyphenated top-level commands (e.g., "clean-interactive") that used to exist
//     as separate entry points.
//   - Flag-style modifiers (e.g., "-i", "--prune") passed to top-level commands to
//     toggle modes or behaviors.
//
// Why this matters: post-v6 the CLI avoids guessing user intent to keep routing
// predictable, safe, and testable. When legacy-like usage is detected, we fail fast
// with a clear error to nudge users toward the unified subcommands (see: "ggc help <command>").
// Detection is intentionally conservative: we flag a hyphen in the top-level command
// or any leading "-" arguments up to a literal "--" separator; everything after "--"
// is treated as data.
func isLegacyLike(args []string) bool {
	if len(args) == 0 {
		return false
	}
	// A hyphenated top-level command is only legacy-like when it is not
	// a registered first-class command (e.g. "debug-keys").
	if strings.Contains(args[0], "-") {
		if _, ok := commandregistry.Find(args[0]); !ok {
			return true
		}
	}
	// A standalone flag token as the top-level command is always legacy-like.
	if strings.HasPrefix(args[0], "-") {
		return true
	}
	// Any flag-style argument (starts with '-' or '--') is considered legacy-like
	for _, a := range args[1:] {
		if a == "--" { // treat everything after "--" as data
			break
		}
		if strings.HasPrefix(a, "-") {
			return true
		}
	}
	return false
}

// waitForContinue waits for user input to continue
func (c *Cmd) waitForContinue() {
	fmt.Print("\nPress Enter to continue...")
	_, _ = fmt.Scanln()
}
