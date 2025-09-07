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

### 1. Update Documentation and UI Files
**⚠️ CRITICAL: When adding new commands or changing command names/syntax, ALL of the following files MUST be updated:**

#### Core Documentation:
- **README.md**: Update command table, usage examples, and descriptions
- **cmd/templates/help.go**: Update help message templates

#### Interactive UI:
- **cmd/interactive.go**: Update the `commands` array with new command entries

#### Shell Completion Scripts:
- **tools/completions/ggc.bash**: Add/modify command and subcommand completions
- **tools/completions/ggc.fish**: Add/modify command and subcommand completions
- **tools/completions/ggc.zsh**: Add/modify command and subcommand completions

**📋 Checklist for Command Changes:**
- [ ] README.md command table updated
- [ ] cmd/templates/help.go help templates updated
- [ ] cmd/interactive.go commands array updated
- [ ] tools/completions/ggc.bash completions updated
- [ ] tools/completions/ggc.fish completions updated
- [ ] tools/completions/ggc.zsh completions updated
- [ ] All tests pass (`make test`)
- [ ] No lint errors (`make lint`)

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
