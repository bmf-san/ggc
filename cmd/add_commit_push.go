package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func AddCommitPush() {
	// git add .
	addCmd := exec.Command("git", "add", ".")
	addCmd.Stdout = os.Stdout
	addCmd.Stderr = os.Stderr
	if err := addCmd.Run(); err != nil {
		fmt.Printf("Error: failed to add all files: %v\n", err)
		return
	}
	// Enter commit message
	fmt.Print("Enter commit message: ")
	reader := bufio.NewReader(os.Stdin)
	msg, _ := reader.ReadString('\n')
	msg = strings.TrimSpace(msg)
	if msg == "" {
		fmt.Println("Cancelled.")
		return
	}
	// git commit
	commitCmd := exec.Command("git", "commit", "-m", msg)
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr
	if err := commitCmd.Run(); err != nil {
		fmt.Printf("Error: failed to commit: %v\n", err)
		return
	}
	// Get current branch name
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOut, err := branchCmd.Output()
	if err != nil {
		fmt.Printf("Error: failed to get branch name: %v\n", err)
		return
	}
	branch := strings.TrimSpace(string(branchOut))
	// git push
	pushCmd := exec.Command("git", "push", "origin", branch)
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	if err := pushCmd.Run(); err != nil {
		fmt.Printf("Error: failed to push: %v\n", err)
		return
	}
	fmt.Println("add→commit→push done")
}
