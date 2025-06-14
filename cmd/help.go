package cmd

import "fmt"

func ShowHelp() {
	fmt.Print(`gcl: A Go-based CLI tool to streamline Git operations

Usage:
  gcl <command> [subcommand] [options]

Main Commands:
  gcl add <file>              Stage file(s)
  gcl branch current          Show current branch name
  gcl branch checkout         Interactive branch switch
  gcl push current            Push current branch
  gcl push force              Force push current branch
  gcl pull current            Pull current branch
  gcl pull rebase             Pull with rebase
  gcl log simple              Show simple log
  gcl log graph               Show log with graph
  gcl commit allow-empty      Create empty commit
  gcl commit tmp              Temporary commit
  gcl fetch --prune           Fetch with prune
  gcl clean files             Clean files
  gcl clean dirs              Clean directories
  gcl reset clean             Reset and clean
  gcl commit-push             Interactive add/commit/push

Examples:
  gcl add .
  gcl branch current
  gcl branch checkout
  gcl push current
  gcl push force
  gcl pull current
  gcl pull rebase
  gcl log simple
  gcl log graph
  gcl commit allow-empty
  gcl commit tmp
  gcl fetch --prune
  gcl clean files
  gcl clean dirs
  gcl reset clean
  gcl commit-push
`)
}
