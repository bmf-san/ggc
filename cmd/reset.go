package cmd

import (
	"fmt"

	"github.com/bmf-san/gcl/git"
)

func Reset(args []string) {
	if len(args) > 0 && args[0] == "clean" {
		err := git.ResetClean()
		if err != nil {
			fmt.Println("エラー:", err)
			return
		}
		fmt.Println("reset --hard HEAD + clean -fd を実行しました")
		return
	}
	ShowResetHelp()
}

func ShowResetHelp() {
	fmt.Println("使用例: gcl reset clean")
}
