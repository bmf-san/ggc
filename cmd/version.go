package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

// VersionGetter is a function type for getting version info
type VersionGetter func() (version, commit, date string)

var getVersionInfo VersionGetter

// SetVersionGetter sets the version getter function
func SetVersionGetter(getter VersionGetter) {
	getVersionInfo = getter
}

// Versioneer handles status operations.
type Versioneer struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
}

// NewVersioneer creates a new Versioneer instance.
func NewVersioneer() *Versioneer {
	return &Versioneer{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
	}
}

// Version returns the ggc version with the given arguments.
func (v *Versioneer) Version(args []string) {
	if len(args) == 0 {
		version, commit, date := "dev", "none", "unknown"
		if getVersionInfo != nil {
			version, commit, date = getVersionInfo()
		}
		_, _ = fmt.Fprintf(v.outputWriter, "ggc version %s\n", version)
		_, _ = fmt.Fprintf(v.outputWriter, "commit: %s\n", commit)
		_, _ = fmt.Fprintf(v.outputWriter, "built: %s\n", date)
		_, _ = fmt.Fprintf(v.outputWriter, "os/arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	} else {
		v.helper.ShowVersionHelp()
		return
	}
}
