package git

import (
	"os/exec"
)

func PullCurrentBranch() error {
	branch, err := GetCurrentBranch()
	if err != nil {
		return err
	}
	cmd := exec.Command("git", "pull", "origin", branch)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func PullRebaseCurrentBranch() error {
	branch, err := GetCurrentBranch()
	if err != nil {
		return err
	}
	cmd := exec.Command("git", "pull", "--rebase", "origin", branch)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
