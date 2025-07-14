package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
	"runtime"

	"github.com/bmf-san/ggc/config"
)

// Versioneer handles version operations.
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
		configManager := config.NewConfigManager()
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
		_, _ = fmt.Fprintf(v.outputWriter, "ggc version %s\n", loadedConfig.Meta.Version)
		_, _ = fmt.Fprintf(v.outputWriter, "commit: %s\n", loadedConfig.Meta.Commit)
		_, _ = fmt.Fprintf(v.outputWriter, "built: %s\n", loadedConfig.Meta.CreatedAt)
		_, _ = fmt.Fprintf(v.outputWriter, "config version: %s\n", loadedConfig.Meta.ConfigVersion)
		_, _ = fmt.Fprintf(v.outputWriter, "os/arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	} else {
		v.helper.ShowVersionHelp()
		return
	}
}
