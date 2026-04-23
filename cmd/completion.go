package cmd

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// embeddedCompletions ships the shell completion scripts inside the ggc
// binary so `ggc completion install <shell>` works with a stock Homebrew
// install (no access to the source tree required).
//
//go:embed completions/ggc.bash completions/ggc.zsh completions/ggc.fish
var embeddedCompletions embed.FS

// Completer handles the `ggc completion ...` subcommand.
type Completer struct {
	outputWriter io.Writer
	userHomeDir  func() (string, error)
	helper       *Helper
}

// NewCompleter returns a Completer writing to stdout.
func NewCompleter() *Completer {
	return &Completer{
		outputWriter: os.Stdout,
		userHomeDir:  os.UserHomeDir,
		helper:       NewHelper(),
	}
}

// Completion dispatches the subcommand.
func (c *Completer) Completion(args []string) {
	if len(args) == 0 {
		c.helper.outputWriter = c.outputWriter
		c.helper.ShowCompletionHelp()
		return
	}
	switch args[0] {
	case "bash", "zsh", "fish":
		c.print(args[0])
	case "install":
		if len(args) < 2 {
			_, _ = fmt.Fprintln(c.outputWriter, "usage: ggc completion install <bash|zsh|fish>")
			return
		}
		c.install(args[1])
	default:
		c.helper.outputWriter = c.outputWriter
		c.helper.ShowCompletionHelp()
	}
}

// print emits the embedded completion script for the given shell to stdout.
func (c *Completer) print(shell string) {
	data, err := embeddedCompletions.ReadFile("completions/ggc." + shell)
	if err != nil {
		_, _ = fmt.Fprintf(c.outputWriter, "unknown shell: %s\n", shell)
		return
	}
	_, _ = c.outputWriter.Write(data)
}

// install writes the completion script to the conventional location for
// the given shell. It picks the first writable target among a
// deliberately short list of well-known locations so the behaviour is
// predictable across distros and package managers.
func (c *Completer) install(shell string) {
	home, err := c.userHomeDir()
	if err != nil {
		_, _ = fmt.Fprintf(c.outputWriter, "cannot resolve home directory: %v\n", err)
		return
	}
	target, ok := c.targetPath(shell, home)
	if !ok {
		_, _ = fmt.Fprintf(c.outputWriter, "unknown shell: %s (supported: bash, zsh, fish)\n", shell)
		return
	}
	data, err := embeddedCompletions.ReadFile("completions/ggc." + shell)
	if err != nil {
		_, _ = fmt.Fprintf(c.outputWriter, "no embedded completion for %s: %v\n", shell, err)
		return
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		_, _ = fmt.Fprintf(c.outputWriter, "failed to create %s: %v\n", filepath.Dir(target), err)
		return
	}
	if err := os.WriteFile(target, data, 0o644); err != nil {
		_, _ = fmt.Fprintf(c.outputWriter, "failed to write %s: %v\n", target, err)
		return
	}
	_, _ = fmt.Fprintf(c.outputWriter, "installed %s completion to %s\n", shell, target)
	c.printReloadHint(shell)
}

// targetPath returns the canonical per-user install path for a shell.
// Returning a per-user location (rather than /etc/... or /usr/share/...)
// avoids surprising sudo requirements; users who want a system-wide
// install can pipe `ggc completion <shell>` into the path of their choice.
func (c *Completer) targetPath(shell, home string) (string, bool) {
	switch shell {
	case "bash":
		return filepath.Join(home, ".local/share/bash-completion/completions/ggc"), true
	case "zsh":
		return filepath.Join(home, ".zsh/completions/_ggc"), true
	case "fish":
		return filepath.Join(home, ".config/fish/completions/ggc.fish"), true
	default:
		return "", false
	}
}

// printReloadHint tells the user the one manual step they still need:
// loading the new completion in their current shell session.
func (c *Completer) printReloadHint(shell string) {
	switch shell {
	case "bash":
		_, _ = fmt.Fprintln(c.outputWriter, "Restart your shell or `source ~/.bashrc` to activate it.")
	case "zsh":
		_, _ = fmt.Fprintln(c.outputWriter,
			"Ensure ~/.zsh/completions is on $fpath (e.g. add `fpath=(~/.zsh/completions $fpath)` to ~/.zshrc) and restart your shell.")
	case "fish":
		_, _ = fmt.Fprintln(c.outputWriter, "Fish will pick up the new completion in any new session.")
	}
}
