package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bmf-san/ggc/git"
)

type Brancher struct {
	GetCurrentBranch  func() (string, error)
	ListLocalBranches func() ([]string, error)
	execCommand       func(name string, arg ...string) *exec.Cmd
	inputReader       *bufio.Reader
}

func NewBrancher() *Brancher {
	return &Brancher{
		GetCurrentBranch:  git.GetCurrentBranch,
		ListLocalBranches: git.ListLocalBranches,
		execCommand:       exec.Command,
		inputReader:       bufio.NewReader(os.Stdin),
	}
}

func (b *Brancher) Branch(args []string) {
	if len(args) == 1 && args[0] == "current" {
		branch, err := b.GetCurrentBranch()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(branch)
		return
	}
	if len(args) == 1 && args[0] == "checkout" {
		b.branchCheckout()
		return
	}
	if len(args) == 1 && args[0] == "checkout-remote" {
		branchCheckoutRemote()
		return
	}
	if len(args) == 1 && args[0] == "delete" {
		branchDelete()
		return
	}
	if len(args) == 1 && args[0] == "delete-merged" {
		branchDeleteMerged()
		return
	}
	if len(args) == 1 && args[0] == "create" {
		branchCreate()
		return
	}
	ShowBranchHelp()
}

func (b *Brancher) branchCheckout() {
	branches, err := b.ListLocalBranches()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(branches) == 0 {
		fmt.Println("No local branches found.")
		return
	}
	fmt.Println("Local branches:")
	for i, br := range branches {
		fmt.Printf("[%d] %s\n", i+1, br)
	}
	fmt.Print("Enter the number to checkout: ")
	input, _ := b.inputReader.ReadString('\n')
	input = strings.TrimSpace(input)
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(branches) {
		fmt.Println("Invalid number.")
		return
	}
	branch := branches[idx-1]
	cmd := b.execCommand("git", "checkout", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}

func branchCheckoutRemote() {
	branches, err := git.ListRemoteBranches()
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
	branches, err := git.ListLocalBranches()
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
	current, err := git.GetCurrentBranch()
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

func branchCreate() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter new branch name: ")
	branchName, _ := reader.ReadString('\n')
	branchName = strings.TrimSpace(branchName)
	if branchName == "" {
		fmt.Println("Branch name is empty. Cancelled.")
		return
	}
	cmd := exec.Command("git", "checkout", "-b", branchName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Branch '%s' created and checked out.\n", branchName)
}

func ShowBranchHelp() {
	fmt.Println("Usage: ggc branch current | ggc branch checkout | ggc branch checkout-remote | ggc branch create | ggc branch delete | ggc branch delete-merged")
}
