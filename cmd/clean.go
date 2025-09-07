// Package cmd provides command implementations for the ggc CLI tool.
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bmf-san/ggc/v5/git"
)

// Cleaner provides functionality for the clean command.
type Cleaner struct {
	gitClient    git.Clienter
	outputWriter io.Writer
	inputReader  *bufio.Reader
	helper       *Helper
}

// NewCleaner creates a new Cleaner.
func NewCleaner(client git.Clienter) *Cleaner {
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
	files, err := c.getCleanableFiles()
	if err != nil {
		_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
		return
	}
	if len(files) == 0 {
		_, _ = fmt.Fprintln(c.outputWriter, "No files to clean.")
		return
	}

	c.runInteractiveCleanLoop(files)
}

// getCleanableFiles retrieves the list of files that can be cleaned
func (c *Cleaner) getCleanableFiles() ([]string, error) {
	out, err := c.gitClient.CleanDryRun()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	files := []string{}
	for _, line := range lines {
		if strings.HasPrefix(line, "Would remove ") {
			files = append(files, strings.TrimPrefix(line, "Would remove "))
		}
	}
	return files, nil
}

// runInteractiveCleanLoop runs the interactive selection loop
func (c *Cleaner) runInteractiveCleanLoop(files []string) {
	for {
		c.displayFileSelection(files)
		input, _ := c.inputReader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			_, _ = fmt.Fprintln(c.outputWriter, "Canceled.")
			return
		}
		if c.handleSpecialCommands(input, files) {
			return
		}
		if c.handleFileSelection(input, files) {
			return
		}
	}
}

// displayFileSelection shows the file selection interface
func (c *Cleaner) displayFileSelection(files []string) {
	_, _ = fmt.Fprintln(c.outputWriter, "\033[1;36mSelect files to delete by number (space separated, all: select all, none: deselect all, e.g. 1 3 5):\033[0m")
	for i, f := range files {
		_, _ = fmt.Fprintf(c.outputWriter, "  [\033[1;33m%d\033[0m] %s\n", i+1, f)
	}
	_, _ = fmt.Fprint(c.outputWriter, "> ")
}

// handleSpecialCommands processes "all" and "none" commands
func (c *Cleaner) handleSpecialCommands(input string, files []string) bool {
	if input == "all" {
		if err := c.gitClient.CleanFilesForce(files); err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
			return true
		}
		_, _ = fmt.Fprintln(c.outputWriter, "Selected files deleted.")
		return true
	}
	if input == "none" {
		return false // Continue loop
	}
	return false
}

// handleFileSelection processes numeric file selection
func (c *Cleaner) handleFileSelection(input string, files []string) bool {
	selectedFiles, valid := c.parseFileIndices(input, files)
	if !valid {
		return false // Continue loop
	}
	if len(selectedFiles) == 0 {
		_, _ = fmt.Fprintln(c.outputWriter, "\033[1;33mNothing selected.\033[0m")
		return false // Continue loop
	}

	return c.confirmAndDelete(selectedFiles)
}

// parseFileIndices parses user input into selected files
func (c *Cleaner) parseFileIndices(input string, files []string) ([]string, bool) {
	indices := strings.Fields(input)
	selectedFiles := []string{}

	for _, idx := range indices {
		n, err := strconv.Atoi(idx)
		if err != nil || n < 1 || n > len(files) {
			_, _ = fmt.Fprintf(c.outputWriter, "\033[1;31mInvalid number: %s\033[0m\n", idx)
			return nil, false
		}
		selectedFiles = append(selectedFiles, files[n-1])
	}
	return selectedFiles, true
}

// confirmAndDelete confirms deletion and executes it
func (c *Cleaner) confirmAndDelete(selectedFiles []string) bool {
	_, _ = fmt.Fprintf(c.outputWriter, "\033[1;32mSelected files: %v\033[0m\n", selectedFiles)
	_, _ = fmt.Fprint(c.outputWriter, "Delete these files? (y/n): ")

	ans, _ := c.inputReader.ReadString('\n')
	ans = strings.TrimSpace(ans)

	if ans == "y" || ans == "Y" {
		if err := c.gitClient.CleanFilesForce(selectedFiles); err != nil {
			_, _ = fmt.Fprintf(c.outputWriter, "Error: %v\n", err)
			return true
		}
		_, _ = fmt.Fprintln(c.outputWriter, "Selected files deleted.")
		return true
	}
	return false // Continue loop
}
