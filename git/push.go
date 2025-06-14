package git

import (
	"os"
	"os/exec"
)

var execCommand = exec.Command
var getCurrentBranch = GetCurrentBranch

func PushCurrentBranch() error {
	branch, err := getCurrentBranch()
	if err != nil {
		return err
	}
	cmd := execCommand("git", "push", "origin", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func PushForceCurrentBranch() error {
	branch, err := getCurrentBranch()
	if err != nil {
		return err
	}
	cmd := execCommand("git", "push", "--force", "origin", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
