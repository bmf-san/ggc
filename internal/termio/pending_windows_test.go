//go:build windows

package termio

import (
	"os"
	"testing"
)

func TestPendingInputWindows(t *testing.T) {
	// Test with an invalid handle - should return error or 0
	n, err := PendingInput(0)
	// Invalid handle may return an error or 0, both are acceptable
	if err != nil {
		t.Logf("PendingInput with invalid handle returned error (expected): %v", err)
	} else if n != 0 {
		t.Logf("PendingInput with invalid handle returned %d", n)
	}
}

func TestPendingInputWithStdin(t *testing.T) {
	// Test with actual stdin handle
	// Note: This test may behave differently depending on whether
	// stdin is connected to a console or redirected
	fd := os.Stdin.Fd()
	n, err := PendingInput(fd)
	// We don't assert specific values since behavior depends on console state
	// Just verify the function doesn't panic
	t.Logf("PendingInput(stdin): n=%d, err=%v", n, err)
}

func TestPendingInputHookWithPendingEvents(t *testing.T) {
	// Test using hook to simulate pending input
	restore := SetPendingInputFunc(func(fd uintptr) (int, error) {
		return 1, nil // Simulate pending input
	})
	defer restore()

	n, err := PendingInput(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 pending event, got %d", n)
	}
}

func TestPendingInputHookWithNoPendingEvents(t *testing.T) {
	// Test using hook to simulate no pending input
	restore := SetPendingInputFunc(func(fd uintptr) (int, error) {
		return 0, nil // Simulate no pending input
	})
	defer restore()

	n, err := PendingInput(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 pending events, got %d", n)
	}
}
