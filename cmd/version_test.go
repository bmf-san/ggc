package cmd

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v7/config"
	"github.com/bmf-san/ggc/v7/internal/testutil"
)

func TestVersioner_Version(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		expectedOutput []string
	}{
		{
			name: "version no args with default values",
			args: []string{},
			expectedOutput: []string{
				"ggc version",
				"commit:",
				"built:",
				"os/arch:",
			},
		},
		{
			name: "version no args with custom version info",
			args: []string{},
			expectedOutput: []string{
				"ggc version",
				"commit:",
				"built:",
				"os/arch:",
			},
		},
		{
			name: "version with args shows help",
			args: []string{"help"},
			expectedOutput: []string{
				"Usage:",
			},
		},
		{
			name: "version with multiple args shows help",
			args: []string{"invalid", "args"},
			expectedOutput: []string{
				"Usage:",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock getVersionInfo function
			SetVersionGetter(func() (string, string) {
				return "v1.0.0", "abc123"
			})

			var buf bytes.Buffer
			v := &Versioner{
				gitClient:    testutil.NewMockGitClient(),
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

// Tests for getVersionString and getCommitString functions (edge cases)
func TestVersioner_GetVersionString_EdgeCases(t *testing.T) {
	versioner := &Versioner{}

	tests := []struct {
		input    string
		expected string
	}{
		{"v1.0.0", "v1.0.0"},
		{"dev", "dev"},
		{"", "(devel)"},
	}

	for _, tt := range tests {
		result := versioner.getVersionString(tt.input)
		if result != tt.expected {
			t.Errorf("getVersionString(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestVersioner_GetCommitString_EdgeCases(t *testing.T) {
	versioner := &Versioner{}

	tests := []struct {
		input    string
		expected string
	}{
		{"abc123", "abc123"},
		{"1a2b3c4d5e6f", "1a2b3c4d5e6f"},
		{"", "unknown"},
	}

	for _, tt := range tests {
		result := versioner.getCommitString(tt.input)
		if result != tt.expected {
			t.Errorf("getCommitString(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestVersioner_ShouldUpdateVersion(t *testing.T) {
	versioner := &Versioner{}

	tests := []struct {
		name           string
		newVersion     string
		currentVersion string
		expected       bool
	}{
		{name: "empty new version", newVersion: "", currentVersion: "v1.0.0", expected: false},
		{name: "current dev", newVersion: "v1.0.1", currentVersion: "dev", expected: true},
		{name: "current empty", newVersion: "v1.0.1", currentVersion: "", expected: true},
		{name: "upgrade patch", newVersion: "v1.0.2", currentVersion: "v1.0.1", expected: true},
		{name: "same version", newVersion: "v1.0.1", currentVersion: "v1.0.1", expected: false},
		{name: "downgrade", newVersion: "v1.0.1", currentVersion: "v1.1.0", expected: false},
		{name: "non semver new", newVersion: "snapshot", currentVersion: "v1.1.0", expected: true},
		{name: "non semver current", newVersion: "v1.1.0", currentVersion: "snapshot", expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := versioner.shouldUpdateVersion(tt.newVersion, tt.currentVersion)
			if got != tt.expected {
				t.Errorf("shouldUpdateVersion(%q, %q) = %t, want %t", tt.newVersion, tt.currentVersion, got, tt.expected)
			}
		})
	}
}

func TestShouldUpdateToNewerVersion(t *testing.T) {
	tests := []struct {
		name     string
		new      string
		current  string
		expected bool
	}{
		{name: "semver upgrade", new: "v2.0.0", current: "v1.9.9", expected: true},
		{name: "semver downgrade", new: "v1.4.0", current: "v1.5.0", expected: false},
		{name: "identical", new: "v1.2.3", current: "v1.2.3", expected: false},
		{name: "handles prerelease", new: "v1.2.3-beta", current: "v1.2.3", expected: false},
		{name: "fallback for non semver", new: "snapshot", current: "v1.2.3", expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldUpdateToNewerVersion(tt.new, tt.current)
			if got != tt.expected {
				t.Errorf("shouldUpdateToNewerVersion(%q, %q) = %t, want %t", tt.new, tt.current, got, tt.expected)
			}
		})
	}
}

func TestCompareSemanticVersions(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
		ok       bool
	}{
		{name: "v1 greater", v1: "v1.2.0", v2: "v1.1.9", expected: 1, ok: true},
		{name: "v2 greater", v1: "v1.0.0", v2: "v1.0.1", expected: -1, ok: true},
		{name: "equal", v1: "v1.0.0", v2: "1.0.0", expected: 0, ok: true},
		{name: "non semver", v1: "snapshot", v2: "v1.0.0", expected: 0, ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := compareSemanticVersions(tt.v1, tt.v2)
			if got != tt.expected || ok != tt.ok {
				t.Errorf("compareSemanticVersions(%q, %q) = (%d, %t), want (%d, %t)", tt.v1, tt.v2, got, ok, tt.expected, tt.ok)
			}
		})
	}
}

func TestVersioner_ShouldForceUpdateFromBuild(t *testing.T) {
	versioner := &Versioner{}
	loaded := &config.Config{}
	loaded.Meta.Version = "v1.0.0"
	loaded.Meta.Commit = "abc123"

	tests := []struct {
		name         string
		buildVersion string
		buildCommit  string
		expected     bool
	}{
		{name: "no differences", buildVersion: "v1.0.0", buildCommit: "abc123", expected: false},
		{name: "version differs", buildVersion: "v1.0.1", buildCommit: "abc123", expected: true},
		{name: "commit differs", buildVersion: "v1.0.0", buildCommit: "def456", expected: true},
		{name: "empty build info", buildVersion: "", buildCommit: "", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := versioner.shouldForceUpdateFromBuild(tt.buildVersion, tt.buildCommit, loaded)
			if got != tt.expected {
				t.Errorf("shouldForceUpdateFromBuild(%q, %q) = %t, want %t", tt.buildVersion, tt.buildCommit, got, tt.expected)
			}
		})
	}
}
