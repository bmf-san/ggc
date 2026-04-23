package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/bmf-san/ggc/v8/internal/config"

	"go.yaml.in/yaml/v3"
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
		// Keep the helper's writer in sync so callers/tests that redirect
		// d.outputWriter see the help output too.
		d.helper.outputWriter = d.outputWriter
		d.helper.ShowDoctorHelp()
		return
	}
	results := []diagResult{
		d.checkGoRuntime(),
		d.checkGitBinary(),
		d.checkGgcOnPATH(),
		d.checkGgcConfig(),
		d.checkCompletions("bash"),
		d.checkCompletions("zsh"),
		d.checkCompletions("fish"),
		d.checkTerm(),
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
	// Invoke the resolved path so the reported binary is exactly the one
	// we measured, even if PATH changes between LookPath and exec.
	out, err := d.execCommand(path, "--version").Output()
	if err != nil {
		return diagResult{name: "git binary", ok: false, detail: fmt.Sprintf("%s: %v", path, err)}
	}
	trimmed := strings.TrimSpace(string(out))
	if major, minor, ok := parseGitVersion(trimmed); ok {
		if major < minGitMajor || (major == minGitMajor && minor < minGitMinor) {
			return diagResult{
				name:   "git binary",
				ok:     false,
				warn:   true,
				detail: fmt.Sprintf("%s (%s) is older than the recommended %d.%d; some subcommands may not work", path, trimmed, minGitMajor, minGitMinor),
			}
		}
	}
	return diagResult{name: "git binary", ok: true, detail: fmt.Sprintf("%s (%s)", path, trimmed)}
}

// minGit{Major,Minor} is the lowest Git version we actively test against.
// Older Git ships without the porcelain flags several ggc subcommands rely on.
const (
	minGitMajor = 2
	minGitMinor = 30
)

// parseGitVersion extracts the major+minor number from a `git --version`
// line such as "git version 2.44.0" or "git version 2.39.3 (Apple Git-146)".
func parseGitVersion(s string) (int, int, bool) {
	const prefix = "git version "
	if !strings.HasPrefix(s, prefix) {
		return 0, 0, false
	}
	rest := strings.TrimPrefix(s, prefix)
	if idx := strings.IndexAny(rest, " \t"); idx != -1 {
		rest = rest[:idx]
	}
	parts := strings.SplitN(rest, ".", 3)
	if len(parts) < 2 {
		return 0, 0, false
	}
	major, err1 := strconv.Atoi(parts[0])
	minor, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	return major, minor, true
}

// checkGgcOnPATH verifies that a user invoking `ggc` from an arbitrary
// directory gets the same binary that is currently running. Shadowing by
// an old Homebrew install or a stale `go install` copy is a common silent
// failure mode.
func (d *Doctor) checkGgcOnPATH() diagResult {
	self, selfErr := os.Executable()
	resolved, lookErr := d.lookPath("ggc")
	switch {
	case lookErr != nil && selfErr == nil:
		return diagResult{
			name:   "ggc on PATH",
			ok:     false,
			warn:   true,
			detail: fmt.Sprintf("running %s but 'ggc' is not on PATH; add its directory to PATH to use shell completions", self),
		}
	case lookErr != nil:
		return diagResult{name: "ggc on PATH", ok: false, warn: true, detail: "'ggc' not found on PATH"}
	case selfErr != nil:
		return diagResult{name: "ggc on PATH", ok: true, detail: resolved}
	}
	// Compare by resolved symlink target so /opt/homebrew/bin/ggc ->
	// /opt/homebrew/Cellar/ggc/…/bin/ggc still counts as a match.
	selfReal, _ := filepath.EvalSymlinks(self)
	resolvedReal, _ := filepath.EvalSymlinks(resolved)
	if selfReal != "" && resolvedReal != "" && selfReal != resolvedReal {
		return diagResult{
			name:   "ggc on PATH",
			ok:     false,
			warn:   true,
			detail: fmt.Sprintf("running %s but PATH resolves 'ggc' to %s; an older install may shadow this one", self, resolved),
		}
	}
	return diagResult{name: "ggc on PATH", ok: true, detail: resolved}
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
	// Parse the YAML directly: config.NewConfigManager requires a non-nil
	// git.ConfigOps (getDefaultConfig calls methods on it), and the doctor
	// runs at diagnose time without that dependency wired in. We only need
	// to know whether the file is a syntactically valid ggc config.
	if err := parseConfigFile(found); err != nil {
		return diagResult{name: "ggc config", ok: false, detail: fmt.Sprintf("%s: %v", found, err)}
	}
	return diagResult{name: "ggc config", ok: true, detail: fmt.Sprintf("%s loaded", found)}
}

// parseConfigFile validates that the given path is a YAML file that matches
// the ggc config schema. It does not apply defaults or perform any git
// lookups, which makes it safe to call from the doctor without wiring up a
// full config.Manager.
func parseConfigFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var cfg config.Config
	return yaml.Unmarshal(data, &cfg)
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
	case "fish":
		candidates = []string{
			"/usr/share/fish/vendor_completions.d/ggc.fish",
			"/usr/local/share/fish/vendor_completions.d/ggc.fish",
			"/opt/homebrew/share/fish/vendor_completions.d/ggc.fish",
			filepath.Join(home, ".config/fish/completions/ggc.fish"),
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
		detail: "not installed in a well-known location; try `ggc completion install " + shell + "`",
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

// checkTerm warns when $TERM looks like something the interactive TUI
// cannot fully drive (dumb terminal, unset, or vt52-level).
func (d *Doctor) checkTerm() diagResult {
	term := os.Getenv("TERM")
	switch term {
	case "":
		return diagResult{
			name:   "TERM",
			ok:     false,
			warn:   true,
			detail: "$TERM is unset; interactive mode may render incorrectly",
		}
	case "dumb":
		return diagResult{
			name:   "TERM",
			ok:     false,
			warn:   true,
			detail: "$TERM=dumb; interactive mode will not work (one-shot subcommands are fine)",
		}
	}
	return diagResult{name: "TERM", ok: true, detail: term}
}
