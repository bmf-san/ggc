package cmd

import (
	"bytes"
	"testing"

	"github.com/bmf-san/ggc/v5/internal/testutil"
)

func TestStatuser_Constructor(t *testing.T) {
	mockClient := testutil.NewMockGitClient()
	statuser := NewStatuser(mockClient)

	if statuser == nil {
		t.Fatal("Expected NewStatuser to return a non-nil Statuser")
	}
	if statuser.gitClient == nil {
		t.Error("Expected gitClient to be set")
	}
	if statuser.outputWriter == nil {
		t.Error("Expected outputWriter to be set")
	}
	if statuser.helper == nil {
		t.Error("Expected helper to be set")
	}
}

func TestStatuser_Status(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "no args - show full status",
			args: []string{},
		},
		{
			name: "short status",
			args: []string{"short"},
		},
		{
			name: "invalid argument",
			args: []string{"invalid"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			mockClient := testutil.NewMockGitClient()

			statuser := &Statuser{
				gitClient:    mockClient,
				outputWriter: buf,
				helper:       NewHelper(),
			}

			statuser.Status(tt.args)

			// Verify that the function executed without panic and produced output
			output := buf.String()
			
			// Status commands should always produce some output
			if len(output) == 0 {
				t.Errorf("Expected status output for args %v, got empty string", tt.args)
			}
			
			// Verify output content based on arguments
			switch {
			case len(tt.args) == 0:
				// Default status should show repository status information
				if len(output) < 10 { // Minimum reasonable status output
					t.Errorf("Expected detailed status output, got: %s", output)
				}
			case len(tt.args) > 0 && tt.args[0] == "short":
				// Short status should be more concise but still informative
				if len(output) == 0 {
					t.Errorf("Expected short status output, got empty string")
				}
			case len(tt.args) > 0 && tt.args[0] == "invalid":
				// Invalid arguments should produce error or help output
				if len(output) < 5 {
					t.Errorf("Expected error or help output for invalid args, got: %s", output)
				}
			}
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
