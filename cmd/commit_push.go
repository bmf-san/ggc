package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Show interactive add/commit/push
func CommitPushInteractive() {
	// 1. Get changed files with git status --porcelain
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: failed to get git status: %v\n", err)
		return
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		fmt.Println("No changed files.")
		return
	}

	// 2. Show file list with numbers
	files := []string{}
	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		files = append(files, strings.TrimSpace(line[2:]))
	}
	if len(files) == 0 {
		fmt.Println("No files to stage.")
		return
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\033[1;36mSelect files to add by number (space separated, all: select all, none: deselect all, e.g. 1 3 5):\033[0m")
		for i, f := range files {
			fmt.Printf("  [\033[1;33m%d\033[0m] %s\n", i+1, f)
		}
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Println("Cancelled.")
			return
		}
		if input == "all" {
			addArgs := append([]string{"add"}, files...)
			addCmd := exec.Command("git", addArgs...)
			addCmd.Stdout = os.Stdout
			addCmd.Stderr = os.Stderr
			if err := addCmd.Run(); err != nil {
				fmt.Printf("Error: failed to add files: %v\n", err)
				return
			}
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
			if err != nil || n < 1 || n > len(files) {
				fmt.Printf("\033[1;31mInvalid number: %s\033[0m\n", idx)
				valid = false
				break
			}
			tmp = append(tmp, files[n-1])
		}
		if !valid {
			continue
		}
		if len(tmp) == 0 {
			fmt.Println("\033[1;33mNothing selected.\033[0m")
			continue
		}
		// Confirmation prompt
		fmt.Printf("\033[1;32mSelected files: %v\033[0m\n", tmp)
		fmt.Print("Add these files? (y/n): ")
		ans, _ := reader.ReadString('\n')
		ans = strings.TrimSpace(ans)
		if ans == "y" || ans == "Y" {
			addArgs := append([]string{"add"}, tmp...)
			addCmd := exec.Command("git", addArgs...)
			addCmd.Stdout = os.Stdout
			addCmd.Stderr = os.Stderr
			if err := addCmd.Run(); err != nil {
				fmt.Printf("Error: failed to add files: %v\n", err)
				return
			}
			break
		}
	}

	fmt.Print("\n\r")
	fmt.Print("Enter commit message: ")
	msg, _ := reader.ReadString('\n')
	msg = strings.TrimSpace(msg)
	if msg == "" {
		fmt.Println("Cancelled.")
		return
	}

	// 5. Run git commit
	commitCmd := exec.Command("git", "commit", "-m", msg)
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr
	if err := commitCmd.Run(); err != nil {
		fmt.Printf("Error: failed to commit: %v\n", err)
		return
	}

	// 6. Get current branch name
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOut, err := branchCmd.Output()
	if err != nil {
		fmt.Printf("Error: failed to get branch name: %v\n", err)
		return
	}
	branch := strings.TrimSpace(string(branchOut))

	// 7. Run git push
	pushCmd := exec.Command("git", "push", "origin", branch)
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	if err := pushCmd.Run(); err != nil {
		fmt.Printf("Error: failed to push: %v\n", err)
		return
	}

	fmt.Println("Done!")
}
