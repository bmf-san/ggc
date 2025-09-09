# Improve test coverage, refactor git package structure, and simplify config package

## Description of Changes

This Pull Request significantly improves test coverage across multiple packages, refactors the git package structure for better organization and consistency, and simplifies the config package by eliminating over-engineered abstractions while improving testability.

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

#### Config Package Simplification & Testability Enhancement
**Major Refactoring Journey:**

1. **Initial State**: Direct OS calls in functions → Hard to test
2. **Over-Engineering Phase**: Complex filesystem abstractions (fileReader, fileWriter, memoryFileOps) → 240+ lines of test-only code
3. **Simplification**: Eliminated all test-only abstractions → Real temp directories with t.TempDir()
4. **Final Solution**: Function-level dependency injection → Perfect balance

**Current Design:**
```go
// Public API (unchanged)
func (cm *Manager) Load() error
func (cm *Manager) Save() error

// Testable API (new)
func (cm *Manager) LoadWithFileOps(fileOps FileOps) error
func (cm *Manager) SaveWithFileOps(fileOps FileOps) error
```

**Benefits Achieved:**
- ✅ **Testability**: MockFileOps with in-memory filesystem
- ✅ **Simplicity**: No Manager struct complexity
- ✅ **Compatibility**: Public API unchanged
- ✅ **Performance**: Fast, isolated tests
- ✅ **Code Quality**: 191 lines removed, cleaner design

#### Enhanced Test Quality
- Added comprehensive error case coverage for all git package functions
- Improved router package with alias functionality tests
- Enhanced config package with detailed configuration tests
- Improved cmd package tests with meaningful assertions (eliminated "panic-only" tests)
- Achieved 100% coverage for: git/diff.go, git/status.go, git/stash.go, git/rebase.go, git/rev-parse.go, git/rev-list.go, git/ls-files.go

#### Technical Improvements
- Excluded testutil package from coverage reports (updated Makefile)
- Standardized mock usage with `testutil.NewMockGitClient()`
- Utilized table-driven tests for efficient test management
- Added comprehensive documentation for test utilities
- Applied YAGNI principle: removed unused interfaces (GitConfigExecutor, Validator)
- Implemented function-level dependency injection for OS operations

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

4. **Untestable Code**: Direct OS calls in config package functions made testing difficult and led to side effects.

5. **Over-Engineering**: Initial attempts at abstraction created unnecessary complexity with test-only code that violated YAGNI principles.

6. **Meaningless Tests**: Many cmd package tests only checked for "no panic" without verifying actual functionality.

### Impact on Users

- **Zero Breaking Changes**: All user-facing commands remain exactly the same
- **Improved Reliability**: Higher test coverage means more stable code
- **Better Maintainability**: Cleaner code structure for future development
- **Enhanced Performance**: Faster, more reliable config operations with better error handling

### Files Changed Summary

**New Files (11):**
- 7 new test files in cmd/ package
- 2 new git/ package files (rev-list.go, ls-files.go)
- 2 new test files for git/ package

**Modified Files (15):**
- Enhanced existing test files with error cases and meaningful assertions
- Updated Makefile for proper coverage calculation
- Moved functions to appropriate files in git/ package
- Updated existing git/ package files
- Major config package refactoring with function-level dependency injection
- Improved cmd package tests with output verification
- Removed unused interfaces and over-engineered abstractions

**Deleted Files (4):**
- git/util.go (functions moved to appropriate files)
- git/util_test.go (tests moved to appropriate files)
- config/filesystem.go (over-engineered abstraction removed)
- config/memory_filesystem.go (test-only complexity removed)

**Key Refactoring Principles Applied:**
- **YAGNI**: Don't build abstractions until you actually need them
- **Function-level DI**: Inject dependencies at function level, not struct level
- **Interface Segregation**: Minimal, focused interfaces
- **Testability**: Easy to test without complex setup

This refactoring significantly improves the codebase quality while maintaining full backward compatibility for all users.
