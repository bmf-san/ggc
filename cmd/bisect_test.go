package cmd

import (
	"bytes"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v8/internal/testutil"
)

type mockBisectClient struct {
	testutil.MockGitClient
	called  bool
	gotName string
	gotArgs []string
	err     error
}

func (m *mockBisectClient) RunGit(name string, args []string) error {
	m.called = true
	m.gotName = name
	m.gotArgs = slices.Clone(args)
	return m.err
}

func TestBisector_Bisect_Help(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBisectClient{}
	b := NewBisector(mockClient)
	b.outputWriter = &buf
	b.helper.outputWriter = &buf

	b.Bisect(nil)

	if mockClient.called {
		t.Fatal("RunGit should not be called for help")
	}
	if !strings.Contains(buf.String(), "ggc bisect") {
		t.Fatalf("expected bisect help output, got: %q", buf.String())
	}
}

func TestBisector_Bisect_Start_RequiresBadAndGood(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBisectClient{}
	b := NewBisector(mockClient)
	b.outputWriter = &buf

	b.Bisect([]string{"start", "bad-ref-only"})

	if mockClient.called {
		t.Fatal("RunGit should not be called when args are invalid")
	}
	if !strings.Contains(buf.String(), "Usage: ggc bisect start <bad> <good>") {
		t.Fatalf("expected usage output, got: %q", buf.String())
	}
}

func TestBisector_Bisect_Start_ForwardsToGit(t *testing.T) {
	mockClient := &mockBisectClient{}
	b := NewBisector(mockClient)
	b.outputWriter = &bytes.Buffer{}

	b.Bisect([]string{"start", "deadbeef", "v1.0.0"})

	if !mockClient.called {
		t.Fatal("expected RunGit to be called")
	}
	if mockClient.gotName != "bisect" {
		t.Fatalf("expected bisect command, got %q", mockClient.gotName)
	}
	if !slices.Equal(mockClient.gotArgs, []string{"start", "deadbeef", "v1.0.0"}) {
		t.Fatalf("unexpected args: %v", mockClient.gotArgs)
	}
}

func TestBisector_Bisect_Run_RequiresScriptOrCommand(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBisectClient{}
	b := NewBisector(mockClient)
	b.outputWriter = &buf

	b.Bisect([]string{"run"})

	if mockClient.called {
		t.Fatal("RunGit should not be called when run target is missing")
	}
	if !strings.Contains(buf.String(), "Usage: ggc bisect run <script-or-command>") {
		t.Fatalf("expected usage output, got: %q", buf.String())
	}
}

func TestBisector_Bisect_Run_ForwardsToGit(t *testing.T) {
	mockClient := &mockBisectClient{}
	b := NewBisector(mockClient)
	b.outputWriter = &bytes.Buffer{}

	b.Bisect([]string{"run", "./scripts/test.sh", "--fast"})

	if !mockClient.called {
		t.Fatal("expected RunGit to be called")
	}
	if mockClient.gotName != "bisect" {
		t.Fatalf("expected bisect command, got %q", mockClient.gotName)
	}
	if !slices.Equal(mockClient.gotArgs, []string{"run", "./scripts/test.sh", "--fast"}) {
		t.Fatalf("unexpected args: %v", mockClient.gotArgs)
	}
}

func TestBisector_Bisect_FallbackSubcommand_ForwardsToGit(t *testing.T) {
	mockClient := &mockBisectClient{}
	b := NewBisector(mockClient)
	b.outputWriter = &bytes.Buffer{}

	b.Bisect([]string{"bad"})

	if !mockClient.called {
		t.Fatal("expected RunGit to be called")
	}
	if mockClient.gotName != "bisect" {
		t.Fatalf("expected bisect command, got %q", mockClient.gotName)
	}
	if !slices.Equal(mockClient.gotArgs, []string{"bad"}) {
		t.Fatalf("unexpected args: %v", mockClient.gotArgs)
	}
}

func TestBisector_Bisect_ForwardError(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBisectClient{err: errors.New("bisect failed")}
	b := NewBisector(mockClient)
	b.outputWriter = &buf

	b.Bisect([]string{"bad"})

	if !strings.Contains(buf.String(), "Error: bisect failed") {
		t.Fatalf("expected error output, got: %q", buf.String())
	}
}
