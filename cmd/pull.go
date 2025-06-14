package cmd

import (
	"fmt"

	"github.com/bmf-san/ggc/git"
)

func Pull(args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "current":
			err := git.PullCurrentBranch()
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		case "rebase":
			err := git.PullRebaseCurrentBranch()
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
	}
	ShowPullHelp()
}

func ShowPullHelp() {
	fmt.Println("Usage: ggc pull current | ggc pull rebase")
}
