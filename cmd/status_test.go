package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v6/git"
)

type mockStatusInfoReader struct {
	currentBranch        string
	upstreamBranch       string
	aheadBehindCount     string
	statusWithColor      string
	statusShortWithColor string
}

func (m *mockStatusInfoReader) GetCurrentBranch() (string, error) {
	if m.currentBranch == "" {
		return "main", nil
	}
	return m.currentBranch, nil
}
func (m *mockStatusInfoReader) GetUpstreamBranchName(string) (string, error) {
	if m.upstreamBranch == "" {
		return "origin/main", nil
	}
	return m.upstreamBranch, nil
}
func (m *mockStatusInfoReader) GetAheadBehindCount(string, string) (string, error) {
	if m.aheadBehindCount == "" {
		return "0 0", nil
	}
	return m.aheadBehindCount, nil
}
func (m *mockStatusInfoReader) StatusWithColor() (string, error) {
	return m.statusWithColor, nil
}
func (m *mockStatusInfoReader) StatusShortWithColor() (string, error) {
	return m.statusShortWithColor, nil
}

var _ git.StatusInfoReader = (*mockStatusInfoReader)(nil)

func TestStatuser_Constructor(t *testing.T) {
	mockClient := &mockStatusInfoReader{}
	statuser := NewStatuser(mockClient)

	if statuser == nil {
		t.Fatal("Expected NewStatuser to return a non-nil Statuser")
	}
	if statuser != nil && statuser.gitClient == nil {
		t.Error("Expected gitClient to be set")
	}
	if statuser != nil && statuser.outputWriter == nil {
		t.Error("Expected outputWriter to be set")
	}
	if statuser != nil && statuser.helper == nil {
		t.Error("Expected helper to be set")
	}
}

func TestStatuser_Status(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		shouldShowHelp bool
	}{
		{
			name:           "no args - show full status",
			args:           []string{},
			expectedOutput: "On branch main", // Should show branch info
			shouldShowHelp: false,
		},
		{
			name:           "short status",
			args:           []string{"short"},
			expectedOutput: "", // Mock client returns empty for StatusShortWithColor
			shouldShowHelp: false,
		},
		{
			name:           "invalid argument - should show help",
			args:           []string{"invalid"},
			expectedOutput: "Usage: ggc status [command]",
			shouldShowHelp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := &mockStatusInfoReader{}

			statuser := &Statuser{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			// Set helper's output writer to capture help output
			statuser.helper.outputWriter = buf

			statuser.Status(tt.args)

			output := buf.String()

			// Verify expected behavior
			if tt.shouldShowHelp {
				if !strings.Contains(output, tt.expectedOutput) {
					t.Errorf("Expected help output containing '%s', got: %s", tt.expectedOutput, output)
				}
			} else {
				if tt.expectedOutput == "" {
					// For short status, mock returns empty string - this is expected
					// We verify the command executed without error
					if strings.Contains(output, "Error:") {
						t.Errorf("Unexpected error in status operation: %s", output)
					}
				} else {
					// For full status, should show branch info
					if !strings.Contains(output, tt.expectedOutput) {
						t.Errorf("Expected output containing '%s', got: %s", tt.expectedOutput, output)
					}
				}
			}

			// Verify no panic occurred
			if t.Failed() {
				t.Logf("Command args: %v", tt.args)
				t.Logf("Full output: %s", output)
			}
		})
	}
}

func TestStatuser_StatusOperations(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		testFunc func(*testing.T, *Statuser, *bytes.Buffer)
	}{
		{
			name: "full status shows branch and upstream info",
			args: []string{},
			testFunc: func(t *testing.T, statuser *Statuser, buf *bytes.Buffer) {
				output := buf.String()
				// Should show current branch
				if !strings.Contains(output, "On branch") {
					t.Errorf("Expected branch info in full status, got: %s", output)
				}
				// Mock client returns "main" for GetCurrentBranch
				if !strings.Contains(output, "main") {
					t.Errorf("Expected branch name 'main' in output, got: %s", output)
				}
				// Should not show error
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in full status: %s", output)
				}
			},
		},
		{
			name: "short status calls StatusShortWithColor",
			args: []string{"short"},
			testFunc: func(t *testing.T, statuser *Statuser, buf *bytes.Buffer) {
				output := buf.String()
				// Mock client doesn't return errors for StatusShortWithColor
				if strings.Contains(output, "Error:") {
					t.Errorf("Unexpected error in short status: %s", output)
				}
				// Should not show branch info (that's for full status only)
				if strings.Contains(output, "On branch") {
					t.Errorf("Short status should not show branch info, got: %s", output)
				}
			},
		},
		{
			name: "invalid argument shows help",
			args: []string{"invalid"},
			testFunc: func(t *testing.T, statuser *Statuser, buf *bytes.Buffer) {
				output := buf.String()
				// Should show help
				if !strings.Contains(output, "Usage:") {
					t.Errorf("Expected help output for invalid argument, got: %s", output)
				}
				// Should not attempt status operations
				if strings.Contains(output, "On branch") {
					t.Errorf("Invalid argument should not show status, got: %s", output)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := &mockStatusInfoReader{}

			statuser := &Statuser{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			statuser.helper.outputWriter = buf

			statuser.Status(tt.args)
			tt.testFunc(t, statuser, buf)
		})
	}
}

func TestStatuser_UpstreamStatus(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T, *Statuser, *bytes.Buffer)
	}{
		{
			name: "full status includes upstream information",
			testFunc: func(t *testing.T, statuser *Statuser, buf *bytes.Buffer) {
				statuser.Status([]string{})
				output := buf.String()

				// Should show branch info
				if !strings.Contains(output, "On branch main") {
					t.Errorf("Expected branch info, got: %s", output)
				}

				// Mock client returns empty for upstream operations, so no upstream info expected
				// But should not show error
				if strings.Contains(output, "Error getting current branch:") {
					t.Errorf("Unexpected branch error: %s", output)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := &mockStatusInfoReader{}

			statuser := &Statuser{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}
			statuser.helper.outputWriter = buf

			tt.testFunc(t, statuser, buf)
		})
	}
}

func TestStatuser_FormatMethods(t *testing.T) {
	statuser := &Statuser{
		outputWriter: &bytes.Buffer{},
		helper:       NewHelper(),
	}

	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "formatUpToDate formats correctly",
			testFunc: func(t *testing.T) {
				result := statuser.formatUpToDate("origin/main")
				expected := "Your branch is up to date with 'origin/main'"
				if result != expected {
					t.Errorf("Expected %q, got %q", expected, result)
				}
			},
		},
		{
			name: "formatAheadBehind handles up-to-date case",
			testFunc: func(t *testing.T) {
				result := statuser.formatAheadBehind("origin/main", "0", "0")
				expected := "Your branch is up to date with 'origin/main'"
				if result != expected {
					t.Errorf("Expected %q, got %q", expected, result)
				}
			},
		},
		{
			name: "formatAheadBehind handles ahead case",
			testFunc: func(t *testing.T) {
				result := statuser.formatAheadBehind("origin/main", "2", "0")
				expected := "Your branch is ahead of 'origin/main' by 2 commit(s)"
				if result != expected {
					t.Errorf("Expected %q, got %q", expected, result)
				}
			},
		},
		{
			name: "formatAheadBehind handles behind case",
			testFunc: func(t *testing.T) {
				result := statuser.formatAheadBehind("origin/main", "0", "3")
				expected := "Your branch is behind 'origin/main' by 3 commit(s)"
				if result != expected {
					t.Errorf("Expected %q, got %q", expected, result)
				}
			},
		},
		{
			name: "formatAheadBehind handles diverged case",
			testFunc: func(t *testing.T) {
				result := statuser.formatAheadBehind("origin/main", "2", "3")
				expected := "Your branch and 'origin/main' have diverged,\nand have 2 and 3 different commits each, respectively"
				if result != expected {
					t.Errorf("Expected %q, got %q", expected, result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func TestParseCounts(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedAhead  string
		expectedBehind string
		expectedOK     bool
	}{
		{
			name:           "valid counts",
			input:          "2 3",
			expectedAhead:  "2",
			expectedBehind: "3",
			expectedOK:     true,
		},
		{
			name:           "zero counts",
			input:          "0 0",
			expectedAhead:  "0",
			expectedBehind: "0",
			expectedOK:     true,
		},
		{
			name:       "invalid format - single number",
			input:      "2",
			expectedOK: false,
		},
		{
			name:       "invalid format - empty",
			input:      "",
			expectedOK: false,
		},
		{
			name:       "invalid format - too many numbers",
			input:      "1 2 3",
			expectedOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ahead, behind, ok := parseCounts(tt.input)
			if ok != tt.expectedOK {
				t.Errorf("Expected ok=%v, got %v", tt.expectedOK, ok)
			}
			if tt.expectedOK {
				if ahead != tt.expectedAhead {
					t.Errorf("Expected ahead=%q, got %q", tt.expectedAhead, ahead)
				}
				if behind != tt.expectedBehind {
					t.Errorf("Expected behind=%q, got %q", tt.expectedBehind, behind)
				}
			}
		})
	}
}
