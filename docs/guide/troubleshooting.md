# Troubleshooting

## `ggc doctor`

The fastest way to diagnose a ggc installation:

```bash
ggc doctor
```

It prints a checklist like:

```
[OK  ] Go runtime: go1.25.0 (darwin/arm64)
[OK  ] git binary: /usr/bin/git (git version 2.46.0)
[OK  ] ggc config: /Users/you/.config/ggc/config.yaml loaded
[WARN] bash completions: not installed in a well-known location
[OK  ] zsh completions: /opt/homebrew/share/zsh/site-functions/_ggc
[OK  ] stdin TTY: stdin is a TTY

Everything looks good.
```

- `[OK  ]` — nothing to do
- `[WARN]` — usable but suboptimal (e.g. completions not picked up by your shell)
- `[FAIL]` — ggc cannot work until this is fixed (e.g. `git` not in `$PATH`)

## Verbose error messages

ggc prints a compact error by default. To see the exact git command that produced the failure, set `GGC_VERBOSE=1`:

```bash
GGC_VERBOSE=1 ggc pull
```

## Reporting a bug

Please paste the output of `ggc doctor` and the verbose error into the issue. Without those two the maintainers usually can't reproduce the problem.

## Opening an issue

See the [issue forms](https://github.com/bmf-san/ggc/issues/new/choose).
