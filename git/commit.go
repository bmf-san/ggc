package git

import (
	"os"
)

func CommitAllowEmpty() error {
	cmd := execCommand("git", "commit", "--allow-empty", "-m", "empty commit")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CommitTmp() error {
	cmd := execCommand("git", "commit", "-m", "tmp")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
