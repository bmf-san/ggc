package cmd

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func newTestDoctor(out *bytes.Buffer) *Doctor {
	d := NewDoctor()
	d.outputWriter = out
	return d
}

func TestDoctor_GoRuntime_AlwaysOK(t *testing.T) {
	d := newTestDoctor(&bytes.Buffer{})
	r := d.checkGoRuntime()
	if !r.ok {
		t.Fatalf("Go runtime check should always be OK, got %+v", r)
	}
	if r.detail == "" {
		t.Fatal("Go runtime detail should not be empty")
	}
}

func TestDoctor_GitBinary_NotFound(t *testing.T) {
	d := newTestDoctor(&bytes.Buffer{})
	d.lookPath = func(string) (string, error) { return "", errors.New("not found") }
	r := d.checkGitBinary()
	if r.ok || !strings.Contains(r.detail, "not found") {
		t.Fatalf("want not-OK with 'not found', got %+v", r)
	}
}

func TestDoctor_GitBinary_Works(t *testing.T) {
	d := newTestDoctor(&bytes.Buffer{})
	d.lookPath = func(string) (string, error) { return "/usr/bin/git", nil }
	// Spawn a command that prints a predictable line and succeeds.
	d.execCommand = func(_ string, _ ...string) *exec.Cmd {
		return exec.Command("echo", "git version 2.99.9")
	}
	r := d.checkGitBinary()
	if !r.ok || !strings.Contains(r.detail, "2.99.9") {
		t.Fatalf("want OK with version, got %+v", r)
	}
}

func TestDoctor_Config_NoFile(t *testing.T) {
	tmp := t.TempDir()
	d := newTestDoctor(&bytes.Buffer{})
	d.userHomeDir = func() (string, error) { return tmp, nil }
	r := d.checkGgcConfig()
	if !r.ok {
		t.Fatalf("missing config should be OK (defaults), got %+v", r)
	}
	if !strings.Contains(r.detail, "no config yet") {
		t.Fatalf("unexpected detail: %s", r.detail)
	}
}

func TestDoctor_Config_ValidFileLoads(t *testing.T) {
	tmp := t.TempDir()
	// Exercise the ~/.ggcconfig.yaml candidate.
	path := filepath.Join(tmp, ".ggcconfig.yaml")
	if err := os.WriteFile(path, []byte("meta:\n  version: v8.0.0\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	d := newTestDoctor(&bytes.Buffer{})
	d.userHomeDir = func() (string, error) { return tmp, nil }
	r := d.checkGgcConfig()
	if !r.ok {
		t.Fatalf("valid config should be OK, got %+v", r)
	}
	if !strings.Contains(r.detail, "loaded") {
		t.Fatalf("detail should mention loaded, got %q", r.detail)
	}
}

func TestDoctor_Config_InvalidYAMLFails(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, ".ggcconfig.yaml")
	// Emit YAML that the parser rejects outright.
	if err := os.WriteFile(path, []byte("meta: [this is not a map\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	d := newTestDoctor(&bytes.Buffer{})
	d.userHomeDir = func() (string, error) { return tmp, nil }
	r := d.checkGgcConfig()
	if r.ok {
		t.Fatalf("malformed config should be FAIL, got %+v", r)
	}
	if !strings.Contains(r.detail, path) {
		t.Fatalf("detail should mention the failing path, got %q", r.detail)
	}
}

func TestDoctor_Completions_MissingIsWarn(t *testing.T) {
	tmp := t.TempDir()
	d := newTestDoctor(&bytes.Buffer{})
	d.userHomeDir = func() (string, error) { return tmp, nil }
	r := d.checkCompletions("zsh")
	// Missing completion is WARN (not a hard failure).
	if r.ok || !r.warn {
		t.Fatalf("missing completions should be WARN, got %+v", r)
	}
}

func TestDoctor_TTY_NonTTY(t *testing.T) {
	d := newTestDoctor(&bytes.Buffer{})
	// A regular file stat has no ModeCharDevice.
	f, err := os.CreateTemp(t.TempDir(), "stdin")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	d.stdinStat = func() (os.FileInfo, error) { return os.Stat(f.Name()) }
	r := d.checkTTY()
	if !r.ok || !strings.Contains(r.detail, "not a TTY") {
		t.Fatalf("want OK with 'not a TTY', got %+v", r)
	}
}

func TestDoctor_FullReport_NoHardFailures(t *testing.T) {
	var buf bytes.Buffer
	d := newTestDoctor(&buf)
	tmp := t.TempDir()
	d.userHomeDir = func() (string, error) { return tmp, nil }
	d.lookPath = func(string) (string, error) { return "/usr/bin/git", nil }
	d.execCommand = func(_ string, _ ...string) *exec.Cmd {
		return exec.Command("echo", "git version 2.99.9")
	}
	// Pipe-like stat (non-TTY is not a failure).
	f, _ := os.CreateTemp(tmp, "stdin")
	defer func() { _ = f.Close() }()
	d.stdinStat = func() (os.FileInfo, error) { return os.Stat(f.Name()) }

	d.Doctor(nil)
	out := buf.String()
	if strings.Contains(out, "hard failure") {
		t.Fatalf("should not report hard failures, got:\n%s", out)
	}
	if !strings.Contains(out, "Everything looks good.") {
		t.Fatalf("expected success footer, got:\n%s", out)
	}
}
