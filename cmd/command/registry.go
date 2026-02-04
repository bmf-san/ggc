// Package command provides centralized command registry and metadata management.
package command

import (
	"fmt"
	"strings"
)

// Registry manages command metadata.
type Registry struct {
	commands []Info
}

// NewRegistry creates a new Registry with default commands.
func NewRegistry() *Registry {
	return &Registry{
		commands: defaultCommands(),
	}
}

// NewRegistryWith creates a Registry with custom commands (for testing).
func NewRegistryWith(commands []Info) *Registry {
	return &Registry{commands: commands}
}

// All returns a defensive copy of all commands.
func (r *Registry) All() []Info {
	out := make([]Info, len(r.commands))
	for i := range r.commands {
		out[i] = (&r.commands[i]).clone()
	}
	return out
}

// Find returns the command metadata by name.
func (r *Registry) Find(name string) (Info, bool) {
	for i := range r.commands {
		if strings.EqualFold(r.commands[i].Name, name) {
			return (&r.commands[i]).clone(), true
		}
	}
	return Info{}, false
}

// VisibleCommands returns non-hidden commands.
func (r *Registry) VisibleCommands() []Info {
	var out []Info
	for i := range r.commands {
		if r.commands[i].Hidden {
			continue
		}
		out = append(out, (&r.commands[i]).clone())
	}
	return out
}

// Validate ensures registry consistency.
func (r *Registry) Validate() error {
	return Validate(r.commands)
}

// defaultCommands returns the default command set.
func defaultCommands() []Info {
	return []Info{
		{
			Name:      "help",
			Category:  CategoryBasics,
			Summary:   "Show help information for commands",
			Usage:     []string{"ggc help", "ggc help <command>"},
			Examples:  []string{"ggc help", "ggc help branch"},
			HandlerID: "help",
			Subcommands: []SubcommandInfo{
				{
					Name:    "help",
					Summary: "Show main help message",
					Usage:   []string{"ggc help"},
				},
				{
					Name:    "help <command>",
					Summary: "Show help for a specific command",
					Usage:   []string{"ggc help branch"},
				},
			},
		},
		{
			Name:     "add",
			Category: CategoryBasics,
			Summary:  "Stage changes for the next commit",
			Usage:    []string{"ggc add <file>", "ggc add .", "ggc add interactive", "ggc add patch"},
			Examples: []string{
				"ggc add file.txt   # Add a specific file",
				"ggc add .          # Add all changes to index",
				"ggc add interactive  # Add changes interactively",
				"ggc add patch        # Add changes interactively (patch mode)",
			},
			HandlerID: "add",
			Subcommands: []SubcommandInfo{
				{
					Name:    "add <file>",
					Summary: "Add a specific file to the index",
					Usage:   []string{"ggc add README.md"},
				},
				{
					Name:    "add .",
					Summary: "Add all changes to the index",
					Usage:   []string{"ggc add ."},
				},
				{
					Name:    "add interactive",
					Summary: "Add changes interactively",
					Usage:   []string{"ggc add interactive"},
				},
				{
					Name:    "add patch",
					Summary: "Add changes interactively (patch mode)",
					Usage:   []string{"ggc add patch"},
				},
			},
		},
		{
			Name:     "branch",
			Category: CategoryBranch,
			Summary:  "List, create, and manage branches",
			Usage:    []string{"ggc branch <subcommand>"},
			Examples: []string{
				"ggc branch current                # Show current branch",
				"ggc branch checkout               # Switch to an existing branch",
				"ggc branch checkout remote        # Create and checkout a local branch from the remote",
				"ggc branch create feature/login   # Create and checkout new branch",
				"ggc branch delete feature/login   # Delete local branch",
				"ggc branch delete merged          # Delete local merged branch",
				"ggc branch rename old new         # Rename a branch",
				"ggc branch move feature abc123    # Move branch to specified commit",
				"ggc branch set upstream feature origin/feature  # Set upstream branch",
				"ggc branch info feature           # Show detailed branch information",
				"ggc branch list verbose           # Show detailed branch listing",
				"ggc branch sort date              # List branches sorted by date",
				"ggc branch contains abc123        # Show branches containing a commit",
			},
			HandlerID: "branch",
			Subcommands: []SubcommandInfo{
				{Name: "branch current", Summary: "Show current branch name", Usage: []string{"ggc branch current"}},
				{Name: "branch checkout", Summary: "Switch to an existing branch", Usage: []string{"ggc branch checkout"}},
				{Name: "branch checkout remote", Summary: "Create and checkout a local branch from the remote", Usage: []string{"ggc branch checkout remote"}},
				{Name: "branch create", Summary: "Create and checkout a new branch", Usage: []string{"ggc branch create feature/login"}},
				{Name: "branch delete", Summary: "Delete local branch", Usage: []string{"ggc branch delete feature/login"}, Examples: []string{
					"ggc branch delete feature/123          # Delete a branch",
					"ggc branch delete feature/123 --force  # Force delete a branch",
				}},
				{Name: "branch delete merged", Summary: "Delete local merged branch", Usage: []string{"ggc branch delete merged"}},
				{Name: "branch rename <old> <new>", Summary: "Rename a branch", Usage: []string{"ggc branch rename old new"}},
				{Name: "branch move <branch> <commit>", Summary: "Move branch to specified commit", Usage: []string{"ggc branch move feature abc123"}},
				{Name: "branch set upstream <branch> <upstream>", Summary: "Set upstream for a branch", Usage: []string{"ggc branch set upstream feature origin/feature"}},
				{Name: "branch info <branch>", Summary: "Show detailed branch information", Usage: []string{"ggc branch info feature"}},
				{Name: "branch list verbose", Summary: "Show detailed branch listing", Usage: []string{"ggc branch list verbose"}},
				{Name: "branch list local", Summary: "List local branches", Usage: []string{"ggc branch list local"}},
				{Name: "branch list remote", Summary: "List remote branches", Usage: []string{"ggc branch list remote"}},
				{Name: "branch sort [date|name]", Summary: "List branches sorted by date or name", Usage: []string{"ggc branch sort date"}},
				{Name: "branch contains <commit>", Summary: "Show branches containing a commit", Usage: []string{"ggc branch contains abc123"}},
			},
		},
		{
			Name:     "push",
			Category: CategoryRemote,
			Summary:  "Update remote branches",
			Usage:    []string{"ggc push current", "ggc push force"},
			Examples: []string{
				"ggc push current  # Push current branch to remote",
				"ggc push force    # Force push current branch",
			},
			HandlerID: "push",
			Subcommands: []SubcommandInfo{
				{Name: "push current", Summary: "Push current branch to remote repository", Usage: []string{"ggc push current"}},
				{Name: "push force", Summary: "Force push current branch", Usage: []string{"ggc push force"}},
			},
		},
		{
			Name:     "pull",
			Category: CategoryRemote,
			Summary:  "Fetch and integrate from the remote",
			Usage:    []string{"ggc pull current", "ggc pull rebase"},
			Examples: []string{
				"ggc pull current  # Pull current branch from remote",
				"ggc pull rebase   # Pull with rebase",
			},
			HandlerID: "pull",
			Subcommands: []SubcommandInfo{
				{Name: "pull current", Summary: "Pull current branch from remote repository", Usage: []string{"ggc pull current"}},
				{Name: "pull rebase", Summary: "Pull and rebase", Usage: []string{"ggc pull rebase"}},
			},
		},
		{
			Name:     "log",
			Category: CategoryCommit,
			Summary:  "Inspect commit history",
			Usage:    []string{"ggc log simple", "ggc log graph"},
			Examples: []string{
				"ggc log simple  # Show commit logs in a simple format",
				"ggc log graph   # Show commit logs with a graph",
			},
			HandlerID: "log",
			Subcommands: []SubcommandInfo{
				{Name: "log simple", Summary: "Show simple historical log", Usage: []string{"ggc log simple"}},
				{Name: "log graph", Summary: "Show log with graph", Usage: []string{"ggc log graph"}},
			},
		},
		{
			Name:     "commit",
			Category: CategoryCommit,
			Summary:  "Create commits from staged changes",
			Usage:    []string{"ggc commit <message>", "ggc commit amend", "ggc commit allow empty"},
			Examples: []string{
				"ggc commit \"Update docs\"        # Create commit with a message",
				"ggc commit allow empty            # Create an empty commit",
				"ggc commit amend                  # Amend previous commit (editor)",
				"ggc commit amend no-edit          # Amend without editing commit message",
			},
			HandlerID: "commit",
			Subcommands: []SubcommandInfo{
				{Name: "commit <message>", Summary: "Create commit with a message", Usage: []string{"ggc commit \"Add feature\""}},
				{Name: "commit allow empty", Summary: "Create an empty commit", Usage: []string{"ggc commit allow empty"}},
				{Name: "commit amend", Summary: "Amend previous commit (editor)", Usage: []string{"ggc commit amend"}},
				{Name: "commit amend no-edit", Summary: "Amend without editing commit message", Usage: []string{"ggc commit amend no-edit"}},
			},
		},
		{
			Name:     "fetch",
			Category: CategoryRemote,
			Summary:  "Download objects and refs from remotes",
			Usage:    []string{"ggc fetch", "ggc fetch prune"},
			Examples: []string{
				"ggc fetch prune   # Fetch and remove stale remote-tracking references",
			},
			HandlerID: "fetch",
			Subcommands: []SubcommandInfo{
				{Name: "fetch", Summary: "Fetch from the remote", Usage: []string{"ggc fetch"}},
				{Name: "fetch prune", Summary: "Fetch and clean stale references", Usage: []string{"ggc fetch prune"}},
			},
		},
		{
			Name:     "tag",
			Category: CategoryTag,
			Summary:  "Create, list, and manage tags",
			Usage:    []string{"ggc tag list", "ggc tag annotated <tag> <message>", "ggc tag delete <tag>", "ggc tag show <tag>", "ggc tag push [<remote> <tag>]", "ggc tag create <tag>"},
			Examples: []string{
				"ggc tag                                   # List all tags",
				"ggc tag list                              # List all tags (sorted)",
				"ggc tag list v1.*                         # List tags matching pattern",
				"ggc tag create v1.0.0                     # Create tag",
				"ggc tag create v1.0.0 abc123              # Tag specific commit",
				"ggc tag annotated v1.0.0 'Release notes'  # Create annotated tag",
				"ggc tag delete v1.0.0                     # Delete tag",
				"ggc tag push                              # Push all tags to origin",
				"ggc tag push origin v1.0.0                # Push specific tag (remote first)",
				"ggc tag show v1.0.0                       # Show tag information",
			},
			HandlerID: "tag",
			Subcommands: []SubcommandInfo{
				{Name: "tag list", Summary: "List all tags", Usage: []string{"ggc tag list"}},
				{Name: "tag annotated <tag> <message>", Summary: "Create annotated tag", Usage: []string{"ggc tag annotated v1.0.0 \"Release\""}},
				{Name: "tag delete <tag>", Summary: "Delete tag", Usage: []string{"ggc tag delete v1.0.0"}},
				{Name: "tag show <tag>", Summary: "Show tag information", Usage: []string{"ggc tag show v1.0.0"}},
				{Name: "tag push", Summary: "Push tags to remote", Usage: []string{"ggc tag push", "ggc tag push <remote> <tag>"}},
				{Name: "tag create <tag>", Summary: "Create tag", Usage: []string{"ggc tag create v1.0.1"}},
			},
		},
		{
			Name:     "config",
			Category: CategoryConfig,
			Summary:  "Get and set ggc configuration",
			Usage:    []string{"ggc config list", "ggc config get <key>", "ggc config set <key> <value>"},
			Examples: []string{
				"ggc config list                  # List all configuration values",
				"ggc config get <key>             # Get a config value by key path (e.g., 'ui.color')",
				"ggc config set <key> <value>     # Set a config value by key path",
			},
			HandlerID: "config",
			Subcommands: []SubcommandInfo{
				{Name: "config list", Summary: "List all configuration", Usage: []string{"ggc config list"}},
				{Name: "config get <key>", Summary: "Get a specific config value", Usage: []string{"ggc config get core.editor"}},
				{Name: "config set <key> <value>", Summary: "Set a configuration value", Usage: []string{"ggc config set core.editor vim"}},
			},
		},
		{
			Name:     "hook",
			Category: CategoryHook,
			Summary:  "Manage Git hooks",
			Usage:    []string{"ggc hook <subcommand>"},
			Examples: []string{
				"ggc hook list                    # List all hooks and their status",
				"ggc hook install <hook>          # Install a hook",
				"ggc hook enable <hook>           # Make a hook executable",
				"ggc hook disable <hook>          # Make a hook non-executable",
				"ggc hook uninstall <hook>        # Remove a hook",
				"ggc hook edit <hook>             # Edit a hook",
			},
			HandlerID: "hook",
			Subcommands: []SubcommandInfo{
				{Name: "hook list", Summary: "List all hooks", Usage: []string{"ggc hook list"}},
				{Name: "hook install <hook>", Summary: "Install a hook", Usage: []string{"ggc hook install pre-commit"}},
				{Name: "hook enable <hook>", Summary: "Enable a hook", Usage: []string{"ggc hook enable pre-commit"}},
				{Name: "hook disable <hook>", Summary: "Disable a hook", Usage: []string{"ggc hook disable pre-commit"}},
				{Name: "hook uninstall <hook>", Summary: "Uninstall an existing hook", Usage: []string{"ggc hook uninstall pre-commit"}},
				{Name: "hook edit <hook>", Summary: "Edit a hook's contents", Usage: []string{"ggc hook edit pre-commit"}},
			},
		},
		{
			Name:     "diff",
			Category: CategoryDiff,
			Summary:  "Inspect changes between commits, the index, and the working tree",
			Usage: []string{
				"ggc diff [staged|unstaged|head] [--stat|--name-only|--name-status] [<commit>|<commit1> <commit2>] [--] [<path>...]",
			},
			Examples: []string{
				"ggc diff --stat                     # Show staged + unstaged changes with summary",
				"ggc diff staged cmd/diff.go         # Diff staged changes for a file",
				"ggc diff abc123 def456              # Compare two commits",
				"ggc diff abc123 cmd/diff.go         # Compare commit to working tree for a path",
				"ggc diff -- cmd/deleted_file.go     # Diff a path using -- for disambiguation",
			},
			HandlerID: "diff",
			Subcommands: []SubcommandInfo{
				{Name: "diff", Summary: "Show changes (git diff HEAD)", Usage: []string{"ggc diff"}},
				{Name: "diff unstaged", Summary: "Show unstaged changes", Usage: []string{"ggc diff unstaged"}},
				{Name: "diff staged", Summary: "Show staged changes", Usage: []string{"ggc diff staged"}},
				{Name: "diff head", Summary: "Alias for default diff against HEAD", Usage: []string{"ggc diff head"}},
			},
		},
		{
			Name:     "version",
			Category: CategoryUtility,
			Summary:  "Display current ggc version",
			Usage:    []string{"ggc version"},
			Examples: []string{
				"ggc version   # Shows build time, latest commit and version number",
			},
			HandlerID: "version",
		},
		{
			Name:     "debug-keys",
			Category: CategoryUtility,
			Summary:  "Debug keybinding issues and capture raw key sequences",
			Usage: []string{
				"ggc debug-keys",
				"ggc debug-keys raw",
				"ggc debug-keys raw <file>",
			},
			Examples: []string{
				"ggc debug-keys                 # Show active keybindings",
				"ggc debug-keys raw             # Capture key sequences interactively",
				"ggc debug-keys raw keys.txt    # Capture and save to keys.txt",
			},
			HandlerID: "debug-keys",
			Subcommands: []SubcommandInfo{
				{
					Name:    "debug-keys",
					Summary: "Show current keybindings",
					Usage:   []string{"ggc debug-keys"},
				},
				{
					Name:    "debug-keys raw",
					Summary: "Capture key sequences interactively",
					Usage:   []string{"ggc debug-keys raw"},
				},
				{
					Name:    "debug-keys raw <file>",
					Summary: "Capture key sequences and save them to a file",
					Usage:   []string{"ggc debug-keys raw keys.txt"},
				},
			},
		},
		{
			Name:     "clean",
			Category: CategoryCleanup,
			Summary:  "Remove untracked files and directories",
			Usage:    []string{"ggc clean files", "ggc clean dirs", "ggc clean interactive"},
			Examples: []string{
				"ggc clean files       # Clean untracked files",
				"ggc clean dirs        # Clean untracked directories",
				"ggc clean interactive # Clean files interactively",
			},
			HandlerID: "clean",
			Subcommands: []SubcommandInfo{
				{Name: "clean files", Summary: "Clean untracked files", Usage: []string{"ggc clean files"}},
				{Name: "clean dirs", Summary: "Clean untracked directories", Usage: []string{"ggc clean dirs"}},
				{Name: "clean interactive", Summary: "Clean files interactively", Usage: []string{"ggc clean interactive"}},
			},
		},
		{
			Name:     "stash",
			Category: CategoryStash,
			Summary:  "Save and reapply work-in-progress changes",
			Usage:    []string{"ggc stash <subcommand>"},
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
				"ggc stash create                       # Create stash and return object name",
				"ggc stash store <object>               # Store stash object",
			},
			HandlerID: "stash",
			Subcommands: []SubcommandInfo{
				{Name: "stash", Summary: "Stash current changes", Usage: []string{"ggc stash"}},
				{Name: "stash list", Summary: "List all stashes", Usage: []string{"ggc stash list"}},
				{Name: "stash show", Summary: "Show changes in stash", Usage: []string{"ggc stash show"}},
				{Name: "stash show <stash>", Summary: "Show changes in specific stash", Usage: []string{"ggc stash show stash@{1}"}},
				{Name: "stash apply", Summary: "Apply stash without removing it", Usage: []string{"ggc stash apply"}},
				{Name: "stash apply <stash>", Summary: "Apply specific stash without removing it", Usage: []string{"ggc stash apply stash@{1}"}},
				{Name: "stash pop", Summary: "Apply and remove the latest stash", Usage: []string{"ggc stash pop"}},
				{Name: "stash pop <stash>", Summary: "Apply and remove specific stash", Usage: []string{"ggc stash pop stash@{1}"}},
				{Name: "stash drop", Summary: "Remove the latest stash", Usage: []string{"ggc stash drop"}},
				{Name: "stash drop <stash>", Summary: "Remove specific stash", Usage: []string{"ggc stash drop stash@{1}"}},
				{Name: "stash branch <branch>", Summary: "Create branch from stash", Usage: []string{"ggc stash branch feature"}},
				{Name: "stash branch <branch> <stash>", Summary: "Create branch from specific stash", Usage: []string{"ggc stash branch feature stash@{1}"}},
				{Name: "stash push", Summary: "Save changes to new stash", Usage: []string{"ggc stash push"}},
				{Name: "stash push -m <message>", Summary: "Save changes to new stash with message", Usage: []string{"ggc stash push -m \"WIP\""}},
				{Name: "stash save <message>", Summary: "Save changes to new stash with message", Usage: []string{"ggc stash save \"WIP\""}},
				{Name: "stash clear", Summary: "Remove all stashes", Usage: []string{"ggc stash clear"}},
				{Name: "stash create", Summary: "Create stash and return object name", Usage: []string{"ggc stash create"}},
				{Name: "stash store <object>", Summary: "Store stash object", Usage: []string{"ggc stash store 1234abcd"}},
			},
		},
		{
			Name:     "status",
			Category: CategoryStatus,
			Summary:  "Show working tree status",
			Usage:    []string{"ggc status", "ggc status short"},
			Examples: []string{
				"ggc status        # Full detailed status output",
				"ggc status short  # Short, concise output (porcelain format)",
			},
			HandlerID: "status",
			Subcommands: []SubcommandInfo{
				{Name: "status", Summary: "Show working tree status", Usage: []string{"ggc status"}},
				{Name: "status short", Summary: "Show concise status (porcelain format)", Usage: []string{"ggc status short"}},
			},
		},
		{
			Name:     "rebase",
			Category: CategoryRebase,
			Summary:  "Reapply commits on top of another base tip",
			Usage:    []string{"ggc rebase <subcommand>"},
			Examples: []string{
				"ggc rebase interactive  # Interactive rebase",
				"ggc rebase main         # Rebase current branch onto 'main'",
				"ggc rebase continue     # Continue an in-progress rebase",
				"ggc rebase abort        # Abort an in-progress rebase",
				"ggc rebase skip         # Skip current patch and continue",
			},
			HandlerID: "rebase",
			Subcommands: []SubcommandInfo{
				{Name: "rebase interactive", Summary: "Interactive rebase", Usage: []string{"ggc rebase interactive"}},
				{Name: "rebase <upstream>", Summary: "Rebase current branch onto <upstream>", Usage: []string{"ggc rebase main"}},
				{Name: "rebase continue", Summary: "Continue an in-progress rebase", Usage: []string{"ggc rebase continue"}},
				{Name: "rebase abort", Summary: "Abort an in-progress rebase", Usage: []string{"ggc rebase abort"}},
				{Name: "rebase skip", Summary: "Skip current patch and continue", Usage: []string{"ggc rebase skip"}},
			},
		},
		{
			Name:      "remote",
			Category:  CategoryRemote,
			Summary:   "Manage remotes",
			Usage:     []string{"ggc remote list", "ggc remote add <name> <url>", "ggc remote remove <name>", "ggc remote set-url <name> <url>"},
			Examples:  []string{"ggc remote list", "ggc remote add origin git@github.com:user/repo.git"},
			HandlerID: "remote",
			Subcommands: []SubcommandInfo{
				{Name: "remote list", Summary: "List all remote repositories", Usage: []string{"ggc remote list"}},
				{Name: "remote add <name> <url>", Summary: "Add remote repository", Usage: []string{"ggc remote add upstream git@github.com:user/repo.git"}},
				{Name: "remote remove <name>", Summary: "Remove remote repository", Usage: []string{"ggc remote remove upstream"}},
				{Name: "remote set-url <name> <url>", Summary: "Change remote URL", Usage: []string{"ggc remote set-url origin git@github.com:user/new.git"}},
			},
		},
		{
			Name:      "restore",
			Category:  CategoryCleanup,
			Summary:   "Restore files in working tree or staging area",
			Usage:     []string{"ggc restore <file>", "ggc restore .", "ggc restore staged <file>", "ggc restore staged .", "ggc restore <commit> <file>"},
			Examples:  []string{"ggc restore staged .", "ggc restore main README.md"},
			HandlerID: "restore",
			Subcommands: []SubcommandInfo{
				{Name: "restore <file>", Summary: "Restore file in working directory from index", Usage: []string{"ggc restore README.md"}},
				{Name: "restore .", Summary: "Restore all files in working directory from index", Usage: []string{"ggc restore ."}},
				{Name: "restore staged <file>", Summary: "Unstage file (restore from HEAD to index)", Usage: []string{"ggc restore staged README.md"}},
				{Name: "restore staged .", Summary: "Unstage all files", Usage: []string{"ggc restore staged ."}},
				{Name: "restore <commit> <file>", Summary: "Restore file from specific commit", Usage: []string{"ggc restore HEAD~1 README.md"}},
			},
		},
		{
			Name:      "quit",
			Category:  CategoryUtility,
			Summary:   "Exit interactive mode",
			Usage:     []string{"quit"},
			Examples:  []string{"quit"},
			HandlerID: "quit",
		},
	}
}

// DefaultRegistry is the singleton registry containing all built-in commands.
var DefaultRegistry = NewRegistry()

// Validate ensures the provided command metadata is internally consistent.
func Validate(commands []Info) error {
	seen := make(map[string]struct{})
	for i := range commands {
		cmd := &commands[i]
		if err := validateCommand(cmd, seen); err != nil {
			return err
		}
	}
	return nil
}

func validateCommand(cmd *Info, seen map[string]struct{}) error {
	if strings.TrimSpace(cmd.Name) == "" {
		return fmt.Errorf("command name cannot be empty")
	}
	key := strings.ToLower(cmd.Name)
	if _, ok := seen[key]; ok {
		return fmt.Errorf("duplicate command name: %s", cmd.Name)
	}
	seen[key] = struct{}{}
	if strings.TrimSpace(cmd.Summary) == "" {
		return fmt.Errorf("command summary missing for %s", cmd.Name)
	}
	if !cmd.Hidden && strings.TrimSpace(cmd.HandlerID) == "" {
		return fmt.Errorf("handler ID missing for %s", cmd.Name)
	}

	return validateSubcommands(cmd)
}

func validateSubcommands(cmd *Info) error {
	subSeen := make(map[string]struct{})
	for _, sub := range cmd.Subcommands {
		if strings.TrimSpace(sub.Name) == "" {
			return fmt.Errorf("subcommand name cannot be empty for %s", cmd.Name)
		}
		subKey := strings.ToLower(sub.Name)
		if _, ok := subSeen[subKey]; ok {
			return fmt.Errorf("duplicate subcommand %s under %s", sub.Name, cmd.Name)
		}
		subSeen[subKey] = struct{}{}
		if strings.TrimSpace(sub.Summary) == "" {
			return fmt.Errorf("subcommand summary missing for %s -> %s", cmd.Name, sub.Name)
		}
	}
	return nil
}
