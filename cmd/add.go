package cmd

import (
	"fmt"
	"os/exec"
)

func Add(args []string) {
	if len(args) == 0 {
		fmt.Println("使用例: gcl add <file>")
		return
	}
	cmd := exec.Command("git", append([]string{"add"}, args...)...)
	err := cmd.Run()
	if err != nil {
		fmt.Println("エラー:", err)
	}
}
