package cmd

import (
	"io"
	"os"

	"github.com/bmf-san/ggc/v8/internal/git"
)

// Bisector handles git bisect operations.
type Bisector struct {
	gitClient    git.PassthroughOps
	outputWriter io.Writer
	helper       *Helper
}

// NewBisector creates a new Bisector instance.
func NewBisector(client git.PassthroughOps) *Bisector {
	return &Bisector{
		gitClient:    client,
		outputWriter: os.Stdout,
		helper:       NewHelper(),
	}
}

// Bisect executes git bisect commands.
//
// Supported guided flows:
//   - ggc bisect start <bad> <good>
//   - ggc bisect run <script-or-command>
//
// Other bisect subcommands are forwarded to git as-is.
func (b *Bisector) Bisect(args []string) {
	if len(args) == 0 || args[0] == "help" {
		b.helper.ShowPassthroughHelp("bisect")
		return
	}

	switch args[0] {
	case "start":
		b.start(args[1:])
	case "run":
		b.run(args[1:])
	default:
		b.forward(args)
	}
}

func (b *Bisector) start(args []string) {
	if len(args) < 2 {
		WriteLine(b.outputWriter, "Usage: ggc bisect start <bad> <good>")
		return
	}
	b.forward(append([]string{"start"}, args...))
}

func (b *Bisector) run(args []string) {
	if len(args) == 0 {
		WriteLine(b.outputWriter, "Usage: ggc bisect run <script-or-command>")
		return
	}
	b.forward(append([]string{"run"}, args...))
}

func (b *Bisector) forward(args []string) {
	if err := b.gitClient.RunGit("bisect", args); err != nil {
		WriteError(b.outputWriter, err)
	}
}
