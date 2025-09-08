package cmd

import (
	"testing"
)

// TestIsCommitLikeStrict tests the isCommitLikeStrict function
func TestIsCommitLikeStrict(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"HEAD", "HEAD", true},
		{"HEAD^", "HEAD^", true},
		{"HEAD~1", "HEAD~1", true},
		{"HEAD@{0}", "HEAD@{0}", true},
		{"refs/heads/main", "refs/heads/main", true},
		{"origin/main", "origin/main", true},
		{"abc123def", "abc123def", true}, // hex object name
		{"invalid", "invalid", false},
		{"", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isCommitLikeStrict(tt.input)
			if result != tt.expected {
				t.Errorf("isCommitLikeStrict(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsHEADVariant tests the isHEADVariant function
func TestIsHEADVariant(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"HEAD", "HEAD", true},
		{"HEAD^", "HEAD^", true},
		{"HEAD^1", "HEAD^1", true},
		{"HEAD~1", "HEAD~1", true},
		{"HEAD~10", "HEAD~10", true},
		{"HEAD@{0}", "HEAD@{0}", true},
		{"HEAD@{yesterday}", "HEAD@{yesterday}", true},
		{"main", "main", false},
		{"origin/main", "origin/main", false},
		{"", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isHEADVariant(tt.input)
			if result != tt.expected {
				t.Errorf("isHEADVariant(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsExplicitRef tests the isExplicitRef function
func TestIsExplicitRef(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"refs/heads/main", "refs/heads/main", true},
		{"refs/heads/feature", "refs/heads/feature", true},
		{"refs/remotes/origin/main", "refs/remotes/origin/main", true},
		{"refs/tags/v1.0.0", "refs/tags/v1.0.0", true},
		{"main", "main", false},
		{"origin/main", "origin/main", false},
		{"", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isExplicitRef(tt.input)
			if result != tt.expected {
				t.Errorf("isExplicitRef(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsRemoteRef tests the isRemoteRef function
func TestIsRemoteRef(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"origin/main", "origin/main", true},
		{"origin/feature", "origin/feature", true},
		{"origin/develop", "origin/develop", true},
		{"main", "main", false},
		{"refs/heads/main", "refs/heads/main", false},
		{"upstream/main", "upstream/main", false}, // only origin/ is checked
		{"", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRemoteRef(tt.input)
			if result != tt.expected {
				t.Errorf("isRemoteRef(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsHexObjectName tests the isHexObjectName function
func TestIsHexObjectName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"short hex", "abc1234", true},  // 7 chars minimum
		{"full hex", "abc123def456789012345678901234567890abcd", true},
		{"mixed case", "AbC123DeF", true},
		{"with invalid char", "abc123g", false},
		{"too short", "abc12", false},  // 5 chars < 7
		{"empty", "", false},
		{"non-hex", "hello", false},
		{"numbers only", "1234567", true},  // 7 chars minimum
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isHexObjectName(tt.input)
			if result != tt.expected {
				t.Errorf("isHexObjectName(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsHexChar tests the isHexChar function
func TestIsHexChar(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		{"0", '0', true},
		{"9", '9', true},
		{"a", 'a', true},
		{"f", 'f', true},
		{"A", 'A', true},
		{"F", 'F', true},
		{"g", 'g', false},
		{"G", 'G', false},
		{"z", 'z', false},
		{"space", ' ', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isHexChar(tt.input)
			if result != tt.expected {
				t.Errorf("isHexChar(%c) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
