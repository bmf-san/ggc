package cmd

import (
	"fmt"

	"github.com/bmf-san/ggc/git"
)

type Resetter struct {
	ResetClean func() error
}

func NewResetter() *Resetter {
	return &Resetter{
		ResetClean: git.ResetClean,
	}
}

func (r *Resetter) Reset(args []string) {
	if len(args) > 0 && args[0] == "clean" {
		err := r.ResetClean()
		if err != nil {
			fmt.Println("Error:", err)
		}
		return
	}
	ShowResetHelp()
}

func ShowResetHelp() {
	fmt.Println("Usage: ggc reset clean")
}
