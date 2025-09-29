# Contribution Guide (CONTRIBUTING.md)

## Introduction

Thank you for your interest in contributing to this repository! Bug reports, feature requests, and pull requests (PRs) are all welcome.

## Issues
- Please use [Issues](./issues) for bug reports, feature requests, or questions.
- Include as much detail as possible: steps to reproduce, expected/actual behavior, environment info (OS, Go version, etc.).

## Pull Requests (PR)
- Do not commit directly to the `main` branch. Always create a new branch and submit a PR.
- Each PR should focus on a single purpose (e.g., separate bug fixes and new features).
- Run `make fmt` and `make lint` and `make test` before submitting, and ensure there are no errors.
- Please manually test major commands as well.

## Coding Standards
- Go 1.25 or later, standard library only
- Pass linting (golangci-lint) and static analysis

## Implementation Guidelines
When implementing new features or modifying existing ones, please ensure to:

### 1. Update Registry, Documentation, and UI Files
**When adding new commands or changing command names/syntax, update the centralized registry and any surfaces that depend on it.**

#### Command Registry:
- **cmd/command/registry.go**: Add or modify `CommandInfo` entries (usage, examples, handler IDs, visibility)
  - Set `Hidden: true` for experimental/internal commands you do not want exposed via `help` or interactive search.

#### Documentation:
- **Auto-generated**: Run `make docs` to update the README.md command table from the registry

#### Shell Completion Scripts:
- **Auto-generated**: Run `make docs` (or `make completions`) to regenerate the Bash/Zsh/Fish completion scripts from the registry.
- **Do not edit** files under `tools/completions/` manually—changes will be overwritten by the generator.

**📋 Checklist for Command Changes:**
- [ ] cmd/command/registry.go entry added/updated (usage, examples, handler)
- [ ] Run `make docs` to update README.md and regenerate shell completions
- [ ] All tests pass (`make test`)
- [ ] No lint errors (`make lint`)

## Adding New Commands

ggc uses a centralized command registry system that eliminates the need to update multiple files when adding commands. Here's the streamlined workflow:

### 1. Add to Registry
Edit `cmd/command/registry.go` and add your command entry:

```go
{
    Name:      "mycommand",
    Category:  command.CategoryUtility,
    Summary:   "Does something useful",
    Usage:     []string{"ggc mycommand", "ggc mycommand --help"},
    Examples:  []string{"ggc mycommand", "ggc mycommand file.txt"},
    HandlerID: "mycommand",
    Subcommands: []command.SubcommandInfo{
        {
            Name:    "mycommand subaction",
            Summary: "Performs a sub-action",
            Usage:   []string{"ggc mycommand subaction"},
        },
    },
},
```

### 2. Implement Handler
Add your handler function to the appropriate `cmd/*.go` file and register it in `cmd/cmd.go`:

```go
// In cmd/cmd.go handlers map
"mycommand": func(args []string) { cmd.MyCommand(args) },
```

### 3. Update Documentation
```bash
make docs  # Updates README.md command table automatically
```

### 2. Follow existing code patterns:
   - Place command implementations in appropriate files under `cmd/`
   - Add corresponding test files
   - Use consistent error handling and output formatting

### 3. Consider user experience:
   - Provide clear, helpful error messages
   - Add examples in help text
   - Ensure command behavior is intuitive

## Testing
- macOS/Linux/WSL2 are recommended environments
- Use `make build` to build the binary, `make test` to run tests
- Update tests when adding or modifying features
- Add test cases for error scenarios

## Internal Design: Segmented Git Interfaces

To reduce mock surface area and improve maintainability, the `git` package defines small, focused interfaces that represent cohesive slices of functionality. For example:

```
// git/interfaces.go
type DiffReader interface {
    Diff() (string, error)
    DiffStaged() (string, error)
    DiffHead() (string, error)
    DiffWith(args []string) (string, error)
}
```

Guidelines:
- Commands in `cmd/` should depend on the smallest interface they need (e.g., `git.DiffReader` for `cmd/diff.go`).
- The concrete `git.Client` implements the full set of operations and automatically satisfies these smaller interfaces.
- In tests, prefer defining minimal mocks that satisfy only the required small interface instead of a large, catch‑all client surface.

This approach follows the Interface Segregation Principle and helps avoid updating large mock types when unrelated functionality changes.

### Interface Style Guide

- Prefer narrow, role-based interfaces over a single large one.
- Use clear suffixes to communicate intent:
  - `Reader` for read-only queries (e.g., `DiffReader`, `StatusReader`, `BranchReader`).
  - `Writer` for mutating operations (e.g., `CommitWriter`, `BranchWriter`).
  - Use simple nouns for single-purpose operations (e.g., `Pusher`, `Puller`).
- Compose small interfaces for command needs instead of expanding them:
  - Example composite used by status:

```
type StatusInfoReader interface {
    StatusReader
    BranchUpstreamReader
}
```

- Define interfaces in `git/interfaces.go`; implement behavior in `git/*.go` on `git.Client`.
- Constructors in `cmd/*` should accept the smallest interface required, for example:

```
// cmd/diff.go
type Differ struct { gitClient git.DiffReader /* ... */ }
func NewDiffer(c git.DiffReader) *Differ { /* ... */ }

// cmd/commit.go
type Committer struct { gitClient git.CommitWriter /* ... */ }
func NewCommitter(c git.CommitWriter) *Committer { /* ... */ }
```

- Tests should create minimal mocks that satisfy only the specific interface required by the command being tested.

## Command Design Guidelines

### Naming Conventions
1. Struct Names: Use `-er` suffix consistently (e.g., `Brancher`, `Committer`)
2. Field Names: Match struct names in lowercase (e.g., `brancher`, `committer`)
3. Function Names: Use descriptive verbs (e.g., `GetCurrentBranch`, `ListLocalBranches`)

### Command Structure
1. Format: `ggc <command> <subcommand> [modifier] [arguments]`
2. No Option Flags: Use subcommands instead of `-flag` or `--flag`
3. No Hyphens: Use spaces to separate words (e.g., `clean interactive`)
4. Hierarchical: Group related functionality under common commands

### Examples
- ✅ `ggc restore staged <file>`
- ✅ `ggc commit allow empty`
- ✅ `ggc branch set upstream <branch> <upstream>`
- ❌ `ggc restore --staged <file>`
- ❌ `ggc commit --allow-empty`
- ❌ `ggc branch set-upstream <branch> <upstream>`

### When Adding New Commands
1. Follow the established hierarchy
2. Use descriptive subcommand names
3. Update all related files:
   - Command implementation
   - Help templates
   - Interactive mode commands
   - Completion scripts
   - Documentation

## Other
- If unsure, please open an Issue for discussion first!
- Documentation and README improvements are also welcome.

---

Thank you for your cooperation!
