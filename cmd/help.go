// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/bmf-san/ggc/cmd/templates"
)

// Helper provides help message functionality.
type Helper struct {
	outputWriter io.Writer
}

// NewHelper creates a new Helper.
func NewHelper() *Helper {
	return &Helper{
		outputWriter: os.Stdout,
	}
}

// ShowHelp shows the main help message.
func (h *Helper) ShowHelp() {
	_, _ = fmt.Fprint(h.outputWriter, templates.RenderMainHelp())
}

// ShowCommandHelp shows help message for a command.
func (h *Helper) ShowCommandHelp(data templates.HelpData) {
	helpMsg, err := templates.RenderCommandHelp(data)
	if err != nil {
		_, _ = fmt.Fprintf(h.outputWriter, "Error: %v\n", err)
		return
	}
	_, _ = fmt.Fprint(h.outputWriter, helpMsg)
}

// ShowAddHelp shows help message for add command.
func (h *Helper) ShowAddHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc add <file> | ggc add -p",
		Description: "Add file contents to the index",
		Examples: []string{
			"ggc add file.txt   # Add a specific file",
			"ggc add -p         # Add changes interactively",
		},
	})
}

// ShowBranchHelp shows help message for branch command.
func (h *Helper) ShowBranchHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc branch <command>",
		Description: "List, create, or delete branches",
		Examples: []string{
			"ggc branch current     # Show current branch",
			"ggc branch checkout    # Checkout existing branch",
			"ggc branch create      # Create and checkout new branch",
			"ggc branch delete      # Delete a branch",
		},
	})
}

// ShowCleanHelp shows help message for clean command.
func (h *Helper) ShowCleanHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc clean <command>",
		Description: "Clean untracked files and directories",
		Examples: []string{
			"ggc clean files    # Clean untracked files",
			"ggc clean dirs     # Clean untracked directories",
		},
	})
}

// ShowCommitHelp shows help message for commit command.
func (h *Helper) ShowCommitHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc commit <message>",
		Description: "Commit staged changes",
		Examples: []string{
			"ggc commit allow-empty   # Create empty commit",
			"ggc commit tmp           # Create temporary commit",
		},
	})
}

// ShowLogHelp shows help message for log command.
func (h *Helper) ShowLogHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc log <command>",
		Description: "Show commit logs",
		Examples: []string{
			"ggc log simple   # Show commit logs in a simple format",
			"ggc log graph    # Show commit logs with a graph",
		},
	})
}

// ShowPullHelp shows help message for pull command.
func (h *Helper) ShowPullHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc pull <command>",
		Description: "Pull changes from remote",
		Examples: []string{
			"ggc pull current   # Pull current branch from remote",
			"ggc pull rebase    # Pull with rebase",
		},
	})
}

// ShowPushHelp shows help message for push command.
func (h *Helper) ShowPushHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc push <command>",
		Description: "Push changes to remote",
		Examples: []string{
			"ggc push current   # Push current branch to remote",
			"ggc push force     # Force push current branch to remote",
		},
	})
}

// ShowRemoteHelp shows help message for remote command.
func (h *Helper) ShowRemoteHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc remote <command>",
		Description: "Manage set of tracked repositories",
		Examples: []string{
			"ggc remote list              # List remote repositories",
			"ggc remote add name url      # Add a remote repository",
			"ggc remote remove name       # Remove a remote repository",
			"ggc remote set-url name url  # Change remote repository URL",
		},
	})
}

// ShowStashHelp shows help message for stash command.
func (h *Helper) ShowStashHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc stash [command]",
		Description: "Stash changes",
		Examples: []string{
			"ggc stash          # Stash current changes",
			"ggc stash pop      # Apply and remove the latest stash",
			"ggc stash drop     # Remove the latest stash",
		},
	})
}

// ShowStatusHelp shows help message for status command.
func (h *Helper) ShowStatusHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc status [command]",
		Description: "Show the working tree status",
		Examples: []string{
			"ggc status           # Full detailed status output",
			"ggc status short     # Short, concise output (porcelain format)",
		},
	})
}

// ShowRebaseHelp shows help message for rebase command.
func (h *Helper) ShowRebaseHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc rebase",
		Description: "Rebase current branch",
		Examples: []string{
			"ggc rebase   # Rebase current branch",
		},
	})
}

// ShowResetHelp shows help message for reset command.
func (h *Helper) ShowResetHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc reset",
		Description: "Reset and clean",
		Examples: []string{
			"ggc reset   # Reset and clean",
		},
	})
}

// ShowStashPullPopHelp shows help message for stash-pull-pop command.
func (h *Helper) ShowStashPullPopHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc stash-pull-pop",
		Description: "Stash changes, pull from remote, and pop stashed changes",
		Examples: []string{
			"ggc stash-pull-pop   # Stash changes, pull from remote, and pop stashed changes",
		},
	})
}

// ShowResetCleanHelp shows help message for reset-clean command.
func (h *Helper) ShowResetCleanHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc reset-clean",
		Description: "Reset to HEAD and clean untracked files and directories",
		Examples: []string{
			"ggc reset-clean   # Reset to HEAD and clean untracked files and directories",
		},
	})
}

// ShowListBranchesHelp displays help for the list branches command.
func (h *Helper) ShowListBranchesHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc branch list-local | ggc branch list-remote",
		Description: "List local or remote branches",
		Examples: []string{
			"ggc branch list-local    # List local branches",
			"ggc branch list-remote   # List remote branches",
		},
	})
}

// ShowDeleteBranchHelp displays help for the delete branch command.
func (h *Helper) ShowDeleteBranchHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc branch delete <branch-name> [--force]",
		Description: "Delete a branch",
		Examples: []string{
			"ggc branch delete feature/123          # Delete a branch",
			"ggc branch delete feature/123 --force  # Force delete a branch",
		},
	})
}

// ShowDeleteMergedBranchHelp displays help for the delete merged branch command.
func (h *Helper) ShowDeleteMergedBranchHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc branch delete-merged",
		Description: "Delete merged branches",
		Examples: []string{
			"ggc branch delete-merged   # Delete all merged branches",
		},
	})
}

// ShowDiffHelp displays help for the git diff command.
func (h *Helper) ShowDiffHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc diff [options]",
		Description: "Show changes between commits, commit and working tree, etc",
		Examples: []string{
			"ggc diff           # Diff all changes (unstaged and staged)",
			"ggc diff staged    # Diff only staged changes",
			"ggc diff unstaged  # Diff only unstaged changes",
		},
	})
}

// ShowFetchHelp shows help message for fetch command.
func (h *Helper) ShowFetchHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc fetch [options]",
		Description: "Download objects and refs from another repository",
		Examples: []string{
			"ggc fetch --prune   # Fetch and remove any remote-tracking references that no longer exist on the remote",
		},
	})
}
