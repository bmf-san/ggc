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

func ShowBranchHelp() {
	fmt.Println("使用例: gcl branch current | gcl branch checkout")
}
