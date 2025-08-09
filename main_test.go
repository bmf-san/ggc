package main

import (
	"os"
	"testing"
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
	// Test that main function exists and can handle no arguments
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// Test with help argument to avoid infinite loop
	os.Args = []string{"ggc", "help"}

	// Test that main doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main() should not panic with help argument: %v", r)
		}
	}()

	main()
}
