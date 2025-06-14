package cmd

import (
	"fmt"

	"github.com/bmf-san/gcl/git"
)

func Fetch(args []string) {
	if len(args) > 0 && args[0] == "--prune" {
		err := git.FetchPrune()
		if err != nil {
			fmt.Println("エラー:", err)
		}
		return
	}
	ShowFetchHelp()
}

func ShowFetchHelp() {
	fmt.Println("使用例: gcl fetch --prune")
}
