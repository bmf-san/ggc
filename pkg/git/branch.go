// Package git provides a high-level interface to git commands.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// LocalBranchLister provides only local branch listing.
type LocalBranchLister interface {
	ListLocalBranches() ([]string, error)
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

// BranchInfo contains rich information about a branch.
type BranchInfo struct {
	Name            string
	IsCurrentBranch bool
	Upstream        string
	AheadBehind     string // e.g. "ahead 2, behind 1"
	LastCommitSHA   string
	LastCommitMsg   string
}

func splitBranchLines(out []byte) []string {
	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return []string{}
	}
	return strings.Split(trimmed, "\n")
}

func normalizeBranchName(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", fmt.Errorf("branch name cannot be empty")
	}
	cmd := exec.Command("git", "check-ref-format", "--branch", trimmed)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("invalid branch name %q: %w", trimmed, err)
	}
	return trimmed, nil
}

// ValidateBranchName checks whether the provided name is a valid git branch name.
func ValidateBranchName(name string) error {
	_, err := normalizeBranchName(name)
	return err
}

// ListLocalBranches lists local branches.
func (c *Client) ListLocalBranches() ([]string, error) {
	cmd := c.execCommand("git", "branch", "--format", "%(refname:short)")
	out, err := cmd.Output()
	if err != nil {
		return nil, NewOpError("list local branches", "git branch --format %(refname:short)", err)
	}
	lines := splitBranchLines(out)
	return lines, nil
}

// ListRemoteBranches lists remote branches.
func (c *Client) ListRemoteBranches() ([]string, error) {
	cmd := c.execCommand("git", "branch", "-r", "--format", "%(refname:short)")
	out, err := cmd.Output()
	if err != nil {
		return nil, NewOpError("list remote branches", "git branch -r --format %(refname:short)", err)
	}
	lines := splitBranchLines(out)
	// Exclude HEAD references such as origin/HEAD -> origin/main
	filtered := []string{}
	for _, l := range lines {
		if strings.Contains(l, "->") {
			continue
		}
		filtered = append(filtered, strings.TrimSpace(l))
	}
	return filtered, nil
}

// CheckoutBranch checks out an existing branch.
func (c *Client) CheckoutBranch(name string) error {
	cmd := c.execCommand("git", "checkout", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("checkout branch", "git checkout "+name, err)
	}
	return nil
}

// CheckoutNewBranch creates a new branch and checks it out.
func (c *Client) CheckoutNewBranch(name string) error {
	normalized, err := normalizeBranchName(name)
	if err != nil {
		return err
	}

	cmd := c.execCommand("git", "checkout", "-b", normalized)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("checkout new branch", fmt.Sprintf("git checkout -b %s", normalized), err)
	}
	return nil
}

// CheckoutNewBranchFromRemote creates a new local branch tracking a remote branch.
func (c *Client) CheckoutNewBranchFromRemote(localBranch, remoteBranch string) error {
	normalizedLocal, err := normalizeBranchName(localBranch)
	if err != nil {
		return err
	}

	cmd := c.execCommand("git", "checkout", "-b", normalizedLocal, "--track", remoteBranch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("checkout new branch from remote", fmt.Sprintf("git checkout -b %s --track %s", normalizedLocal, remoteBranch), err)
	}
	return nil
}

// DeleteBranch deletes a branch.
func (c *Client) DeleteBranch(name string) error {
	normalized, err := normalizeBranchName(name)
	if err != nil {
		return err
	}

	cmd := c.execCommand("git", "branch", "-d", normalized)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("delete branch", "git branch -d "+normalized, err)
	}
	return nil
}

// ListMergedBranches lists branches that have been merged.
func (c *Client) ListMergedBranches() ([]string, error) {
	cmd := c.execCommand("git", "branch", "--merged")
	out, err := cmd.Output()
	if err != nil {
		return nil, NewOpError("list merged branches", "git branch --merged", err)
	}

	branches := splitBranchLines(out)
	result := []string{}
	for _, branch := range branches {
		branch = strings.TrimSpace(branch)
		if branch != "" && !strings.HasPrefix(branch, "*") {
			result = append(result, branch)
		}
	}
	return result, nil
}

// RenameBranch renames a branch (git branch -m <old> <new>).
func (c *Client) RenameBranch(old, newName string) error {
	trimmedOld := strings.TrimSpace(old)
	if trimmedOld == "" {
		return fmt.Errorf("branch name cannot be empty")
	}

	normalizedNew, err := normalizeBranchName(newName)
	if err != nil {
		return err
	}

	cmd := c.execCommand("git", "branch", "-m", trimmedOld, normalizedNew)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("rename branch", fmt.Sprintf("git branch -m %s %s", trimmedOld, normalizedNew), err)
	}
	return nil
}

// MoveBranch moves a branch pointer to a specific commit (git branch -f <branch> <commit>).
func (c *Client) MoveBranch(branch, commit string) error {
	normalized, err := normalizeBranchName(branch)
	if err != nil {
		return err
	}
	trimmedCommit := strings.TrimSpace(commit)
	if trimmedCommit == "" {
		return fmt.Errorf("commit cannot be empty")
	}

	cmd := c.execCommand("git", "branch", "-f", normalized, trimmedCommit)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("move branch", fmt.Sprintf("git branch -f %s %s", normalized, trimmedCommit), err)
	}
	return nil
}

// SetUpstreamBranch sets upstream for a branch (git branch -u <upstream> <branch>).
func (c *Client) SetUpstreamBranch(branch, upstream string) error {
	normalizedBranch, err := normalizeBranchName(branch)
	if err != nil {
		return err
	}
	trimmedUpstream := strings.TrimSpace(upstream)
	if trimmedUpstream == "" {
		return fmt.Errorf("upstream branch cannot be empty")
	}

	cmd := c.execCommand("git", "branch", "-u", trimmedUpstream, normalizedBranch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewOpError("set upstream branch", fmt.Sprintf("git branch -u %s %s", trimmedUpstream, normalizedBranch), err)
	}
	return nil
}

// ListBranchesVerbose lists branches with verbose info (parses `git branch -vv`).
func (c *Client) ListBranchesVerbose() ([]BranchInfo, error) {
	cmd := c.execCommand("git", "branch", "-vv")
	out, err := cmd.Output()
	if err != nil {
		return nil, NewOpError("list branches verbose", "git branch -vv", err)
	}
	lines := splitBranchLines(out)
	infos := make([]BranchInfo, 0, len(lines))
	for _, line := range lines {
		l := strings.TrimRight(line, "\r\n")
		if strings.TrimSpace(l) == "" {
			continue
		}
		info := parseBranchVVLine(l)
		if info.Name != "" {
			infos = append(infos, info)
		}
	}
	return infos, nil
}

// GetBranchInfo returns BranchInfo for a specific branch using ListBranchesVerbose.
func (c *Client) GetBranchInfo(branch string) (*BranchInfo, error) {
	infos, err := c.ListBranchesVerbose()
	if err != nil {
		return nil, err
	}

	for _, bi := range infos {
		if bi.Name == branch {
			b := bi
			return &b, nil
		}
	}

	// As a fallback, try to build minimal info with separate commands
	return c.buildBranchInfoFallback(branch)
}

// buildBranchInfoFallback builds BranchInfo using individual git commands
func (c *Client) buildBranchInfoFallback(branch string) (*BranchInfo, error) {
	sha, err := c.getBranchSHA(branch)
	if err != nil {
		return nil, err
	}

	msg, err := c.getBranchLastCommitMsg(branch)
	if err != nil {
		return nil, err
	}

	current, _ := c.GetCurrentBranch()
	upstream, _ := c.GetUpstreamBranchName(branch)
	aheadBehind := c.calculateAheadBehind(branch, upstream)

	bi := BranchInfo{
		Name:            branch,
		IsCurrentBranch: branch == current,
		Upstream:        upstream,
		AheadBehind:     aheadBehind,
		LastCommitSHA:   sha,
		LastCommitMsg:   msg,
	}
	return &bi, nil
}

// getBranchSHA gets the SHA for a branch
func (c *Client) getBranchSHA(branch string) (string, error) {
	shaCmd := c.execCommand("git", "rev-parse", "--short", branch)
	shaOut, shaErr := shaCmd.Output()
	if shaErr != nil {
		return "", NewOpError("get branch info", fmt.Sprintf("git rev-parse --short %s", branch), shaErr)
	}
	return strings.TrimSpace(string(shaOut)), nil
}

// getBranchLastCommitMsg gets the last commit message for a branch
func (c *Client) getBranchLastCommitMsg(branch string) (string, error) {
	msgCmd := c.execCommand("git", "log", "-1", "--pretty=%s", branch)
	msgOut, msgErr := msgCmd.Output()
	if msgErr != nil {
		return "", NewOpError("get branch info", fmt.Sprintf("git log -1 --pretty=%%s %s", branch), msgErr)
	}
	return strings.TrimSpace(string(msgOut)), nil
}

// calculateAheadBehind calculates ahead/behind status for a branch vs upstream
func (c *Client) calculateAheadBehind(branch, upstream string) string {
	if upstream == "" {
		return ""
	}
	ab, err := c.GetAheadBehindCount(branch, upstream)
	if err != nil {
		return ""
	}
	ahead, behind, ok := parseAheadBehind(strings.TrimSpace(ab))
	if !ok || (ahead == "0" && behind == "0") {
		return ""
	}
	return formatAheadBehind(ahead, behind)
}

func parseAheadBehind(s string) (string, string, bool) {
	parts := strings.Split(s, "\t")
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func formatAheadBehind(ahead, behind string) string {
	if ahead != "0" && behind != "0" {
		return fmt.Sprintf("ahead %s, behind %s", ahead, behind)
	}
	if ahead != "0" {
		return fmt.Sprintf("ahead %s", ahead)
	}
	if behind != "0" {
		return fmt.Sprintf("behind %s", behind)
	}
	return ""
}

// SortBranches lists branches sorted by the specified key ("name" or "date").
func (c *Client) SortBranches(by string) ([]string, error) {
	var sortKey string
	switch strings.ToLower(strings.TrimSpace(by)) {
	case "date":
		// Newest first
		sortKey = "-committerdate"
	case "name", "":
		sortKey = "refname"
	default:
		sortKey = by // pass-through to git for flexibility
	}
	cmd := c.execCommand("git", "branch", "--sort="+sortKey, "--format", "%(refname:short)")
	out, err := cmd.Output()
	if err != nil {
		return nil, NewOpError("sort branches", fmt.Sprintf("git branch --sort=%s --format %%(refname:short)", sortKey), err)
	}
	return splitBranchLines(out), nil
}

// BranchesContaining lists branches containing a given commit.
func (c *Client) BranchesContaining(commit string) ([]string, error) {
	cmd := c.execCommand("git", "branch", "--contains", commit)
	out, err := cmd.Output()
	if err != nil {
		return nil, NewOpError("branches containing commit", "git branch --contains "+commit, err)
	}
	lines := splitBranchLines(out)
	res := []string{}
	for _, l := range lines {
		l = strings.TrimSpace(l)
		name := strings.TrimSpace(strings.TrimPrefix(l, "* "))
		if name != "" {
			res = append(res, name)
		}
	}
	return res, nil
}

// parseBranchVVLine parses a single line of `git branch -vv` output into BranchInfo.
func parseBranchVVLine(line string) BranchInfo {
	// Example lines:
	// * main    1a2b3c4 [origin/main: ahead 2, behind 1] Commit message here
	//   feature 5d6e7f8 [origin/feature] Another message
	//   local    abcdef0 Commit without upstream
	info := BranchInfo{}
	l := line
	if strings.HasPrefix(l, "*") {
		info.IsCurrentBranch = true
		l = strings.TrimSpace(l[1:])
	} else {
		l = strings.TrimSpace(l)
	}

	// Split into tokens first to grab name and sha
	fields := strings.Fields(l)
	if len(fields) < 2 {
		return info
	}
	info.Name = fields[0]
	info.LastCommitSHA = fields[1]

	// Remainder after name and sha
	remainder := strings.TrimSpace(strings.TrimPrefix(l, info.Name))
	remainder = strings.TrimSpace(strings.TrimPrefix(remainder, info.LastCommitSHA))

	// Check for upstream bracket
	upstream := ""
	aheadBehind := ""
	if strings.HasPrefix(strings.TrimSpace(remainder), "[") {
		// Extract [ ... ]
		rb := remainder
		endIdx := strings.Index(rb, "]")
		if endIdx > 1 {
			inside := strings.TrimSpace(rb[1:endIdx])
			// Formats: "origin/main: ahead 2, behind 1" or "origin/main" or "gone"
			if inside != "" {
				// If contains ':', split
				if idx := strings.Index(inside, ":"); idx != -1 {
					upstream = strings.TrimSpace(inside[:idx])
					aheadBehind = strings.TrimSpace(inside[idx+1:])
				} else {
					upstream = strings.TrimSpace(inside)
				}
			}
			// Remainder after bracket
			if endIdx+1 < len(rb) {
				remainder = strings.TrimSpace(rb[endIdx+1:])
			} else {
				remainder = ""
			}
		}
	}
	info.Upstream = upstream
	info.AheadBehind = aheadBehind
	info.LastCommitMsg = strings.TrimSpace(remainder)
	return info
}

// Ensure unused import exec referenced when building tests using helperCommand
var _ = exec.Cmd{}
