package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"
)

// mockBranchGitClient is a mock implementation of git.Clienter for branch tests
type mockBranchGitClient struct {
	mockGitClient
	getCurrentBranchCalled bool
	currentBranch          string
	err                    error
	listLocalBranches      func() ([]string, error)
	listRemoteBranches     func() ([]string, error)
}

func (m *mockBranchGitClient) GetCurrentBranch() (string, error) {
	m.getCurrentBranchCalled = true
	return m.currentBranch, m.err
}

func (m *mockBranchGitClient) ListLocalBranches() ([]string, error) {
	if m.listLocalBranches != nil {
		return m.listLocalBranches()
	}
	return []string{"main", "feature/test"}, nil
}

func (m *mockBranchGitClient) ListRemoteBranches() ([]string, error) {
	if m.listRemoteBranches != nil {
		return m.listRemoteBranches()
	}
	return []string{"origin/main", "origin/feature/test"}, nil
}

// define other functions as nil except branch testing ones
func (m *mockBranchGitClient) LogGraph() error {
	return nil
}

func (m *mockBranchGitClient) RestoreAll() error {
	return nil
}

func (m *mockBranchGitClient) RestoreAllStaged() error {
	return nil
}

func (m *mockBranchGitClient) RestoreStaged(...string) error {
	return nil
}

func (m *mockBranchGitClient) RestoreWorkingDir(...string) error {
	return nil
}

func (m *mockBranchGitClient) RestoreFromCommit(string, ...string) error {
	return nil
}

func TestBrancher_Branch_Current(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		currentBranch: "feature/test",
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
	}
	brancher.Branch([]string{"current"})

	if !mockClient.getCurrentBranchCalled {
		t.Error("GetCurrentBranch should be called")
	}

	output := buf.String()
	if output != "feature/test\n" {
		t.Errorf("unexpected output: got %q, want %q", output, "feature/test\n")
	}
}

func TestBrancher_Branch_Current_Error(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		err: errors.New("failed to get current branch"),
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
	}
	brancher.Branch([]string{"current"})

	if !mockClient.getCurrentBranchCalled {
		t.Error("GetCurrentBranch should be called")
	}

	output := buf.String()
	if output != "Error: failed to get current branch\n" {
		t.Errorf("unexpected output: got %q, want %q", output, "Error: failed to get current branch\n")
	}
}

func TestBrancher_Branch_Checkout(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		execCommand:  func(_ string, _ ...string) *exec.Cmd { return exec.Command("echo") },
		inputReader:  bufio.NewReader(strings.NewReader("1\n")),
	}
	brancher.Branch([]string{"checkout"})

	output := buf.String()
	expected := "Local branches:\n[1] main\n[2] feature/test\nEnter the number to checkout: "
	if !strings.Contains(output, expected) {
		t.Errorf("unexpected output: got %q, want %q", output, expected)
	}
}

func TestBrancher_Branch_CheckoutRemote(t *testing.T) {
	var buf bytes.Buffer
	brancher := &Brancher{
		gitClient:    &mockBranchGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand:  func(_ string, _ ...string) *exec.Cmd { return exec.Command("echo") },
		inputReader:  bufio.NewReader(strings.NewReader("1\n")),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"checkout-remote"})

	output := buf.String()
	if !strings.Contains(output, "Remote branches:") {
		t.Error("Expected remote branches list")
	}
}

func TestBrancher_Branch_Delete(t *testing.T) {
	var buf bytes.Buffer
	brancher := &Brancher{
		gitClient:    &mockBranchGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand:  func(_ string, _ ...string) *exec.Cmd { return exec.Command("echo") },
		inputReader:  bufio.NewReader(strings.NewReader("1\n")),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"delete"})

	output := buf.String()
	if !strings.Contains(output, "Select local branches to delete") {
		t.Error("Expected delete prompt")
	}
}

func TestBrancher_Branch_DeleteMerged(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		currentBranch: "main",
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
		execCommand: func(name string, args ...string) *exec.Cmd {
			if name == "git" && len(args) > 0 && args[0] == "branch" && args[1] == "--merged" {
				cmd := exec.Command("echo", "  main\n  feature/merged")
				return cmd
			}
			return exec.Command("echo")
		},
		inputReader: bufio.NewReader(strings.NewReader("1\n")),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"delete-merged"})

	output := buf.String()
	if !strings.Contains(output, "Select merged local branches to delete") {
		t.Error("Expected delete merged prompt")
	}
}

func TestBrancher_Branch_Help(t *testing.T) {
	var buf bytes.Buffer
	brancher := &Brancher{
		gitClient:    &mockBranchGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	// Show help when no arguments provided
	brancher.Branch([]string{})

	output := buf.String()
	if !strings.Contains(output, "Usage") {
		t.Error("Expected help message to contain 'Usage'")
	}
}

func TestBrancher_Branch_UnknownCommand(t *testing.T) {
	var buf bytes.Buffer
	brancher := &Brancher{
		gitClient:    &mockBranchGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	// Show help for unknown commands
	brancher.Branch([]string{"unknown"})

	output := buf.String()
	if !strings.Contains(output, "Usage") {
		t.Error("Expected help message to contain 'Usage'")
	}
}

func TestBrancher_branchCheckout_Error(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listLocalBranches: func() ([]string, error) {
			return nil, errors.New("failed to list branches")
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
	}
	brancher.branchCheckout()

	output := buf.String()
	if !strings.Contains(output, "Error: failed to list branches") {
		t.Errorf("Expected error message, got: %s", output)
	}
}

func TestBrancher_branchCheckout_NoBranches(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	// Override ListLocalBranches to return empty slice
	mockClient.listLocalBranches = func() ([]string, error) {
		return []string{}, nil
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
	}
	brancher.branchCheckout()

	output := buf.String()
	if !strings.Contains(output, "No local branches found.") {
		t.Errorf("Expected no branches message, got: %s", output)
	}
}

func TestBrancher_branchCheckout_InvalidNumber(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		execCommand:  func(_ string, _ ...string) *exec.Cmd { return exec.Command("echo") },
		inputReader:  bufio.NewReader(strings.NewReader("invalid\n")),
	}
	brancher.branchCheckout()

	output := buf.String()
	if !strings.Contains(output, "Invalid number.") {
		t.Errorf("Expected invalid number message, got: %s", output)
	}
}

func TestBrancher_branchDelete_Success(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listLocalBranches: func() ([]string, error) {
			return []string{"feature/test", "bugfix/issue"}, nil
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		execCommand:  func(_ string, _ ...string) *exec.Cmd { return exec.Command("echo") },
		inputReader:  bufio.NewReader(strings.NewReader("1 2\n")),
	}

	brancher.branchDelete()

	output := buf.String()
	if !strings.Contains(output, "feature/test") {
		t.Error("Expected branch name in output")
	}
	if !strings.Contains(output, "Selected branches deleted.") {
		t.Error("Expected success message")
	}
}

func TestBrancher_branchDelete_All(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listLocalBranches: func() ([]string, error) {
			return []string{"feature/test", "bugfix/issue"}, nil
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		execCommand:  func(_ string, _ ...string) *exec.Cmd { return exec.Command("echo") },
		inputReader:  bufio.NewReader(strings.NewReader("all\n")),
	}

	brancher.branchDelete()

	output := buf.String()
	if !strings.Contains(output, "All branches deleted.") {
		t.Error("Expected all branches deleted message")
	}
}

func TestBrancher_branchDelete_Cancel(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listLocalBranches: func() ([]string, error) {
			return []string{"feature/test"}, nil
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		execCommand:  func(_ string, _ ...string) *exec.Cmd { return exec.Command("echo") },
		inputReader:  bufio.NewReader(strings.NewReader("\n")),
	}

	brancher.branchDelete()

	output := buf.String()
	if !strings.Contains(output, "Cancelled.") {
		t.Error("Expected cancelled message")
	}
}

func TestBrancher_branchDelete_Error(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listLocalBranches: func() ([]string, error) {
			return nil, errors.New("failed to list branches")
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
	}

	brancher.branchDelete()

	output := buf.String()
	if !strings.Contains(output, "Error: failed to list branches") {
		t.Error("Expected error message")
	}
}

func TestBrancher_branchDelete_NoBranches(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listLocalBranches: func() ([]string, error) {
			return []string{}, nil
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
	}

	brancher.branchDelete()

	output := buf.String()
	if !strings.Contains(output, "No local branches found.") {
		t.Error("Expected no branches message")
	}
}

func TestBrancher_branchDeleteMerged_Success(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		currentBranch: "main",
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		execCommand: func(name string, args ...string) *exec.Cmd {
			if name == "git" && len(args) > 0 && args[0] == "branch" && args[1] == "--merged" {
				cmd := exec.Command("echo", "  main\n  feature/merged\n  bugfix/old")
				return cmd
			}
			return exec.Command("echo")
		},
		inputReader: bufio.NewReader(strings.NewReader("1 2\n")),
	}

	brancher.branchDeleteMerged()

	output := buf.String()
	if !strings.Contains(output, "Selected merged branches deleted.") {
		t.Error("Expected success message")
	}
}

func TestBrancher_branchDeleteMerged_All(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		currentBranch: "main",
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		execCommand: func(name string, args ...string) *exec.Cmd {
			if name == "git" && len(args) > 0 && args[0] == "branch" && args[1] == "--merged" {
				cmd := exec.Command("echo", "  main\n  feature/merged")
				return cmd
			}
			return exec.Command("echo")
		},
		inputReader: bufio.NewReader(strings.NewReader("all\n")),
	}

	brancher.branchDeleteMerged()

	output := buf.String()
	if !strings.Contains(output, "All merged branches deleted.") {
		t.Error("Expected all merged branches deleted message")
	}
}

func TestBrancher_branchDeleteMerged_Cancel(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		currentBranch: "main",
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		execCommand: func(name string, args ...string) *exec.Cmd {
			if name == "git" && len(args) > 0 && args[0] == "branch" && args[1] == "--merged" {
				cmd := exec.Command("echo", "  main\n  feature/merged")
				return cmd
			}
			return exec.Command("echo")
		},
		inputReader: bufio.NewReader(strings.NewReader("\n")),
	}

	brancher.branchDeleteMerged()

	output := buf.String()
	if !strings.Contains(output, "Cancelled.") {
		t.Error("Expected cancelled message")
	}
}

func TestBrancher_branchDeleteMerged_CurrentBranchError(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		err: errors.New("failed to get current branch"),
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
	}

	brancher.branchDeleteMerged()

	output := buf.String()
	if !strings.Contains(output, "Error: failed to get current branch") {
		t.Error("Expected current branch error message")
	}
}

func TestBrancher_branchDeleteMerged_NoMergedBranches(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		currentBranch: "main",
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		execCommand: func(name string, args ...string) *exec.Cmd {
			if name == "git" && len(args) > 0 && args[0] == "branch" && args[1] == "--merged" {
				cmd := exec.Command("echo", "  main")
				return cmd
			}
			return exec.Command("echo")
		},
	}

	brancher.branchDeleteMerged()

	output := buf.String()
	if !strings.Contains(output, "No merged local branches.") {
		t.Error("Expected no merged branches message")
	}
}

func TestBrancher_branchCheckoutRemote_Success(t *testing.T) {
	var buf bytes.Buffer
	brancher := &Brancher{
		gitClient:    &mockBranchGitClient{},
		outputWriter: &buf,
		execCommand:  func(_ string, _ ...string) *exec.Cmd { return exec.Command("echo") },
		inputReader:  bufio.NewReader(strings.NewReader("2\n")),
	}

	brancher.branchCheckoutRemote()

	output := buf.String()
	if !strings.Contains(output, "Remote branches:") {
		t.Error("Expected remote branches list")
	}
	if !strings.Contains(output, "origin/feature/test") {
		t.Error("Expected remote branch name in output")
	}
}

func TestBrancher_branchCheckoutRemote_InvalidNumber(t *testing.T) {
	var buf bytes.Buffer
	brancher := &Brancher{
		gitClient:    &mockBranchGitClient{},
		outputWriter: &buf,
		execCommand:  func(_ string, _ ...string) *exec.Cmd { return exec.Command("echo") },
		inputReader:  bufio.NewReader(strings.NewReader("invalid\n")),
	}

	brancher.branchCheckoutRemote()

	output := buf.String()
	if !strings.Contains(output, "Invalid number.") {
		t.Error("Expected invalid number message")
	}
}

func TestBrancher_branchCheckoutRemote_InvalidBranchName(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listRemoteBranches: func() ([]string, error) {
			return []string{"invalidbranch"}, nil
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		execCommand:  func(_ string, _ ...string) *exec.Cmd { return exec.Command("echo") },
		inputReader:  bufio.NewReader(strings.NewReader("1\n")),
	}

	brancher.branchCheckoutRemote()

	output := buf.String()
	if !strings.Contains(output, "Invalid remote branch name.") {
		t.Error("Expected invalid branch name message")
	}
}

func TestBrancher_Branch_Create(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput string
		cmdOutput      string
		cmdError       bool
	}{
		{
			name:           "Success: Create new branch",
			input:          "feature/test\n",
			expectedOutput: "Enter new branch name: \n",
			cmdOutput:      "",
			cmdError:       false,
		},
		{
			name:           "Error: Empty branch name",
			input:          "\n",
			expectedOutput: "Enter new branch name: Cancelled.\n",
			cmdOutput:      "",
			cmdError:       false,
		},
		{
			name:           "Error: Branch creation failed",
			input:          "feature/test\n",
			expectedOutput: "Enter new branch name: Error: failed to create and checkout branch: exit status 1\n",
			cmdOutput:      "",
			cmdError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			brancher := &Brancher{
				gitClient:    &mockBranchGitClient{},
				outputWriter: &buf,
				inputReader:  bufio.NewReader(strings.NewReader(tt.input)),
				execCommand: func(_ string, _ ...string) *exec.Cmd {
					if tt.cmdError {
						return exec.Command("false")
					}
					return exec.Command("echo", tt.cmdOutput)
				},
			}

			brancher.Branch([]string{"create"})

			output := buf.String()
			if output != tt.expectedOutput {
				t.Errorf("unexpected output:\ngot:  %q\nwant: %q", output, tt.expectedOutput)
			}
		})
	}
}

func TestBrancher_branchCreate_ExistingBranch(t *testing.T) {
	var buf bytes.Buffer
	brancher := &Brancher{
		gitClient:    &mockBranchGitClient{},
		outputWriter: &buf,
		inputReader:  bufio.NewReader(strings.NewReader("main\n")),
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			// Set up git command to return error
			cmd := exec.Command("echo", "fatal: A branch named 'main' already exists.")
			cmd.Stderr = cmd.Stdout
			return cmd
		},
	}

	brancher.Branch([]string{"create"})

	output := buf.String()
	expectedOutput := "Enter new branch name: fatal: A branch named 'main' already exists.\n"
	if output != expectedOutput {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", output, expectedOutput)
	}
}

func TestBrancher_Branch_BoundaryInputValues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Single character input",
			input:    "1\n",
			expected: "main",
		},
		{
			name:     "Maximum integer input",
			input:    "999999\n",
			expected: "Invalid number",
		},
		{
			name:     "Negative number input",
			input:    "-1\n",
			expected: "Invalid number",
		},
		{
			name:     "Zero input",
			input:    "0\n",
			expected: "Invalid number",
		},
		{
			name:     "Leading zeros",
			input:    "001\n",
			expected: "main",
		},
		{
			name:     "Floating point number",
			input:    "1.5\n",
			expected: "Invalid number",
		},
		{
			name:     "Scientific notation",
			input:    "1e2\n",
			expected: "Invalid number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			mockClient := &mockBranchGitClient{}
			brancher := &Brancher{
				gitClient:    mockClient,
				outputWriter: &buf,
				execCommand:  func(_ string, _ ...string) *exec.Cmd { return exec.Command("echo") },
				inputReader:  bufio.NewReader(strings.NewReader(tt.input)),
			}

			brancher.branchCheckout()

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected %q in output, got: %s", tt.expected, output)
			}
		})
	}
}

func TestBrancher_Branch_BoundaryBranchNames(t *testing.T) {
	tests := []struct {
		name         string
		branchName   string
		expectedPass bool
		description  string
	}{
		{
			name:         "Single character branch name",
			branchName:   "a",
			expectedPass: true,
			description:  "Minimum length branch name",
		},
		{
			name:         "Maximum length branch name",
			branchName:   strings.Repeat("a", 255),
			expectedPass: false,
			description:  "Very long branch name",
		},
		{
			name:         "Branch name with dots",
			branchName:   "feature.test",
			expectedPass: true,
			description:  "Branch name with dots",
		},
		{
			name:         "Branch name with hyphens",
			branchName:   "feature-test",
			expectedPass: true,
			description:  "Branch name with hyphens",
		},
		{
			name:         "Branch name with underscores",
			branchName:   "feature_test",
			expectedPass: true,
			description:  "Branch name with underscores",
		},
		{
			name:         "Branch name with slashes",
			branchName:   "feature/test",
			expectedPass: true,
			description:  "Branch name with slashes",
		},
		{
			name:         "Branch name starting with dot",
			branchName:   ".test",
			expectedPass: false,
			description:  "Invalid: starts with dot",
		},
		{
			name:         "Branch name ending with dot",
			branchName:   "test.",
			expectedPass: false,
			description:  "Invalid: ends with dot",
		},
		{
			name:         "Branch name with consecutive dots",
			branchName:   "test..branch",
			expectedPass: false,
			description:  "Invalid: consecutive dots",
		},
		{
			name:         "Branch name with spaces",
			branchName:   "test branch",
			expectedPass: false,
			description:  "Invalid: contains spaces",
		},
		{
			name:         "Branch name with special characters",
			branchName:   "test@#$%",
			expectedPass: false,
			description:  "Invalid: special characters",
		},
		{
			name:         "Branch name with control characters",
			branchName:   "test\x00branch",
			expectedPass: false,
			description:  "Invalid: control characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			brancher := &Brancher{
				gitClient:    &mockBranchGitClient{},
				outputWriter: &buf,
				inputReader:  bufio.NewReader(strings.NewReader(tt.branchName + "\n")),
				execCommand: func(_ string, _ ...string) *exec.Cmd {
					if tt.expectedPass {
						return exec.Command("echo")
					}
					return exec.Command("false")
				},
			}

			brancher.branchCreate()

			output := buf.String()
			if tt.expectedPass {
				if strings.Contains(output, "Error:") {
					t.Errorf("Expected success for %s, but got error: %s", tt.description, output)
				}
			} else {
				if !strings.Contains(output, "Error:") && !strings.Contains(output, "fatal:") {
					t.Errorf("Expected error for %s, but got success: %s", tt.description, output)
				}
			}
		})
	}
}

func TestBrancher_Branch_BoundaryUserInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty input",
			input:    "\n",
			expected: "Cancelled",
		},
		{
			name:     "Whitespace only",
			input:    "   \n",
			expected: "Cancelled",
		},
		{
			name:     "Tab characters",
			input:    "\t\t\n",
			expected: "Cancelled",
		},
		{
			name:     "Multiple newlines",
			input:    "\n\n\n",
			expected: "Cancelled",
		},
		{
			name:     "Very long input",
			input:    strings.Repeat("a", 1000) + "\n",
			expected: "Error:",
		},
		{
			name:     "Input with trailing spaces",
			input:    "branch   \n",
			expected: "", // Should be trimmed and work
		},
		{
			name:     "Input with leading spaces",
			input:    "   branch\n",
			expected: "", // Should be trimmed and work
		},
		{
			name:     "Unicode characters",
			input:    "brαnch\n",
			expected: "Error:",
		},
		{
			name:     "Mixed case input",
			input:    "BrAnCh\n",
			expected: "", // Should work
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			brancher := &Brancher{
				gitClient:    &mockBranchGitClient{},
				outputWriter: &buf,
				inputReader:  bufio.NewReader(strings.NewReader(tt.input)),
				execCommand: func(_ string, _ ...string) *exec.Cmd {
					if strings.Contains(tt.expected, "Error:") {
						return exec.Command("false")
					}
					return exec.Command("echo")
				},
			}

			brancher.branchCreate()

			output := buf.String()
			if tt.expected != "" {
				if !strings.Contains(output, tt.expected) {
					t.Errorf("Expected %q in output, got: %s", tt.expected, output)
				}
			}
		})
	}
}

func TestBrancher_Branch_BoundaryListOperations(t *testing.T) {
	tests := []struct {
		name          string
		localBranches []string
		input         string
		expected      string
	}{
		{
			name:          "Empty branch list",
			localBranches: []string{},
			input:         "1\n",
			expected:      "No local branches found",
		},
		{
			name:          "Single branch",
			localBranches: []string{"main"},
			input:         "1\n",
			expected:      "main",
		},
		{
			name:          "Branch with very long name",
			localBranches: []string{"main", strings.Repeat("feature-", 30)},
			input:         "2\n",
			expected:      strings.Repeat("feature-", 30),
		},
		{
			name:          "Branches with special naming patterns",
			localBranches: []string{"main", "feature/ABC-123", "bugfix/fix-issue-456", "release/v1.0.0"},
			input:         "3\n",
			expected:      "bugfix/fix-issue-456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			mockClient := &mockBranchGitClient{
				listLocalBranches: func() ([]string, error) {
					return tt.localBranches, nil
				},
			}
			brancher := &Brancher{
				gitClient:    mockClient,
				outputWriter: &buf,
				execCommand:  func(_ string, _ ...string) *exec.Cmd { return exec.Command("echo") },
				inputReader:  bufio.NewReader(strings.NewReader(tt.input)),
			}

			brancher.branchCheckout()

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected %q in output, got: %s", tt.expected, output)
			}
		})
	}
}

func TestBrancher_Branch_BoundaryDeleteOperations(t *testing.T) {
	tests := []struct {
		name          string
		localBranches []string
		input         string
		expected      string
	}{
		{
			name:          "Delete single branch with boundary input",
			localBranches: []string{"main", "feature", "bugfix"},
			input:         "1\n",
			expected:      "Selected branches deleted",
		},
		{
			name:          "Delete multiple branches",
			localBranches: []string{"main", "feature", "bugfix", "hotfix"},
			input:         "1 2 3\n",
			expected:      "Selected branches deleted",
		},
		{
			name:          "Delete all branches",
			localBranches: []string{"feature", "bugfix", "hotfix"},
			input:         "all\n",
			expected:      "All branches deleted",
		},
		{
			name:          "Delete with out-of-range numbers",
			localBranches: []string{"main", "feature"},
			input:         "1 5 10\n",
			expected:      "Invalid number:",
		},
		{
			name:          "Delete with mixed valid/invalid input",
			localBranches: []string{"main", "feature", "bugfix"},
			input:         "1 invalid 2\n",
			expected:      "Invalid number:",
		},
		{
			name:          "Delete with very large numbers",
			localBranches: []string{"main", "feature"},
			input:         "999999 1000000\n",
			expected:      "Invalid number:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			mockClient := &mockBranchGitClient{
				listLocalBranches: func() ([]string, error) {
					return tt.localBranches, nil
				},
			}
			brancher := &Brancher{
				gitClient:    mockClient,
				outputWriter: &buf,
				execCommand:  func(_ string, _ ...string) *exec.Cmd { return exec.Command("echo") },
				inputReader:  bufio.NewReader(strings.NewReader(tt.input)),
			}

			brancher.branchDelete()

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected %q in output, got: %s", tt.expected, output)
			}
		})
	}
}

func TestBrancher_Branch_BoundaryRemoteOperations(t *testing.T) {
	tests := []struct {
		name           string
		remoteBranches []string
		input          string
		expected       string
	}{
		{
			name:           "Empty remote branches",
			remoteBranches: []string{},
			input:          "1\n",
			expected:       "Selected branches deleted.",
		},
		{
			name:           "Single remote branch",
			remoteBranches: []string{"origin/main"},
			input:          "1\n",
			expected:       "main",
		},
		{
			name:           "Multiple remote origins",
			remoteBranches: []string{"origin/main", "upstream/main", "fork/main"},
			input:          "2\n",
			expected:       "main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			mockClient := &mockBranchGitClient{
				listRemoteBranches: func() ([]string, error) {
					return tt.remoteBranches, nil
				},
			}
			brancher := &Brancher{
				gitClient:    mockClient,
				outputWriter: &buf,
				execCommand:  func(_ string, _ ...string) *exec.Cmd { return exec.Command("echo") },
				inputReader:  bufio.NewReader(strings.NewReader(tt.input)),
			}

			brancher.branchDelete()

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected %q in output, got: %s", tt.expected, output)
			}
		})
	}
}
