//go:build !windows

package termio

import (
	"os"
	"testing"
)

func mustClose(t *testing.T, f *os.File, name string) {
	t.Helper()
	if err := f.Close(); err != nil {
		t.Fatalf("close %s failed: %v", name, err)
	}
}

func TestPendingInputPipe(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	t.Cleanup(func() {
		mustClose(t, r, "pipe reader")
		mustClose(t, w, "pipe writer")
	})

	fd := r.Fd()

	n, err := PendingInput(fd)
	if err != nil {
		t.Fatalf("PendingInput before write returned error: %v", err)
	}
	if n != 0 {
		t.Fatalf("PendingInput before write returned %d, want 0", n)
	}

	if _, err := w.Write([]byte("x")); err != nil {
		t.Fatalf("write to pipe failed: %v", err)
	}

	n, err = PendingInput(fd)
	if err != nil {
		t.Fatalf("PendingInput after write returned error: %v", err)
	}
	if n != 1 {
		t.Fatalf("PendingInput after write returned %d, want 1", n)
	}

	var buf [1]byte
	if _, err := r.Read(buf[:]); err != nil {
		t.Fatalf("read from pipe failed: %v", err)
	}

	n, err = PendingInput(fd)
	if err != nil {
		t.Fatalf("PendingInput after drain returned error: %v", err)
	}
	if n != 0 {
		t.Fatalf("PendingInput after drain returned %d, want 0", n)
	}
}
