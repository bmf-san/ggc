package main

import (
	"io"
	"os"
	"testing"
)

// BenchmarkStartup_Version exercises the full RunApp path for the lightest
// command (`ggc version`). This benchmark is the baseline for any future
// startup-time optimization work: if a change makes it noticeably faster or
// slower, benchstat will show it on the PR.
//
// The benchmark swaps stdout/stderr for a drain pipe so terminal I/O does
// not dominate the measurement, and sandboxes HOME/XDG_CONFIG_HOME so that
// `go test -bench` never reads or writes the developer's real config file.
func BenchmarkStartup_Version(b *testing.B) {
	sandboxHome(b)
	silence := redirectStdio(b)
	defer silence()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := RunApp([]string{"version"}); err != nil {
			b.Fatalf("RunApp(version) failed: %v", err)
		}
	}
}

// BenchmarkStartup_Help exercises the `ggc help` path which touches the full
// command registry and the help renderer.
func BenchmarkStartup_Help(b *testing.B) {
	sandboxHome(b)
	silence := redirectStdio(b)
	defer silence()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := RunApp([]string{"help"}); err != nil {
			b.Fatalf("RunApp(help) failed: %v", err)
		}
	}
}

// sandboxHome points HOME and XDG_CONFIG_HOME at a per-benchmark temp
// directory so that reading/writing the ggc config file in RunApp never
// touches the developer's real ~/.ggcconfig.yaml.
func sandboxHome(tb testing.TB) {
	tb.Helper()
	tmp := tb.TempDir()
	tb.Setenv("HOME", tmp)
	tb.Setenv("XDG_CONFIG_HOME", tmp)
}

// redirectStdio points os.Stdout and os.Stderr at a pipe that is drained on
// a goroutine, so writers never block and the benchmark's output does not
// bleed into the test runner's log. Writes to the real os.Stdout/os.Stderr
// already have observable side effects, so no additional "keepalive"
// allocation is needed to prevent the compiler from eliding them.
func redirectStdio(tb testing.TB) func() {
	tb.Helper()
	origStdout := os.Stdout
	origStderr := os.Stderr

	// Use real *os.File pipes because some downstream code checks
	// (*os.File).Stat() and panics on a regular *bytes.Buffer.
	r, w, err := os.Pipe()
	if err != nil {
		tb.Fatal(err)
	}
	os.Stdout = w
	os.Stderr = w

	// Drain the pipe on a goroutine so writers never block.
	done := make(chan struct{})
	go func() {
		_, _ = io.Copy(io.Discard, r)
		close(done)
	}()

	return func() {
		_ = w.Close()
		<-done
		_ = r.Close()
		os.Stdout = origStdout
		os.Stderr = origStderr
	}
}
