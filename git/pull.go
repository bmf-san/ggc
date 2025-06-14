package git

import (
	"os"
)

func PullCurrentBranch() error {
	branch, err := getCurrentBranch()
	if err != nil {
		return err
	}
	cmd := execCommand("git", "pull", "origin", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func PullRebaseCurrentBranch() error {
	branch, err := getCurrentBranch()
	if err != nil {
		return err
	}
	cmd := execCommand("git", "pull", "--rebase", "origin", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
