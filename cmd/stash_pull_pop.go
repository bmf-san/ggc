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
		fmt.Printf("エラー: git stash に失敗しました: %v\n", err)
		return
	}
	// 現在のブランチ名取得
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOut, err := branchCmd.Output()
	if err != nil {
		fmt.Printf("エラー: ブランチ名の取得に失敗しました: %v\n", err)
		return
	}
	branch := strings.TrimSpace(string(branchOut))
	// git pull
	pullCmd := exec.Command("git", "pull", "origin", branch)
	pullCmd.Stdout = nil
	pullCmd.Stderr = nil
	if err := pullCmd.Run(); err != nil {
		fmt.Printf("エラー: git pull に失敗しました: %v\n", err)
		return
	}
	// git stash pop
	popCmd := exec.Command("git", "stash", "pop")
	popCmd.Stdout = nil
	popCmd.Stderr = nil
	if err := popCmd.Run(); err != nil {
		fmt.Printf("エラー: git stash pop に失敗しました: %v\n", err)
		return
	}
	fmt.Println("stash→pull→pop 完了")
}
