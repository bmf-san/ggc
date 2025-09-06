package cmd

import (
	"fmt"
	"strings"

	"github.com/bmf-san/ggc/v5/git"
)

// Completer handles dynamic completion for subcommands/args
type Completer struct {
	gitClient git.Clienter
}

// NewCompleter creates a new Completer.
var NewCompleter = func(client git.Clienter) *Completer {
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
		c.completeBranch(args)
	case "files":
		c.completeFiles()
	default:
		// Other completions can be added in the future
	}
}

func (c *Completer) completeBranch(args []string) {
	if len(args) == 1 {
		subs := []string{
			"current",
			"checkout",
			"checkout-remote",
			"create",
			"delete",
			"delete-merged",
			"rename",
			"move",
			"set-upstream",
			"info",
			"list",
			"list --verbose",
			"sort",
			"contains",
		}
		for _, s := range subs {
			fmt.Println(s)
		}
		return
	}
	branches, err := c.gitClient.ListLocalBranches()
	if err != nil {
		return
	}
	for _, b := range branches {
		fmt.Println(b)
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
