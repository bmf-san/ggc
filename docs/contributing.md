# Contributing

We love contributions! See the canonical [`CONTRIBUTING.md`](https://github.com/bmf-san/ggc/blob/main/CONTRIBUTING.md) in the repository for the full guide. In short:

1. Fork and clone the repo.
2. `make test` before you push.
3. Follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages and PR titles — the CI enforces it.
4. Every third-party GitHub Action must be pinned to a full commit SHA with a trailing `# vX.Y.Z` comment, per OSSF Scorecard.

## Building the docs site locally

```bash
pip install mkdocs-material
mkdocs serve
```

Then open <http://127.0.0.1:8000>.
