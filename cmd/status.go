package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bmf-san/ggc/v7/pkg/git"
)

// Statuser handles status operations.
type Statuser struct {
	outputWriter io.Writer
	helper       *Helper
	gitClient    git.StatusInfoReader
}

// NewStatuser creates a new Statuser instance.
func NewStatuser(client git.StatusInfoReader) *Statuser {
	return &Statuser{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		gitClient:    client,
	}
}

// getUpstreamStatus gets the upstream tracking status ahead/behind counts.
func (s *Statuser) getUpstreamStatus(branch string) string {
	// Check if upstream exists
	upstream, err := s.gitClient.GetUpstreamBranchName(branch)
	if err != nil {
		return ""
	}
	output, err := s.gitClient.GetAheadBehindCount(branch, upstream)
	if err != nil {
		return s.formatUpToDate(upstream)
	}
	ahead, behind, ok := parseCounts(output)
	if !ok {
		return s.formatUpToDate(upstream)
	}
	return s.formatAheadBehind(upstream, ahead, behind)
}

func parseCounts(output string) (string, string, bool) {
	counts := strings.Fields(strings.TrimSpace(output))
	if len(counts) != 2 {
		return "", "", false
	}
	return counts[0], counts[1], true
}

func (s *Statuser) formatUpToDate(upstream string) string {
	return fmt.Sprintf("Your branch is up to date with '%s'", upstream)
}

func (s *Statuser) formatAheadBehind(upstream, ahead, behind string) string {
	switch {
	case ahead == "0" && behind == "0":
		return s.formatUpToDate(upstream)
	case ahead != "0" && behind == "0":
		return fmt.Sprintf("Your branch is ahead of '%s' by %s commit(s)", upstream, ahead)
	case ahead == "0" && behind != "0":
		return fmt.Sprintf("Your branch is behind '%s' by %s commit(s)", upstream, behind)
	default:
		return fmt.Sprintf("Your branch and '%s' have diverged,\nand have %s and %s different commits each, respectively", upstream, ahead, behind)
	}
}

// Status executes git status with the given arguments.
func (s *Statuser) Status(args []string) {
	if len(args) == 0 {
		// Show status with color and branch info
		branch, err := s.gitClient.GetCurrentBranch()
		if err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error getting current branch: %v\n", err)
			return
		}

		upstreamStatus := s.getUpstreamStatus(branch)

		_, _ = fmt.Fprintf(s.outputWriter, "On branch %s\n", branch)
		if upstreamStatus != "" {
			_, _ = fmt.Fprintf(s.outputWriter, "%s\n", upstreamStatus)
		}
		_, _ = fmt.Fprintf(s.outputWriter, "\n")

		if output, err := s.gitClient.StatusWithColor(); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: %v\n", err)
		} else {
			_, _ = fmt.Fprint(s.outputWriter, output)
		}
		return
	}

	switch args[0] {
	case "short":
		if output, err := s.gitClient.StatusShortWithColor(); err != nil {
			_, _ = fmt.Fprintf(s.outputWriter, "Error: %v\n", err)
		} else {
			_, _ = fmt.Fprint(s.outputWriter, output)
		}
		return
	default:
		s.helper.ShowStatusHelp()
		return
	}
}
