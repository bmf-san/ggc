//go:build !windows

package prompt

import (
	"bytes"
	"errors"
	"io"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/creack/pty"
)

func TestInputCtrlCCancelsWithoutEnter(t *testing.T) {
	master, slave, err := pty.Open()
	if err != nil {
		t.Skipf("pty open not available: %v", err)
	}

	t.Cleanup(func() {
		_ = master.Close()
		_ = slave.Close()
	})

	prompter := New(slave, slave).(*StandardPrompter)

	type result struct {
		line     string
		canceled bool
		err      error
	}

	resultCh := make(chan result, 1)
	go func() {
		line, canceled, inputErr := prompter.Input("branch: ")
		resultCh <- result{line: line, canceled: canceled, err: inputErr}
	}()

	waitForPrompt(t, master, "branch: ")

	if _, err := master.Write([]byte{3}); err != nil {
		t.Fatalf("failed to send ctrl+c: %v", err)
	}

	select {
	case res := <-resultCh:
		if res.err != nil {
			t.Fatalf("Input returned error: %v", res.err)
		}
		if !res.canceled {
			t.Fatal("expected canceled to be true")
		}
		if res.line != "" {
			t.Fatalf("expected empty line, got %q", res.line)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Input did not cancel after ctrl+c within timeout")
	}
}

func waitForPrompt(t *testing.T, master *os.File, prompt string) {
	t.Helper()

	var buf bytes.Buffer
	deadline := time.Now().Add(2 * time.Second)

	if err := syscall.SetNonblock(int(master.Fd()), true); err != nil {
		t.Fatalf("failed to set non-blocking mode: %v", err)
	}

	for time.Now().Before(deadline) {
		tmp := make([]byte, 128)
		n, err := master.Read(tmp)
		if n > 0 {
			buf.Write(tmp[:n])
			if bytes.Contains(buf.Bytes(), []byte(prompt)) {
				return
			}
		}

		if err != nil {
			if errors.Is(err, syscall.EAGAIN) {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			if errors.Is(err, io.EOF) {
				break
			}
			t.Fatalf("reading prompt: %v", err)
		}
	}

	t.Fatalf("prompt %q not observed; captured %q", prompt, buf.String())
}
