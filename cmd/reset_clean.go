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
		fmt.Printf("エラー: git reset --hard HEAD に失敗しました: %v\n", err)
		return
	}
	cleanCmd := exec.Command("git", "clean", "-fd")
	cleanCmd.Stdout = nil
	cleanCmd.Stderr = nil
	if err := cleanCmd.Run(); err != nil {
		fmt.Printf("エラー: git clean -fd に失敗しました: %v\n", err)
		return
	}
	fmt.Println("reset --hard HEAD→clean -fd 完了")
}
