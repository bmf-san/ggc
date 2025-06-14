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
