package cmd

import (
	"fmt"
	"io"
	"os"
)

type Helper struct {
	writer io.Writer
}

func NewHelper() *Helper {
	return &Helper{writer: os.Stdout}
}

func (h *Helper) ShowHelp() {
	if _, err := fmt.Fprint(h.writer, `ggc: A Go-based CLI tool to streamline Git operations

Usage:
  ggc <command> [subcommand] [options]

Main Commands:
  ggc add <file>              Stage file(s)
  ggc branch current          Show current branch name
  ggc branch checkout         Interactive branch switch
  ggc push current            Push current branch
  ggc push force              Force push current branch
  ggc pull current            Pull current branch
  ggc pull rebase             Pull with rebase
  ggc log simple              Show simple log
  ggc log graph               Show log with graph
  ggc commit allow-empty      Create empty commit
  ggc commit tmp              Temporary commit
  ggc fetch --prune           Fetch with prune
  ggc clean files             Clean files
  ggc clean dirs              Clean directories
  ggc reset clean             Reset and clean
  ggc commit-push             Interactive add/commit/push

Examples:
  ggc add .
  ggc branch current
  ggc branch checkout
  ggc push current
  ggc push force
  ggc pull current
  ggc pull rebase
  ggc log simple
  ggc log graph
  ggc commit allow-empty
  ggc commit tmp
  ggc fetch --prune
  ggc clean files
  ggc clean dirs
  ggc reset clean
  ggc commit-push
`); err != nil {
		// ignore error (for test/lint)
		_ = err
	}
}

// For backward compatibility
func ShowHelp() {
	NewHelper().ShowHelp()
}
