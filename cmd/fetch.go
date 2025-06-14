package cmd

import (
	"fmt"

	"github.com/bmf-san/ggc/git"
)

func Fetch(args []string) {
	if len(args) > 0 && args[0] == "--prune" {
		err := git.FetchPrune()
		if err != nil {
			fmt.Println("Error:", err)
		}
		return
	}
	ShowFetchHelp()
}

func ShowFetchHelp() {
	fmt.Println("Usage: ggc fetch --prune")
}
