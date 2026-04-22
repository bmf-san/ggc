# Configuration & aliases

ggc reads configuration from one of:

1. `$XDG_CONFIG_HOME/ggc/config.yaml`
2. `~/.config/ggc/config.yaml` (default on Linux/macOS)
3. `~/.ggcconfig.yaml` (legacy path; still respected for older installs)

The first file that exists wins. If none exists, built-in defaults are used.

## Anatomy

```yaml
meta:
  version: v8.3.0     # auto-maintained
  commit: abcdef1     # auto-maintained
  created-at: 2025-01-15_12:34:56

ui:
  color: true

git:
  default-remote: origin
  default-branch: main

aliases:
  ship: status && commit amend --no-edit && push force
  cleanup: branch delete merged
```

## Aliases

An alias is a named sequence of `ggc` commands separated by `&&`. Anything you can type in the prompt you can put behind an alias.

```bash
ggc ship
```

runs the three commands above in order, stopping at the first failure.

See [alias validation](https://github.com/bmf-san/ggc/blob/main/internal/config/alias_validate.go) for the exact grammar.

## Keybindings

The interactive prompt's keybindings are defined in code and not yet user-configurable. See [#issue tracker](https://github.com/bmf-san/ggc/issues) for progress on making them user-override-able.

## Editing

```bash
ggc config edit   # opens the config file in $EDITOR
ggc config path   # prints the resolved path
```
