package termio

import "testing"

func TestPendingInputHookOverride(t *testing.T) {
	const (
		fakeFD   = uintptr(123)
		fakeSize = 7
	)

	calls := 0
	restore := SetPendingInputFunc(func(fd uintptr) (int, error) {
		calls++
		if fd != fakeFD {
			t.Fatalf("pendingInput called with fd %d, want %d", fd, fakeFD)
		}
		return fakeSize, nil
	})
	t.Cleanup(restore)

	size, err := PendingInput(fakeFD)
	if err != nil {
		t.Fatalf("PendingInput returned error: %v", err)
	}
	if size != fakeSize {
		t.Fatalf("PendingInput returned %d, want %d", size, fakeSize)
	}
	if calls != 1 {
		t.Fatalf("pendingInput called %d times, want 1", calls)
	}

	restore()

	stubHits := 0
	restore2 := SetPendingInputFunc(func(fd uintptr) (int, error) {
		stubHits++
		return 0, nil
	})
	t.Cleanup(restore2)

	if _, err := PendingInput(0); err != nil {
		t.Fatalf("PendingInput after restore returned error: %v", err)
	}
	if stubHits != 1 {
		t.Fatalf("pendingInput after restore called %d times, want 1", stubHits)
	}
}
