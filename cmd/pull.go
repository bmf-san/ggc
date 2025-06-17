package cmd

import (
	"fmt"

	"github.com/bmf-san/ggc/git"
)

type Puller struct {
	PullCurrentBranch       func() error
	PullRebaseCurrentBranch func() error
}

func NewPuller() *Puller {
	return &Puller{
		PullCurrentBranch:       git.PullCurrentBranch,
		PullRebaseCurrentBranch: git.PullRebaseCurrentBranch,
	}
}

func (p *Puller) Pull(args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "current":
			err := p.PullCurrentBranch()
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		case "rebase":
			err := p.PullRebaseCurrentBranch()
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

// 旧インターフェース維持用ラッパー
// func Pull(args []string) {
// 	NewPuller().Pull(args)
// }
