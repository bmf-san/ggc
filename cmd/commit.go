package cmd

import (
	"fmt"

	"github.com/bmf-san/gcl/git"
)

func Commit(args []string) {
	if len(args) > 0 && args[0] == "allow-empty" {
		err := git.CommitAllowEmpty()
		if err != nil {
			fmt.Println("エラー:", err)
			return
		}
		fmt.Println("空コミットを作成しました")
		return
	}
	ShowCommitHelp()
}

func ShowCommitHelp() {
	fmt.Println("使用例: gg commit allow-empty")
}
