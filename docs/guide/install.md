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
