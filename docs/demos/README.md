# Demo Automation

This directory houses automated VHS scripts used to generate GIFs demonstrating key ggc workflows.

## Structure
- `scripts/`: Charmbracelet VHS `.tape` definitions.
- `generated/`: Output GIFs produced by running the scripts.
- `workspaces/`: Temporary git fixtures created during generation (ignored from git).

## Available Scenarios
- `cli-workflow` → `docs/demos/generated/cli-workflow.gif`
- `interactive-overview` → `docs/demos/generated/interactive-overview.gif`
- `branch-management` → `docs/demos/generated/branch-management.gif`
- `stash-cycle` → `docs/demos/generated/stash-cycle.gif`

Run `make demos` (or `make docs`) after installing VHS and its dependencies to refresh all assets.
