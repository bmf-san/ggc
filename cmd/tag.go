package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Tagger handles tagging operations.
type Tagger struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
}

// NewTagger creates a new Tagger instance.
func NewTagger() *Tagger {
	return &Tagger{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
	}
}

// Tag executes git tag operations with the given arguments.
func (t *Tagger) Tag(args []string) {
	var cmd *exec.Cmd
	
	if len(args) == 0 {
		cmd = t.execCommand("git", "tag")
	} else {
		switch args[0] {
		case "list", "l":
			t.listTags(args[1:])
			return
		case "create", "c":
			t.createTag(args[1:])
			return
		case "delete", "d":
			t.deleteTag(args[1:])
			return
		case "annotated", "a":
			t.createAnnotatedTag(args[1:])
			return
		case "push":
			t.pushTags(args[1:])
			return
		case "show":
			t.showTag(args[1:])
			return
		default:
			t.helper.ShowTagHelp()
			return
		}
	}
	
	if err := t.runCommand(cmd); err != nil {
		_, _ = fmt.Fprintf(t.outputWriter, "Error: %v\n", err)
		return
	}
}

// listTags lists tags with optional pattern matching
func (t *Tagger) listTags(args []string) {
	var cmd *exec.Cmd
	
	if len(args) == 0 {
		cmd = t.execCommand("git", "tag", "--sort=-version:refname")
	} else {
		// List tags matching pattern
		gitArgs := append([]string{"tag", "--sort=-version:refname", "-l"}, args...)
		cmd = t.execCommand("git", gitArgs...)
	}
	
	if err := t.runCommand(cmd); err != nil {
		_, _ = fmt.Fprintf(t.outputWriter, "Error listing tags: %v\n", err)
	}
}

// createTag creates a lightweight tag
func (t *Tagger) createTag(args []string) {
	if len(args) == 0 {
		_, _ = fmt.Fprintf(t.outputWriter, "Error: tag name is required\n")
		return
	}
	
	tagName := args[0]
	var cmd *exec.Cmd
	
	if len(args) > 1 {
		// tag specific commit
		cmd = t.execCommand("git", "tag", tagName, args[1])
	} else {
		// tag current commit
		cmd = t.execCommand("git", "tag", tagName)
	}
	
	if err := t.runCommand(cmd); err != nil {
		_, _ = fmt.Fprintf(t.outputWriter, "Error creating tag: %v\n", err)
		return
	}
	
	_, _ = fmt.Fprintf(t.outputWriter, "Tag '%s' created successfully\n", tagName)
}

// deleteTag deletes one or more tags
func (t *Tagger) deleteTag(args []string) {
	if len(args) == 0 {
		_, _ = fmt.Fprintf(t.outputWriter, "Error: tag name(s) required\n")
		return
	}
	
	for _, tagName := range args {
		cmd := t.execCommand("git", "tag", "-d", tagName)
		
		if err := t.runCommand(cmd); err != nil {
			_, _ = fmt.Fprintf(t.outputWriter, "Error deleting tag '%s': %v\n", tagName, err)
		} else {
			_, _ = fmt.Fprintf(t.outputWriter, "Tag '%s' deleted successfully\n", tagName)
		}
	}
}

// pushTags pushes tags to remote repository
func (t *Tagger) pushTags(args []string) {
	var cmd *exec.Cmd
	
	if len(args) == 0 {
		// push all tags
		cmd = t.execCommand("git", "push", "origin", "--tags")
	} else {
		// push specific tag
		tagName := args[0]
		remote := "origin"
		if len(args) > 1 {
			remote = args[1]
		}
		cmd = t.execCommand("git", "push", remote, tagName)
	}
	
	if err := t.runCommand(cmd); err != nil {
		_, _ = fmt.Fprintf(t.outputWriter, "Error pushing tags: %v\n", err)
		return
	}
	
	_, _ = fmt.Fprintf(t.outputWriter, "Tags pushed successfully\n")
}

// showTag shows information about a specific tag
func (t *Tagger) showTag(args []string) {
	if len(args) == 0 {
		_, _ = fmt.Fprintf(t.outputWriter, "Error: tag name is required\n")
		return
	}
	
	tagName := args[0]
	cmd := t.execCommand("git", "show", tagName)
	
	if err := t.runCommand(cmd); err != nil {
		_, _ = fmt.Fprintf(t.outputWriter, "Error showing tag '%s': %v\n", tagName, err)
	}
}

// runCommand executes a command and pipes output to the writer
func (t *Tagger) runCommand(cmd *exec.Cmd) error {
	cmd.Stdout = t.outputWriter
	cmd.Stderr = t.outputWriter
	return cmd.Run()
}

// GetLatestTag returns the latest tag name
func (t *Tagger) GetLatestTag() (string, error) {
	cmd := t.execCommand("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// TagExists checks if a tag exists
func (t *Tagger) TagExists(tagName string) bool {
	cmd := t.execCommand("git", "tag", "-l", tagName)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == tagName
}

// GetTagCommit returns the commit hash for a given tag
func (t *Tagger) GetTagCommit(tagName string) (string, error) {
	cmd := t.execCommand("git", "rev-list", "-n", "1", tagName)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// createAnnotatedTag creates an annotated tag with message
func (t *Tagger) createAnnotatedTag(args []string) {
	if len(args) == 0 {
		_, _ = fmt.Fprintf(t.outputWriter, "Error: tag name is required\n")
		return
	}
	
	tagName := args[0]
	var cmd *exec.Cmd
	
	if len(args) > 1 {
		// Use provided message
		message := strings.Join(args[1:], " ")
		cmd = t.execCommand("git", "tag", "-a", tagName, "-m", message)
	} else {
		// Open editor for message
		cmd = t.execCommand("git", "tag", "-a", tagName)
	}
	
	if err := t.runCommand(cmd); err != nil {
		_, _ = fmt.Fprintf(t.outputWriter, "Error creating annotated tag: %v\n", err)
		return
	}
	
	_, _ = fmt.Fprintf(t.outputWriter, "Annotated tag '%s' created successfully\n", tagName)
}
