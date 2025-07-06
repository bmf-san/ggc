// Package templates provides templates for help messages.
package templates

import (
	"bytes"
	"text/template"
)

// HelpData contains data for help message templates.
type HelpData struct {
	Usage       string
	Description string
	Examples    []string
}

// Templates for help messages.
var (
	mainHelpTemplate = `ggc: A Go-based CLI tool to streamline Git operations

Usage:
  ggc <command> [subcommand] [options]

Main Commands:
  ggc add <file>              Stage file(s)
  ggc add .                   Stage all changes
  ggc add -p                  Stage changes interactively
  ggc add-commit-push         Add, commit, and push all at once
  ggc branch current          Show current branch name
  ggc branch checkout         Interactive branch switch
  ggc branch checkout-remote  Create and checkout new local branch from remote
  ggc branch delete          Interactive delete of local branches
  ggc branch delete-merged   Interactive delete of merged local branches
  ggc clean files             Clean files
  ggc clean dirs              Clean directories
  ggc clean-interactive       Interactive file cleaning
  ggc commit allow-empty      Create empty commit
  ggc commit tmp              Temporary commit
  ggc commit-push-interactive Interactive add/commit/push
  ggc complete <shell>        Generate shell completion script (bash|zsh)
  ggc fetch --prune          Fetch and remove stale remote-tracking branches
  ggc diff                    Show changes between commits, commit and working tree
  ggc log simple              Show simple log
  ggc log graph               Show log with graph
  ggc pull current            Pull current branch
  ggc pull rebase             Pull with rebase
  ggc pull-rebase-push        Pull with rebase and push all at once
  ggc push current            Push current branch
  ggc push force              Force push current branch
  ggc rebase                  Rebase current branch
  ggc remote list             Show remotes
  ggc remote add <n> <url>    Add remote
  ggc remote remove <n>       Remove remote
  ggc remote set-url <n> <url> Change remote URL
  ggc reset                   Reset and clean
  ggc reset-clean            Reset to HEAD and clean untracked files
  ggc stash                   Stash changes
  ggc stash-pull-pop          Stash, pull, and pop all at once
  ggc status                  Show the working tree status
`

	commandHelpTemplate = `Usage: {{.Usage}}

Description:
  {{.Description}}

Examples:
{{range .Examples}}  {{.}}
{{end}}
`
)

// RenderMainHelp renders the main help message.
func RenderMainHelp() string {
	return mainHelpTemplate
}

// RenderCommandHelp renders help message for a specific command.
func RenderCommandHelp(data HelpData) (string, error) {
	tmpl, err := template.New("commandHelp").Parse(commandHelpTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
