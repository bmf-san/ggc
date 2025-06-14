package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Show interactive add/commit/push
func CommitPushInteractive() {
	// 1. git status --porcelain で変更ファイル一覧取得
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("エラー: git status の取得に失敗しました: %v\n", err)
		return
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		fmt.Println("変更されたファイルはありません")
		return
	}

	// 2. ファイル一覧を番号付きで表示
	files := []string{}
	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		files = append(files, strings.TrimSpace(line[2:]))
	}
	if len(files) == 0 {
		fmt.Println("ステージング可能なファイルはありません")
		return
	}
	reader := bufio.NewReader(os.Stdin)
	selected := []string{}
	for {
		fmt.Println("\033[1;36m追加するファイルを番号で選択（スペース区切り, all:全選択, none:全解除, 例: 1 3 5）:\033[0m")
		for i, f := range files {
			fmt.Printf("  [\033[1;33m%d\033[0m] %s\n", i+1, f)
		}
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Println("キャンセルしました")
			return
		}
		if input == "all" {
			selected = files
			break
		}
		if input == "none" {
			selected = []string{}
			continue
		}
		indices := strings.Fields(input)
		tmp := []string{}
		valid := true
		for _, idx := range indices {
			n, err := strconv.Atoi(idx)
			if err != nil || n < 1 || n > len(files) {
				fmt.Printf("\033[1;31m無効な番号: %s\033[0m\n", idx)
				valid = false
				break
			}
			tmp = append(tmp, files[n-1])
		}
		if !valid {
			continue
		}
		selected = tmp
		if len(selected) == 0 {
			fmt.Println("\033[1;33m何も選択されませんでした\033[0m")
			continue
		}
		// 確認プロンプト
		fmt.Printf("\033[1;32m選択したファイル: %v\033[0m\n", selected)
		fmt.Print("このファイルでaddしますか？ (y/n): ")
		ans, _ := reader.ReadString('\n')
		ans = strings.TrimSpace(ans)
		if ans == "y" || ans == "Y" {
			break
		}
	}

	// 3. git add 実行
	addArgs := append([]string{"add"}, selected...)
	addCmd := exec.Command("git", addArgs...)
	addCmd.Stdout = os.Stdout
	addCmd.Stderr = os.Stderr
	if err := addCmd.Run(); err != nil {
		fmt.Printf("エラー: git add に失敗しました: %v\n", err)
		return
	}

	// 4. コミットメッセージ入力
	fmt.Print("コミットメッセージを入力してください: ")
	msg, _ := reader.ReadString('\n')
	msg = strings.TrimSpace(msg)
	if msg == "" {
		fmt.Println("キャンセルしました")
		return
	}

	// 5. git commit 実行
	commitCmd := exec.Command("git", "commit", "-m", msg)
	commitCmd.Stdout = os.Stdout
	commitCmd.Stderr = os.Stderr
	if err := commitCmd.Run(); err != nil {
		fmt.Printf("エラー: git commit に失敗しました: %v\n", err)
		return
	}

	// 6. 現在のブランチ名取得
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOut, err := branchCmd.Output()
	if err != nil {
		fmt.Printf("エラー: ブランチ名の取得に失敗しました: %v\n", err)
		return
	}
	branch := strings.TrimSpace(string(branchOut))

	// 7. git push 実行
	pushCmd := exec.Command("git", "push", "origin", branch)
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	if err := pushCmd.Run(); err != nil {
		fmt.Printf("エラー: git push に失敗しました: %v\n", err)
		return
	}

	fmt.Println("完了しました！")
}
