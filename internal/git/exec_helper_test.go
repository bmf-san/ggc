package git

import (
	"encoding/base64"
	"os"
	"os/exec"
	"testing"
)

// The Unix shorthands used by these mocks (echo, true, false) don't exist on
// Windows, and an echo.exe found on PATH may not support "-n". Route every
// mock through a helper process instead: the test binary re-executes itself,
// prints (or fails) exactly what the caller asked for, and behaves identically
// on every platform. This is the pattern documented in os/exec's Example.
const (
	helperEnv     = "GGC_HELPER_PROCESS"
	helperModeEnv = "GGC_HELPER_MODE"
	helperModeRun = "-test.run=^TestHelperProcess$"
)

// echoCommand returns a command that writes output to stdout verbatim (no
// trailing newline) and exits 0.
func echoCommand(output string) *exec.Cmd {
	cmd := exec.Command(os.Args[0], helperModeRun, "--")
	cmd.Env = append(os.Environ(),
		helperEnv+"=1",
		helperModeEnv+"="+base64.StdEncoding.EncodeToString([]byte(output)))
	return cmd
}

// successCommand returns a command that exits 0 without output.
func successCommand() *exec.Cmd {
	return echoCommand("")
}

// failingCommand returns a command that exits 1 without output.
func failingCommand() *exec.Cmd {
	cmd := exec.Command(os.Args[0], helperModeRun, "--")
	cmd.Env = append(os.Environ(),
		helperEnv+"=1",
		helperModeEnv+"=fail")
	return cmd
}

// TestHelperProcess is the re-executed helper behind echoCommand,
// successCommand and failingCommand. It is a no-op unless re-invoked with
// GGC_HELPER_PROCESS set by one of those constructors.
func TestHelperProcess(t *testing.T) {
	if os.Getenv(helperEnv) != "1" {
		return
	}
	mode := ""
	for _, e := range os.Environ() {
		if len(e) > len(helperModeEnv)+1 && e[:len(helperModeEnv)+1] == helperModeEnv+"=" {
			mode = e[len(helperModeEnv)+1:]
			break
		}
	}
	if mode == "fail" {
		os.Exit(1)
	}
	out, err := base64.StdEncoding.DecodeString(mode)
	if err != nil {
		os.Exit(2)
	}
	_, _ = os.Stdout.Write(out)
	os.Exit(0)
}
