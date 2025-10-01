//go:build windows

package termio

import "testing"

func TestPendingInputStub(t *testing.T) {
	n, err := PendingInput(0)
	if err != nil {
		t.Fatalf("PendingInput returned error: %v", err)
	}
	if n != 0 {
		t.Fatalf("PendingInput returned %d, want 0", n)
	}
}
