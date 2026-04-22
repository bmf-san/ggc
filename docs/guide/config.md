# Configuration & aliases

ggc reads configuration from one of:

1. `$XDG_CONFIG_HOME/ggc/config.yaml`
2. `~/.config/ggc/config.yaml` (default on Linux/macOS)
3. `~/.ggcconfig.yaml` (legacy path; still respected for older installs)

The first file that exists wins. If none exists, built-in defaults are used. On Windows the path resolves via `%APPDATA%\ggc\config.yaml`.

## Anatomy

```yaml
meta:
  version: v8.3.0     # auto-maintained by ggc
  commit: abcdef1     # auto-maintained by ggc
  created-at: 2025-01-15_12:34:56

ui:
  color: true

git:
  default-remote: origin
  default-branch: main

aliases:
  ship: status && commit amend --no-edit && push force
  cleanup: branch delete merged

interactive:
  profile: default   # one of: default | emacs | vi | readline
```

`meta.*` is rewritten by ggc on startup; don't edit it by hand.

## Aliases

An alias is a named sequence of `ggc` commands separated by `&&`. Anything you can type in the prompt you can put behind an alias.

```bash
ggc ship
```

runs the three commands above in order, stopping at the first failure.

### Placeholders

Numeric placeholders let an alias accept arguments:

```yaml
aliases:
  commit-msg: "commit -m '{0}'"
  feature:    ["branch checkout {0}", "push {0}"]
```

Then:

```bash
ggc commit-msg "Fix bug"
ggc feature main
```

See the [alias validation grammar](https://github.com/bmf-san/ggc/blob/main/internal/config/alias_validate.go) for the exact rules (nesting, escaping, reserved names).

## Keybindings

### Profiles

Pick a profile in one line:

```yaml
interactive:
  profile: emacs
```

| Profile     | Feel                                      |
|-------------|-------------------------------------------|
| `default`   | Curated ggc defaults (recommended)        |
| `emacs`     | GNU readline-style                        |
| `vi`        | Modal, `hjkl` navigation                  |
| `readline`  | Strict GNU readline compatibility         |

### Overriding individual keys

Under `interactive.keybindings` you can override any binding by its logical name:

```yaml
interactive:
  keybindings:
    move_up: "ctrl+p"
    move_down: "ctrl+n"
    accept: "enter"
    cancel: "esc"
    toggle_workflow: "ctrl+w"
```

### Supported key notations

Three equivalent forms are accepted; pick whichever reads best:

| Notation     | Example               |
|--------------|-----------------------|
| `ctrl+`      | `ctrl+p`              |
| `C-`         | `C-p`                 |
| raw caret    | `^P`                  |

`alt+` / `M-` and `shift+` work the same way. Function keys (`f1`..`f12`), arrow keys (`up`, `down`, `left`, `right`), and named keys (`enter`, `esc`, `tab`, `space`, `backspace`) are all recognized.

### Layered overrides

The config is evaluated in this order, later layers winning:

1. `interactive.profile` baseline
2. `interactive.keybindings` (global)
3. `interactive.<os>` — `darwin` / `linux` / `windows`
4. `interactive.contexts.<ctx>` — e.g. `contexts.picker`, `contexts.search`
5. `interactive.terminals.<term>` — e.g. `terminals.alacritty`, `terminals.iterm2`

Example: use emacs everywhere, but tweak `move_up` on macOS only:

```yaml
interactive:
  profile: emacs
  darwin:
    keybindings:
      move_up: "ctrl+p"
```

### Inspecting the resolved keymap

```bash
ggc config list
ggc config get interactive.keybindings
ggc config get interactive.darwin.keybindings
```

If you're unsure what key your terminal is sending, run:

```bash
ggc debug-keys
```

and press keys — it prints the raw escape sequences.

## Editing

```bash
ggc config edit   # opens the config file in $EDITOR
ggc config path   # prints the resolved path
ggc config list   # print the fully-merged config
```

## tmux

Under tmux, most terminals mangle the modifier prefix unless `xterm-keys` is on. Add to `~/.tmux.conf`:

```tmux
set -g xterm-keys on
```

Then reload tmux. Without this, `ctrl+arrow` / `alt+...` will arrive as bare arrow keys.
