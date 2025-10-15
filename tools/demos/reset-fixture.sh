#!/usr/bin/env bash

set -euo pipefail

usage() {
	printf 'Usage: %s <scenario>\n' "${0##*/}" >&2
	exit 1
}

SCENARIO=${1:-}
if [[ -z "$SCENARIO" ]]; then
	usage
fi

REPO_ROOT=$(git rev-parse --show-toplevel)
WORKSPACE_DIR="$REPO_ROOT/docs/demos/workspaces/$SCENARIO"
rm -rf "$WORKSPACE_DIR"
mkdir -p "$WORKSPACE_DIR"
cd "$WORKSPACE_DIR"

# Initialize isolated git repository for demo captures
git init -b main >/dev/null
git config user.email "demo@example.com"
git config user.name "GGC Demo"

cat <<'DOC' > README.md
# GGC Demo Repository

This fixture repository is generated automatically for VHS demo recordings.
It contains a small commit history and predictable changes for CLI showcase.
DOC

git add README.md
git commit -m "docs: add demo readme" >/dev/null

mkdir -p docs
cat <<'DOC' > docs/workflow.md
# Demo Workflow

This document tracks the planned git workflow showcased in the GIFs.
DOC

git add docs/workflow.md
git commit -m "docs: add workflow notes" >/dev/null

cat <<'SRC' > app.go
package main

import "fmt"

func main() {
	fmt.Println("Hello from the demo repo")
}
SRC

git add app.go
git commit -m "feat: add sample application" >/dev/null

case "$SCENARIO" in
	cli-workflow)
		cat <<'DOC' > docs/changelog.md
# Changelog

- chore: tidy docs
DOC
		git add docs/changelog.md
		printf '\n## Upcoming release\n- Polish docs\n' >> README.md
		;;
	interactive-overview)
		cat <<'DOC' > docs/checklist.md
# Release Checklist

- [ ] Update README
- [ ] Run tests
- [ ] Tag release
DOC
		git add docs/checklist.md
		printf '\nChecklist captured for interactive demo.\n' >> README.md
		;;
	branch-management)
		git checkout -b feature/onboarding >/dev/null
		cat <<'DOC' > docs/branch-playbook.md
# Branch Playbook

- Keep main stable
- Cut feature branches per task
DOC
		git add docs/branch-playbook.md
		git commit -m "docs: add branch playbook" >/dev/null
		git checkout main >/dev/null
		printf '\n## Branch planning\n- Prepare new docs branch\n' >> README.md
		;;
	stash-cycle)
		printf '\n## Draft release notes\n- [ ] Update changelog\n' >> README.md
		echo "// TODO: refactor output handling" >> app.go
		printf '\nPending checklist update during stash demo.\n' >> docs/workflow.md
		;;
	*)
		echo "Unknown scenario: $SCENARIO" >&2
		exit 1
		;;
esac

# Leave tracked files with unstaged modifications for demos
if command -v perl >/dev/null 2>&1; then
	perl -0pi -e 's/demo/Demo/g' README.md
	perl -0pi -e 's/demo/Demo/g' docs/workflow.md
	perl -0pi -e 's/demo/Demo/g' docs/changelog.md 2>/dev/null || true
	perl -0pi -e 's/demo/Demo/g' docs/checklist.md 2>/dev/null || true
	perl -0pi -e 's/demo/Demo/g' docs/branch-playbook.md 2>/dev/null || true
else
	sed 's/demo/Demo/g' README.md > README.tmp && mv README.tmp README.md
	sed 's/demo/Demo/g' docs/workflow.md > workflow.tmp && mv workflow.tmp docs/workflow.md
	if [[ -f docs/changelog.md ]]; then
		sed 's/demo/Demo/g' docs/changelog.md > changelog.tmp && mv changelog.tmp docs/changelog.md
	fi
	if [[ -f docs/checklist.md ]]; then
		sed 's/demo/Demo/g' docs/checklist.md > checklist.tmp && mv checklist.tmp docs/checklist.md
	fi
	if [[ -f docs/branch-playbook.md ]]; then
		sed 's/demo/Demo/g' docs/branch-playbook.md > branch.tmp && mv branch.tmp docs/branch-playbook.md
	fi
fi

printf '\nfunc Version() string { return "0.1.0" }\n' >> app.go
printf 'Fixture prepared at %s\n' "$WORKSPACE_DIR"
printf 'Staged files:\n'
git diff --cached --name-only
printf '\nWorking tree changes:\n'
git status --short
