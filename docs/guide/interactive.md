# Interactive mode

Running `ggc` with no arguments drops you into the interactive prompt. There are two modes:

- **Search mode** (default) — fuzzy-search over all `ggc` commands and run one.
- **Workflow mode** — build a queue of commands, then execute them in sequence with one keystroke.

## Search mode

From the prompt:

- type a few letters → the command list narrows via fuzzy match
- <kbd>Enter</kbd> — execute the highlighted command
- <kbd>Tab</kbd> — add the highlighted command to the workflow queue and stay in search
- <kbd>↑</kbd>/<kbd>↓</kbd> or <kbd>Ctrl</kbd>+<kbd>P</kbd>/<kbd>Ctrl</kbd>+<kbd>N</kbd> — move selection
- <kbd>Ctrl</kbd>+<kbd>C</kbd> — cancel the current input
- <kbd>Ctrl</kbd>+<kbd>D</kbd> — exit

### Fuzzy pickers

Commands that take a branch, file, or stash entry open a nested picker using the same keys (<kbd>↑</kbd>/<kbd>↓</kbd>, <kbd>Enter</kbd> to accept, <kbd>Esc</kbd> to cancel).

## Workflow mode

Workflow mode turns the interactive prompt into a command pipeline builder. Typical use: stage → commit → push in one go without re-typing anything.

1. In search mode, highlight a command and press <kbd>Tab</kbd>. The command is appended to the workflow queue and the prompt stays in search mode so you can add the next one.
2. Press <kbd>Ctrl</kbd>+<kbd>T</kbd> to switch to workflow view. You see the queued commands in order.
3. In workflow view: <kbd>x</kbd> runs the queue (execution stops on the first failure), <kbd>n</kbd> creates a new workflow, <kbd>d</kbd> / <kbd>Ctrl</kbd>+<kbd>D</kbd> deletes the active workflow, <kbd>Ctrl</kbd>+<kbd>N</kbd>/<kbd>Ctrl</kbd>+<kbd>P</kbd> cycles between workflows.
4. <kbd>Ctrl</kbd>+<kbd>T</kbd> again returns to search mode without clearing the queue; <kbd>c</kbd> clears the active workflow.

Commands with placeholders (e.g. aliases like `commit-msg: "commit -m '{0}'"`) will prompt for the placeholder value when they run, not when they're queued.

## Keybinding profiles

The interactive prompt ships with four profiles:

- `default` — curated ggc defaults
- `emacs` — GNU readline-style
- `vi` — modal, `hjkl` navigation
- `readline` — strict readline compatibility

Pick one in `~/.config/ggc/config.yaml`:

```yaml
interactive:
  profile: emacs
```

Fine-grained overrides (per-OS, per-context, per-terminal, custom key combos) are documented in [Configuration & aliases → Keybindings](config.md#keybindings).

## Exiting

From search mode: <kbd>Ctrl</kbd>+<kbd>D</kbd> or type `quit` + <kbd>Enter</kbd>. `quit` only works inside interactive mode; invoking `ggc quit` from a shell is a no-op.
