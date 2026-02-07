//go:build windows

package termio

import (
	"syscall"
	"unsafe"
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	peekConsoleInput = kernel32.NewProc("PeekConsoleInputW")
)

// inputRecord represents a Windows INPUT_RECORD structure.
// We only need the EventType field to filter for keyboard events.
type inputRecord struct {
	EventType uint16
	_         uint16 // padding
	Event     [16]byte
}

// keyEvent is the Windows KEY_EVENT constant.
const keyEvent = 0x0001

// pendingInput returns whether there are pending keyboard input events on Windows.
// This uses PeekConsoleInputW to check for KEY_EVENT type events only,
// filtering out mouse, window resize, and other non-keyboard events.
func pendingInput(fd uintptr) (int, error) {
	var inputRecords [16]inputRecord
	var numEventsRead uint32

	ret, _, err := peekConsoleInput.Call(
		fd,
		uintptr(unsafe.Pointer(&inputRecords[0])),
		uintptr(len(inputRecords)),
		uintptr(unsafe.Pointer(&numEventsRead)),
	)

	if ret == 0 {
		return 0, err
	}

	// Count only keyboard events.
	// Return 1 if any keyboard event is found, for consistency with Unix
	// implementation which also returns 1 (not actual count) when input is available.
	for i := uint32(0); i < numEventsRead; i++ {
		if inputRecords[i].EventType == keyEvent {
			return 1, nil
		}
	}

	return 0, nil
}
