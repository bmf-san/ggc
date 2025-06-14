package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bmf-san/gcl/git"
)

func Branch(args []string) {
	if len(args) == 1 && args[0] == "current" {
		branch, err := git.GetCurrentBranch()
		if err != nil {
			fmt.Println("エラー:", err)
			return
		}
		fmt.Println(branch)
		return
	}
	if len(args) == 1 && args[0] == "checkout" {
		branchCheckout()
		return
	}
	if len(args) == 1 && args[0] == "checkout-remote" {
		branchCheckoutRemote()
		return
	}
	if len(args) == 1 && args[0] == "delete" {
		branchDelete()
		return
	}
	if len(args) == 1 && args[0] == "delete-merged" {
		branchDeleteMerged()
		return
	}
	ShowBranchHelp()
}

func branchCheckout() {
	branches, err := git.ListLocalBranches()
	if err != nil {
		fmt.Println("エラー:", err)
		return
	}
	if len(branches) == 0 {
		fmt.Println("ローカルブランチが見つかりません")
		return
	}
	fmt.Println("ローカルブランチ一覧:")
	for i, b := range branches {
		fmt.Printf("[%d] %s\n", i+1, b)
	}
	fmt.Print("チェックアウトする番号を入力してください: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(branches) {
		fmt.Println("無効な番号です")
		return
	}
	branch := branches[idx-1]
	cmd := exec.Command("git", "checkout", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("エラー:", err)
	}
}

func branchCheckoutRemote() {
	branches, err := git.ListRemoteBranches()
	if err != nil {
		fmt.Println("エラー:", err)
		return
	}
	if len(branches) == 0 {
		fmt.Println("リモートブランチが見つかりません")
		return
	}
	fmt.Println("リモートブランチ一覧:")
	for i, b := range branches {
		fmt.Printf("[%d] %s\n", i+1, b)
	}
	fmt.Print("チェックアウトする番号を入力してください: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(branches) {
		fmt.Println("無効な番号です")
		return
	}
	remoteBranch := branches[idx-1]
	// origin/feature/foo → feature/foo
	parts := strings.SplitN(remoteBranch, "/", 2)
	if len(parts) != 2 {
		fmt.Println("無効なリモートブランチ名です")
		return
	}
	localBranch := parts[1]
	cmd := exec.Command("git", "checkout", "-b", localBranch, "--track", remoteBranch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("エラー:", err)
	}
}

func branchDelete() {
	branches, err := git.ListLocalBranches()
	if err != nil {
		fmt.Println("エラー:", err)
		return
	}
	if len(branches) == 0 {
		fmt.Println("ローカルブランチが見つかりません")
		return
	}
	fmt.Println("削除するローカルブランチを番号で選択（スペース区切り, 例: 1 3 5）:")
	for i, b := range branches {
		fmt.Printf("[%d] %s\n", i+1, b)
	}
	fmt.Print("> ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Println("キャンセルしました")
		return
	}
	indices := strings.Fields(input)
	selected := []string{}
	for _, idx := range indices {
		n, err := strconv.Atoi(idx)
		if err != nil || n < 1 || n > len(branches) {
			fmt.Printf("無効な番号: %s\n", idx)
			return
		}
		selected = append(selected, branches[n-1])
	}
	if len(selected) == 0 {
		fmt.Println("何も選択されませんでした")
		return
	}
	for _, b := range selected {
		cmd := exec.Command("git", "branch", "-d", b)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("エラー: %s の削除に失敗しました: %v\n", b, err)
		}
	}
	fmt.Println("選択したブランチを削除しました")
}

func branchDeleteMerged() {
	// 現在のブランチを取得
	current, err := git.GetCurrentBranch()
	if err != nil {
		fmt.Println("エラー: 現在のブランチ取得に失敗:", err)
		return
	}
	// マージ済みローカルブランチ一覧取得（現在のブランチ以外）
	cmd := exec.Command("git", "branch", "--merged")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("エラー: マージ済みブランチ取得に失敗:", err)
		return
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	branches := []string{}
	for _, l := range lines {
		b := strings.TrimSpace(strings.TrimPrefix(l, "* "))
		if b != "" && b != current {
			branches = append(branches, b)
		}
	}
	if len(branches) == 0 {
		fmt.Println("マージ済みローカルブランチはありません")
		return
	}
	fmt.Println("削除するマージ済みローカルブランチを番号で選択（スペース区切り, 例: 1 3 5）:")
	for i, b := range branches {
		fmt.Printf("[%d] %s\n", i+1, b)
	}
	fmt.Print("> ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Println("キャンセルしました")
		return
	}
	indices := strings.Fields(input)
	selected := []string{}
	for _, idx := range indices {
		n, err := strconv.Atoi(idx)
		if err != nil || n < 1 || n > len(branches) {
			fmt.Printf("無効な番号: %s\n", idx)
			return
		}
		selected = append(selected, branches[n-1])
	}
	if len(selected) == 0 {
		fmt.Println("何も選択されませんでした")
		return
	}
	for _, b := range selected {
		cmd := exec.Command("git", "branch", "-d", b)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("エラー: %s の削除に失敗しました: %v\n", b, err)
		}
	}
	fmt.Println("選択したマージ済みブランチを削除しました")
}

func ShowBranchHelp() {
	fmt.Println("使用例: gcl branch current | gcl branch checkout | gcl branch checkout-remote | gcl branch delete | gcl branch delete-merged")
}
