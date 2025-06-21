package cmd

import (
	"fmt"

	"github.com/bmf-san/ggc/git"
)

type Pusher struct {
	PushCurrentBranch      func() error
	PushForceCurrentBranch func() error
}

func NewPusher() *Pusher {
	return &Pusher{
		PushCurrentBranch:      git.PushCurrentBranch,
		PushForceCurrentBranch: git.PushForceCurrentBranch,
	}
}

func (p *Pusher) Push(args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "current":
			err := p.PushCurrentBranch()
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		case "force":
			err := p.PushForceCurrentBranch()
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
	}
	ShowPushHelp()
}

func ShowPushHelp() {
	fmt.Println("Usage: ggc push current | ggc push force")
}
