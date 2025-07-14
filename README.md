# ggc
[![GitHub release](https://img.shields.io/github/release/bmf-san/ggc.svg)](https://github.com/bmf-san/ggc/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/bmf-san/ggc)](https://goreportcard.com/report/github.com/bmf-san/ggc)
[![codecov](https://codecov.io/gh/bmf-san/ggc/branch/main/graph/badge.svg)](https://codecov.io/gh/bmf-san/ggc)
[![GitHub license](https://img.shields.io/github/license/bmf-san/ggc)](https://github.com/bmf-san/ggc/blob/main/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/bmf-san/ggc.svg)](https://pkg.go.dev/github.com/bmf-san/ggc)
[![Sourcegraph](https://sourcegraph.com/github.com/bmf-san/ggc/-/badge.svg)](https://sourcegraph.com/github.com/bmf-san/ggc?badge)


A Go Git CLI.

<img src="./docs/icon.png" alt="ggc" title="ggc" width="100px">

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

### Pre-compiled Binaries (Recommended)

Pre-compiled binaries are available for multiple platforms and architectures. This is the fastest way to get started with `ggc`.

#### Download from Releases
Visit the [Releases](https://github.com/bmf-san/ggc/releases/) page to download the latest binary for your platform:

#### Supported Platforms:

- macOS: `darwin_amd64` (Intel), `darwin_arm64` (Apple Silicon)
- Linux: `linux_amd64`, `linux_arm64`
- Windows: `windows_amd64`

### Quick Install with Script

The easiest way to install `ggc` is using the provided installation script:

```sh
# Download and run the installation script
curl -sSL https://raw.githubusercontent.com/bmf-san/ggc/main/install.sh | bash
```

Or download and run it manually:

```sh
# Download the script
curl -O https://raw.githubusercontent.com/bmf-san/ggc/main/install.sh

# Make it executable
chmod +x install.sh

# Run the script
./install.sh
```

The script will:
- Detect your operating system and architecture
- Download the appropriate binary for your system
- Install using `go install`, manual install fallback
- Verify the installation

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

# Run formatter
make fmt

# Run tests
make test

# Run linter
make lint

# Run tests with coverage
make cover

# Run tests and lint
make test-and-lint

# Build with go build
make build

# Build and run with version info
make run
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

- `config` - Manage ggc configuration file
  - `config list` - List all variables in ggc configuration file
  - `config get <key>` - Get configuration variable from key
  - `config set <key> <value>` - Set configuration variable from key and value

- `hook` - Manage ggc hooks (integrates with `git hooks`)
  - `hook list` -	List all hooks
  - `hook install <hook>` - Adds a new hook script from a sample or blank template
  - `hook enable <hook>` - Makes the hook executable
  - `hook disable <hook>` - Makes the hook non-executable
  - `hook uninstall <hook>` - Removes the hook file
  - `hook edit <hook>` - Opens the hook file in your config editor (`default.editor` in `~/.ggcconfig.yaml`)

- `tag` - List all tags
    - `tag list` - List all tags (sorted)
        - `tag list v1.*` - List tags matching pattern
    - `tag create v1.0.0` - Create tag
        - `tag create v1.0.0 abc123` - Tag specific commit
    - `tag annotated v1.0.0 'Release notes'` - Create annotated tag
    - `tag delete v1.0.0` - Delete tag
    - `tag push` - Push all tags to origin
        - `tag push v1.0.0` - Push specific tag
    - `tag show v1.0.0` - Show tag information

- `stash` - Stash changes
  - `stash` - Stash current changes
  - `stash pop` - Apply and remove the latest stash
  - `stash drop` - Remove the latest stash
  - `stash-pull-pop` - Stash changes, pull from remote, and pop stashed changes

- `status` - Show the working tree status
  - `status short` - Show concise output (porcelain format)

- `version` - Show the currently installed ggc version

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

### Fish
Add the following to your `~/.config/fish/config.fish`:
```fish
if test -f (go env GOPATH)/pkg/mod/github.com/bmf-san/ggc@*/tools/completions/ggc.fish
    source (go env GOPATH)/pkg/mod/github.com/bmf-san/ggc@*/tools/completions/ggc.fish
end
```

This setup will automatically find the completion script regardless of the installed version.

# Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) and [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for details.

# Sponsor

If you’d like to support my work, please consider sponsoring me!

[GitHub Sponsors – bmf-san](https://github.com/sponsors/bmf-san)

Or simply giving ⭐ on GitHub is greatly appreciated—it keeps me motivated to maintain and improve the project! :D

# Stargazers
[![Stargazers repo roster for @bmf-san/ggc](https://reporoster.com/stars/bmf-san/ggc)](https://github.com/bmf-san/ggc/stargazers)

# Forkers
[![Forkers repo roster for @bmf-san/ggc](https://reporoster.com/forks/bmf-san/ggc)](https://github.com/bmf-san/ggc/network/members)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
