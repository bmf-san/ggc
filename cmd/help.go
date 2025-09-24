// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/bmf-san/ggc/v6/cmd/templates"
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
	helpMsg, err := templates.RenderMainHelp()
	if err != nil {
		_, _ = fmt.Fprintf(h.outputWriter, "Error: %v\n", err)
		return
	}
	_, _ = fmt.Fprint(h.outputWriter, helpMsg)
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
		Usage:       "ggc add <file> | ggc add interactive | ggc add patch",
		Description: "Add file contents to the index",
		Examples: []string{
			"ggc add file.txt   # Add a specific file",
			"ggc add interactive  # Add changes interactively",
			"ggc add patch        # Add changes interactively (patch mode)",
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
			"ggc branch delete merged              # Delete merged branches",
			"ggc branch rename <old> <new>         # Rename branch",
			"ggc branch move <branch> <commit>     # Move branch pointer",
			"ggc branch set upstream <branch> <up> # Set upstream branch",
			"ggc branch info <branch>              # Show branch details",
			"ggc branch list verbose              # Detailed branch list",
			"ggc branch sort [date|name]           # Sort branches",
			"ggc branch contains <commit>          # Show branches containing commit",
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
			"ggc clean interactive # Clean interactively",
		},
	})
}

// ShowCommitHelp shows help message for commit command.
func (h *Helper) ShowCommitHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc commit <message> | ggc commit amend [no-edit] [<message>] | ggc commit allow empty",
		Description: "Commit staged changes",
		Examples: []string{
			"ggc commit amend                # Amend previous commit (editor)",
			"ggc commit amend no-edit        # Amend without editing commit message",
			"ggc commit amend Updated title  # Amend with new message",
			"ggc commit allow empty          # Create empty commit",
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
			"ggc stash                              # Stash current changes",
			"ggc stash list                         # List all stashes",
			"ggc stash show [stash]                 # Show changes in stash",
			"ggc stash apply [stash]                # Apply stash without removing it",
			"ggc stash pop [stash]                  # Apply and remove stash",
			"ggc stash drop [stash]                 # Remove stash",
			"ggc stash branch <branch> [stash]      # Create branch from stash",
			"ggc stash push [-m message] [files]    # Save changes to new stash",
			"ggc stash save [message]               # Save changes to new stash",
			"ggc stash clear                        # Remove all stashes",
			"ggc stash create [message]             # Create stash and return object name",
			"ggc stash store <object> [message]     # Store stash object",
		},
	})
}

// ShowHookHelp displays help information for hook commands.
func (h *Helper) ShowHookHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc hook [command]",
		Description: "Manage Git hooks",
		Examples: []string{
			"ggc hook list                    # List all hooks and their status",
			"ggc hook install <hook>          # Install a hook from sample or create template",
			"ggc hook enable <hook>           # Make a hook executable",
			"ggc hook edit <hook>             # Edit a hook in default.editor in config",
			"ggc hook disable <hook>          # Make a hook non-executable",
			"ggc hook uninstall <hook>        # Remove a hook",
		},
	})
}

// ShowConfigHelp shows help message for config command.
func (h *Helper) ShowConfigHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc config [command]",
		Description: "Get, set, and list configuration values for ggc (file located at ~/.ggcconfig.yaml or ~/.config/ggc/config.yaml)",
		Examples: []string{
			"ggc config list                  # List all configuration values",
			"ggc config get <key>             # Get a config value by key path (e.g., 'ui.color', 'default.branch')",
			"ggc config set <key> <value>     # Set a config value by key path (e.g., 'ui.color true', 'default.branch main')",
		},
	})
}

// ShowRestoreHelp shows help message for restore command.
func (h *Helper) ShowRestoreHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc restore [command]",
		Description: "Restore working tree files",
		Examples: []string{
			"restore <file>                # Restore file in working directory from index",
			"restore <commit> <file>       # Restore file from specific commit",
			"restore .                     # Restore all files in working directory from index",
			"restore staged <file>       # Unstage file (restore from HEAD to index)",
			"restore staged .            # Unstage all files",
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

// ShowTagHelp shows help message for tag command.
func (h *Helper) ShowTagHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc tag [command] [options]",
		Description: "Create, list, delete and verify tags",
		Examples: []string{
			"ggc tag                                   # List all tags",
			"ggc tag list                              # List all tags (sorted)",
			"ggc tag list v1.*                         # List tags matching pattern",
			"ggc tag create v1.0.0                     # Create tag",
			"ggc tag create v1.0.0 abc123              # Tag specific commit",
			"ggc tag annotated v1.0.0 'Release notes'  # Create annotated tag",
			"ggc tag delete v1.0.0                     # Delete tag",
			"ggc tag push                              # Push all tags to origin",
			"ggc tag push v1.0.0                       # Push specific tag",
			"ggc tag show v1.0.0                       # Show tag information",
		},
	})
}

// ShowVersionHelp shows help message for Version command.
func (h *Helper) ShowVersionHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc version",
		Description: "Show current ggc version",
		Examples: []string{
			"ggc version           # Shows build time, latest commit and version number",
		},
	})
}

// ShowRebaseHelp shows help message for rebase command.
func (h *Helper) ShowRebaseHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc rebase [interactive | <upstream> | continue | abort | skip]",
		Description: "Rebase current branch onto another branch; supports interactive and common workflows",
		Examples: []string{
			"ggc rebase interactive   # Interactive rebase",
			"ggc rebase main          # Rebase current branch onto 'main'",
			"ggc rebase continue      # Continue after resolving conflicts",
			"ggc rebase abort         # Abort the in-progress rebase",
			"ggc rebase skip          # Skip current patch and continue",
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

// ShowListBranchesHelp displays help for the list branches command.
func (h *Helper) ShowListBranchesHelp() {
	h.ShowCommandHelp(templates.HelpData{
		Usage:       "ggc branch list local | ggc branch list remote",
		Description: "List local or remote branches",
		Examples: []string{
			"ggc branch list local    # List local branches",
			"ggc branch list remote   # List remote branches",
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
		Usage:       "ggc branch delete merged",
		Description: "Delete merged branches",
		Examples: []string{
			"ggc branch delete merged   # Delete all merged branches",
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
		Usage:       "ggc fetch [subcommand]",
		Description: "Download objects and refs from another repository",
		Examples: []string{
			"ggc fetch prune   # Fetch and remove any remote-tracking references that no longer exist on the remote",
		},
	})
}
