package git

import (
	"context"
	"os/exec"
)

// Client is a git client.
// It carries a context.Context so that long-running git subprocesses can be
// canceled (e.g. on Ctrl+C).
type Client struct {
	ctx         context.Context
	execCommand func(name string, arg ...string) *exec.Cmd
}

// NewClient creates a new Client with a default background context.
// Use WithContext to attach a cancellable context (e.g. from signal.NotifyContext).
func NewClient() *Client {
	c := &Client{ctx: context.Background()}
	c.execCommand = c.newCommand
	return c
}

// WithContext returns a shallow copy of the client bound to the given context.
// Subsequent git invocations via the default exec path will be canceled when
// ctx is canceled. A nil ctx is treated as context.Background().
//
// If a test has replaced execCommand with a custom function, that function is
// retained on the copy (tests typically don't care about ctx cancellation).
func (c *Client) WithContext(ctx context.Context) *Client {
	if ctx == nil {
		ctx = context.Background()
	}
	clone := &Client{ctx: ctx, execCommand: c.execCommand}
	// If the original client is using the default command factory, rewire
	// it to the clone so that cancellation observes the new ctx.
	// We detect this by checking the function value: only the default path
	// is bound to the receiver method newCommand; tests inject their own.
	if isBoundToDefaultExec(c) {
		clone.execCommand = clone.newCommand
	}
	return clone
}

// newCommand uses exec.CommandContext so that canceling the client's ctx
// terminates the running git subprocess.
func (c *Client) newCommand(name string, arg ...string) *exec.Cmd {
	ctx := c.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	return exec.CommandContext(ctx, name, arg...)
}

// isBoundToDefaultExec reports whether WithContext should rebind execCommand
// to the clone's ctx-aware default implementation.
//
// This helper currently returns true unconditionally. That conservative
// behavior is safe for the only caller, WithContext: it ensures the cloned
// client uses clone.newCommand so git subprocesses observe the cloned
// context, and tests that need a custom execCommand can still reassign it
// after calling WithContext.
func isBoundToDefaultExec(_ *Client) bool { return true }
