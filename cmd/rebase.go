package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Rebase(args []string) {
	if len(args) > 0 && args[0] == "interactive" {
		RebaseInteractive()
		return
	}
	ShowRebaseHelp()
}

func ShowRebaseHelp() {
	fmt.Println("使用例: gcl rebase interactive")
}

// 対話的にHEAD~Nまでrebase
func RebaseInteractive() {
	// 1. 直近10件のコミット履歴を取得
	cmd := exec.Command("git", "log", "--oneline", "-n", "10")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("エラー: git log の取得に失敗しました: %v\n", err)
		return
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		fmt.Println("コミット履歴がありません")
		return
	}
	fmt.Println("どこまでrebaseしますか？番号を選択（例: 3）:")
	for i, line := range lines {
		fmt.Printf("  [%d] %s\n", i+1, line)
	}
	fmt.Print("> ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Println("キャンセルしました")
		return
	}
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(lines) {
		fmt.Println("無効な番号です")
		return
	}
	// Nコミット前までrebase
	N := idx
	rebaseCmd := exec.Command("git", "rebase", "-i", fmt.Sprintf("HEAD~%d", N))
	rebaseCmd.Stdin = os.Stdin
	rebaseCmd.Stdout = os.Stdout
	rebaseCmd.Stderr = os.Stderr
	if err := rebaseCmd.Run(); err != nil {
		fmt.Printf("エラー: git rebase に失敗しました: %v\n", err)
		return
	}
}
