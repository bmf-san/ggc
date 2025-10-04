package git

// This file defines focused, cohesive interfaces that represent
// small slices of Git functionality. Commands in the cmd/ package
// should depend on the smallest interface they need instead of a
// single, monolithic client interface. This reduces test mock burden
// and improves maintainability.

// DiffReader provides read-only diff operations.
// Implemented by Client and any compatible mock in tests.
type DiffReader interface {
	Diff() (string, error)
	DiffStaged() (string, error)
	DiffHead() (string, error)
	DiffWith(args []string) (string, error)
}

// StatusReader provides read-only status output with color support.
type StatusReader interface {
	StatusWithColor() (string, error)
	StatusShortWithColor() (string, error)
}

// BranchUpstreamReader provides information about the current branch and its upstream.
type BranchUpstreamReader interface {
	GetCurrentBranch() (string, error)
	GetUpstreamBranchName(branch string) (string, error)
	GetAheadBehindCount(branch, upstream string) (string, error)
}

// StatusInfoReader is a pragmatic composite for the status command dependencies.
// It avoids pulling in an overly broad client surface area.
type StatusInfoReader interface {
	StatusReader
	BranchUpstreamReader
}

// ConfigOps provides git config get/set operations used by config.Manager and Configurer.
type ConfigOps interface {
	ConfigGetGlobal(key string) (string, error)
	ConfigSetGlobal(key, value string) error
	GetVersion() (string, error)
	GetCommitHash() (string, error)
}

// LocalBranchLister provides only local branch listing.
type LocalBranchLister interface {
	ListLocalBranches() ([]string, error)
}

// FileLister provides repository file listing.
type FileLister interface {
	ListFiles() (string, error)
}

// CommitWriter provides write operations for commits.
type CommitWriter interface {
	Commit(message string) error
	CommitAmend() error
	CommitAmendNoEdit() error
	CommitAmendWithMessage(message string) error
	CommitAllowEmpty() error
}

// Pusher provides push operation.
type Pusher interface {
	Push(force bool) error
}

// Puller provides pull operation.
type Puller interface {
	Pull(rebase bool) error
}

// Stager provides add operations for staging changes.
type Stager interface {
	Add(files ...string) error
	AddInteractive() error
}

// FetchOps provides fetch operation(s).
type FetchOps interface {
	Fetch(prune bool) error
}

// BranchReader provides read-only branch queries.
type BranchReader interface {
	GetCurrentBranch() (string, error)
	ListLocalBranches() ([]string, error)
	ListMergedBranches() ([]string, error)
	ListBranchesVerbose() ([]BranchInfo, error)
	SortBranches(by string) ([]string, error)
	BranchesContaining(commit string) ([]string, error)
	GetBranchInfo(branch string) (*BranchInfo, error)
	ListRemoteBranches() ([]string, error)
	RevParseVerify(ref string) bool
}

// BranchWriter provides branch mutation operations.
type BranchWriter interface {
	CheckoutNewBranch(name string) error
	CheckoutBranch(name string) error
	CheckoutNewBranchFromRemote(localBranch, remoteBranch string) error
	DeleteBranch(name string) error
	RenameBranch(old, newName string) error
	MoveBranch(branch, commit string) error
	SetUpstreamBranch(branch, upstream string) error
}

// BranchOps is a pragmatic composite for the branch command dependencies.
type BranchOps interface {
	BranchReader
	BranchWriter
}

// RemoteManager provides remote repository management operations.
type RemoteManager interface {
	RemoteList() error
	RemoteAdd(name, url string) error
	RemoteRemove(name string) error
	RemoteSetURL(name, url string) error
}

// LogReader provides read-only access to git log output.
type LogReader interface {
	LogSimple() error
	LogGraph() error
}

// CleanOps provides operations used by the clean command.
type CleanOps interface {
	CleanFiles() error
	CleanDirs() error
	CleanDryRun() (string, error)
	CleanFilesForce(files []string) error
}

// ResetOps provides operations used by the reset command.
type ResetOps interface {
	GetCurrentBranch() (string, error)
	ResetHardAndClean() error
	ResetHard(commit string) error
}

// TagOps provides operations used by the tag command.
type TagOps interface {
	// list/show
	TagList(pattern []string) error
	TagShow(name string) error
	// create/delete
	TagCreate(name string, commit string) error
	TagCreateAnnotated(name, message string) error
	TagDelete(names []string) error
	// push
	TagPush(remote, name string) error
	TagPushAll(remote string) error
	// query
	GetLatestTag() (string, error)
	TagExists(name string) bool
	GetTagCommit(name string) (string, error)
}

// RebaseOps provides operations used by the rebase command.
type RebaseOps interface {
	// sequence operations
	RebaseInteractive(commitCount int) error
	Rebase(upstream string) error
	RebaseContinue() error
	RebaseAbort() error
	RebaseSkip() error
	// discovery
	GetCurrentBranch() (string, error)
	GetUpstreamBranch(branch string) (string, error)
	LogOneline(from, to string) (string, error)
	RevParseVerify(ref string) bool
}

// StashOps provides operations used by the stash command.
type StashOps interface {
	Stash() error
	StashList() (string, error)
	StashShow(stash string) error
	StashApply(stash string) error
	StashPop(stash string) error
	StashDrop(stash string) error
	StashClear() error
}

// RestoreOps provides operations used by the restore command.
type RestoreOps interface {
	RestoreWorkingDir(paths ...string) error
	RestoreStaged(paths ...string) error
	RestoreFromCommit(commit string, paths ...string) error
	RevParseVerify(ref string) bool
}
