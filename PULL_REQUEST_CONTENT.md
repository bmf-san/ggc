# Improve test coverage and refactor git package structure

## Description of Changes

This Pull Request significantly improves test coverage across multiple packages and refactors the git package structure for better organization and consistency.

### 🎯 Key Achievements

#### Test Coverage Improvements
- **cmd package**: 54.4% → 60.2% (+5.8%)
- **router package**: 78.3% → 93.5% (+15.2%) 🌟
- **config package**: 78.3% → 82.2% (+3.9%)
- **git package**: 85.2% → 89.9% (+4.7%)
- **Overall project**: ~65% → **69.3%** (+4.3%)

#### New Test Files Created (9 files)
- `cmd/diff_test.go` - Comprehensive tests for Diff command
- `cmd/remote_test.go` - Basic tests for Remote command
- `cmd/status_test.go` - Tests for Status command and utilities
- `cmd/reset_test.go` - Tests for Reset command
- `cmd/restore_test.go` - Tests for Restore command
- `cmd/stash_test.go` - Tests for Stash command
- `cmd/tag_test.go` - Tests for Tag command and utility methods
- `git/rev-list_test.go` - Comprehensive tests for rev-list functions
- `git/ls-files_test.go` - Tests for ls-files functions

#### Git Package Structure Refactoring
**Before (inconsistent design):**
```
util.go (mixed functions)
├── ListFiles() - git ls-files
├── GetUpstreamBranchName() - git rev-parse
└── GetAheadBehindCount() - git rev-list
```

**After (command-based consistent design):**
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
```

#### Enhanced Test Quality
- Added comprehensive error case coverage for all git package functions
- Improved router package with alias functionality tests
- Enhanced config package with detailed configuration tests
- Achieved 100% coverage for: git/diff.go, git/status.go, git/stash.go, git/rebase.go, git/rev-parse.go, git/rev-list.go, git/ls-files.go

#### Technical Improvements
- Excluded testutil package from coverage reports (updated Makefile)
- Standardized mock usage with `testutil.NewMockGitClient()`
- Utilized table-driven tests for efficient test management
- Added comprehensive documentation for test utilities

## Related Issue

This addresses the need for improved test coverage and better code organization in the git package as discussed in internal development.

## Checklist

- [x] I have read the [CONTRIBUTING.md](https://github.com/bmf-san/ggc/blob/main/CONTRIBUTING.md)
- [x] I have added or updated tests (9 new test files, enhanced existing tests)
- [x] I have updated the documentation (if required) - Internal refactoring, no user-facing changes
- [x] Code is formatted with `make fmt`
- [x] Code passes linter checks via `make lint`
- [x] All tests are passing

### Test Results Verification
```bash
=== All tests passing ===
ok github.com/bmf-san/ggc/v5/cmd    1.185s coverage: 60.2%
ok github.com/bmf-san/ggc/v5/config 0.735s coverage: 82.2%
ok github.com/bmf-san/ggc/v5/git    1.598s coverage: 89.9%
ok github.com/bmf-san/ggc/v5/router 1.470s coverage: 93.5%
```

## Screenshots (if appropriate)

Not applicable - this is an internal refactoring and testing improvement.

## Additional Context

### Why This Refactoring Was Needed

1. **Inconsistent Design**: The original `git/util.go` mixed functions from different Git commands (ls-files, rev-parse, rev-list), making it difficult to maintain and understand.

2. **Low Test Coverage**: Several packages had insufficient test coverage, particularly missing error case scenarios.

3. **Code Organization**: Functions were scattered across files without logical grouping by their underlying Git command.

### Impact on Users

- **Zero Breaking Changes**: All user-facing commands remain exactly the same
- **Improved Reliability**: Higher test coverage means more stable code
- **Better Maintainability**: Cleaner code structure for future development

### Files Changed Summary

**New Files (11):**
- 7 new test files in cmd/ package
- 2 new git/ package files (rev-list.go, ls-files.go)
- 2 new test files for git/ package

**Modified Files (10):**
- Enhanced existing test files with error cases
- Updated Makefile for proper coverage calculation
- Moved functions to appropriate files
- Updated existing git/ package files

**Deleted Files (2):**
- git/util.go (functions moved to appropriate files)
- git/util_test.go (tests moved to appropriate files)

This refactoring significantly improves the codebase quality while maintaining full backward compatibility for all users.
