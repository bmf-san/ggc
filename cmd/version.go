package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/bmf-san/ggc/v4/config"
	"github.com/bmf-san/ggc/v4/git"
)

// VersionGetter is a function type for getting version info
type VersionGetter func() (version, commit string)

var getVersionInfo VersionGetter

// SetVersionGetter sets the version getter function
func SetVersionGetter(getter VersionGetter) {
	getVersionInfo = getter
}

// Versioneer handles version operations.
type Versioneer struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
	gitClient    git.Clienter
}

// NewVersioneer creates a new Versioneer instance.
func NewVersioneer(client git.Clienter) *Versioneer {
	return &Versioneer{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
		gitClient:    client,
	}
}

// Version returns the ggc version with the given arguments.
func (v *Versioneer) Version(args []string) {
	if len(args) == 0 {
		configManager := config.NewConfigManagerWithClient(v.gitClient)
		configManager.LoadConfig()

		loadedConfig := configManager.GetConfig()

		if loadedConfig.Meta.CreatedAt == "" {
			createdAt := time.Now().UTC().Format("2006-01-02_15:04:05")
			if err := configManager.Set("meta.created-at", createdAt); err != nil {
				_, _ = fmt.Fprintf(v.outputWriter, "warn: failed to set created-at: %v\n", err)
			} else {
				loadedConfig = configManager.GetConfig()
			}
		}
		if loadedConfig.Meta.Version == "dev" || loadedConfig.Meta.Commit == "unknown" {
			version, commit := getVersionInfo()
			if err := configManager.Set("meta.version", version); err != nil {
				_, _ = fmt.Fprintf(v.outputWriter, "warn: failed to set version: %v\n", err)
			}
			if err := configManager.Set("meta.commit", commit); err != nil {
				_, _ = fmt.Fprintf(v.outputWriter, "warn: failed to set commit: %v\n", err)
			}
		}

		version := loadedConfig.Meta.Version
		commit := loadedConfig.Meta.Commit
		if version == "" {
			version = "(devel)"
		}
		if commit == "" {
			commit = "unknown"
		}

		_, _ = fmt.Fprintf(v.outputWriter, "ggc version %s\n", version)
		_, _ = fmt.Fprintf(v.outputWriter, "commit: %s\n", commit)
		_, _ = fmt.Fprintf(v.outputWriter, "built: %s\n", loadedConfig.Meta.CreatedAt)
		_, _ = fmt.Fprintf(v.outputWriter, "config version: %s\n", loadedConfig.Meta.ConfigVersion)
		_, _ = fmt.Fprintf(v.outputWriter, "os/arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	} else {
		v.helper.ShowVersionHelp()
		return
	}
}
