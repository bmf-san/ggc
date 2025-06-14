package git

import (
	"os/exec"
)

func PushCurrentBranch() error {
	branch, err := GetCurrentBranch()
	if err != nil {
		return err
	}
	cmd := exec.Command("git", "push", "origin", branch)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
