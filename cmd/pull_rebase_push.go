package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type PullRebasePusher struct {
	execCommand  func(name string, arg ...string) *exec.Cmd
	outputWriter io.Writer
}

func NewPullRebasePusher() *PullRebasePusher {
	return &PullRebasePusher{
		execCommand:  exec.Command,
		outputWriter: os.Stdout,
	}
}

func (prp *PullRebasePusher) PullRebasePush() {
	branch, err := prp.getBranch()
	if err != nil {
		_, _ = fmt.Fprintf(prp.outputWriter, "Error: Failed to get branch name\n%v\n", err)
		return
	}
	_, _ = fmt.Fprintf(prp.outputWriter, "current branch: %s\n", branch)

	if err := prp.gitPull(branch); err != nil {
		_, _ = fmt.Fprintf(prp.outputWriter, "Error: Failed to git pull\n%v\n", err)
		return
	}

	if err := prp.gitRebase("origin/main"); err != nil {
		_, _ = fmt.Fprintf(prp.outputWriter, "Error: Failed to git rebase\n%v\n", err)
		return
	}

	if err := prp.gitPush(branch); err != nil {
		_, _ = fmt.Fprintf(prp.outputWriter, "Error: Failed to git push\n%v\n", err)
		return
	}
	_, _ = fmt.Fprintln(prp.outputWriter, "pull→rebase→push completed")
}

func (prp *PullRebasePusher) getBranch() (string, error) {
	cmd := prp.execCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (prp *PullRebasePusher) gitPull(branch string) error {
	cmd := prp.execCommand("git", "pull", "origin", branch)
	cmd.Stdout = prp.outputWriter
	cmd.Stderr = prp.outputWriter
	return cmd.Run()
}

func (prp *PullRebasePusher) gitRebase(baseBranch string) error {
	cmd := prp.execCommand("git", "rebase", baseBranch)
	cmd.Stdout = prp.outputWriter
	cmd.Stderr = prp.outputWriter
	return cmd.Run()
}

func (prp *PullRebasePusher) gitPush(branch string) error {
	cmd := prp.execCommand("git", "push", "origin", branch, "--force-with-lease")
	cmd.Stdout = prp.outputWriter
	cmd.Stderr = prp.outputWriter
	return cmd.Run()
}
