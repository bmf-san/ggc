package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// BenchmarkStartup_Version exercises the full RunApp path for the lightest
// command (`ggc version`). This benchmark is the baseline for any future
// startup-time optimization work: if a change makes it noticeably faster or
// slower, benchstat will show it on the PR.
//
// The benchmark swaps stdout/stderr for /dev/null so terminal I/O does not
// dominate the measurement.
func BenchmarkStartup_Version(b *testing.B) {
	silence := redirectStdio(b)
	defer silence()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RunApp([]string{"version"})
	}
}

// BenchmarkStartup_Help exercises the `ggc help` path which touches the full
// command registry and the help renderer.
func BenchmarkStartup_Help(b *testing.B) {
	silence := redirectStdio(b)
	defer silence()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RunApp([]string{"help"})
	}
}

// redirectStdio points os.Stdout and os.Stderr at a discard buffer during the
// benchmark and restores them on cleanup. We also drop a discard buffer into
// the returned closer so the compiler cannot elide the writes.
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
		// Touch buf so the optimizer does not drop the allocation.
		_ = bytes.NewBuffer(nil)
	}
}
