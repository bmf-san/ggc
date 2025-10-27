package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bmf-san/ggc/v7/pkg/git"
)

// Tagger handles tagging operations.
type Tagger struct {
	gitClient interface {
		git.TagOps
		git.ConfigOps
	}
	outputWriter io.Writer
	helper       *Helper
	// defaultRemote caches the default remote name to avoid
	// reloading configuration on each tag push.
	defaultRemote string
}

// NewTagger creates a new Tagger instance.
func NewTagger(client interface {
	git.TagOps
	git.ConfigOps
}) *Tagger {
	return &Tagger{
		gitClient:     client,
		outputWriter:  os.Stdout,
		helper:        NewHelper(),
		defaultRemote: "origin",
	}
}

// Tag executes git tag operations with the given arguments.
func (t *Tagger) Tag(args []string) {
	if len(args) == 0 {
		if err := t.gitClient.TagList(nil); err != nil {
			_, _ = fmt.Fprintf(t.outputWriter, "Error: %v\n", err)
		}
		return
	}

	switch args[0] {
	case "list", "l":
		t.listTags(args[1:])
		return
	case "create", "c":
		t.createTag(args[1:])
		return
	case "delete", "d":
		t.deleteTags(args[1:])
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

// listTags lists tags with optional pattern matching
func (t *Tagger) listTags(args []string) {
	if err := t.gitClient.TagList(args); err != nil {
		_, _ = fmt.Fprintf(t.outputWriter, "Error: %v\n", err)
	}
}

// createTag creates a new tag
func (t *Tagger) createTag(args []string) {
	if len(args) == 0 {
		_, _ = fmt.Fprintf(t.outputWriter, "Error: tag name is required\n")
		return
	}

	tagName := args[0]

	if len(args) == 1 {
		// tag current commit
		if err := t.gitClient.TagCreate(tagName, ""); err != nil {
			_, _ = fmt.Fprintf(t.outputWriter, "Error: %v\n", err)
			return
		}
	} else {
		// tag specific commit
		if err := t.gitClient.TagCreate(tagName, args[1]); err != nil {
			_, _ = fmt.Fprintf(t.outputWriter, "Error: %v\n", err)
			return
		}
	}

	_, _ = fmt.Fprintf(t.outputWriter, "Tag '%s' created\n", tagName)
}

// deleteTags deletes one or more tags
func (t *Tagger) deleteTags(args []string) {
	if len(args) == 0 {
		_, _ = fmt.Fprintf(t.outputWriter, "Error: at least one tag name is required\n")
		return
	}

	if err := t.gitClient.TagDelete(args); err != nil {
		_, _ = fmt.Fprintf(t.outputWriter, "Error: %v\n", err)
		return
	}

	for _, tagName := range args {
		_, _ = fmt.Fprintf(t.outputWriter, "Tag '%s' deleted\n", tagName)
	}
}

// pushTags pushes tags to remote
func (t *Tagger) pushTags(args []string) {
	// Use cached default remote; fallback to "origin" if unset
	remote := strings.TrimSpace(t.defaultRemote)
	if remote == "" {
		remote = "origin"
	}

	if len(args) == 0 {
		// push all tags
		if err := t.gitClient.TagPushAll(remote); err != nil {
			_, _ = fmt.Fprintf(t.outputWriter, "Error: %v\n", err)
			return
		}
		_, _ = fmt.Fprintf(t.outputWriter, "All tags pushed to %s\n", remote)
	} else {
		// push specific tag
		var tagName string
		if len(args) == 1 {
			// backwards-compatible: assume single arg is tag name
			tagName = args[0]
		} else {
			// git-compatible ordering: remote first, tag second
			candidate := strings.TrimSpace(args[0])
			if candidate == "" {
				_, _ = fmt.Fprintf(t.outputWriter, "Error: remote name cannot be empty or whitespace\n")
				return
			}
			remote = candidate
			tagName = args[1]
		}
		if err := t.gitClient.TagPush(remote, tagName); err != nil {
			_, _ = fmt.Fprintf(t.outputWriter, "Error: %v\n", err)
			return
		}
		_, _ = fmt.Fprintf(t.outputWriter, "Tag '%s' pushed to %s\n", tagName, remote)
	}
}

// showTag shows information about a tag
func (t *Tagger) showTag(args []string) {
	if len(args) == 0 {
		_, _ = fmt.Fprintf(t.outputWriter, "Error: tag name is required\n")
		return
	}

	tagName := args[0]
	if err := t.gitClient.TagShow(tagName); err != nil {
		_, _ = fmt.Fprintf(t.outputWriter, "Error: %v\n", err)
	}
}

// GetLatestTag gets the latest tag.
func (t *Tagger) GetLatestTag() (string, error) {
	return t.gitClient.GetLatestTag()
}

// TagExists checks if a tag exists.
func (t *Tagger) TagExists(tagName string) bool {
	return t.gitClient.TagExists(tagName)
}

// GetTagCommit gets the commit hash for a tag.
func (t *Tagger) GetTagCommit(tagName string) (string, error) {
	return t.gitClient.GetTagCommit(tagName)
}

// CreateAnnotatedTag creates an annotated tag
func (t *Tagger) CreateAnnotatedTag(args []string) {
	if len(args) == 0 {
		_, _ = fmt.Fprintf(t.outputWriter, "Error: tag name is required\n")
		return
	}

	tagName := args[0]
	if len(args) > 1 {
		// Use provided message
		message := strings.Join(args[1:], " ")
		if err := t.gitClient.TagCreateAnnotated(tagName, message); err != nil {
			_, _ = fmt.Fprintf(t.outputWriter, "Error: %v\n", err)
			return
		}
	} else {
		// Open editor for message
		if err := t.gitClient.TagCreateAnnotated(tagName, ""); err != nil {
			_, _ = fmt.Fprintf(t.outputWriter, "Error: %v\n", err)
			return
		}
	}

	_, _ = fmt.Fprintf(t.outputWriter, "Annotated tag '%s' created\n", tagName)
}
