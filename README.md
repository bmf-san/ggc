# ggc
[![GitHub release](https://img.shields.io/github/release/bmf-san/ggc.svg)](https://github.com/bmf-san/ggc/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/bmf-san/ggc)](https://goreportcard.com/report/github.com/bmf-san/ggc)
[![codecov](https://codecov.io/gh/bmf-san/ggc/branch/main/graph/badge.svg)](https://codecov.io/gh/bmf-san/ggc)
[![GitHub license](https://img.shields.io/github/license/bmf-san/ggc)](https://github.com/bmf-san/ggc/blob/main/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/bmf-san/ggc.svg)](https://pkg.go.dev/github.com/bmf-san/ggc)
[![Sourcegraph](https://sourcegraph.com/github.com/bmf-san/ggc/-/badge.svg)](https://sourcegraph.com/github.com/bmf-san/ggc?badge)


A Go Git CLI.

<img src="https://storage.googleapis.com/gopherizeme.appspot.com/gophers/22bdcabe630eb8f45ed8c740ea665a8345f1d3f6.png" alt="ggc" title="ggc" width="250px">

This logo was created by [gopherize.me](https://gopherize.me/gopher/d654ddf2b81c2b4123684f93071af0cf559eb0b5).

## Overview

ggc is a Git tool written in Go, offering both traditional CLI commands and an interactive interface with incremental search. You can either run subcommands like ggc add directly, or launch the interactive mode by simply typing ggc. Designed to be fast, user-friendly, and extensible.

## Features

- Traditional command-line interface (CLI): Run ggc <command> [args] to execute specific operations directly.
- Interactive interface: Run ggc with no arguments to launch an incremental search UI for command selection.
- Simple commands for common Git operations (add, push, pull, branch, log, etc.)
- Composite commands that combine multiple Git operations
- Interactive UI for branch/file selection and message input
- Implemented using only the Go standard library (+ golang.org/x/term)

## Supported Environments
- OS: macOS (Apple Silicon/Intel) - Verified
- Go version: 1.24 or later recommended
- Dependencies: Go standard library + golang.org/x/term (no extra packages required)
- Requirement: `git` command must be installed

## Installation

### Build with make

```sh
git clone <repository URL>
make build
```

Place the `ggc` binary in a directory included in your PATH.

### Development Setup

For development, you can use the Makefile to install required tools and dependencies:

```sh
# Install all dependencies and tools
make deps

# Run tests
make test

# Run linter
make lint

# Run tests with coverage
make cover

# Run tests and lint
make test-and-lint
```

The Makefile will automatically install required tools like `golangci-lint` using `go install`.

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
- After command execution, results are displayed and you can press Enter to continue
- After viewing results, you return to the command selection screen for continuous use
- Use "quit" command or ctrl+c to exit interactive mode
- All UI and prompts are in English

### Available Commands

- `add` - Add file contents to the index
  - `add <file>` - Add a specific file
  - `add .` - Add all changes
  - `add -p` - Add changes interactively

- `branch` - List, create, or delete branches
  - `branch current` - Show current branch
  - `branch checkout` - Checkout existing branch
  - `branch checkout-remote` - Checkout remote branch
  - `branch create` - Create and checkout new branch
  - `branch delete` - Delete a branch
  - `branch delete-merged` - Delete merged branches
  - `branch list-local` - List local branches
  - `branch list-remote` - List remote branches

- `clean` - Clean untracked files and directories
  - `clean files` - Clean untracked files
  - `clean dirs` - Clean untracked directories

- `commit` - Commit staged changes
  - `commit amend <message>` - Amend to previous commit
  - `commit amend --no-edit` - Amend without editing commit message
  - `commit allow-empty` - Create empty commit
  - `commit tmp` - Create temporary commit

- `diff` - Show changes between commits, commit and working tree, etc.
  - `diff staged` - Show staged changes
  - `diff unstaged` - Show unstaged changes

- `fetch` - Fetch from remote
  - `fetch --prune` - Fetch and prune remote branches

- `log` - Show commit logs
  - `log simple` - Show commit logs in a simple format
  - `log graph` - Show commit logs with a graph

- `pull` - Pull changes from remote
  - `pull current` - Pull current branch from remote
  - `pull rebase` - Pull with rebase

- `push` - Push changes to remote
  - `push current` - Push current branch to remote
  - `push force` - Force push current branch to remote

- `rebase` - Rebase current branch

- `remote` - Manage set of tracked repositories
  - `remote list` - List remote repositories
  - `remote add <name> <url>` - Add a remote repository
  - `remote remove <name>` - Remove a remote repository
  - `remote set-url <name> <url>` - Change remote repository URL

- `reset` - Reset and clean
  - `reset-clean` - Reset to HEAD and clean untracked files and directories

- `stash` - Stash changes
  - `stash` - Stash current changes
  - `stash pop` - Apply and remove the latest stash
  - `stash drop` - Remove the latest stash
  - `stash-pull-pop` - Stash changes, pull from remote, and pop stashed changes

- `status` - Show the working tree status
  - `status short` - Show concise output (porcelain format)

- `add-commit-push` - Add all changes, commit, and push in one command
- `commit-push-interactive` - Commit and push interactively
- `pull-rebase-push` - Pull with rebase and push in one command

## Directory Structure

```
main.go                  # Entry point
router/                  # Command routing logic
cmd/                     # Command entry handlers
git/                     # Git operation wrappers
```

## Shell Completion

### Bash
Add the following to your `~/.bash_profile` or `~/.bashrc`:
```bash
if [ -f "$(go env GOPATH)/pkg/mod/github.com/bmf-san/ggc@*/tools/completions/ggc.bash" ]; then
  . "$(go env GOPATH)"/pkg/mod/github.com/bmf-san/ggc@*/tools/completions/ggc.bash
fi
```

### Zsh
Add the following to your `~/.zshrc`:
```zsh
if [ -f "$(go env GOPATH)/pkg/mod/github.com/bmf-san/ggc@*/tools/completions/ggc.bash" ]; then
  . "$(go env GOPATH)"/pkg/mod/github.com/bmf-san/ggc@*/tools/completions/ggc.bash
fi
```

This setup will automatically find the completion script regardless of the installed version.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
