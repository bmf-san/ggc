package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

func PullRebasePush() {
	// Get current branch name
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOut, err := branchCmd.Output()
	if err != nil {
		fmt.Printf("Error: Failed to get branch name: %v\n", err)
		return
	}
	branch := strings.TrimSpace(string(branchOut))
	// git pull
	pullCmd := exec.Command("git", "pull", "origin", branch)
	pullCmd.Stdout = nil
	pullCmd.Stderr = nil
	if err := pullCmd.Run(); err != nil {
		fmt.Printf("Error: Failed to git pull: %v\n", err)
		return
	}
	// git rebase origin/main
	rebaseCmd := exec.Command("git", "rebase", "origin/main")
	rebaseCmd.Stdout = nil
	rebaseCmd.Stderr = nil
	if err := rebaseCmd.Run(); err != nil {
		fmt.Printf("Error: Failed to git rebase: %v\n", err)
		return
	}
	// git push
	pushCmd := exec.Command("git", "push", "origin", branch)
	pushCmd.Stdout = nil
	pushCmd.Stderr = nil
	if err := pushCmd.Run(); err != nil {
		fmt.Printf("Error: Failed to git push: %v\n", err)
		return
	}
	fmt.Println("pull→rebase→push completed")
}
