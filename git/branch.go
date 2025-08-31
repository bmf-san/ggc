// Package git provides a high-level interface to git commands.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ListLocalBranches lists local branches.
func (c *Client) ListLocalBranches() ([]string, error) {
	cmd := c.execCommand("git", "branch", "--format", "%(refname:short)")
	out, err := cmd.Output()
	if err != nil {
		return nil, NewError("list local branches", "git branch --format %(refname:short)", err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	return lines, nil
}

// ListRemoteBranches lists remote branches.
func (c *Client) ListRemoteBranches() ([]string, error) {
	cmd := c.execCommand("git", "branch", "-r", "--format", "%(refname:short)")
	out, err := cmd.Output()
	if err != nil {
		return nil, NewError("list remote branches", "git branch -r --format %(refname:short)", err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
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
		return NewError("checkout branch", "git checkout "+name, err)
	}
	return nil
}

// CheckoutNewBranchFromRemote creates a new local branch tracking a remote branch.
func (c *Client) CheckoutNewBranchFromRemote(localBranch, remoteBranch string) error {
	cmd := c.execCommand("git", "checkout", "-b", localBranch, "--track", remoteBranch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("checkout new branch from remote", fmt.Sprintf("git checkout -b %s --track %s", localBranch, remoteBranch), err)
	}
	return nil
}

// DeleteBranch deletes a branch.
func (c *Client) DeleteBranch(name string) error {
	cmd := c.execCommand("git", "branch", "-d", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("delete branch", "git branch -d "+name, err)
	}
	return nil
}

// ListMergedBranches lists branches that have been merged.
func (c *Client) ListMergedBranches() ([]string, error) {
	cmd := c.execCommand("git", "branch", "--merged")
	out, err := cmd.Output()
	if err != nil {
		return nil, NewError("list merged branches", "git branch --merged", err)
	}

	branches := strings.Split(strings.TrimSpace(string(out)), "\n")
	var result []string
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
	cmd := c.execCommand("git", "branch", "-m", old, newName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("rename branch", fmt.Sprintf("git branch -m %s %s", old, newName), err)
	}
	return nil
}

// MoveBranch moves a branch pointer to a specific commit (git branch -f <branch> <commit>).
func (c *Client) MoveBranch(branch, commit string) error {
	cmd := c.execCommand("git", "branch", "-f", branch, commit)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("move branch", fmt.Sprintf("git branch -f %s %s", branch, commit), err)
	}
	return nil
}

// SetUpstreamBranch sets upstream for a branch (git branch -u <upstream> <branch>).
func (c *Client) SetUpstreamBranch(branch, upstream string) error {
	cmd := c.execCommand("git", "branch", "-u", upstream, branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return NewError("set upstream branch", fmt.Sprintf("git branch -u %s %s", upstream, branch), err)
	}
	return nil
}

// ListBranchesVerbose lists branches with verbose info (parses `git branch -vv`).
func (c *Client) ListBranchesVerbose() ([]BranchInfo, error) {
	cmd := c.execCommand("git", "branch", "-vv")
	out, err := cmd.Output()
	if err != nil {
		return nil, NewError("list branches verbose", "git branch -vv", err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
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
	shaCmd := c.execCommand("git", "rev-parse", "--short", branch)
	shaOut, shaErr := shaCmd.Output()
	if shaErr != nil {
		return nil, NewError("get branch info", fmt.Sprintf("git rev-parse --short %s", branch), shaErr)
	}
	msgCmd := c.execCommand("git", "log", "-1", "--pretty=%s", branch)
	msgOut, msgErr := msgCmd.Output()
	if msgErr != nil {
		return nil, NewError("get branch info", fmt.Sprintf("git log -1 --pretty=%%s %s", branch), msgErr)
	}
	current, _ := c.GetCurrentBranch()
	upstream, _ := c.GetUpstreamBranchName(branch)
	aheadBehind := ""
	if upstream != "" {
		if ab, err := c.GetAheadBehindCount(branch, upstream); err == nil {
			// rev-list --left-right --count returns "<ahead>\t<behind>"
			parts := strings.Split(strings.TrimSpace(ab), "\t")
			if len(parts) == 2 && (parts[0] != "0" || parts[1] != "0") {
				if parts[0] != "0" && parts[1] != "0" {
					aheadBehind = fmt.Sprintf("ahead %s, behind %s", parts[0], parts[1])
				} else if parts[0] != "0" {
					aheadBehind = fmt.Sprintf("ahead %s", parts[0])
				} else if parts[1] != "0" {
					aheadBehind = fmt.Sprintf("behind %s", parts[1])
				}
			}
		}
	}
	bi := BranchInfo{
		Name:            branch,
		IsCurrentBranch: branch == current,
		Upstream:        upstream,
		AheadBehind:     aheadBehind,
		LastCommitSHA:   strings.TrimSpace(string(shaOut)),
		LastCommitMsg:   strings.TrimSpace(string(msgOut)),
	}
	return &bi, nil
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
		return nil, NewError("sort branches", fmt.Sprintf("git branch --sort=%s --format %%(refname:short)", sortKey), err)
	}
	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return []string{}, nil
	}
	return strings.Split(trimmed, "\n"), nil
}

// BranchesContaining lists branches containing a given commit.
func (c *Client) BranchesContaining(commit string) ([]string, error) {
	cmd := c.execCommand("git", "branch", "--contains", commit)
	out, err := cmd.Output()
	if err != nil {
		return nil, NewError("branches containing commit", "git branch --contains "+commit, err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	res := []string{}
	for _, l := range lines {
		name := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(l), "* "))
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
