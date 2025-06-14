package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func AddCommitPush() {
	// git add .
	addCmd := exec.Command("git", "add", ".")
	addCmd.Stdout = os.Stdout
	addCmd.Stderr = os.Stderr
	if err := addCmd.Run(); err != nil {
		fmt.Printf("エラー: git add . に失敗しました: %v\n", err)
		return
	}
	// コミットメッセージ入力
	fmt.Print("コミットメッセージを入力してください: ")
	reader := bufio.NewReader(os.Stdin)
	msg, _ := reader.ReadString('\n')
	msg = strings.TrimSpace(msg)
	if msg == "" {
		fmt.Println("キャンセルしました")
		return
	}
	// git commit
	commitCmd := exec.Command("git", "commit", "-m", msg)
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr
	if err := commitCmd.Run(); err != nil {
		fmt.Printf("エラー: git commit に失敗しました: %v\n", err)
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
	// git push
	pushCmd := exec.Command("git", "push", "origin", branch)
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	if err := pushCmd.Run(); err != nil {
		fmt.Printf("エラー: git push に失敗しました: %v\n", err)
		return
	}
	fmt.Println("add→commit→push 完了")
}
