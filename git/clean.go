package git

import (
	"os"
)

func CleanFiles() error {
	cmd := execCommand("git", "clean", "-f")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CleanDirs() error {
	cmd := execCommand("git", "clean", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
