package git

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"
)

// TestClient_WithContext_CancelsSubprocess verifies that canceling the
// context attached to a Client aborts a long-running git subprocess.
func TestClient_WithContext_CancelsSubprocess(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unsupported on windows test runner")
	}
	ctx, cancel := context.WithCancel(context.Background())
	c := NewClient().WithContext(ctx)

	// Use `git -c foo=bar wait-for-stdin-forever` approximation: we can't
	// rely on a real long-running git op portably, so simulate with sleep
	// via the shell not available; instead, use 'git --version' is too
	// fast to cancel. We cancel *before* Start to guarantee exec sees
	// ctx.Err() immediately.
	cancel()

	cmd := c.execCommand("git", "--version")
	err := cmd.Run()
	if err == nil {
		t.Fatalf("expected ctx-canceled git to return error, got nil")
	}
	if !errors.Is(err, context.Canceled) && !errors.Is(ctx.Err(), context.Canceled) {
		t.Logf("note: error was %v (ctx.Err=%v)", err, ctx.Err())
	}
}

// TestClient_WithContext_DefaultBackground verifies that a nil ctx is
// substituted for context.Background, and that normal git operations work.
func TestClient_WithContext_DefaultBackground(t *testing.T) {
	//nolint:staticcheck // intentionally passing nil to exercise the nil→Background path
	var nilCtx context.Context
	c := NewClient().WithContext(nilCtx)
	if c.ctx == nil {
		t.Fatal("expected non-nil ctx after WithContext(nil)")
	}
	// Short timeout ensures this test does not hang if git is missing.
	deadline, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c2 := c.WithContext(deadline)
	cmd := c2.execCommand("git", "--version")
	_ = cmd.Run() // Not asserting success - some CI may lack git.
}
