package cmd

import (
	"fmt"
	"os/exec"
)

func Stash(args []string) {
	if len(args) > 0 && args[0] == "trash" {
		// git add .
		addCmd := exec.Command("git", "add", ".")
		addCmd.Stdout = nil
		addCmd.Stderr = nil
		if err := addCmd.Run(); err != nil {
			fmt.Printf("エラー: git add . に失敗しました: %v\n", err)
			return
		}
		// git stash
		stashCmd := exec.Command("git", "stash")
		stashCmd.Stdout = nil
		stashCmd.Stderr = nil
		if err := stashCmd.Run(); err != nil {
			fmt.Printf("エラー: git stash に失敗しました: %v\n", err)
			return
		}
		fmt.Println("add . → stash を実行しました")
		return
	}
	ShowStashHelp()
}

func ShowStashHelp() {
	fmt.Println("使用例: gcl stash trash")
}
