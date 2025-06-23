package cmd

import (
	"bytes"
	"os/exec"
	"testing"
)

type noopCmd struct {
	cmd *exec.Cmd
}

func newNoopCmd() *noopCmd {
	cmd := exec.Command("true")
	return &noopCmd{cmd: cmd}
}

func (c *noopCmd) Run() error {
	return nil
}

func TestRemoteer_Remote_List(t *testing.T) {
	called := false
	var buf bytes.Buffer
	remoteer := &Remoteer{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			called = true
			return newNoopCmd().cmd
		},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	remoteer.helper.outputWriter = &buf

	remoteer.Remote([]string{"list"})

	if !called {
		t.Error("execCommand is not called in list subcommand")
	}
}

func TestRemoteer_Remote_Add(t *testing.T) {
	called := false
	var buf bytes.Buffer
	remoteer := &Remoteer{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			called = true
			return newNoopCmd().cmd
		},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	remoteer.helper.outputWriter = &buf

	remoteer.Remote([]string{"add", "origin", "https://example.com"})

	if !called {
		t.Error("execCommand is not called in add subcommand")
	}
	output := buf.String()
	if output != "Remote 'origin' added\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}

func TestRemoteer_Remote_Remove(t *testing.T) {
	called := false
	var buf bytes.Buffer
	remoteer := &Remoteer{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			called = true
			return newNoopCmd().cmd
		},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	remoteer.helper.outputWriter = &buf

	remoteer.Remote([]string{"remove", "origin"})

	if !called {
		t.Error("execCommand is not called in remove subcommand")
	}
	output := buf.String()
	if output != "Remote 'origin' removed\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}

func TestRemoteer_Remote_SetURL(t *testing.T) {
	called := false
	var buf bytes.Buffer
	remoteer := &Remoteer{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			called = true
			return newNoopCmd().cmd
		},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	remoteer.helper.outputWriter = &buf

	remoteer.Remote([]string{"set-url", "origin", "https://example.com"})

	if !called {
		t.Error("execCommand is not called in set-url subcommand")
	}
	output := buf.String()
	if output != "Remote 'origin' URL updated\n" {
		t.Errorf("unexpected output: got %q", output)
	}
}

func TestRemoteer_Remote_Help(t *testing.T) {
	var buf bytes.Buffer
	remoteer := &Remoteer{
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return newNoopCmd().cmd
		},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	remoteer.helper.outputWriter = &buf

	remoteer.Remote([]string{"unknown"})

	output := buf.String()
	if output == "" || output[:5] != "Usage" {
		t.Errorf("Usage should be displayed, but got: %s", output)
	}
}
