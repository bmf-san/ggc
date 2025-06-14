# Contribution Guide (CONTRIBUTING.md)

## Introduction

Thank you for your interest in contributing to this repository! Bug reports, feature requests, and pull requests (PRs) are all welcome.

## Issues
- Please use [Issues](./issues) for bug reports, feature requests, or questions.
- Include as much detail as possible: steps to reproduce, expected/actual behavior, environment info (OS, Go version, etc.).

## Pull Requests (PR)
- Do not commit directly to the `main` branch. Always create a new branch and submit a PR.
- Each PR should focus on a single purpose (e.g., separate bug fixes and new features).
- Run `make lint` and `make test` before submitting, and ensure there are no errors.
- Please manually test major commands as well.

## Coding Standards
- Go 1.21 or later, standard library only
- Pass linting (golangci-lint) and static analysis

## Testing
- macOS/Linux/WSL2 are recommended environments
- Use `make build` to build the binary, `make test` to run tests
- Update the completion script (`tools/completions/gcl.bash`) as needed

## Other
- If unsure, please open an Issue for discussion first!
- Documentation and README improvements are also welcome.

---

Thank you for your cooperation!