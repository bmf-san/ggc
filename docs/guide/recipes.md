# Recipes

Task-oriented walkthroughs that combine ggc subcommands into typical Git workflows. Each recipe is a few lines you can copy-paste or build into a workflow (interactive mode → <kbd>Tab</kbd> → <kbd>Ctrl</kbd>+<kbd>T</kbd>).

For the full reference see [Commands](commands.md); for the interactive flow see [Interactive mode](interactive.md).

## Start a feature branch

```bash
ggc fetch prune           # tidy stale remote refs
ggc branch create         # prompts for a name, creates and checks out
# ... do work ...
ggc status                # what's changed
ggc add .
ggc commit "feat: widget supports dark mode"
ggc push current          # first push sets upstream automatically
```

## Amend the last commit before pushing

```bash
ggc add .
ggc commit amend no-edit  # re-use the current message
# or
ggc commit amend          # reopens the editor
```

Reach for `commit amend` only on commits that have **not** been pushed. For published branches, prefer a fixup + autosquash (below).

## Fixup + autosquash

You noticed a typo three commits ago on your feature branch:

```bash
ggc log graph              # find the <commit> you want to patch
# ... make the typo fix ...
ggc add .
ggc commit fixup <commit>
ggc rebase autosquash      # squashes the fixup into the right commit
```

The branch history stays clean; reviewers only see the amended commit.

## Sync a long-running branch with main

```bash
git switch main
ggc pull current
git switch feature/widget
ggc rebase main            # or: ggc rebase interactive for cherry-picking
# If conflicts:
ggc rebase continue        # after resolving
# or abort:
ggc rebase abort
```

## Clean up after a merged PR

```bash
ggc fetch prune
ggc branch delete merged   # removes local branches already merged into the default branch
```

## Unstage / undo

```bash
ggc restore staged <file>   # unstage a single file
ggc restore staged .        # unstage everything
ggc restore <file>          # discard working-tree changes for a file
ggc reset hard <commit>     # rewind HEAD and working tree to <commit>
```

## Stash a WIP and come back to it

```bash
ggc stash push -m "WIP: migration"
# ... work on the other thing ...
ggc stash list
ggc stash pop               # latest stash
# or, to apply a specific stash without removing it:
ggc stash apply <stash>
```

## Tag a release

```bash
ggc tag create v1.2.0
ggc tag push                # push all local tags
# for an annotated tag with a message:
ggc tag annotated v1.2.0 "First stable release"
```

## Inspect before committing

```bash
ggc diff unstaged           # what you'd lose with `restore .`
ggc diff staged             # what will go into the next commit
ggc log simple              # quick scroll-friendly log
ggc log graph               # visual graph across branches
```

## Build your own workflow

Drop into the fuzzy picker (`ggc` with no args), search each subcommand you want, press <kbd>Tab</kbd> to queue it, and <kbd>Ctrl</kbd>+<kbd>T</kbd> to run the full pipeline. Save the sequence as an alias in `~/.ggcconfig.yaml` — see [Configuration & aliases](config.md).
