package cmd

import (
	"fmt"

	"github.com/bmf-san/gcl/git"
)

func Commit(args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "allow-empty":
			err := git.CommitAllowEmpty()
			if err != nil {
				fmt.Println("エラー:", err)
				return
			}
			fmt.Println("空コミットを作成しました")
			return
		case "tmp":
			err := git.CommitTmp()
			if err != nil {
				fmt.Println("エラー:", err)
				return
			}
			fmt.Println("一時コミット(tmp)を作成しました")
			return
		}
	}
	ShowCommitHelp()
}

func ShowCommitHelp() {
	fmt.Println("使用例: gcl commit allow-empty | gcl commit tmp")
}
