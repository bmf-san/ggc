package cmd

import (
	"fmt"

	"github.com/bmf-san/gcl/git"
)

func Branch(args []string) {
	if len(args) > 0 && args[0] == "current" {
		branch, err := git.GetCurrentBranch()
		if err != nil {
			fmt.Println("エラー:", err)
			return
		}
		fmt.Println(branch)
		return
	}
	ShowBranchHelp()
}

func ShowBranchHelp() {
	fmt.Println("使用例: gg branch current")
}
