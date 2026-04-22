# ggc

**ggc** is an interactive Git CLI written in Go. It gives you short, memorable subcommands and a fuzzy-finder UI for the everyday workflows that `git` makes too verbose.

```console
$ ggc
> status       show what's changed
  commit       commit staged changes
  branch       switch, create, delete branches
  rebase       interactive rebase helper
  stash        push/pop/apply/clear the stash
  ...
```

## Why ggc?

- **Type less, do more.** `ggc s` stages + commits in one step, `ggc branch` opens a fuzzy picker over local branches, `ggc rebase i 3` starts an interactive rebase 3 commits back.
- **Unified syntax.** No `-`/`--` flag soup. Every command is a verb followed by plain words.
- **Scripts stay scripts.** `ggc` is a thin layer over `git`: anything you can't express in `ggc` you can always fall back to.
- **Safe by default.** Destructive operations ask for confirmation unless you pass `--yes`.

## Get it

See [Installation](guide/install.md). Quickest path:

```bash
go install github.com/bmf-san/ggc/v8@latest
```

## Next steps

1. [Quick start](guide/quickstart.md) — the 10-minute tour
2. [Commands](guide/commands.md) — reference of every `ggc` command
3. [Interactive mode](guide/interactive.md) — fuzzy finders and keybindings
4. [Configuration & aliases](guide/config.md) — `~/.config/ggc/config.yaml`
5. [Troubleshooting](guide/troubleshooting.md) — `ggc doctor` and common issues
