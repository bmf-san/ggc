package git

import (
	"os"
	"os/exec"
)

func LogSimple() error {
	cmd := exec.Command("git", "log", "--oneline")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
