# Installation

## Quick install (macOS / Linux)

The fastest way — the install script picks the right binary for your OS/arch, installs it under `/usr/local/bin/`, and verifies it runs:

```bash
curl -sSL https://raw.githubusercontent.com/bmf-san/ggc/main/install.sh | bash
```

## Homebrew (macOS / Linux)

```bash
brew install ggc
# later:
brew upgrade ggc
```

Formula: <https://formulae.brew.sh/formula/ggc>.

## `go install`

Requires Go 1.25 or newer:

```bash
go install github.com/bmf-san/ggc/v8@latest
```

The binary lands in `$GOBIN` (usually `$HOME/go/bin`). Make sure it's on your `PATH`.

!!! warning "No version metadata with `go install`"
    Release notes, commit hash and build date are baked in via ldflags during `make build` / CI. `go install` skips that, so `ggc version` will only print `dev`. Prefer the script, Homebrew, or pre-built binaries if you care about version info.

## Pre-built binaries

Download the archive for your OS/arch from the [releases page](https://github.com/bmf-san/ggc/releases) and drop `ggc` on your `PATH`.

Supported targets:

- macOS: `darwin_amd64`, `darwin_arm64`, and a universal binary (one file for both) starting from v8.3.0
- Linux: `linux_amd64`, `linux_arm64`
- Windows: `windows_amd64`

## Windows

Windows binaries are published to the [releases page](https://github.com/bmf-san/ggc/releases) as `ggc_Windows_x86_64.zip`. Steps:

1. Download and unzip the archive.
2. Move `ggc.exe` to a folder on your `PATH` (e.g. `%USERPROFILE%\bin`).
3. In PowerShell, check: `ggc version`.

If you use Git Bash or WSL the Linux instructions above also work unchanged.

## Build from source

```bash
git clone https://github.com/bmf-san/ggc.git
cd ggc
make build
sudo mv ggc /usr/local/bin/
```

`make build` stamps the binary with the current tag / commit so `ggc version` reports real values.

## Shell completions

After installing, generate and install the completion script for your shell:

=== "bash"

    ```bash
    ggc completion bash | sudo tee /etc/bash_completion.d/ggc
    ```

=== "zsh"

    ```bash
    ggc completion zsh > ~/.zsh/completions/_ggc
    # add '~/.zsh/completions' to fpath in your .zshrc
    ```

=== "fish"

    ```bash
    ggc completion fish > ~/.config/fish/completions/ggc.fish
    ```

=== "PowerShell"

    ```powershell
    ggc completion powershell | Out-String | Invoke-Expression
    ```

## Verify

```bash
ggc doctor
```

Should print a green checklist. If anything is `WARN` or `FAIL`, the output explains what to fix. See [Troubleshooting](troubleshooting.md) for details.
# Installation

## Homebrew (macOS / Linux)

```bash
brew install bmf-san/tap/ggc
```

## Go

Requires Go 1.25 or newer:

```bash
go install github.com/bmf-san/ggc/v8@latest
```

## Pre-built binaries

Download the archive for your OS/arch from the [releases page](https://github.com/bmf-san/ggc/releases) and drop `ggc` on your `$PATH`.

macOS universal binaries (one file that runs on both Intel and Apple Silicon) are published starting with v8.3.0.

## Shell completions

After installing, generate and install the completion script for your shell:

=== "bash"

    ```bash
    ggc completion bash | sudo tee /etc/bash_completion.d/ggc
    ```

=== "zsh"

    ```bash
    ggc completion zsh > ~/.zsh/completions/_ggc
    # add '~/.zsh/completions' to fpath in your .zshrc
    ```

=== "fish"

    ```bash
    ggc completion fish > ~/.config/fish/completions/ggc.fish
    ```

## Verify

```bash
ggc doctor
```

Should print a green checklist. If something is WARN or FAIL, the output explains what to fix.
