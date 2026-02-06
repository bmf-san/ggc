package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bmf-san/ggc/v7/pkg/git"
)

// Differ handles git diff operations.
type Differ struct {
	gitClient    git.DiffReader
	outputWriter io.Writer
	helper       *Helper
}

// NewDiffer creates a new Differ instance.
func NewDiffer(client git.DiffReader) *Differ {
	return &Differ{
		gitClient:    client,
		outputWriter: os.Stdout,
		helper:       NewHelper(),
	}
}

type diffMode int

const (
	diffModeDefault diffMode = iota
	diffModeUnstaged
	diffModeStaged
	diffModeHead
)

func (m diffMode) String() string {
	switch m {
	case diffModeUnstaged:
		return "unstaged"
	case diffModeStaged:
		return "staged"
	case diffModeHead:
		return "head"
	default:
		return "default"
	}
}

type diffOptions struct {
	mode       diffMode
	commits    []string
	paths      []string
	stat       bool
	nameOnly   bool
	nameStatus bool
}

type diffUsageError struct {
	message string
}

func (e *diffUsageError) Error() string {
	return e.message
}

func newDiffUsageError(message string) error {
	return &diffUsageError{message: message}
}

// Diff executes git diff with the given arguments.
func (d *Differ) Diff(args []string) {
	if d.helper != nil && d.helper.outputWriter != d.outputWriter {
		d.helper.outputWriter = d.outputWriter
	}

	opts, err := parseDiffArgs(args, defaultPathExists)
	if err != nil {
		var usageErr *diffUsageError
		if errors.As(err, &usageErr) {
			if usageErr.message != "" {
				WriteErrorf(d.outputWriter, "%s", usageErr.message)
			}
			d.helper.ShowDiffHelp()
			return
		}

		WriteError(d.outputWriter, err)
		return
	}

	gitArgs := buildDiffArgs(opts)
	output, err := d.gitClient.DiffWith(gitArgs)
	if err != nil {
		WriteError(d.outputWriter, err)
		return
	}

	_, _ = fmt.Fprint(d.outputWriter, output)
}

func parseDiffArgs(args []string, pathExists func(string) bool) (*diffOptions, error) {
	if pathExists == nil {
		pathExists = func(string) bool { return false }
	}

	state := newDiffParseState()
	for _, arg := range args {
		if err := state.consume(arg); err != nil {
			return nil, err
		}
	}

	commits, prePaths, err := classifyDiffArgs(state.positional, pathExists)
	if err != nil {
		return nil, err
	}

	state.opts.commits = commits
	state.opts.paths = append(append([]string(nil), prePaths...), state.pathsAfterDash...)

	if err := state.validateModes(); err != nil {
		return nil, err
	}

	return state.opts, nil
}

type diffParseState struct {
	opts           *diffOptions
	positional     []string
	pathsAfterDash []string
	modeSet        bool
	seenDoubleDash bool
}

func newDiffParseState() *diffParseState {
	return &diffParseState{opts: &diffOptions{mode: diffModeDefault}}
}

func (s *diffParseState) consume(arg string) error {
	if !s.seenDoubleDash && arg == "--" {
		s.seenDoubleDash = true
		return nil
	}

	if s.seenDoubleDash {
		s.pathsAfterDash = append(s.pathsAfterDash, arg)
		return nil
	}

	return s.handleBeforeDoubleDash(arg)
}

func (s *diffParseState) handleBeforeDoubleDash(arg string) error {
	switch arg {
	case "unstaged", "staged", "head":
		return s.setMode(arg)
	case "--stat":
		s.opts.stat = true
		return nil
	case "--name-only":
		if s.opts.nameStatus {
			return newDiffUsageError("--name-only cannot be combined with --name-status")
		}
		s.opts.nameOnly = true
		return nil
	case "--name-status":
		if s.opts.nameOnly {
			return newDiffUsageError("--name-only cannot be combined with --name-status")
		}
		s.opts.nameStatus = true
		return nil
	}

	if strings.HasPrefix(arg, "--") {
		return newDiffUsageError(fmt.Sprintf("unknown flag %s", arg))
	}

	s.positional = append(s.positional, arg)
	return nil
}

func (s *diffParseState) setMode(mode string) error {
	if s.modeSet {
		return newDiffUsageError("multiple diff modes specified")
	}

	s.modeSet = true
	s.opts.mode = mapMode(mode)
	return nil
}

func mapMode(mode string) diffMode {
	switch mode {
	case "unstaged":
		return diffModeUnstaged
	case "staged":
		return diffModeStaged
	case "head":
		return diffModeHead
	default:
		return diffModeDefault
	}
}

func (s *diffParseState) validateModes() error {
	if (s.opts.mode == diffModeStaged || s.opts.mode == diffModeUnstaged || s.opts.mode == diffModeHead) && len(s.opts.commits) > 0 {
		return newDiffUsageError(fmt.Sprintf("%s mode does not accept commit arguments (but allows path arguments)", s.opts.mode.String()))
	}
	return nil
}

func classifyDiffArgs(tokens []string, pathExists func(string) bool) ([]string, []string, error) {
	if len(tokens) == 0 {
		return nil, nil, nil
	}

	maxCommits := len(tokens)
	if maxCommits > 2 {
		maxCommits = 2
	}

	for commitsCount := 0; commitsCount <= maxCommits; commitsCount++ {
		commitsCandidate := append([]string(nil), tokens[:commitsCount]...)
		pathsCandidate := tokens[commitsCount:]

		if pathsAreValid(pathsCandidate, pathExists) {
			return commitsCandidate, append([]string(nil), pathsCandidate...), nil
		}
	}

	if len(tokens) > 2 {
		return nil, nil, newDiffUsageError("too many commit arguments (maximum two)")
	}

	return nil, nil, newDiffUsageError("unable to determine commit and path arguments")
}

func pathsAreValid(paths []string, pathExists func(string) bool) bool {
	for _, p := range paths {
		if p == "" {
			return false
		}
		if pathExists(p) {
			continue
		}
		return false
	}
	return true
}

func buildDiffArgs(opts *diffOptions) []string {
	args := append([]string(nil), modeArg(opts)...)
	args = append(args, summaryArgs(opts)...)
	args = append(args, commitArgs(opts)...)

	if len(opts.paths) > 0 {
		args = append(args, "--")
		args = append(args, opts.paths...)
	}

	return args
}

func modeArg(opts *diffOptions) []string {
	switch opts.mode {
	case diffModeStaged:
		return []string{"--staged"}
	default:
		return nil
	}
}

func summaryArgs(opts *diffOptions) []string {
	var out []string
	if opts.stat {
		out = append(out, "--stat")
	}
	if opts.nameOnly {
		out = append(out, "--name-only")
	}
	if opts.nameStatus {
		out = append(out, "--name-status")
	}
	return out
}

func commitArgs(opts *diffOptions) []string {
	if len(opts.commits) > 0 {
		return append([]string(nil), opts.commits...)
	}

	switch opts.mode {
	case diffModeDefault, diffModeHead:
		return []string{"HEAD"}
	default:
		return nil
	}
}

func defaultPathExists(path string) bool {
	if path == "" {
		return false
	}
	if strings.HasPrefix(path, "--") {
		return false
	}

	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
