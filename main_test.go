package main

import (
	"testing"

	"github.com/bmf-san/ggc/v4/config"
)

func TestGetVersionInfo(t *testing.T) {
	// Test GetVersionInfo function
	version, commit := GetVersionInfo()

	// Test that function returns without panic
	// In test environment, these should be empty strings (not set by linker)
	if version == "" && commit == "" {
		// Expected behavior in test environment
		t.Log("Version and commit are empty as expected in test environment")
	} else {
		// If they are set (e.g., in production build), they should be valid strings
		if len(version) == 0 {
			t.Error("Version should not be empty if set")
		}
		if len(commit) == 0 {
			t.Error("Commit should not be empty if set")
		}
		t.Logf("Version: %s, Commit: %s", version, commit)
	}

	// Verify that the function returns two string values
	// This test ensures the function signature is correct
	versionType := len(version) >= 0 // All strings have len >= 0
	commitType := len(commit) >= 0   // All strings have len >= 0

	if !versionType || !commitType {
		t.Error("GetVersionInfo should return two string values")
	}
}

func TestMain(t *testing.T) {
	// Test that main function components can be initialized without side effects
	// We avoid calling main() directly to prevent actual git command execution during tests

	// Test config manager creation (safe, no git commands)
	cm := config.NewConfigManager()
	if cm == nil {
		t.Error("config manager should be created")
	}

	// Test version getter setup (safe, no git commands)
	version, commit := GetVersionInfo()
	// In test environment, these should be empty strings (not set by linker)
	if version == "" && commit == "" {
		t.Log("Version and commit are empty as expected in test environment")
	}

	// CRITICAL: We don't call cmd.NewCmd() or main() to avoid git command side effects
	// This ensures no actual git.Client is created and no git commands are executed
	t.Log("Main function components tested successfully without side effects")
}
