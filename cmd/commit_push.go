package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type CommitPusher struct {
	execCommand  func(name string, arg ...string) *exec.Cmd
	inputReader  *bufio.Reader
	outputWriter io.Writer
}

func NewCommitPusher() *CommitPusher {
	return &CommitPusher{
		execCommand:  exec.Command,
		inputReader:  bufio.NewReader(os.Stdin),
		outputWriter: os.Stdout,
	}
}

func (c *CommitPusher) CommitPushInteractive() {
	cmd := c.execCommand("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		if _, err := fmt.Fprintf(c.outputWriter, "Error: failed to get git status: %v\n", err); err != nil {
			_ = err
		}
		return
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		if _, err := fmt.Fprintln(c.outputWriter, "No changed files."); err != nil {
			_ = err
		}
		return
	}
	files := []string{}
	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		files = append(files, strings.TrimSpace(line[2:]))
	}
	if len(files) == 0 {
		if _, err := fmt.Fprintln(c.outputWriter, "No files to stage."); err != nil {
			_ = err
		}
		return
	}
	reader := c.inputReader
	for {
		if _, err := fmt.Fprintln(c.outputWriter, "\033[1;36mSelect files to add by number (space separated, all: select all, none: deselect all, e.g. 1 3 5):\033[0m"); err != nil {
			_ = err
		}
		for i, f := range files {
			if _, err := fmt.Fprintf(c.outputWriter, "  [\033[1;33m%d\033[0m] %s\n", i+1, f); err != nil {
				_ = err
			}
		}
		if _, err := fmt.Fprint(c.outputWriter, "> "); err != nil {
			_ = err
		}
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			if _, err := fmt.Fprintln(c.outputWriter, "Cancelled."); err != nil {
				_ = err
			}
			return
		}
		if input == "all" {
			addArgs := append([]string{"add"}, files...)
			addCmd := c.execCommand("git", addArgs...)
			addCmd.Stdout = c.outputWriter
			addCmd.Stderr = c.outputWriter
			if err := addCmd.Run(); err != nil {
				if _, err := fmt.Fprintf(c.outputWriter, "Error: failed to add files: %v\n", err); err != nil {
					_ = err
				}
				return
			}
			break
		}
		if input == "none" {
			continue
		}
		indices := strings.Fields(input)
		tmp := []string{}
		valid := true
		for _, idx := range indices {
			n, err := strconv.Atoi(idx)
			if err != nil || n < 1 || n > len(files) {
				if _, err := fmt.Fprintf(c.outputWriter, "\033[1;31mInvalid number: %s\033[0m\n", idx); err != nil {
					_ = err
				}
				valid = false
				break
			}
			tmp = append(tmp, files[n-1])
		}
		if !valid {
			continue
		}
		if len(tmp) == 0 {
			if _, err := fmt.Fprintln(c.outputWriter, "\033[1;33mNothing selected.\033[0m"); err != nil {
				_ = err
			}
			continue
		}
		if _, err := fmt.Fprintf(c.outputWriter, "\033[1;32mSelected files: %v\033[0m\n", tmp); err != nil {
			_ = err
		}
		if _, err := fmt.Fprint(c.outputWriter, "Add these files? (y/n): "); err != nil {
			_ = err
		}
		ans, _ := reader.ReadString('\n')
		ans = strings.TrimSpace(ans)
		if ans == "y" || ans == "Y" {
			addArgs := append([]string{"add"}, tmp...)
			addCmd := c.execCommand("git", addArgs...)
			addCmd.Stdout = c.outputWriter
			addCmd.Stderr = c.outputWriter
			if err := addCmd.Run(); err != nil {
				if _, err := fmt.Fprintf(c.outputWriter, "Error: failed to add files: %v\n", err); err != nil {
					_ = err
				}
				return
			}
			break
		}
	}
	if _, err := fmt.Fprint(c.outputWriter, "\n\r"); err != nil {
		_ = err
	}
	if _, err := fmt.Fprint(c.outputWriter, "Enter commit message: "); err != nil {
		_ = err
	}
	msg, _ := c.inputReader.ReadString('\n')
	msg = strings.TrimSpace(msg)
	if msg == "" {
		if _, err := fmt.Fprintln(c.outputWriter, "Cancelled."); err != nil {
			_ = err
		}
		return
	}
	commitCmd := c.execCommand("git", "commit", "-m", msg)
	commitCmd.Stdout = c.outputWriter
	commitCmd.Stderr = c.outputWriter
	if err := commitCmd.Run(); err != nil {
		if _, err := fmt.Fprintf(c.outputWriter, "Error: failed to commit: %v\n", err); err != nil {
			_ = err
		}
		return
	}
	branchCmd := c.execCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOut, err := branchCmd.Output()
	if err != nil {
		if _, err := fmt.Fprintf(c.outputWriter, "Error: failed to get branch name: %v\n", err); err != nil {
			_ = err
		}
		return
	}
	branch := strings.TrimSpace(string(branchOut))
	pushCmd := c.execCommand("git", "push", "origin", branch)
	pushCmd.Stdout = c.outputWriter
	pushCmd.Stderr = c.outputWriter
	if err := pushCmd.Run(); err != nil {
		if _, err := fmt.Fprintf(c.outputWriter, "Error: failed to push: %v\n", err); err != nil {
			_ = err
		}
		return
	}
	if _, err := fmt.Fprintln(c.outputWriter, "Done!"); err != nil {
		_ = err
	}
}

// For backward compatibility
func CommitPushInteractive() {
	NewCommitPusher().CommitPushInteractive()
}
