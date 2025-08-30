package cmd

import (
	"fmt"
	"strings"

	"github.com/bmf-san/ggc/v4/git"
)

// Completer handles dynamic completion for subcommands/args
type Completer struct {
	gitClient git.Clienter
}

// NewCompleter creates a new Completer.
var NewCompleter = func() *Completer {
	return &Completer{
		gitClient: git.NewClient(),
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
			subs := []string{"current", "checkout", "checkout-remote", "create", "delete", "delete-merged"}
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
		// Get list of files managed by git ls-files
		out, err := c.gitClient.ListFiles()
		if err != nil {
			return
		}
		files := strings.Split(strings.TrimSpace(out), "\n")
		for _, f := range files {
			fmt.Println(f)
		}
	default:
		// Other completions can be added in the future
	}
}
