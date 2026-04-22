// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	commandregistry "github.com/bmf-san/ggc/v8/cmd/command"
	"github.com/bmf-san/ggc/v8/internal/config"
	"github.com/bmf-san/ggc/v8/internal/git"
	"github.com/bmf-san/ggc/v8/internal/interactive"
)

// Interactive mode command constants.
// These are used to signal special states when returning from interactive mode.
const (
	interactiveQuitCommand     = "quit"
	interactiveWorkflowCommand = "workflow-executed"
)

// Cmd represents the command-line interface.
type Cmd struct {
	registry      *commandregistry.Registry
	configManager *config.Manager
	gitClient     git.StatusInfoReader
	outputWriter  io.Writer
	helper        *Helper
	brancher      *Brancher
	committer     *Committer
	logger        *Logger
	puller        *Puller
	pusher        *Pusher
	resetter      *Resetter
	cleaner       *Cleaner
	adder         *Adder
	remoter       *Remoter
	rebaser       *Rebaser
	stasher       *Stasher
	configurer    *Configurer
	hooker        *Hooker
	tagger        *Tagger
	statuser      *Statuser
	versioner     *Versioner
	differ        *Differ
	restorer      *Restorer
	fetcher       *Fetcher
	cmdRouter     *commandRouter
	debugger      *Debugger
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

// NewCmd creates a new Cmd with the provided git client and config manager.
// The config manager is required for alias resolution and other configuration features.
// It returns an error if the command registry is inconsistent (a developer error
// indicating a command is registered without a handler). Callers that are sure
// the registry is valid may ignore the error.
func NewCmd(client GitDeps, cm *config.Manager) (*Cmd, error) {
	registry := commandregistry.NewRegistry()

	// Populate the alias validator whitelist so that internal/config does not
	// need to import cmd/command (fixes dependency inversion).
	all := registry.All()
	names := make([]string, len(all))
	for i := range all {
		names[i] = all[i].Name
	}
	config.SetValidCommandNames(names)

	tagger := NewTagger(client)
	// Inline default-remote configuration to avoid a post-construction setter.
	if cm != nil {
		if r := strings.TrimSpace(cm.GetConfig().Git.DefaultRemote); r != "" {
			tagger.defaultRemote = r
		}
	}

	cmd := &Cmd{
		registry:      registry,
		configManager: cm,
		gitClient:     client,
		outputWriter:  os.Stdout,
		helper:        NewHelper(registry),
		brancher:      NewBrancher(client),
		committer:     NewCommitter(client),
		logger:        NewLogger(client),
		puller:        NewPuller(client),
		pusher:        NewPusher(client),
		resetter:      NewResetter(client),
		cleaner:       NewCleaner(client),
		adder:         NewAdder(client),
		remoter:       NewRemoter(client),
		rebaser:       NewRebaser(client),
		stasher:       NewStasher(client),
		configurer:    NewConfigurer(client),
		hooker:        NewHooker(client),
		tagger:        tagger,
		statuser:      NewStatuser(client),
		versioner:     NewVersioner(client).withConfigManager(cm),
		differ:        NewDiffer(client),
		restorer:      NewRestorer(client),
		fetcher:       NewFetcher(client),
		debugger:      NewDebugger(),
	}
	router, err := newCommandRouter(cmd)
	if err != nil {
		return nil, err
	}
	cmd.cmdRouter = router
	return cmd, nil
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

// buildInteractiveCommands converts the command registry into the flat list of
// CommandInfo entries consumed by the interactive UI. This keeps the cmd layer
// as the sole owner of registry knowledge so that internal/interactive has no
// dependency on cmd/command.
func buildInteractiveCommands(registry *commandregistry.Registry) []interactive.CommandInfo {
	var list []interactive.CommandInfo
	allCmds := registry.All()
	for i := range allCmds {
		if allCmds[i].Hidden {
			continue
		}
		if len(allCmds[i].Subcommands) == 0 {
			list = append(list, interactive.CommandInfo{Command: allCmds[i].Name, Description: allCmds[i].Summary})
			continue
		}
		for j := range allCmds[i].Subcommands {
			if allCmds[i].Subcommands[j].Hidden {
				continue
			}
			list = append(list, interactive.CommandInfo{Command: allCmds[i].Subcommands[j].Name, Description: allCmds[i].Subcommands[j].Summary})
		}
	}
	return list
}

// Interactive starts the interactive UI mode.
func (c *Cmd) Interactive() {
	// Set up global Ctrl+C handling without introducing a reset window
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer signal.Stop(sigChan)
	go func() {
		<-sigChan
		_, _ = fmt.Fprintln(c.outputWriter, "\nExiting...")
		signal.Stop(sigChan)
		signal.Reset(os.Interrupt)
		os.Exit(0)
	}()

	// Create persistent UI instance to preserve state; pass already-loaded
	// config so NewUI does not perform a second config load (Problem H fix).
	ui := interactive.NewUI(c.gitClient, buildInteractiveCommands(c.registry), c.configManager.GetConfig(), c)

	for {
		args := ui.Run()
		if args == nil {
			break
		}

		// Check for "quit" command
		if len(args) >= 2 && args[1] == interactiveQuitCommand {
			break
		}

		// Check for workflow execution
		if len(args) >= 2 && args[1] == interactiveWorkflowCommand {
			// Workflow was executed, wait for user to continue
			c.waitForContinue()
			ui.ResetToSearchMode()
			continue
		}

		if err := c.Route(args[1:]); err != nil {
			_, _ = fmt.Fprintln(c.outputWriter, "Error:", err)
		}

		// Wait for user to continue
		c.waitForContinue()
		ui.ResetToSearchMode()
	}
}

// Route routes the command to the appropriate handler based on args.
// It returns an error if the command is not recognized.
func (c *Cmd) Route(args []string) error {
	if len(args) == 0 {
		c.Help(nil)
		return nil
	}

	return c.routeCommand(args[0], args[1:])
}

// routeCommand routes to the appropriate command handler
func (c *Cmd) routeCommand(cmd string, args []string) error {

	if c.cmdRouter.route(cmd, args) {
		return nil
	}

	return fmt.Errorf("unknown command: %q", cmd)
}

// waitForContinue waits for user input to continue
func (c *Cmd) waitForContinue() {
	_, _ = fmt.Fprint(c.outputWriter, "\nPress Enter to continue...")
	_, _ = fmt.Scanln()
}
