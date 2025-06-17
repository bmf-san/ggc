package cmd

import (
	"fmt"

	"github.com/bmf-san/ggc/git"
)

type Committer struct {
	CommitAllowEmpty func() error
	CommitTmp        func() error
}

func NewCommitter() *Committer {
	return &Committer{
		CommitAllowEmpty: git.CommitAllowEmpty,
		CommitTmp:        git.CommitTmp,
	}
}

func (c *Committer) Commit(args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "allow-empty":
			err := c.CommitAllowEmpty()
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		case "tmp":
			err := c.CommitTmp()
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
	}
	ShowCommitHelp()
}

func ShowCommitHelp() {
	fmt.Println("Usage: ggc commit allow-empty | ggc commit tmp")
}
