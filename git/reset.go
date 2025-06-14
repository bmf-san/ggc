package git

import (
	"os"
)

func ResetClean() error {
	cmd1 := execCommand("git", "reset", "--hard", "HEAD")
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	if err := cmd1.Run(); err != nil {
		return err
	}
	cmd2 := execCommand("git", "clean", "-fd")
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	return cmd2.Run()
}
