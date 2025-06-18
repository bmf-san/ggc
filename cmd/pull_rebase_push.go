package cmd

import (
	"fmt"
	"io"
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
		outputWriter: nil,
	}
}

func (p *PullRebasePusher) PullRebasePush() {
	branchCmd := p.execCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOut, err := branchCmd.Output()
	if err != nil {
		if _, err := fmt.Fprintf(p.outputWriter, "Error: Failed to get branch name: %v\n", err); err != nil {
			_ = err
		}
		return
	}
	branch := strings.TrimSpace(string(branchOut))
	pullCmd := p.execCommand("git", "pull", "origin", branch)
	pullCmd.Stdout = p.outputWriter
	pullCmd.Stderr = p.outputWriter
	if err := pullCmd.Run(); err != nil {
		if _, err := fmt.Fprintf(p.outputWriter, "Error: Failed to git pull: %v\n", err); err != nil {
			_ = err
		}
		return
	}
	rebaseCmd := p.execCommand("git", "rebase", "origin/main")
	rebaseCmd.Stdout = p.outputWriter
	rebaseCmd.Stderr = p.outputWriter
	if err := rebaseCmd.Run(); err != nil {
		if _, err := fmt.Fprintf(p.outputWriter, "Error: Failed to git rebase: %v\n", err); err != nil {
			_ = err
		}
		return
	}
	pushCmd := p.execCommand("git", "push", "origin", branch)
	pushCmd.Stdout = p.outputWriter
	pushCmd.Stderr = p.outputWriter
	if err := pushCmd.Run(); err != nil {
		if _, err := fmt.Fprintf(p.outputWriter, "Error: Failed to git push: %v\n", err); err != nil {
			_ = err
		}
		return
	}
	if _, err := fmt.Fprintln(p.outputWriter, "pull→rebase→push completed"); err != nil {
		_ = err
	}
}

// 既存互換用
func PullRebasePush() {
	p := NewPullRebasePusher()
	p.outputWriter = nil // 既存通り（os.Stdoutに出力）
	p.PullRebasePush()
}
