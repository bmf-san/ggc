# gcl

## Overview

`gcl` is a CLI tool written in Go to streamline Git operations. It aims to be a maintainable and extensible alternative to shell scripts and aliases, using only the Go standard library with minimal dependencies.

## Features
- Simple commands for common Git operations (add, push, pull, branch, log, etc.)
- Composite commands that combine multiple Git operations
- Interactive UI for branch/file selection and message input
- Implemented using only the Go standard library

## Supported Environments
- OS: macOS (Apple Silicon/Intel), Linux, WSL2 (Windows Subsystem for Linux)
- Go version: 1.21 or later recommended
- Dependencies: Go standard library only (no extra packages required)
- Requirement: `git` command must be installed

## Installation

### Build with make

```sh
git clone <repository URL>
make build
```

Place the `gcl` binary in a directory included in your PATH.

### Global install with go install

```sh
go install github.com/bmf-san/gcl@latest
```

- The `gcl` binary will be installed to `$GOBIN` (usually `$HOME/go/bin`).
- If `$GOBIN` is in your `PATH`, you can use `gcl` from anywhere.
- If not, add it to your `PATH`:

```sh
export PATH=$PATH:$(go env GOBIN)
# or
export PATH=$PATH:$HOME/go/bin
```

## Usage

```sh
gcl <command> [subcommand] [options]
```

### Main Command Examples

|     gcl Command Example      |       Actual git Command       |              Description               |
| --------------------------- | ------------------------------ | -------------------------------------- |
| gcl add <file>              | git add <file>                 | Stage file(s)                          |
| gcl branch current          | git rev-parse --abbrev-ref HEAD| Show current branch name               |
| gcl branch checkout         | git branch ... → git checkout <selected> | Interactive branch switch     |
| gcl branch checkout-remote  | git branch -r ... → git checkout -b ... --track ... | Create and checkout new local branch from remote |
| gcl push current            | git push origin <branch>        | Push current branch                    |
| gcl push force              | git push --force origin <branch>| Force push current branch              |
| gcl pull current            | git pull origin <branch>        | Pull current branch                    |
| gcl pull rebase             | git pull --rebase origin <branch>| Pull with rebase                      |
| gcl log simple              | git log --oneline               | Show simple log                        |
| gcl log graph               | git log --graph                 | Show log with graph                    |
| gcl commit allow-empty      | git commit --allow-empty -m ... | Create empty commit                    |
| gcl commit tmp              | git commit -m "tmp"             | Temporary commit                       |
| gcl fetch --prune           | git fetch --prune               | Fetch with prune                       |
| gcl clean files             | git clean -f                    | Clean files                            |
| gcl clean dirs              | git clean -d                    | Clean directories                      |
| gcl reset clean             | git reset --hard HEAD; git clean -fd | Reset and clean                   |
| gcl commit-push             | git add ... → git commit ... → git push | Interactive add/commit/push      |
| gcl clean interactive       | git clean -nd → git clean -f -- <selected> | Interactive file selection and clean |
| gcl stash trash             | git add . → git stash           | Add all changes and stash              |
| gcl rebase interactive      | git log ... → git rebase -i HEAD~N | Interactive rebase up to HEAD~N   |
| gcl branch delete           | git branch ... → git branch -d <selected> | Interactive delete local branches |
| gcl branch delete-merged    | git branch --merged ... → git branch -d <selected> | Interactive delete merged local branches |
| gcl remote list             | git remote -v                   | Show remotes                           |
| gcl remote add <name> <url> | git remote add <name> <url>     | Add remote                             |
| gcl remote remove <name>    | git remote remove <name>        | Remove remote                          |
| gcl remote set-url <name> <url> | git remote set-url <name> <url> | Change remote URL                  |
| gcl add-commit-push         | git add . → git commit ... → git push | Add, commit, and push all at once |
| gcl pull-rebase-push        | git pull → git rebase origin/main → git push | Pull, rebase, and push all at once |
| gcl stash-pull-pop          | git stash → git pull → git stash pop | Stash, pull, and pop all at once  |
| gcl reset-clean             | git reset --hard HEAD → git clean -fd | Reset and clean all at once        |

## Directory Structure

```
main.go                  # Entry point
router/                  # Command routing logic
cmd/                     # Command entry handlers
git/                     # Git operation wrappers
```

## Completion Script

A bash completion script is available at `tools/completions/gcl.bash`.

### How to Enable (bash/zsh)

```sh
# For bash
source /path/to/gcl/tools/completions/gcl.bash
# For zsh, you can also use source
```

- Add the above to your `.bashrc` or `.zshrc` to enable completion automatically on terminal startup.
- Subcommand completion is supported.

## Future Plans
- Custom configuration via `.gclconfig`
- Mock implementation for testing
- More composite commands and interactive UI

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.