//go:build windows

package termio

// pendingInput is a stub on Windows; console APIs require a different strategy.
func pendingInput(uintptr) (int, error) {
	return 0, nil
}
