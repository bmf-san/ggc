package git

import (
	"os/exec"
)

func ResetClean() error {
	cmd1 := exec.Command("git", "reset", "--hard", "HEAD")
	cmd1.Stdout = nil
	cmd1.Stderr = nil
	if err := cmd1.Run(); err != nil {
		return err
	}
	cmd2 := exec.Command("git", "clean", "-fd")
	cmd2.Stdout = nil
	cmd2.Stderr = nil
	return cmd2.Run()
}
