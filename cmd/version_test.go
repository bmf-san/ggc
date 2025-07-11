package cmd

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestVersioneer_Version(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		versionGetter  VersionGetter
		expectedOutput []string
	}{
		{
			name: "version no args with default values",
			args: []string{},
			versionGetter: nil,
			expectedOutput: []string{
				"ggc version dev",
				"commit: none",
				"built: unknown",
				"os/arch:",
			},
		},
		{
			name: "version no args with custom version info",
			args: []string{},
			versionGetter: func() (version, commit, date string) {
				return "v1.0.0", "abc123", "2024-01-01"
			},
			expectedOutput: []string{
				"ggc version v1.0.0",
				"commit: abc123",
				"built: 2024-01-01",
				"os/arch:",
			},
		},
		{
			name: "version with args shows help",
			args: []string{"help"},
			versionGetter: func() (version, commit, date string) {
				return "v1.0.0", "abc123", "2024-01-01"
			},
			expectedOutput: []string{
				"Usage:",
			},
		},
		{
			name: "version with multiple args shows help",
			args: []string{"invalid", "args"},
			versionGetter: nil,
			expectedOutput: []string{
				"Usage:",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			originalGetter := getVersionInfo
			SetVersionGetter(tc.versionGetter)
			defer SetVersionGetter(originalGetter)
			
			v := &Versioneer{
				outputWriter: &buf,
				helper:       NewHelper(),
				execCommand:  exec.Command,
			}
			
			v.helper.outputWriter = &buf
			v.Version(tc.args)
			
			output := buf.String()
			for _, expected := range tc.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("expected output to contain %q, got %q", expected, output)
				}
			}
		})
	}
}

// Helper test to verify version getter functionality
func TestSetVersionGetter(t *testing.T) {
	originalGetter := getVersionInfo
	defer SetVersionGetter(originalGetter)
	
	customGetter := func() (version, commit, date string) {
		return "test-version", "test-commit", "test-date"
	}
	
	SetVersionGetter(customGetter)
	
	if getVersionInfo == nil {
		t.Error("expected getVersionInfo to be set")
	}
	
	version, commit, date := getVersionInfo()
	if version != "test-version" || commit != "test-commit" || date != "test-date" {
		t.Errorf("expected version info (test-version, test-commit, test-date), got (%s, %s, %s)", 
			version, commit, date)
	}
}
