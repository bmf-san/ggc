# Mapping Your Terminal

Terminals differ in the escape sequences they send for special keys. Use the built-in debug mode to discover what your terminal emits, then bind those sequences.

## Quick Guide

1) Run: `ggc interactive --debug-keys`
2) Press the key you want to map (Fn key, Ctrl+Arrow, Option+Arrow, etc.)
3) Observe stderr logs like:

```
[esc] 0x5B
[csi] 0x31
[csi] 0x35
[csi] 0x41
```

The bytes after ESC are the “tail”. In this example the tail is `[15A`. Lowercased, you would reference it as `esc:[15a` in config. Typically you’ll see final bytes like `A/B/C/D` or `~`.

Add a binding using the raw tail (without the leading ESC):

```yaml
interactive:
  keybindings:
    move_up: ["esc:[1;5A"]  # example: Ctrl+Right/Left often use 1;5C / 1;5D
```

## Function Keys

You can use convenience tokens `f2`..`f12` that map to common sequences. If your terminal uses different tails, prefer explicit `esc:` entries.

Common tails (no leading ESC):

- F2: `oq`, `[12~`;  F3: `or`, `[13~`;  F4: `os`, `[14~`
- F5: `[15~`; F6: `[17~`; F7: `[18~`; F8: `[19~`
- F9: `[20~`; F10: `[21~`; F11: `[23~`; F12: `[24~`

## macOS Option/Alt Keys

- Many macOS terminals can be configured to send ESC+<letter> for Option keys. In iTerm2, set “Left/Right Option Key” to “Esc+”. The bindings here accept `alt+<letter>` tokens and raw `esc:` sequences.
- Alt+Backspace is enabled by default to delete the previous word; you can disable it by not including it in your config, or explicitly bind other options.

## tmux Notes

- Ensure `set -g xterm-keys on` in your `.tmux.conf` for richer key reporting (Ctrl/Alt modifiers on arrows, etc.).
- tmux can alter or prefix sequences. Use `ggc interactive --debug-keys` inside tmux to capture the tails as actually received by ggc, then bind with `esc:<tail>`.
- If Alt keys are intercepted by the OS/terminal, configure your terminal to send ESC for Option/Alt and ensure tmux passes it through.
