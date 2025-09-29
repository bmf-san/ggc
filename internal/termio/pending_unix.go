//go:build !windows

package termio

import "golang.org/x/sys/unix"

func pendingInput(fd uintptr) (int, error) {
	pollFds := []unix.PollFd{{Fd: int32(fd), Events: unix.POLLIN}}
	n, err := unix.Poll(pollFds, 0)
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, nil
	}
	if pollFds[0].Revents&unix.POLLIN != 0 {
		return 1, nil
	}
	return 0, nil
}
