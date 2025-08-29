// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bmf-san/ggc/v3/git"
)

// Brancher provides functionality for the branch command.
type Brancher struct {
	gitClient    git.Clienter
	execCommand  func(name string, arg ...string) *exec.Cmd
	inputReader  *bufio.Reader
	outputWriter io.Writer
	helper       *Helper
}

// NewBrancher creates a new Brancher.
func NewBrancher() *Brancher {
	return &Brancher{
		gitClient:    git.NewClient(),
		execCommand:  exec.Command,
		inputReader:  bufio.NewReader(os.Stdin),
		outputWriter: os.Stdout,
		helper:       NewHelper(),
	}
}

// Branch executes the branch command with the given arguments.
func (b *Brancher) Branch(args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "current":
			branch, err := b.gitClient.GetCurrentBranch()
			if err != nil {
				_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
				return
			}
			_, _ = fmt.Fprintln(b.outputWriter, branch)
			return
		case "checkout":
			b.branchCheckout()
			return
		case "checkout-remote":
			b.branchCheckoutRemote()
			return
		case "create":
			b.branchCreate()
			return
		case "delete":
			b.branchDelete()
			return
		case "delete-merged":
			b.branchDeleteMerged()
			return
		}
	}
	b.helper.ShowBranchHelp()
}

func (b *Brancher) branchCheckout() {
	branches, err := b.gitClient.ListLocalBranches()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	if len(branches) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No local branches found.")
		return
	}
	_, _ = fmt.Fprintln(b.outputWriter, "Local branches:")
	for i, br := range branches {
		_, _ = fmt.Fprintf(b.outputWriter, "[%d] %s\n", i+1, br)
	}
	_, _ = fmt.Fprint(b.outputWriter, "Enter the number to checkout: ")
	input, _ := b.inputReader.ReadString('\n')
	input = strings.TrimSpace(input)
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(branches) {
		_, _ = fmt.Fprintln(b.outputWriter, "Invalid number.")
		return
	}
	branch := branches[idx-1]
	cmd := b.execCommand("git", "checkout", branch)
	cmd.Stdout = b.outputWriter
	cmd.Stderr = b.outputWriter
	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
	}
}

func (b *Brancher) branchCheckoutRemote() {
	branches, err := b.gitClient.ListRemoteBranches()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	if len(branches) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No remote branches found.")
		return
	}
	_, _ = fmt.Fprintln(b.outputWriter, "Remote branches:")
	for i, br := range branches {
		_, _ = fmt.Fprintf(b.outputWriter, "[%d] %s\n", i+1, br)
	}
	_, _ = fmt.Fprint(b.outputWriter, "Enter the number to checkout: ")
	input, _ := b.inputReader.ReadString('\n')
	input = strings.TrimSpace(input)
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(branches) {
		_, _ = fmt.Fprintln(b.outputWriter, "Invalid number.")
		return
	}
	remoteBranch := branches[idx-1]
	// origin/feature/foo → feature/foo
	parts := strings.SplitN(remoteBranch, "/", 2)
	if len(parts) != 2 {
		_, _ = fmt.Fprintln(b.outputWriter, "Invalid remote branch name.")
		return
	}
	localBranch := parts[1]
	cmd := b.execCommand("git", "checkout", "-b", localBranch, "--track", remoteBranch)
	cmd.Stdout = b.outputWriter
	cmd.Stderr = b.outputWriter
	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
	}
}

func (b *Brancher) branchCreate() {
	_, _ = fmt.Fprint(b.outputWriter, "Enter new branch name: ")
	input, _ := b.inputReader.ReadString('\n')
	branchName := strings.TrimSpace(input)
	if branchName == "" {
		_, _ = fmt.Fprintln(b.outputWriter, "Cancelled.")
		return
	}

	cmd := b.execCommand("git", "checkout", "-b", branchName)
	cmd.Stdout = b.outputWriter
	cmd.Stderr = b.outputWriter
	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: failed to create and checkout branch: %v\n", err)
		return
	}
}

func (b *Brancher) branchDelete() {
	branches, err := b.gitClient.ListLocalBranches()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: %v\n", err)
		return
	}
	if len(branches) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No local branches found.")
		return
	}
	for {
		_, _ = fmt.Fprintln(b.outputWriter, "\033[1;36mSelect local branches to delete by number (space separated, all: select all, none: deselect all, e.g. 1 3 5):\033[0m")
		for i, br := range branches {
			_, _ = fmt.Fprintf(b.outputWriter, "  [\033[1;33m%d\033[0m] %s\n", i+1, br)
		}
		_, _ = fmt.Fprint(b.outputWriter, "> ")
		input, _ := b.inputReader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			_, _ = fmt.Fprintln(b.outputWriter, "Cancelled.")
			return
		}
		if input == "all" {
			for _, br := range branches {
				cmd := b.execCommand("git", "branch", "-d", br)
				cmd.Stdout = b.outputWriter
				cmd.Stderr = b.outputWriter
				if err := cmd.Run(); err != nil {
					_, _ = fmt.Fprintf(b.outputWriter, "Error: failed to delete %s: %v\n", br, err)
				}
			}
			_, _ = fmt.Fprintln(b.outputWriter, "All branches deleted.")
			break
		}
		if input == "none" {
			continue
		}
		indices := strings.Fields(input)
		tmp := []string{}
		valid := true
		for _, idx := range indices {
			n, err := strconv.Atoi(idx)
			if err != nil || n < 1 || n > len(branches) {
				_, _ = fmt.Fprintf(b.outputWriter, "\033[1;31mInvalid number: %s\033[0m\n", idx)
				valid = false
				break
			}
			tmp = append(tmp, branches[n-1])
		}
		if !valid {
			continue
		}
		for _, br := range tmp {
			cmd := b.execCommand("git", "branch", "-d", br)
			cmd.Stdout = b.outputWriter
			cmd.Stderr = b.outputWriter
			if err := cmd.Run(); err != nil {
				_, _ = fmt.Fprintf(b.outputWriter, "Error: failed to delete %s: %v\n", br, err)
			}
		}
		_, _ = fmt.Fprintln(b.outputWriter, "Selected branches deleted.")
		break
	}
}

func (b *Brancher) branchDeleteMerged() {
	current, err := b.gitClient.GetCurrentBranch()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: failed to get current branch: %v\n", err)
		return
	}
	cmd := b.execCommand("git", "branch", "--merged")
	out, err := cmd.Output()
	if err != nil {
		_, _ = fmt.Fprintf(b.outputWriter, "Error: failed to get merged branches: %v\n", err)
		return
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	branches := []string{}
	for _, l := range lines {
		br := strings.TrimSpace(strings.TrimPrefix(l, "* "))
		if br != "" && br != current {
			branches = append(branches, br)
		}
	}
	if len(branches) == 0 {
		_, _ = fmt.Fprintln(b.outputWriter, "No merged local branches.")
		return
	}
	for {
		_, _ = fmt.Fprintln(b.outputWriter, "\033[1;36mSelect merged local branches to delete by number (space separated, all: select all, none: deselect all, e.g. 1 3 5):\033[0m")
		for i, br := range branches {
			_, _ = fmt.Fprintf(b.outputWriter, "  [\033[1;33m%d\033[0m] %s\n", i+1, br)
		}
		_, _ = fmt.Fprint(b.outputWriter, "> ")
		input, _ := b.inputReader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			_, _ = fmt.Fprintln(b.outputWriter, "Cancelled.")
			return
		}
		if input == "all" {
			for _, br := range branches {
				cmd := b.execCommand("git", "branch", "-d", br)
				cmd.Stdout = b.outputWriter
				cmd.Stderr = b.outputWriter
				if err := cmd.Run(); err != nil {
					_, _ = fmt.Fprintf(b.outputWriter, "Error: failed to delete %s: %v\n", br, err)
				}
			}
			_, _ = fmt.Fprintln(b.outputWriter, "All merged branches deleted.")
			break
		}
		if input == "none" {
			continue
		}
		indices := strings.Fields(input)
		tmp := []string{}
		valid := true
		for _, idx := range indices {
			n, err := strconv.Atoi(idx)
			if err != nil || n < 1 || n > len(branches) {
				_, _ = fmt.Fprintf(b.outputWriter, "\033[1;31mInvalid number: %s\033[0m\n", idx)
				valid = false
				break
			}
			tmp = append(tmp, branches[n-1])
		}
		if !valid {
			continue
		}
		for _, br := range tmp {
			cmd := b.execCommand("git", "branch", "-d", br)
			cmd.Stdout = b.outputWriter
			cmd.Stderr = b.outputWriter
			if err := cmd.Run(); err != nil {
				_, _ = fmt.Fprintf(b.outputWriter, "Error: failed to delete %s: %v\n", br, err)
			}
		}
		_, _ = fmt.Fprintln(b.outputWriter, "Selected merged branches deleted.")
		break
	}
}
