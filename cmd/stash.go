package cmd

import (
	"fmt"
	"os/exec"
)

type Stasher struct {
	execCommand func(name string, arg ...string) *exec.Cmd
}

func NewStasher() *Stasher {
	return &Stasher{execCommand: exec.Command}
}

func (s *Stasher) Stash(args []string) {
	if len(args) > 0 && args[0] == "trash" {
		addCmd := s.execCommand("git", "add", ".")
		addCmd.Stdout = nil
		addCmd.Stderr = nil
		if err := addCmd.Run(); err != nil {
			fmt.Printf("Error: failed to add all files: %v\n", err)
			return
		}
		stashCmd := s.execCommand("git", "stash")
		stashCmd.Stdout = nil
		stashCmd.Stderr = nil
		if err := stashCmd.Run(); err != nil {
			fmt.Printf("Error: failed to stash: %v\n", err)
			return
		}
		fmt.Println("add . → stash done")
		return
	}
	ShowStashHelp()
}

func ShowStashHelp() {
	fmt.Println("Usage: ggc stash trash")
}

// 旧インターフェース維持用ラッパー
// func Stash(args []string) {
// 	NewStasher().Stash(args)
// }
