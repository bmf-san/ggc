package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bmf-san/gcl/git"
)

func Log(args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "simple":
			err := git.LogSimple()
			if err != nil {
				fmt.Println("エラー:", err)
			}
			return
		case "graph":
			err := logGraph()
			if err != nil {
				fmt.Println("エラー:", err)
			}
			return
		}
	}
	ShowLogHelp()
}

func logGraph() error {
	cmd := exec.Command("git", "log", "--graph")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ShowLogHelp() {
	fmt.Println("使用例: gcl log simple | gcl log graph")
}
