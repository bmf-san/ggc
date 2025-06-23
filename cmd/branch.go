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

	"github.com/bmf-san/ggc/git"
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
			branchCheckoutRemote()
			return
		case "delete":
			branchDelete()
			return
		case "delete-merged":
			branchDeleteMerged()
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

func branchCheckoutRemote() {
	gitClient := git.NewClient()
	branches, err := gitClient.ListRemoteBranches()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(branches) == 0 {
		fmt.Println("No remote branches found.")
		return
	}
	fmt.Println("Remote branches:")
	for i, b := range branches {
		fmt.Printf("[%d] %s\n", i+1, b)
	}
	fmt.Print("Enter the number to checkout: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(branches) {
		fmt.Println("Invalid number.")
		return
	}
	remoteBranch := branches[idx-1]
	// origin/feature/foo → feature/foo
	parts := strings.SplitN(remoteBranch, "/", 2)
	if len(parts) != 2 {
		fmt.Println("Invalid remote branch name.")
		return
	}
	localBranch := parts[1]
	cmd := exec.Command("git", "checkout", "-b", localBranch, "--track", remoteBranch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}

func branchDelete() {
	gitClient := git.NewClient()
	branches, err := gitClient.ListLocalBranches()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(branches) == 0 {
		fmt.Println("No local branches found.")
		return
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\033[1;36mSelect local branches to delete by number (space separated, all: select all, none: deselect all, e.g. 1 3 5):\033[0m")
		for i, b := range branches {
			fmt.Printf("  [\033[1;33m%d\033[0m] %s\n", i+1, b)
		}
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Println("Cancelled.")
			return
		}
		if input == "all" {
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
				fmt.Printf("\033[1;31mInvalid number: %s\033[0m\n", idx)
				valid = false
				break
			}
			tmp = append(tmp, branches[n-1])
		}
		if !valid {
			continue
		}
		for _, b := range tmp {
			cmd := exec.Command("git", "branch", "-d", b)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Printf("Error: failed to delete %s: %v\n", b, err)
			}
		}
		fmt.Println("Selected branches deleted.")
		break
	}
}

func branchDeleteMerged() {
	gitClient := git.NewClient()
	current, err := gitClient.GetCurrentBranch()
	if err != nil {
		fmt.Println("Error: failed to get current branch:", err)
		return
	}
	cmd := exec.Command("git", "branch", "--merged")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error: failed to get merged branches:", err)
		return
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	branches := []string{}
	for _, l := range lines {
		b := strings.TrimSpace(strings.TrimPrefix(l, "* "))
		if b != "" && b != current {
			branches = append(branches, b)
		}
	}
	if len(branches) == 0 {
		fmt.Println("No merged local branches.")
		return
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\033[1;36mSelect merged local branches to delete by number (space separated, all: select all, none: deselect all, e.g. 1 3 5):\033[0m")
		for i, b := range branches {
			fmt.Printf("  [\033[1;33m%d\033[0m] %s\n", i+1, b)
		}
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Println("Cancelled.")
			return
		}
		if input == "all" {
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
				fmt.Printf("\033[1;31mInvalid number: %s\033[0m\n", idx)
				valid = false
				break
			}
			tmp = append(tmp, branches[n-1])
		}
		if !valid {
			continue
		}
		for _, b := range tmp {
			cmd := exec.Command("git", "branch", "-d", b)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Printf("Error: failed to delete %s: %v\n", b, err)
			}
		}
		fmt.Println("Selected merged branches deleted.")
		break
	}
}
