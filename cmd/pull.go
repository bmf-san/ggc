package cmd

import (
	"fmt"

	"github.com/bmf-san/gcl/git"
)

func Pull(args []string) {
	if len(args) > 0 && args[0] == "current" {
		err := git.PullCurrentBranch()
		if err != nil {
			fmt.Println("エラー:", err)
			return
		}
		fmt.Println("現在のブランチをpullしました")
		return
	}
	ShowPullHelp()
}

func ShowPullHelp() {
	fmt.Println("使用例: gg pull current")
}
