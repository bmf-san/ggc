package cmd

import (
	"fmt"

	"github.com/bmf-san/gcl/git"
)

func Push(args []string) {
	if len(args) > 0 && args[0] == "current" {
		err := git.PushCurrentBranch()
		if err != nil {
			fmt.Println("エラー:", err)
			return
		}
		fmt.Println("現在のブランチをpushしました")
		return
	}
	ShowPushHelp()
}

func ShowPushHelp() {
	fmt.Println("使用例: gg push current")
}
