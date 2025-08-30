// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

// Adder provides functionality for the add command.
type Adder struct {
	execCommand func(name string, arg ...string) *exec.Cmd
}

// NewAdder creates a new Adder.
func NewAdder() *Adder {
	return &Adder{execCommand: exec.Command}
}

// Add executes the add command with the given arguments.
func (a *Adder) Add(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: ggc add <file> | ggc add -i|--interactive | ggc add -p")
		return
	}
	if len(args) == 1 && (args[0] == "-i" || args[0] == "--interactive") {
		cmd := a.execCommand("git", "add", "-i")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			fmt.Println("error:", err)
		}
		return
	}
	if len(args) == 1 && args[0] == "-p" {
		cmd := a.execCommand("git", "add", "-p")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			fmt.Println("error:", err)
		}
		return
	}
	cmd := a.execCommand("git", append([]string{"add"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Println("error:", err)
	}
}
