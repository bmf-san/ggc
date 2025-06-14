package git

import (
	"os/exec"
)

func FetchPrune() error {
	cmd := exec.Command("git", "fetch", "--prune")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
