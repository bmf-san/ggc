package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/bmf-san/ggc/git"
)

// Complete handles dynamic completion for subcommands/args
type Completer struct {
	listLocalBranches func() ([]string, error)
	execCommand       func(name string, arg ...string) *exec.Cmd
}

func NewCompleter() *Completer {
	return &Completer{
		listLocalBranches: git.ListLocalBranches,
		execCommand:       exec.Command,
	}
}

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
		branches, err := c.listLocalBranches()
		if err != nil {
			return
		}
		for _, b := range branches {
			fmt.Println(b)
		}
	case "files":
		// Get list of files managed by git ls-files
		cmd := c.execCommand("git", "ls-files")
		out, err := cmd.Output()
		if err != nil {
			return
		}
		files := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, f := range files {
			fmt.Println(f)
		}
	default:
		// Other completions can be added in the future
	}
}
