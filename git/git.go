package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Client is a git client.
type Client struct {
	execCommand          func(name string, arg ...string) *exec.Cmd
	GetCurrentBranchFunc func() (string, error)
}

// Clienter is an interface for a git client.
type Clienter interface {
	// === Repository Information ===
	GetCurrentBranch() (string, error)
	GetBranchName() (string, error)
	GetGitStatus() (string, error)

	// === Status Operations ===
	Status() (string, error)
	StatusShort() (string, error)
	StatusWithColor() (string, error)
	StatusShortWithColor() (string, error)

	// === Staging Operations ===
	Add(files ...string) error
	AddInteractive() error

	// === Commit Operations ===
	Commit(message string) error
	CommitAmend() error
	CommitAmendNoEdit() error
	CommitAmendWithMessage(message string) error
	CommitAllowEmpty() error

	// === Diff Operations ===
	Diff() (string, error)
	DiffStaged() (string, error)
	DiffHead() (string, error)

	// === Branch Operations ===
	ListLocalBranches() ([]string, error)
	ListRemoteBranches() ([]string, error)
	CheckoutNewBranch(name string) error
	CheckoutBranch(name string) error
	CheckoutNewBranchFromRemote(localBranch, remoteBranch string) error
	DeleteBranch(name string) error
	ListMergedBranches() ([]string, error)
	RevParseVerify(ref string) bool

	// === Remote Operations ===
	Push(force bool) error
	Pull(rebase bool) error
	Fetch(prune bool) error
	RemoteList() error
	RemoteAdd(name, url string) error
	RemoteRemove(name string) error
	RemoteSetURL(name, url string) error

	// === Tag Operations ===
	TagList(pattern []string) error
	TagCreate(name string, commit string) error
	TagCreateAnnotated(name, message string) error
	TagDelete(names []string) error
	TagPush(remote, name string) error
	TagPushAll(remote string) error
	TagShow(name string) error
	GetLatestTag() (string, error)
	TagExists(name string) bool
	GetTagCommit(name string) (string, error)

	// === Log Operations ===
	LogSimple() error
	LogGraph() error
	LogOneline(from, to string) (string, error)

	// === Rebase Operations ===
	RebaseInteractive(commitCount int) error
	GetUpstreamBranch(branch string) (string, error)

	// === Stash Operations ===
	Stash() error
	StashList() (string, error)
	StashShow(stash string) error
	StashApply(stash string) error
	StashPop(stash string) error
	StashDrop(stash string) error
	StashClear() error

	// === Restore Operations ===
	RestoreWorkingDir(paths ...string) error
	RestoreStaged(paths ...string) error
	RestoreFromCommit(commit string, paths ...string) error
	RestoreAll() error
	RestoreAllStaged() error

	// === Reset and Clean Operations ===
	ResetHardAndClean() error
	ResetHard(commit string) error
	CleanFiles() error
	CleanDirs() error
	CleanDryRun() (string, error)
	CleanFilesForce(files []string) error

	// === Utility Operations ===
	ListFiles() (string, error)
	GetUpstreamBranchName(branch string) (string, error)
	GetAheadBehindCount(branch, upstream string) (string, error)
}

// NewClient creates a new Client.
func NewClient() *Client {
	return &Client{
		execCommand: exec.Command,
	}
}

// === Repository Information ===

// GetCurrentBranch gets the current branch name.
func (c *Client) GetCurrentBranch() (string, error) {
	if c.GetCurrentBranchFunc != nil {
		return c.GetCurrentBranchFunc()
	}
	cmd := c.execCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("get current branch", "git rev-parse --abbrev-ref HEAD", err)
	}
	branch := strings.TrimSpace(string(out))
	return branch, nil
}

// RevParseVerify checks whether the given ref resolves to a valid object.
// It runs: git rev-parse --verify --quiet <ref>
func (c *Client) RevParseVerify(ref string) bool {
	cmd := c.execCommand("git", "rev-parse", "--verify", "--quiet", ref)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// GetBranchName gets branch name.
func (c *Client) GetBranchName() (string, error) {
	cmd := c.execCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("get branch name", "git rev-parse --abbrev-ref HEAD", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// GetGitStatus gets git status.
func (c *Client) GetGitStatus() (string, error) {
	cmd := c.execCommand("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return "", NewError("get status", "git status --porcelain", err)
	}
	return string(out), nil
}

// === Branch Operations ===

// CheckoutNewBranch creates a new branch and checks it out.
func (c *Client) CheckoutNewBranch(name string) error {
	cmd := c.execCommand("git", "checkout", "-b", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("checkout new branch", fmt.Sprintf("git checkout -b %s", name), err)
	}
	return nil
}
