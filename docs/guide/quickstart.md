# Quick start

This is a 10-minute tour of the commands you'll reach for every day. Each example assumes you're already inside a Git working tree.

## 1. See what's going on

```bash
ggc status
```

Like `git status` but grouped and colored.

```bash
ggc log
ggc diff            # staged + unstaged
ggc diff --staged   # staged only
```

## 2. Stage and commit

```bash
ggc add .               # stage everything
ggc add interactive     # pick hunks interactively
ggc add patch           # patch-mode staging
ggc commit "fix: off-by-one in parser"   # no -m required
ggc commit amend        # amend the last commit
ggc commit amend no-edit
```

## 3. Switch branches

```bash
ggc branch current                # show current branch
ggc branch list local             # list local branches
ggc branch checkout               # list + prompt for a branch
ggc branch create feature/x       # create and switch
ggc branch delete feature/old     # delete a specific branch
ggc branch delete merged          # clean up merged branches
```

## 4. Save work in progress

```bash
ggc stash                 # stash current changes
ggc stash list            # list all stashes
ggc stash show            # show changes in the latest stash
ggc stash pop             # reapply the most recent stash
ggc stash drop            # drop the most recent stash
```

## 5. Rebase

```bash
ggc rebase interactive    # interactive rebase
ggc rebase autosquash     # interactive rebase with --autosquash
ggc rebase main           # rebase current branch onto main
ggc rebase continue       # resume after fixing conflicts
ggc rebase abort          # give up and restore pre-rebase state
```

## 6. Push / pull

```bash
ggc pull
ggc push
ggc push force            # force-with-lease
```

## 7. Tag a release

```bash
ggc tag list                            # list tags
ggc tag create v1.2.0                   # create a lightweight tag
ggc tag annotated v1.2.0 "Release"      # create an annotated tag
ggc tag push                            # push tags to origin
```

## 8. Try interactive mode

Run `ggc` by itself to drop into the fuzzy-search prompt:

```bash
ggc
```

Type a few letters, hit <kbd>Enter</kbd>, and you're off. See [Interactive mode](interactive.md) for the full key list and Workflow mode.

## Where to next?

- Full command reference: [Commands](commands.md)
- Interactive pickers, Workflow mode, keybindings: [Interactive mode](interactive.md)
- Aliases and defaults: [Configuration & aliases](config.md)
