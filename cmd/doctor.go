package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bmf-san/ggc/v8/internal/config"
)

// Doctor inspects the local environment and reports anything that could
// prevent ggc from working correctly.
type Doctor struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
	lookPath     func(string) (string, error)
	userHomeDir  func() (string, error)
	stdinStat    func() (os.FileInfo, error)
}

// NewDoctor creates a new Doctor instance.
func NewDoctor() *Doctor {
	return &Doctor{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
		lookPath:     exec.LookPath,
		userHomeDir:  os.UserHomeDir,
		stdinStat:    func() (os.FileInfo, error) { return os.Stdin.Stat() },
	}
}

// diagResult captures one check's outcome.
type diagResult struct {
	name   string
	ok     bool
	warn   bool
	detail string
}

// Doctor runs diagnostics. Any arg prints help (stub for future --json etc).
func (d *Doctor) Doctor(args []string) {
	if len(args) > 0 {
		d.helper.ShowDoctorHelp()
		return
	}
	results := []diagResult{
		d.checkGoRuntime(),
		d.checkGitBinary(),
		d.checkGgcConfig(),
		d.checkCompletions("bash"),
		d.checkCompletions("zsh"),
		d.checkTTY(),
	}
	d.printReport(results)
}

func (d *Doctor) printReport(results []diagResult) {
	hardFailures := 0
	for _, r := range results {
		prefix := "OK  "
		switch {
		case !r.ok && r.warn:
			prefix = "WARN"
		case !r.ok:
			prefix = "FAIL"
			hardFailures++
		}
		if r.detail == "" {
			_, _ = fmt.Fprintf(d.outputWriter, "[%s] %s\n", prefix, r.name)
		} else {
			_, _ = fmt.Fprintf(d.outputWriter, "[%s] %s: %s\n", prefix, r.name, r.detail)
		}
	}
	if hardFailures > 0 {
		_, _ = fmt.Fprintf(d.outputWriter, "\n%d hard failure(s). ggc may not work correctly.\n", hardFailures)
	} else {
		_, _ = fmt.Fprintln(d.outputWriter, "\nEverything looks good.")
	}
}

func (d *Doctor) checkGoRuntime() diagResult {
	return diagResult{
		name:   "Go runtime",
		ok:     true,
		detail: fmt.Sprintf("%s (%s/%s)", runtime.Version(), runtime.GOOS, runtime.GOARCH),
	}
}

func (d *Doctor) checkGitBinary() diagResult {
	path, err := d.lookPath("git")
	if err != nil {
		return diagResult{name: "git binary", ok: false, detail: "'git' not found on PATH"}
	}
	out, err := d.execCommand("git", "--version").Output()
	if err != nil {
		return diagResult{name: "git binary", ok: false, detail: fmt.Sprintf("%s: %v", path, err)}
	}
	return diagResult{name: "git binary", ok: true, detail: fmt.Sprintf("%s (%s)", path, strings.TrimSpace(string(out)))}
}

// configCandidatePaths mirrors (manager).getConfigPaths without exposing it.
func (d *Doctor) configCandidatePaths() []string {
	home, err := d.userHomeDir()
	if err != nil {
		return nil
	}
	return []string{
		filepath.Join(home, ".ggcconfig.yaml"),
		filepath.Join(home, ".config", "ggc", "config.yaml"),
	}
}

func (d *Doctor) checkGgcConfig() diagResult {
	paths := d.configCandidatePaths()
	if len(paths) == 0 {
		return diagResult{name: "ggc config", ok: false, warn: true, detail: "cannot resolve $HOME"}
	}
	var found string
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			found = p
			break
		}
	}
	if found == "" {
		return diagResult{
			name:   "ggc config",
			ok:     true,
			detail: fmt.Sprintf("no config yet (defaults in use; will be created at %s)", paths[0]),
		}
	}
	m := config.NewConfigManager(nil)
	if err := m.Load(); err != nil {
		return diagResult{name: "ggc config", ok: false, detail: fmt.Sprintf("%s: %v", found, err)}
	}
	return diagResult{name: "ggc config", ok: true, detail: fmt.Sprintf("%s loaded", found)}
}

// checkCompletions looks for an installed completion script in well-known
// locations. Missing is WARN, not FAIL: ggc works without shell completion.
func (d *Doctor) checkCompletions(shell string) diagResult {
	home, err := d.userHomeDir()
	if err != nil {
		return diagResult{name: shell + " completions", ok: false, warn: true, detail: fmt.Sprintf("cannot resolve $HOME: %v", err)}
	}
	var candidates []string
	switch shell {
	case "bash":
		candidates = []string{
			"/etc/bash_completion.d/ggc",
			"/usr/local/etc/bash_completion.d/ggc",
			"/opt/homebrew/etc/bash_completion.d/ggc",
			filepath.Join(home, ".local/share/bash-completion/completions/ggc"),
		}
	case "zsh":
		candidates = []string{
			"/usr/share/zsh/site-functions/_ggc",
			"/usr/local/share/zsh/site-functions/_ggc",
			"/opt/homebrew/share/zsh/site-functions/_ggc",
			filepath.Join(home, ".zsh/completions/_ggc"),
		}
	default:
		return diagResult{name: shell + " completions", ok: true}
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return diagResult{name: shell + " completions", ok: true, detail: c}
		}
	}
	return diagResult{
		name:   shell + " completions",
		ok:     false,
		warn:   true,
		detail: "not installed in a well-known location (see tools/completions/ in the repo)",
	}
}

func (d *Doctor) checkTTY() diagResult {
	fi, err := d.stdinStat()
	if err != nil {
		return diagResult{name: "stdin TTY", ok: false, warn: true, detail: err.Error()}
	}
	if fi.Mode()&os.ModeCharDevice == 0 {
		return diagResult{
			name:   "stdin TTY",
			ok:     true,
			detail: "stdin is not a TTY; interactive mode and `debug-keys raw` will not work in this shell",
		}
	}
	return diagResult{name: "stdin TTY", ok: true, detail: "stdin is a TTY"}
}
