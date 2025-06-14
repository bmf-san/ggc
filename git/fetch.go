package git

import (
	"os"
)

func FetchPrune() error {
	cmd := execCommand("git", "fetch", "--prune")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
