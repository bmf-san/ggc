package cmd

import (
	"fmt"

	"github.com/bmf-san/gcl/git"
)

func Clean(args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "files":
			err := git.CleanFiles()
			if err != nil {
				fmt.Println("エラー:", err)
			}
			return
		case "dirs":
			err := git.CleanDirs()
			if err != nil {
				fmt.Println("エラー:", err)
			}
			return
		}
	}
	ShowCleanHelp()
}

func ShowCleanHelp() {
	fmt.Println("使用例: gcl clean files | gcl clean dirs")
}
