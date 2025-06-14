package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

func Add(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: ggc add <file> | ggc add -p")
		return
	}
	if len(args) == 1 && args[0] == "-p" {
		cmd := exec.Command("git", "add", "-p")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			fmt.Println("error:", err)
		}
		return
	}
	cmd := exec.Command("git", append([]string{"add"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Println("error:", err)
	}
}
