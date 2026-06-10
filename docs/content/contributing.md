---
title: "Contributing"
description: "How to contribute to ggc and build the documentation site locally."
slug: "contributing"
categories:
  - contributing
---

We love contributions! See the canonical [`CONTRIBUTING.md`](https://github.com/bmf-san/ggc/blob/main/CONTRIBUTING.md) in the repository for the full guide. In short:

1. Fork and clone the repo.
2. `make test` before you push.
3. Follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages and PR titles — the CI enforces it.
4. Every third-party GitHub Action must be pinned to a full commit SHA with a trailing `# vX.Y.Z` comment, per OSSF Scorecard.

## Building the docs site locally

The documentation site is built with [gohan](https://github.com/bmf-san/gohan) (a static site generator written in Go) and the [sleyt](https://github.com/bmf-san/sleyt) CSS theme. From the `docs/` directory:

```bash
# Install the gohan CLI once
go install github.com/bmf-san/gohan/cmd/gohan@latest

# Install the sleyt CSS dependency
npm install

# Build the CSS and serve with live reload
npm run serve
```

Then open <http://127.0.0.1:1313>.

To produce a one-off static build into `docs/public/`:

```bash
npm run build
```
