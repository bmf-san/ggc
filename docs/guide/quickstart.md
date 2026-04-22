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
ggc add             # pick files interactively
ggc commit          # opens your $EDITOR like git commit
ggc commit amend    # amend the last commit
```

One-shot:

```bash
ggc s "fix: off-by-one in parser"
```

is equivalent to `git add -A && git commit -m "fix: off-by-one in parser"`.

## 3. Switch branches

```bash
ggc branch                # fuzzy-pick a local branch
ggc branch new feature/x  # create and switch
ggc branch delete         # fuzzy-pick a branch to delete
```

## 4. Save work in progress

```bash
ggc stash              # stash current changes
ggc stash pop          # reapply the most recent stash
ggc stash list         # browse stashes (fuzzy picker)
```

## 5. Rebase

```bash
ggc rebase i 5            # interactive rebase 5 commits back
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
ggc tag                  # list tags
ggc tag create v1.2.0    # create (and sign, if configured)
ggc tag push v1.2.0      # push a single tag to origin
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
