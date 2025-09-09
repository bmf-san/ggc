# Improve test coverage and refactor git package structure

## Description of Changes

This Pull Request significantly improves test coverage across multiple packages and refactors the git package structure for better organization and consistency.

## 🎯 Key Achievements

### Test Coverage Improvements

| Package | Before | After | Improvement |
|---------|--------|-------|-------------|
| **cmd** | 54.4% | **60.2%** | +5.8% |
| **router** | 78.3% | **93.5%** | +15.2% 🌟 |
| **config** | 78.3% | **82.2%** | +3.9% |
| **git** | 85.2% | **89.9%** | +4.7% |
| **Overall** | ~65% | **69.3%** | +4.3% |

### New Test Files Created

#### cmd package (7 new test files)
- `cmd/diff_test.go` - Comprehensive tests for Diff command
- `cmd/remote_test.go` - Basic tests for Remote command
- `cmd/status_test.go` - Tests for Status command and utility functions
- `cmd/reset_test.go` - Tests for Reset command
- `cmd/restore_test.go` - Tests for Restore command
- `cmd/stash_test.go` - Tests for Stash command
- `cmd/tag_test.go` - Tests for Tag command and utility methods

#### git package (supporting new structure)
- `git/rev-list_test.go` - Comprehensive tests for rev-list command functions
- `git/ls-files_test.go` - Tests for ls-files command functions

## 🏗️ Refactoring Details

### Git Package Structure Improvement

**Before: Inconsistent design**
```
util.go (mixed functions)
├── ListFiles() - git ls-files
├── GetUpstreamBranchName() - git rev-parse
└── GetAheadBehindCount() - git rev-list

rev-parse.go
├── GetCurrentBranch()
├── GetBranchName()
├── RevParseVerify()
└── GetCommitHash()

tag.go
└── GetTagCommit() - git rev-list (scattered)
```

**After: Command-based consistent design**
```
rev-parse.go (all rev-parse commands)
├── GetCurrentBranch()
├── GetBranchName()
├── RevParseVerify()
├── GetCommitHash()
└── GetUpstreamBranchName() ← moved

rev-list.go (all rev-list commands)
├── GetAheadBehindCount() ← moved
└── GetTagCommit() ← moved

ls-files.go (all ls-files commands)
└── ListFiles() ← moved

util.go → deleted (all functions moved to appropriate files)
```

### Enhanced Test Quality

#### 1. **Comprehensive Error Case Coverage**
```go
// Normal case test
func TestClient_Diff(t *testing.T) { /* ... */ }

// Error case test (newly added)
func TestClient_Diff_Error(t *testing.T) {
    client := &Client{
        execCommand: func(name string, args ...string) *exec.Cmd {
            return exec.Command("false") // failing command
        },
    }
    _, err := client.Diff()
    if err == nil {
        t.Error("Expected Diff to return an error")
    }
}
```

#### 2. **Complete Alias Functionality Testing (router)**
```go
func TestRouter_WithAliases(t *testing.T) {
    tests := []struct {
        name     string
        alias    config.Alias
        args     []string
        expected []string
    }{
        {
            name: "SimpleAlias",
            alias: config.Alias{
                Type:  config.SimpleAlias,
                Value: "status --short",
            },
            args:     []string{},
            expected: []string{"status", "--short"},
        },
        // ... other test cases
    }
}
```

#### 3. **Detailed Configuration Management Testing (config)**
```go
func TestFlattenMapDirect(t *testing.T) {
    input := map[string]interface{}{
        "aliases": map[string]interface{}{
            "st": "status",
            "co": "checkout",
        },
    }
    result := flattenMap(input, "")
    // detailed verification...
}
```

## 🔧 Technical Improvements

### 1. **Proper Exclusion of testutil Package**
```makefile
# Makefile
cover:
	go test $$(go list ./... | grep -v testutil) -coverprofile=coverage.out
	go tool cover -func=coverage.out
```

### 2. **Mock Usage and Standardization**
```go
// Unified use of testutil.NewMockGitClient()
mockClient := testutil.NewMockGitClient()
tagger := NewTagger(mockClient)
```

### 3. **Table-Driven Test Implementation**
```go
tests := []struct {
    name           string
    args           []string
    expectedOutput string
    wantArgs       []string
}{
    // Efficiently manage multiple test cases
}
```

## ✅ Quality Assurance

### Test Execution Results
```bash
=== All tests passing ===
ok github.com/bmf-san/ggc/v5/cmd    1.185s coverage: 60.2%
ok github.com/bmf-san/ggc/v5/config 0.735s coverage: 82.2%
ok github.com/bmf-san/ggc/v5/git    1.598s coverage: 89.9%
ok github.com/bmf-san/ggc/v5/router 1.470s coverage: 93.5%
```

### Functions Achieving 100% Coverage
- **git/diff.go**: All functions 100%
- **git/status.go**: All functions 100%
- **git/stash.go**: All functions 100%
- **git/rebase.go**: All functions 100%
- **git/rev-parse.go**: All functions 100%
- **git/rev-list.go**: All functions 100%
- **git/ls-files.go**: All functions 100%

## Related Issue

This addresses the need for improved test coverage and better code organization in the git package as discussed in internal development.

## Checklist

- [x] I have read the [CONTRIBUTING.md](https://github.com/bmf-san/ggc/blob/main/CONTRIBUTING.md)
- [x] I have added or updated tests (9 new test files, enhanced existing tests)
- [x] I have updated the documentation (if required) - Internal refactoring, no user-facing changes
- [x] Code is formatted with `make fmt`
- [x] Code passes linter checks via `make lint`
- [x] All tests are passing

## Screenshots (if appropriate)

Not applicable - this is an internal refactoring and testing improvement.

## Additional Context

### Files Changed Summary

#### New Files Created (11 files)
- 7 new test files in cmd/ package
- 2 new git/ package files (rev-list.go, ls-files.go)
- 2 new test files for git/ package

#### Modified Files (10 files)
- `Makefile` - testutil exclusion settings
- `router/router_test.go` - Added alias tests
- `config/config_test.go` - Added configuration function tests
- `git/rev-parse.go` - Added GetUpstreamBranchName
- `git/rev-parse_test.go` - Added new function tests
- `git/diff_test.go` - Added error case tests
- `git/status_test.go` - Added error case tests
- `git/stash_test.go` - Added error case tests
- `git/rebase_test.go` - Added comprehensive tests
- `git/tag.go` - Removed GetTagCommit function
- `git/tag_test.go` - Removed duplicate tests
- `internal/testutil/git_client.go` - Added documentation

#### Deleted Files (2 files)
- `git/util.go` - Functions moved to appropriate files
- `git/util_test.go` - Corresponding tests also moved

### Why This Refactoring Was Needed

1. **Inconsistent Design**: The original `git/util.go` mixed functions from different Git commands (ls-files, rev-parse, rev-list), making it difficult to maintain and understand.

2. **Low Test Coverage**: Several packages had insufficient test coverage, particularly missing error case scenarios.

3. **Code Organization**: Functions were scattered across files without logical grouping by their underlying Git command.

### Impact on Users

- **Zero Breaking Changes**: All user-facing commands remain exactly the same
- **Improved Reliability**: Higher test coverage means more stable code
- **Better Maintainability**: Cleaner code structure for future development

This refactoring significantly improves the codebase quality while maintaining full backward compatibility for all users.
