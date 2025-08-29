// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bmf-san/ggc/v4/git"
)

// Cleaner provides functionality for the clean command.
type Cleaner struct {
	gitClient    git.Clienter
	outputWriter io.Writer
	inputReader  *bufio.Reader
	helper       *Helper
}

// NewCleaner creates a new Cleaner.
func NewCleaner() *Cleaner {
	return NewCleanerWithClient(git.NewClient())
}

// NewCleanerWithClient creates a new Cleaner with the specified git client.
func NewCleanerWithClient(client git.Clienter) *Cleaner {
	c := &Cleaner{
		gitClient:    client,
		outputWriter: os.Stdout,
		inputReader:  bufio.NewReader(os.Stdin),
		helper:       NewHelper(),
	}
	c.helper.outputWriter = c.outputWriter
	return c
}

// Clean executes the clean command with the given arguments.
func (c *Cleaner) Clean(args []string) {
	if len(args) == 0 {
		c.helper.ShowCleanHelp()
		return
	}

	switch args[0] {
	case "files":
		if err := c.gitClient.CleanFiles(); err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		}
	case "dirs":
		if err := c.gitClient.CleanDirs(); err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		}
	default:
		c.helper.ShowCleanHelp()
	}
}

// CleanInteractive interactively selects files to clean.
func (c *Cleaner) CleanInteractive() {
	out, err := c.gitClient.CleanDryRun()
	if err != nil {
		_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		return
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	files := []string{}
	for _, line := range lines {
		if strings.HasPrefix(line, "Would remove ") {
			files = append(files, strings.TrimPrefix(line, "Would remove "))
		}
	}
	if len(files) == 0 {
		_, _ = fmt.Fprintln(c.outputWriter, "No files to clean.")
		return
	}

	for {
		_, _ = fmt.Fprintln(c.outputWriter, "\033[1;36mSelect files to delete by number (space separated, all: select all, none: deselect all, e.g. 1 3 5):\033[0m")
		for i, f := range files {
			_, _ = fmt.Fprintf(c.outputWriter, "  [\033[1;33m%d\033[0m] %s\n", i+1, f)
		}
		_, _ = fmt.Fprint(c.outputWriter, "> ")
		input, _ := c.inputReader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			_, _ = fmt.Fprintln(c.outputWriter, "Cancelled.")
			return
		}
		if input == "all" {
			if err := c.gitClient.CleanFilesForce(files); err != nil {
				_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
				return
			}
			_, _ = fmt.Fprintln(c.outputWriter, "Selected files deleted.")
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
				_, _ = fmt.Fprintf(c.outputWriter, "\033[1;31mInvalid number: %s\033[0m\n", idx)
				valid = false
				break
			}
			tmp = append(tmp, files[n-1])
		}
		if !valid {
			continue
		}
		if len(tmp) == 0 {
			_, _ = fmt.Fprintln(c.outputWriter, "\033[1;33mNothing selected.\033[0m")
			continue
		}
		_, _ = fmt.Fprintf(c.outputWriter, "\033[1;32mSelected files: %v\033[0m\n", tmp)
		_, _ = fmt.Fprint(c.outputWriter, "Delete these files? (y/n): ")
		ans, _ := c.inputReader.ReadString('\n')
		ans = strings.TrimSpace(ans)
		if ans == "y" || ans == "Y" {
			if err := c.gitClient.CleanFilesForce(tmp); err != nil {
				_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
				return
			}
			_, _ = fmt.Fprintln(c.outputWriter, "Selected files deleted.")
			break
		}
	}
}
