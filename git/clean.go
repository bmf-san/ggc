package git

import (
	"os/exec"
)

func CleanFiles() error {
	cmd := exec.Command("git", "clean", "-f")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func CleanDirs() error {
	cmd := exec.Command("git", "clean", "-d")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
