package command

// expansion returns command definitions added as part of the
// "expand git command coverage" effort (see issue #428). Most of these
// commands are thin pass-throughs to `git <name>` and rely on the ggc
// registry only for discovery, help, and shell-completion plumbing.
func expansion() []Info {
	return []Info{
		// --- Tier 1 ---
		{
			Name:     "switch",
			Category: CategoryBranch,
			Summary:  "Switch branches",
			Usage:    []string{"ggc switch [<options>] <branch>"},
			Examples: []string{
				"ggc switch main                       # Switch to an existing branch",
				"ggc switch -c feature/login           # Create and switch to a new branch",
				"ggc switch -C feature/login          # Force-create and switch",
				"ggc switch --detach HEAD~3            # Detached checkout",
				"ggc switch -                          # Switch back to the previous branch",
			},
			Subcommands: []SubcommandInfo{
				{Name: "switch <branch>", Summary: "Switch to an existing branch", Usage: []string{"ggc switch main"}},
				{Name: "switch -c <branch>", Summary: "Create and switch to a new branch", Usage: []string{"ggc switch -c feature/login"}},
				{Name: "switch --detach <ref>", Summary: "Detached checkout at a ref", Usage: []string{"ggc switch --detach HEAD~3"}},
			},
		},
		{
			Name:     "checkout",
			Category: CategoryBranch,
			Summary:  "Switch branches or restore working tree files",
			Usage:    []string{"ggc checkout [<options>] [<branch>|<commit>] [--] [<path>...]"},
			Examples: []string{
				"ggc checkout main                     # Switch to an existing branch",
				"ggc checkout -b feature/login         # Create and switch to a new branch",
				"ggc checkout -- path/to/file.go       # Discard working-tree changes to a file",
				"ggc checkout HEAD~1 -- path/file.go   # Restore a file from a specific commit",
			},
		},
		{
			Name:     "merge",
			Category: CategoryBranch,
			Summary:  "Join two or more development histories together",
			Usage:    []string{"ggc merge [<options>] [<commit>...]"},
			Examples: []string{
				"ggc merge feature/login               # Merge a branch into the current branch",
				"ggc merge --no-ff feature/login       # Force a merge commit",
				"ggc merge --squash feature/login      # Squash all commits into the index",
				"ggc merge --abort                     # Abort an in-progress merge",
				"ggc merge --continue                  # Continue an in-progress merge",
			},
		},
		{
			Name:     "cherry-pick",
			Category: CategoryCommit,
			Summary:  "Apply the changes introduced by some existing commits",
			Usage:    []string{"ggc cherry-pick [<options>] <commit>..."},
			Examples: []string{
				"ggc cherry-pick abc1234               # Apply a single commit",
				"ggc cherry-pick -x abc1234            # Apply and append \"(cherry picked from ...)\"",
				"ggc cherry-pick A..B                  # Apply a range of commits",
				"ggc cherry-pick --continue            # Continue after resolving conflicts",
				"ggc cherry-pick --abort               # Abort the in-progress cherry-pick",
			},
		},
		{
			Name:     "revert",
			Category: CategoryCommit,
			Summary:  "Revert some existing commits",
			Usage:    []string{"ggc revert [<options>] <commit>..."},
			Examples: []string{
				"ggc revert HEAD                       # Revert the latest commit",
				"ggc revert --no-edit abc1234          # Revert without editing the message",
				"ggc revert -n abc1234                 # Revert without committing (stage only)",
				"ggc revert --continue                 # Continue after resolving conflicts",
				"ggc revert --abort                    # Abort the in-progress revert",
			},
		},
		{
			Name:     "blame",
			Category: CategoryBasics,
			Summary:  "Show what revision and author last modified each line of a file",
			Usage:    []string{"ggc blame [<options>] <file>"},
			Examples: []string{
				"ggc blame README.md                   # Show line authorship for a file",
				"ggc blame -L 10,20 README.md          # Limit blame to specific lines",
				"ggc blame -C -C README.md             # Detect copy/move across files",
			},
		},
		// --- Tier 2 ---
		{
			Name:     "worktree",
			Category: CategoryBranch,
			Summary:  "Manage multiple working trees",
			Usage:    []string{"ggc worktree <subcommand> [<options>]"},
			Examples: []string{
				"ggc worktree list                     # List linked working trees",
				"ggc worktree add ../wt-feat feature   # Add a new working tree",
				"ggc worktree remove ../wt-feat        # Remove a linked working tree",
				"ggc worktree prune                    # Prune stale worktree metadata",
			},
		},
		{
			Name:     "bisect",
			Category: CategoryUtility,
			Summary:  "Use binary search to find the commit that introduced a bug",
			Usage:    []string{"ggc bisect <subcommand> [<options>]"},
			Examples: []string{
				"ggc bisect start                      # Start a new bisect session",
				"ggc bisect bad                        # Mark current commit as bad",
				"ggc bisect good v1.0.0                # Mark a known-good commit",
				"ggc bisect reset                      # Finish bisecting",
			},
		},
		{
			Name:     "reflog",
			Category: CategoryUtility,
			Summary:  "Manage reflog information (recovery aid)",
			Usage:    []string{"ggc reflog [<subcommand>] [<options>] [<ref>]"},
			Examples: []string{
				"ggc reflog                            # Show HEAD reflog",
				"ggc reflog show main                  # Show reflog for a specific ref",
				"ggc reflog expire --expire=now --all  # Aggressively expire reflog entries",
			},
		},
		{
			Name:     "format-patch",
			Category: CategoryUtility,
			Summary:  "Prepare patches for e-mail submission",
			Usage:    []string{"ggc format-patch [<options>] <commit-range>"},
			Examples: []string{
				"ggc format-patch -1 HEAD              # Produce a patch for the latest commit",
				"ggc format-patch origin/main..HEAD    # Produce patches for a branch",
			},
		},
		{
			Name:     "am",
			Category: CategoryUtility,
			Summary:  "Apply a series of patches from a mailbox",
			Usage:    []string{"ggc am [<options>] [<mailbox>...]"},
			Examples: []string{
				"ggc am 0001-fix-bug.patch             # Apply a single patch",
				"ggc am --continue                     # Continue after resolving conflicts",
				"ggc am --abort                        # Abort the in-progress am",
			},
		},
		{
			Name:     "sparse-checkout",
			Category: CategoryUtility,
			Summary:  "Reduce the working tree to a subset of tracked files",
			Usage:    []string{"ggc sparse-checkout <subcommand> [<options>]"},
			Examples: []string{
				"ggc sparse-checkout init --cone       # Enable sparse-checkout in cone mode",
				"ggc sparse-checkout set src docs      # Limit working tree to these paths",
				"ggc sparse-checkout list              # Show currently checked-out paths",
				"ggc sparse-checkout disable           # Disable sparse-checkout",
			},
		},
		{
			Name:     "mv",
			Category: CategoryBasics,
			Summary:  "Move or rename a file, directory, or symlink",
			Usage:    []string{"ggc mv [<options>] <source>... <destination>"},
			Examples: []string{
				"ggc mv old.go new.go                  # Rename a tracked file",
				"ggc mv -k a.go b.go pkg/              # Skip move when destination is in the way",
			},
		},
		{
			Name:     "rm",
			Category: CategoryBasics,
			Summary:  "Remove files from the working tree and the index",
			Usage:    []string{"ggc rm [<options>] <file>..."},
			Examples: []string{
				"ggc rm old.go                         # Stage removal of a tracked file",
				"ggc rm --cached secret.env            # Stop tracking but keep the file on disk",
				"ggc rm -r build/                      # Remove a directory recursively",
			},
		},
		{
			Name:     "submodule",
			Category: CategoryUtility,
			Summary:  "Initialize, update, or inspect submodules",
			Usage:    []string{"ggc submodule <subcommand> [<options>]"},
			Examples: []string{
				"ggc submodule status                  # Show submodule status",
				"ggc submodule update --init           # Initialize and update submodules",
				"ggc submodule foreach git status      # Run a command in each submodule",
			},
		},
		// --- Tier 3 ---
		{
			Name:     "describe",
			Category: CategoryUtility,
			Summary:  "Give an object a human-readable name based on an available ref",
			Usage:    []string{"ggc describe [<options>] [<commit>]"},
			Examples: []string{
				"ggc describe                          # Describe current HEAD",
				"ggc describe --tags                   # Use any tag, not just annotated ones",
				"ggc describe --always --dirty         # Always emit a string; mark dirty trees",
			},
		},
		{
			Name:     "range-diff",
			Category: CategoryDiff,
			Summary:  "Compare two commit ranges (e.g. before and after a rebase)",
			Usage:    []string{"ggc range-diff <range1> <range2>"},
			Examples: []string{
				"ggc range-diff main..@{u} main..HEAD  # Compare upstream vs. local rewrite",
				"ggc range-diff abc..def 123..456      # Compare two arbitrary ranges",
			},
		},
		{
			Name:     "grep",
			Category: CategoryBasics,
			Summary:  "Print lines matching a pattern in tracked files",
			Usage:    []string{"ggc grep [<options>] <pattern> [<pathspec>...]"},
			Examples: []string{
				"ggc grep TODO                         # Search tracked files for TODO",
				"ggc grep -n -i fixme                  # Case-insensitive with line numbers",
				"ggc grep -e foo -e bar -- cmd         # Match multiple patterns in cmd/",
			},
		},
		{
			Name:     "notes",
			Category: CategoryUtility,
			Summary:  "Add, read, or edit object notes",
			Usage:    []string{"ggc notes <subcommand> [<options>]"},
			Examples: []string{
				"ggc notes add -m \"reviewed\" HEAD     # Attach a note to HEAD",
				"ggc notes show HEAD                   # Show a note",
				"ggc notes list                        # List notes",
			},
		},
		{
			Name:     "archive",
			Category: CategoryUtility,
			Summary:  "Create an archive of files from a named tree",
			Usage:    []string{"ggc archive [<options>] <tree-ish> [<path>...]"},
			Examples: []string{
				"ggc archive -o out.tar.gz HEAD        # Archive current HEAD to a tarball",
				"ggc archive --format=zip -o v1.zip v1 # Archive a tag as a zip",
			},
		},
		{
			Name:     "shortlog",
			Category: CategoryBasics,
			Summary:  "Summarize git log output grouped by committer",
			Usage:    []string{"ggc shortlog [<options>] [<revision-range>]"},
			Examples: []string{
				"ggc shortlog -sn                      # Summary count by author",
				"ggc shortlog v1.0..HEAD               # Limit to a range",
			},
		},
		{
			Name:     "maintenance",
			Category: CategoryUtility,
			Summary:  "Run scheduled background repository optimizations",
			Usage:    []string{"ggc maintenance <subcommand> [<options>]"},
			Examples: []string{
				"ggc maintenance run                   # Run all enabled tasks once",
				"ggc maintenance start                 # Install scheduled maintenance",
				"ggc maintenance stop                  # Remove scheduled maintenance",
			},
		},
		{
			Name:     "gc",
			Category: CategoryUtility,
			Summary:  "Cleanup unnecessary files and optimize the local repository",
			Usage:    []string{"ggc gc [<options>]"},
			Examples: []string{
				"ggc gc                                # Run a normal gc",
				"ggc gc --aggressive --prune=now       # Aggressively repack and prune",
			},
		},
		{
			Name:     "fsck",
			Category: CategoryUtility,
			Summary:  "Verify the connectivity and validity of objects in the repository",
			Usage:    []string{"ggc fsck [<options>]"},
			Examples: []string{
				"ggc fsck                              # Run a basic fsck",
				"ggc fsck --full --strict              # Comprehensive checks",
			},
		},
		{
			Name:     "prune",
			Category: CategoryUtility,
			Summary:  "Prune all unreachable objects from the object database",
			Usage:    []string{"ggc prune [<options>]"},
			Examples: []string{
				"ggc prune                             # Prune unreachable objects",
				"ggc prune --dry-run                   # Report what would be pruned",
			},
		},
	}
}
