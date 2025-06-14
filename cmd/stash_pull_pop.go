package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

func StashPullPop() {
	// git stash
	stashCmd := exec.Command("git", "stash")
	stashCmd.Stdout = nil
	stashCmd.Stderr = nil
	if err := stashCmd.Run(); err != nil {
		fmt.Printf("error: failed to git stash: %v\n", err)
		return
	}
	// Get current branch name
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOut, err := branchCmd.Output()
	if err != nil {
		fmt.Printf("error: failed to get branch name: %v\n", err)
		return
	}
	branch := strings.TrimSpace(string(branchOut))
	// git pull
	pullCmd := exec.Command("git", "pull", "origin", branch)
	pullCmd.Stdout = nil
	pullCmd.Stderr = nil
	if err := pullCmd.Run(); err != nil {
		fmt.Printf("error: failed to git pull: %v\n", err)
		return
	}
	// git stash pop
	popCmd := exec.Command("git", "stash", "pop")
	popCmd.Stdout = nil
	popCmd.Stderr = nil
	if err := popCmd.Run(); err != nil {
		fmt.Printf("error: failed to git stash pop: %v\n", err)
		return
	}
	fmt.Println("stash→pull→pop done")
}
