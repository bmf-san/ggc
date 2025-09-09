package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/bmf-san/ggc/v5/config"
	"github.com/bmf-san/ggc/v5/git"
)

// VersionGetter is a function type for getting version info
type VersionGetter func() (version, commit string)

var getVersionInfo VersionGetter

// SetVersionGetter sets the version getter function
func SetVersionGetter(getter VersionGetter) {
	getVersionInfo = getter
}

// Versioner handles version operations.
type Versioner struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
	gitClient    git.Clienter
}

// NewVersioner creates a new Versioner instance.
func NewVersioner(client git.Clienter) *Versioner {
	return &Versioner{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
		gitClient:    client,
	}
}

// Version returns the ggc version with the given arguments.
func (v *Versioner) Version(args []string) {
	if len(args) == 0 {
		v.displayVersionInfo()
	} else {
		v.helper.ShowVersionHelp()
	}
}

// displayVersionInfo displays the version information
func (v *Versioner) displayVersionInfo() {
	configManager := config.NewConfigManager(v.gitClient)
	configManager.LoadConfig()
	loadedConfig := configManager.GetConfig()

	v.ensureCreatedAtSet(configManager, loadedConfig)
	v.updateVersionInfoFromBuild(configManager, loadedConfig)
	v.printVersionInfo(loadedConfig)
}

// ensureCreatedAtSet ensures the created-at timestamp is set
func (v *Versioner) ensureCreatedAtSet(configManager *config.Manager, loadedConfig *config.Config) {
	if loadedConfig.Meta.CreatedAt == "" {
		createdAt := time.Now().UTC().Format("2006-01-02_15:04:05")
		if err := configManager.Set("meta.created-at", createdAt); err != nil {
			_, _ = fmt.Fprintf(v.outputWriter, "warn: failed to set created-at: %v\n", err)
		} else {
			*loadedConfig = *configManager.GetConfig()
		}
	}
}

// updateVersionInfoFromBuild updates version info from build info or ldflags
func (v *Versioner) updateVersionInfoFromBuild(configManager *config.Manager, loadedConfig *config.Config) {
	if getVersionInfo == nil {
		return
	}

	newVersion, newCommit := getVersionInfo()
	shouldUpdateVersion := v.shouldUpdateVersion(newVersion, loadedConfig.Meta.Version)
	shouldUpdateCommit := v.shouldUpdateCommit(newCommit, loadedConfig.Meta.Commit)

	if shouldUpdateVersion {
		v.updateConfigValue(configManager, "meta.version", newVersion)
	}
	if shouldUpdateCommit {
		v.updateConfigValue(configManager, "meta.commit", newCommit)
	}
	if shouldUpdateVersion || shouldUpdateCommit {
		*loadedConfig = *configManager.GetConfig()
	}
}

// shouldUpdateVersion determines if version should be updated
func (v *Versioner) shouldUpdateVersion(newVersion, currentVersion string) bool {
	return newVersion != "" && (currentVersion == "dev" || currentVersion != newVersion)
}

// shouldUpdateCommit determines if commit should be updated
func (v *Versioner) shouldUpdateCommit(newCommit, currentCommit string) bool {
	return newCommit != "" && (currentCommit == "unknown" || currentCommit != newCommit)
}

// updateConfigValue updates a config value and handles errors
func (v *Versioner) updateConfigValue(configManager *config.Manager, key, value string) {
	if err := configManager.Set(key, value); err != nil {
		_, _ = fmt.Fprintf(v.outputWriter, "warn: failed to set %s: %v\n", key, err)
	}
}

// printVersionInfo prints the version information
func (v *Versioner) printVersionInfo(loadedConfig *config.Config) {
	version := v.getVersionString(loadedConfig.Meta.Version)
	commit := v.getCommitString(loadedConfig.Meta.Commit)

	_, _ = fmt.Fprintf(v.outputWriter, "ggc version %s\n", version)
	_, _ = fmt.Fprintf(v.outputWriter, "commit: %s\n", commit)
	_, _ = fmt.Fprintf(v.outputWriter, "built: %s\n", loadedConfig.Meta.CreatedAt)
	_, _ = fmt.Fprintf(v.outputWriter, "config version: %s\n", loadedConfig.Meta.ConfigVersion)
	_, _ = fmt.Fprintf(v.outputWriter, "os/arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

// getVersionString returns a formatted version string
func (v *Versioner) getVersionString(version string) string {
	if version == "" {
		return "(devel)"
	}
	return version
}

// getCommitString returns a formatted commit string
func (v *Versioner) getCommitString(commit string) string {
	if commit == "" {
		return "unknown"
	}
	return commit
}
