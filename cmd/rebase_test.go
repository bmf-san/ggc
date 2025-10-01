package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v7/internal/prompt"
)

func TestRebaser_RebaseInteractive_SelectValid(t *testing.T) {
	var buf bytes.Buffer

	mockClient := &mockAddGitClient{}
	r := &Rebaser{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
		prompter:     prompt.New(strings.NewReader("2\n"), &buf),
	}
	r.helper.outputWriter = &buf

	r.RebaseInteractive()

	output := buf.String()
	if !strings.Contains(output, "Current branch: feature/test") {
		t.Errorf("expected branch name, got: %s", output)
	}
	if !strings.Contains(output, "Select number of commits to rebase") {
		t.Errorf("expected prompt, got: %s", output)
	}
	if !strings.Contains(output, "Rebase successful") {
		t.Errorf("expected success message, got: %s", output)
	}
}

func TestRebaser_RebaseInteractive_BranchError(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockAddGitClient{}
	// Override GetCurrentBranch to return error
	mockClient.GetCurrentBranchFunc = func() (string, error) {
		return "", errors.New("failed to get current branch")
	}
	r := &Rebaser{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
		prompter:     prompt.New(strings.NewReader("1\n"), &buf),
	}
	r.helper.outputWriter = &buf

	r.RebaseInteractive()

	output := buf.String()
	if !strings.Contains(output, "Error: failed to get current branch") {
		t.Errorf("expected branch error message, got: %s", output)
	}
}

func TestRebaser_RebaseInteractive_Cancel(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockAddGitClient{}
	r := &Rebaser{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
		prompter:     prompt.New(strings.NewReader("\n"), &buf),
	}
	r.helper.outputWriter = &buf

	r.RebaseInteractive()

	output := buf.String()
	if !strings.Contains(output, "Error: operation canceled") {
		t.Errorf("expected cancellation message, got: %s", output)
	}
}

func TestRebaser_RebaseInteractive_InvalidNumber(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockAddGitClient{}
	r := &Rebaser{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
		prompter:     prompt.New(strings.NewReader("abc\n"), &buf),
	}
	r.helper.outputWriter = &buf

	r.RebaseInteractive()

	output := buf.String()
	if !strings.Contains(output, "Error: invalid number") {
		t.Errorf("expected invalid number message, got: %s", output)
	}
}

func TestRebaser_RebaseInteractive_NoHistory(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockAddGitClient{}
	// Override LogOneline to return empty history
	mockClient.LogOnelineFunc = func(_, _ string) (string, error) {
		return "", nil
	}
	r := &Rebaser{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
		prompter:     prompt.New(strings.NewReader("1\n"), &buf),
	}
	r.helper.outputWriter = &buf

	r.RebaseInteractive()

	output := buf.String()
	if !strings.Contains(output, "Error: no commit history found") {
		t.Errorf("expected no history message, got: %s", output)
	}
}

func TestRebaser_RebaseInteractive_LogError(t *testing.T) {
	var buf bytes.Buffer

	mockClient := &mockAddGitClient{}
	// Override LogOneline to return error
	mockClient.LogOnelineFunc = func(_, _ string) (string, error) {
		return "", errors.New("failed to get git log")
	}
	r := &Rebaser{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
		prompter:     prompt.New(strings.NewReader("1\n"), &buf),
	}
	r.helper.outputWriter = &buf

	r.RebaseInteractive()

	output := buf.String()
	if !strings.Contains(output, "Error: failed to get git log") {
		t.Errorf("expected log error message, got: %s", output)
	}
}

func TestRebaser_RebaseInteractive_RebaseError(t *testing.T) {
	var buf bytes.Buffer

	mockClient := &mockAddGitClient{}
	// Override RebaseInteractive to return error
	mockClient.RebaseInteractiveFunc = func(_ int) error {
		return errors.New("rebase failed")
	}
	r := &Rebaser{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
		prompter:     prompt.New(strings.NewReader("1\n"), &buf),
	}
	r.helper.outputWriter = &buf

	r.RebaseInteractive()

	output := buf.String()
	if !strings.Contains(output, "Error: rebase failed") {
		t.Errorf("expected rebase error message, got: %s", output)
	}
}

func TestRebaser_Rebase(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		expectedOutput string
		mockInput      string
	}{
		{
			name:           "interactive rebase",
			args:           []string{"interactive"},
			expectedOutput: "Current branch:",
			mockInput:      "1\n",
		},
		{
			name:           "no args",
			args:           []string{},
			expectedOutput: "Usage: ggc rebase [interactive | <upstream> | continue | abort | skip]",
		},
		{
			name:           "invalid ref errors",
			args:           []string{"invalid"},
			expectedOutput: "Error: unknown ref",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			mockClient := &mockAddGitClient{}
			if tc.name == "invalid ref errors" {
				mockClient.RevParseVerifyFunc = func(_ string) bool { return false }
			}
			r := &Rebaser{
				gitClient:    mockClient,
				outputWriter: &buf,
				helper:       NewHelper(),
				prompter:     prompt.New(strings.NewReader(tc.mockInput), &buf),
			}
			r.helper.outputWriter = &buf

			r.Rebase(tc.args)

			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}

func TestRebaser_Rebase_Subcommands(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockAddGitClient{}
	r := &Rebaser{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
		prompter:     prompt.New(strings.NewReader(""), &buf),
	}
	r.helper.outputWriter = &buf

	// continue
	buf.Reset()
	mockClient.RebaseContinueCalled = false
	r.Rebase([]string{"continue"})
	if !mockClient.RebaseContinueCalled {
		t.Errorf("expected RebaseContinue to be called")
	}
	if !strings.Contains(buf.String(), "Rebase successful") {
		t.Errorf("expected success message, got: %s", buf.String())
	}

	// abort
	buf.Reset()
	mockClient.RebaseAbortCalled = false
	r.Rebase([]string{"abort"})
	if !mockClient.RebaseAbortCalled {
		t.Errorf("expected RebaseAbort to be called")
	}
	if !strings.Contains(buf.String(), "Rebase aborted") {
		t.Errorf("expected abort message, got: %s", buf.String())
	}

	// skip
	buf.Reset()
	mockClient.RebaseSkipCalled = false
	r.Rebase([]string{"skip"})
	if !mockClient.RebaseSkipCalled {
		t.Errorf("expected RebaseSkip to be called")
	}
	if !strings.Contains(buf.String(), "Rebase successful") {
		t.Errorf("expected success message, got: %s", buf.String())
	}
}

func TestRebaser_Rebase_BasicOnto(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockAddGitClient{}
	mockClient.RevParseVerifyFunc = func(ref string) bool {
		return ref == "main" || ref == "origin/main"
	}
	r := &Rebaser{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
		prompter:     prompt.New(strings.NewReader(""), &buf),
	}
	r.helper.outputWriter = &buf

	buf.Reset()
	r.Rebase([]string{"main"})
	if !mockClient.RebaseCalled || mockClient.RebaseUpstream != "main" {
		t.Errorf("expected Rebase to be called with 'main', got called=%v upstream=%q", mockClient.RebaseCalled, mockClient.RebaseUpstream)
	}
	if !strings.Contains(buf.String(), "Rebase successful") {
		t.Errorf("expected success message, got: %s", buf.String())
	}
}

func TestRebaser_Rebase_InteractiveCancel(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockAddGitClient{}
	// Override LogOneline to return empty history
	mockClient.LogOnelineFunc = func(_, _ string) (string, error) {
		return "", nil
	}
	r := &Rebaser{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
		prompter:     prompt.New(strings.NewReader("\n"), &buf),
	}
	r.helper.outputWriter = &buf

	r.Rebase([]string{"interactive"})

	output := buf.String()
	if !strings.Contains(output, "Error: no commit history found") {
		t.Errorf("expected no history message, got %q", output)
	}
}

func TestRebaser_Rebase_Error(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockAddGitClient{}
	r := &Rebaser{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
		prompter:     prompt.New(strings.NewReader("y\n"), &buf),
	}
	r.helper.outputWriter = &buf

	r.Rebase([]string{"interactive"})

	output := buf.String()
	if !strings.Contains(output, "Error:") {
		t.Errorf("expected error message, got %q", output)
	}
}
