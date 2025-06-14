package cmd

import (
	"fmt"
	"os/exec"
)

func ResetClean() {
	resetCmd := exec.Command("git", "reset", "--hard", "HEAD")
	resetCmd.Stdout = nil
	resetCmd.Stderr = nil
	if err := resetCmd.Run(); err != nil {
		fmt.Printf("Error: git reset --hard HEAD failed: %v\n", err)
		return
	}
	cleanCmd := exec.Command("git", "clean", "-fd")
	cleanCmd.Stdout = nil
	cleanCmd.Stderr = nil
	if err := cleanCmd.Run(); err != nil {
		fmt.Printf("Error: git clean -fd failed: %v\n", err)
		return
	}
	fmt.Println("reset --hard HEAD and clean -fd done")
}
