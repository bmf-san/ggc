package cmd

import (
	"bytes"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v8/internal/testutil"
)

type mockPassthroughClient struct {
	testutil.MockGitClient
	called  bool
	gotName string
	gotArgs []string
	err     error
}

func (m *mockPassthroughClient) RunGit(name string, args []string) error {
	m.called = true
	m.gotName = name
	m.gotArgs = slices.Clone(args)
	return m.err
}

func TestPassthroughCommand_Run(t *testing.T) {
	mock := &mockPassthroughClient{}
	var buf bytes.Buffer
	pc := newPassthroughCommand("cherry-pick", mock)
	pc.outputWriter = &buf
	pc.helper.outputWriter = &buf

	pc.Run([]string{"-x", "abc123"})

	if !mock.called {
		t.Fatal("RunGit was not called")
	}
	if mock.gotName != "cherry-pick" {
		t.Errorf("got name %q, want cherry-pick", mock.gotName)
	}
	if !slices.Equal(mock.gotArgs, []string{"-x", "abc123"}) {
		t.Errorf("got args %v, want [-x abc123]", mock.gotArgs)
	}
}

func TestPassthroughCommand_Help(t *testing.T) {
	mock := &mockPassthroughClient{}
	var buf bytes.Buffer
	pc := newPassthroughCommand("cherry-pick", mock)
	pc.outputWriter = &buf
	pc.helper.outputWriter = &buf

	pc.Run([]string{"help"})

	if mock.called {
		t.Error("RunGit should not be called for help subcommand")
	}
	if !strings.Contains(buf.String(), "cherry-pick") {
		t.Errorf("expected help output to mention cherry-pick, got %q", buf.String())
	}
}

func TestPassthroughCommand_Error(t *testing.T) {
	mock := &mockPassthroughClient{err: errors.New("bad ref")}
	var buf bytes.Buffer
	pc := newPassthroughCommand("revert", mock)
	pc.outputWriter = &buf
	pc.helper.outputWriter = &buf

	pc.Run([]string{"deadbeef"})

	if !strings.Contains(buf.String(), "bad ref") {
		t.Errorf("expected error in output, got %q", buf.String())
	}
}

func TestBuildPassthroughs_CoversAllNames(t *testing.T) {
	mock := &mockPassthroughClient{}
	m := buildPassthroughs(mock)
	if len(m) != len(passthroughCommandNames) {
		t.Fatalf("buildPassthroughs returned %d entries, want %d", len(m), len(passthroughCommandNames))
	}
	for _, name := range passthroughCommandNames {
		pc, ok := m[name]
		if !ok {
			t.Errorf("missing passthrough for %q", name)
			continue
		}
		if pc.name != name {
			t.Errorf("passthrough[%q].name = %q", name, pc.name)
		}
	}
}
