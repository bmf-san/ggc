# Quick start

This is a 10-minute tour of the commands you'll reach for every day. Each example assumes you're already inside a Git working tree.

## 1. See what's going on

```bash
ggc status
```

Like `git status` but grouped and colored.

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

## 4. Rebase

```bash
ggc rebase i 5            # interactive rebase 5 commits back
ggc rebase continue       # resume after fixing conflicts
```

## 5. Push / pull

```bash
ggc pull
ggc push
ggc push force            # force-with-lease
```

## Where to next?

- Full command reference: [Commands](commands.md)
- Interactive pickers and keybindings: [Interactive mode](interactive.md)
- Aliases and defaults: [Configuration & aliases](config.md)
