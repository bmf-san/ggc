# ggc

A Go Git CLI & CUI.

<img src="https://storage.googleapis.com/gopherizeme.appspot.com/gophers/22bdcabe630eb8f45ed8c740ea665a8345f1d3f6.png" alt="ggc" title="ggc" width="250px">

This logo was created by [gopherize.me](https://gopherize.me/gopher/d654ddf2b81c2b4123684f93071af0cf559eb0b5).

## Overview

ggc is a Git tool written in Go, providing both a traditional command-line interface (CLI) and an interactive character user interface (CUI) with incremental search. It is designed to be fast, user-friendly, and extensible. (Go Git CLI & CUI)

## Features
- Traditional command-line interface (CLI): run `ggc <command> [args]` for direct operations
- Interactive character user interface (CUI): run `ggc` with no arguments to launch an incremental search UI for command selection
- Simple commands for common Git operations (add, push, pull, branch, log, etc.)
- Composite commands that combine multiple Git operations
- Interactive UI for branch/file selection and message input
- All prompts and UI are in English
- All prompts and command inputs are always aligned to the left (no terminal right-shift issues)
- Implemented using only the Go standard library (+ golang.org/x/term)

## Supported Environments
- OS: macOS (Apple Silicon/Intel), Linux, WSL2 (Windows Subsystem for Linux)
- Go version: 1.21 or later recommended
- Dependencies: Go standard library + golang.org/x/term (no extra packages required)
- Requirement: `git` command must be installed

## Installation

### Build with make

```sh
git clone <repository URL>
make build
```

Place the `ggc` binary in a directory included in your PATH.

### Global install with go install

```sh
go install github.com/bmf-san/ggc@latest
```

- The `ggc` binary will be installed to `$GOBIN` (usually `$HOME/go/bin`).
- If `$GOBIN` is in your `PATH`, you can use `ggc` from anywhere.
- If not, add it to your `PATH`:

```sh
export PATH=$PATH:$(go env GOBIN)
# or
export PATH=$PATH:$HOME/go/bin
```

## Usage

### Interactive Command Selection (Incremental Search UI)

Just run:

```sh
ggc
```

- Type to filter commands (incremental search)
- Use ctrl+n/ctrl+p to move selection, Enter to execute
- If a command requires arguments (e.g. `<file>`, `<name>`, `<url>`), you will be prompted for input (always left-aligned)
- All UI and prompts are in English

### Main Command Examples

|     ggc Command Example      |       Actual git Command       |              Description               |
| --------------------------- | ------------------------------ | -------------------------------------- |
| ggc add <file>              | git add <file>                 | Stage file(s)                          |
| ggc add-commit-push         | git add . → git commit ... → git push | Add, commit, and push all at once |
| ggc branch current          | git rev-parse --abbrev-ref HEAD| Show current branch name               |
| ggc branch checkout         | git branch ... → git checkout <selected> | Interactive branch switch     |
| ggc branch checkout-remote  | git branch -r ... → git checkout -b <n> --track <remote>/<branch> | Create local branch from remote and switch to it |
| ggc branch delete          | git branch ... → git branch -d <selected> | Interactive delete of non-merged local branches |
| ggc branch delete-merged   | git branch --merged ... → git branch -d <selected> | Interactive delete of already merged local branches |
| ggc clean files             | git clean -f                    | Clean files                            |
| ggc clean dirs              | git clean -d                    | Clean directories                      |
| ggc clean-interactive       | git clean -nd → git clean -f    | Interactive file cleaning              |
| ggc commit allow-empty      | git commit --allow-empty -m ... | Create empty commit                    |
| ggc commit tmp              | git commit -m "tmp"             | Temporary commit                       |
| ggc commit-push-interactive | Interactive add/commit/push     | Select files, commit, and push interactively |
| ggc complete               | bash/zsh completion            | Generate shell completion script        |
| ggc fetch --prune          | git fetch --prune              | Fetch and remove stale remote-tracking branches |
| ggc log simple              | git log --oneline               | Show simple log                        |
| ggc log graph               | git log --graph                 | Show log with graph                    |
| ggc pull current            | git pull origin <branch>        | Pull current branch                    |
| ggc pull rebase             | git pull --rebase origin <branch>| Pull with rebase                      |
| ggc pull-rebase-push        | git pull → git rebase → git push | Pull with rebase and push all at once |
| ggc push current            | git push origin <branch>        | Push current branch                    |
| ggc push force              | git push --force origin <branch>| Force push current branch              |
| ggc rebase                  | git rebase                      | Rebase current branch                  |
| ggc remote list             | git remote -v                   | Show remotes                           |
| ggc remote add <n> <url>    | git remote add <n> <url>       | Add remote                             |
| ggc remote remove <n>       | git remote remove <n>          | Remove remote                          |
| ggc remote set-url <n> <url>| git remote set-url <n> <url>   | Change remote URL                      |
| ggc reset                   | git reset --hard HEAD; git clean -fd | Reset and clean                   |
| ggc stash                   | git stash                       | Stash changes                          |
| ggc stash-pull-pop          | git stash → git pull → git stash pop | Stash, pull, and pop all at once  |

## Directory Structure

```
main.go                  # Entry point
router/                  # Command routing logic
cmd/                     # Command entry handlers
git/                     # Git operation wrappers
```

## Completion Script

A bash completion script is available at `tools/completions/ggc.bash`.

### How to Enable (bash/zsh)

```sh
# For bash
source /path/to/ggc/tools/completions/ggc.bash
# For zsh, you can also use source
```

- Add the above to your `.bashrc` or `.zshrc` to enable completion automatically on terminal startup.
- Subcommand completion is supported.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.