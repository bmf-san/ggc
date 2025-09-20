# Interactive Keybindings: Examples

Practical examples for per‑OS overlays, per‑TERM overlays, and chord recipes.

## Profiles and Layers

Resolution order: defaults → profile → platform (OS) → terminal ($TERM) → your config.

## Per‑OS Overlays

Example: macOS (darwin) tweaks navigation and input context.

```yaml
interactive:
  profile: default
  darwin:
    keybindings:
      move_up: "ctrl+p"
      move_down: "ctrl+n"
    contexts:
      input:
        keybindings:
          clear_line: "ctrl+u"
          delete_word: ["ctrl+w", "alt+backspace"]
```

Linux and Windows are analogous under `interactive.linux` and `interactive.windows`.

## Per‑TERM Overlays

Target specific terminals via `$TERM` value (e.g., `xterm-256color`).

```yaml
interactive:
  terminals:
    xterm-256color:
      keybindings:
        # Bind F5 to delete_to_end (common ESC tail is [15~)
        delete_to_end: ["f5"]
        # Or use raw escape tails if your terminal differs
        move_up: ["esc:[1;5A"]  # Ctrl+Arrow Right/Left variants vary by terminal
      contexts:
        input:
          keybindings:
            # Custom input-only mapping
            move_to_beginning: "ctrl+a"
```

## Chord Recipes

- Two-step delete word (Emacs-ish):

```yaml
interactive:
  chords:
    delete_word: ["ctrl+k", "ctrl+d"]
```

- Alternative with Alt keys or Fn key:

```yaml
interactive:
  chords:
    delete_word:
      - ["alt+b", "alt+backspace"]
      - ["f2", "ctrl+w"]
```

Enable a timeout so the first step applies if the second doesn’t arrive in time:

```yaml
interactive:
  chord-timeout-ms: 800
```

## Export/Import

- Show effective bindings: `ggc config keybindings list`
- Export complete: `ggc config keybindings export`
- Export only differences from your profile: `ggc config keybindings export --delta`
- Import from YAML: `ggc config keybindings import mykeys.yaml` (add `--dry-run` to preview)

