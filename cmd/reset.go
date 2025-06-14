package cmd

import (
	"fmt"

	"github.com/bmf-san/gcl/git"
)

func Reset(args []string) {
	if len(args) > 0 && args[0] == "clean" {
		err := git.ResetClean()
		if err != nil {
			fmt.Println("Error:", err)
		}
		return
	}
	ShowResetHelp()
}

func ShowResetHelp() {
	fmt.Println("Usage: gcl reset clean")
}
