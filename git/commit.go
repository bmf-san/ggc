package git

import (
	"os/exec"
)

func CommitAllowEmpty() error {
	cmd := exec.Command("git", "commit", "--allow-empty", "-m", "empty commit")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func CommitTmp() error {
	cmd := exec.Command("git", "commit", "-m", "tmp")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
