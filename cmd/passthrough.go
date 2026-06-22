package cmd

import (
	"io"
	"os"

	"github.com/bmf-san/ggc/v8/internal/git"
)

// passthroughCommand is a thin wrapper that forwards arguments to the
// underlying `git <name>` invocation. It is used by ggc commands that do
// not need a bespoke argument parser or interactive layer; users get the
// full surface of the underlying git porcelain while still being able to
// discover the command through `ggc help`, completions, and the registry.
type passthroughCommand struct {
	name         string
	gitClient    git.PassthroughOps
	outputWriter io.Writer
	helper       *Helper
}

// newPassthroughCommand creates a pass-through wrapper for the given git
// subcommand name.
func newPassthroughCommand(name string, client git.PassthroughOps) *passthroughCommand {
	return &passthroughCommand{
		name:         name,
		gitClient:    client,
		outputWriter: os.Stdout,
		helper:       NewHelper(),
	}
}

// Run forwards args to `git <name>`. If the first argument is the literal
// string "help", it prints help rendered from the registry instead.
func (p *passthroughCommand) Run(args []string) {
	if len(args) > 0 && args[0] == "help" {
		p.helper.ShowPassthroughHelp(p.name)
		return
	}
	if err := p.gitClient.RunGit(p.name, args); err != nil {
		WriteError(p.outputWriter, err)
	}
}

// passthroughCommandNames lists every ggc command that is implemented as a
// thin pass-through to `git <name>`. Adding a name here, along with a
// matching registry entry, automatically wires it into the router.
var passthroughCommandNames = []string{
	// Tier 1
	"switch",
	"checkout",
	"merge",
	"cherry-pick",
	"revert",
	"blame",
	// Tier 2
	"worktree",
	"reflog",
	"format-patch",
	"am",
	"sparse-checkout",
	"mv",
	"rm",
	"submodule",
	// Tier 3
	"describe",
	"range-diff",
	"grep",
	"notes",
	"archive",
	"shortlog",
	"maintenance",
	"gc",
	"fsck",
	"prune",
}

// buildPassthroughs constructs a map of all pass-through commands keyed by
// their canonical name. The router dispatches to the matching entry by name.
func buildPassthroughs(client git.PassthroughOps) map[string]*passthroughCommand {
	m := make(map[string]*passthroughCommand, len(passthroughCommandNames))
	for _, name := range passthroughCommandNames {
		m[name] = newPassthroughCommand(name, client)
	}
	return m
}
