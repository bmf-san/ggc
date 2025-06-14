package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Rebase(args []string) {
	if len(args) > 0 && args[0] == "interactive" {
		RebaseInteractive()
		return
	}
	ShowRebaseHelp()
}

func ShowRebaseHelp() {
	fmt.Println("Usage: ggc rebase interactive")
}

// Interactively rebase up to HEAD~N
func RebaseInteractive() {
	// 1. Get the last 10 commit logs
	cmd := exec.Command("git", "log", "--oneline", "-n", "10")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("error: failed to get git log: %v\n", err)
		return
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		fmt.Println("No commit history found")
		return
	}
	fmt.Println("Where do you want to rebase up to? Select a number (e.g., 3):")
	for i, line := range lines {
		fmt.Printf("  [%d] %s\n", i+1, line)
	}
	fmt.Print("> ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Println("Cancelled")
		return
	}
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(lines) {
		fmt.Println("Invalid number")
		return
	}
	// N commits before rebase
	N := idx
	rebaseCmd := exec.Command("git", "rebase", "-i", fmt.Sprintf("HEAD~%d", N))
	rebaseCmd.Stdin = os.Stdin
	rebaseCmd.Stdout = os.Stdout
	rebaseCmd.Stderr = os.Stderr
	if err := rebaseCmd.Run(); err != nil {
		fmt.Printf("error: git rebase failed: %v\n", err)
		return
	}
}
