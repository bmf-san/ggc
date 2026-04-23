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

## Docker (ghcr.io)

Multi-arch images (`linux/amd64`, `linux/arm64`) are published to GitHub Container Registry on every tagged release, starting from the next release after v8.3.0.

```bash
docker pull ghcr.io/bmf-san/ggc:latest
# or pin a version
docker pull ghcr.io/bmf-san/ggc:v8.4.0

# run against the current directory
docker run --rm -it -v "$PWD:/work" ghcr.io/bmf-san/ggc:latest status
```

The image is based on `alpine:3.22` and bundles `git`. Runs as an unprivileged user (`ggc`).

## Build from source

```bash
git clone https://github.com/bmf-san/ggc.git
cd ggc
make build
sudo mv ggc /usr/local/bin/
```

`make build` stamps the binary with the current tag / commit so `ggc version` reports real values.

## Shell completions

The completion scripts are embedded in the `ggc` binary. One command per shell:

```bash
ggc completion install bash   # -> ~/.local/share/bash-completion/completions/ggc
ggc completion install zsh    # -> ~/.zsh/completions/_ggc
ggc completion install fish   # -> ~/.config/fish/completions/ggc.fish
```

Restart your shell (or for zsh: make sure `~/.zsh/completions` is on `$fpath`).

### Piping to a custom location

`ggc completion <shell>` prints the script to stdout, so you can redirect it anywhere:

```bash
ggc completion zsh  | sudo tee /usr/local/share/zsh/site-functions/_ggc
ggc completion bash | sudo tee /etc/bash_completion.d/ggc
```

### Reading the pre-built files directly

Source files are also versioned in [`cmd/completions/`](https://github.com/bmf-san/ggc/tree/main/cmd/completions); they are regenerated from the command registry by `make completions`.

## Verify

```bash
ggc doctor
```

## Verifying a downloaded release

Starting with releases signed by cosign keyless, you can verify archives before installing:

```bash
TAG=v8.2.0     # or whichever release
BASE="https://github.com/bmf-san/ggc/releases/download/${TAG}"

curl -sSLO "${BASE}/checksums.txt"
curl -sSLO "${BASE}/checksums.txt.sig"
curl -sSLO "${BASE}/checksums.txt.pem"

cosign verify-blob \
  --certificate checksums.txt.pem \
  --signature   checksums.txt.sig \
  --certificate-identity-regexp "^https://github.com/bmf-san/ggc/\.github/workflows/release\.yml@" \
  --certificate-oidc-issuer     "https://token.actions.githubusercontent.com" \
  checksums.txt

# Then verify the archive you downloaded against the signed checksum file:
sha256sum -c --ignore-missing checksums.txt
```

The certificate is short-lived (Fulcio) and the signature is recorded in the public Rekor transparency log — no long-lived key to rotate.

