package cmd

import (
	"fmt"
	"strings"

	"github.com/bmf-san/ggc/v5/git"
)

// Completer handles dynamic completion for subcommands/args
type Completer struct {
	// Needs only local branches and file listing for completion
	gitClient interface {
		git.LocalBranchLister
		git.FileLister
	}
}

// NewCompleter creates a new Completer.
var NewCompleter = func(client interface {
	git.LocalBranchLister
	git.FileLister
}) *Completer {
	return &Completer{
		gitClient: client,
	}
}

// Complete provides completion for various subcommands.
func (c *Completer) Complete(args []string) {
	if len(args) < 1 {
		return
	}
	switch args[0] {
	case "branch":
		if len(args) == 1 {
			// Suggest subcommand candidates
			subs := []string{
				"current",
				"checkout",
				"create",
				"delete",
				// follow-ups: 'delete merged'
				// Enhanced branch management
				"rename",
				"move",
				// follow-ups: 'set upstream'
				"info",
				"list",
				// follow-ups: 'list verbose'
				"sort",
				"contains",
			}
			for _, s := range subs {
				fmt.Println(s)
			}
			return
		}
		// For the second argument and beyond, suggest local branch names
		branches, err := c.gitClient.ListLocalBranches()
		if err != nil {
			return
		}
		for _, b := range branches {
			fmt.Println(b)
		}
	case "files":
		c.completeFiles()
	default:
		// Other completions can be added in the future
	}
}

func (c *Completer) completeFiles() {
	out, err := c.gitClient.ListFiles()
	if err != nil {
		return
	}
	files := strings.Split(strings.TrimSpace(out), "\n")
	for _, f := range files {
		fmt.Println(f)
	}
}
