// Package templates provides templates for help messages.
package templates

import (
	"bytes"
	"golang.org/x/term"
	"os"
	"text/template"
)

// HelpData contains data for help message templates.
type HelpData struct {
	Logo        string
	Usage       string
	Description string
	Examples    []string
}

// Templates for help messages.
var (
	mainHelpTemplate = `{{.Logo}}ggc: A Go-based CLI tool to streamline Git operations

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
  ggc commit amend <message>  Amend to previous commit
  ggc commit amend --no-edit  Amend without editing commit message
  ggc commit allow-empty      Create empty commit
  ggc commit tmp              Temporary commit
  ggc commit-push-interactive Interactive add/commit/push
  ggc complete <shell>        Generate shell completion script (bash|zsh)
  ggc fetch --prune          Fetch and remove stale remote-tracking branches
  ggc diff                    Show changes between commits, commit and working tree
  ggc tag                     Create, list, and delete tags
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
  ggc version                 Show current ggc version
  ggc reset                   Reset and clean
  ggc reset-clean            Reset to HEAD and clean untracked files
  ggc stash                   Stash changes
  ggc stash-pull-pop          Stash, pull, and pop all at once
  ggc status                  Show the working tree status
`

	commandHelpTemplate = `{{.Logo}}
Usage: {{.Usage}}

Description:
  {{.Description}}

Examples:
{{range .Examples}}  {{.}}
{{end}}
`
)

// MainHelpData contains data for main help message.
type MainHelpData struct {
	Logo string
}

func selectLogo() string {
	if termWidth, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		if termWidth < 50 {
			return SmallLogo
		}
	}
	return Logo
}

// RenderMainHelp renders the main help message.
func RenderMainHelp() (string, error) {
	tmpl, err := template.New("mainHelp").Parse(mainHelpTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	data := MainHelpData{
		Logo: selectLogo(),
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
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
