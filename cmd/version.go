package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bmf-san/ggc/v7/pkg/config"
	"github.com/bmf-san/ggc/v7/pkg/git"
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
	gitClient    git.ConfigOps
}

// NewVersioner creates a new Versioner instance.
func NewVersioner(client git.ConfigOps) *Versioner {
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
	newVersion, newCommit, shouldUpdate := v.resolveBuildUpdates(loadedConfig)
	if !shouldUpdate {
		return
	}

	updated := false
	if newVersion != "" {
		v.updateConfigValue(configManager, "meta.version", newVersion)
		updated = true
	}
	if newCommit != "" {
		v.updateConfigValue(configManager, "meta.commit", newCommit)
		updated = true
	}

	if updated {
		*loadedConfig = *configManager.GetConfig()
	}
}

func (v *Versioner) resolveBuildUpdates(loadedConfig *config.Config) (string, string, bool) {
	if getVersionInfo == nil {
		return "", "", false
	}
	newVersion, newCommit := getVersionInfo()
	if newVersion == "" && newCommit == "" {
		return "", "", false
	}
	forceUpdate := v.shouldForceUpdateFromBuild(newVersion, newCommit, loadedConfig)
	versionUpdate := buildUpdateValue(newVersion, loadedConfig.Meta.Version, forceUpdate, v.shouldUpdateVersion)
	commitUpdate := buildUpdateValue(newCommit, loadedConfig.Meta.Commit, forceUpdate, v.shouldUpdateCommit)
	if versionUpdate == "" && commitUpdate == "" {
		return "", "", false
	}
	return versionUpdate, commitUpdate, true
}

func buildUpdateValue(newValue, currentValue string, force bool, shouldUpdate func(string, string) bool) string {
	if newValue == "" {
		return ""
	}
	if force || shouldUpdate(newValue, currentValue) {
		return newValue
	}
	return ""
}

// shouldUpdateVersion determines if version should be updated
func (v *Versioner) shouldUpdateVersion(newVersion, currentVersion string) bool {
	if newVersion == "" {
		return false
	}

	if currentVersion == "dev" || currentVersion == "" {
		return true
	}

	if newVersion == currentVersion {
		return false
	}

	return shouldUpdateToNewerVersion(newVersion, currentVersion)
}

// shouldUpdateCommit determines if commit should be updated
func (v *Versioner) shouldUpdateCommit(newCommit, currentCommit string) bool {
	return newCommit != "" && (currentCommit == "unknown" || currentCommit != newCommit)
}

// shouldForceUpdateFromBuild determines if config should be updated to match build info even without a semantic upgrade
func (v *Versioner) shouldForceUpdateFromBuild(buildVersion, buildCommit string, loadedConfig *config.Config) bool {
	if loadedConfig == nil {
		return false
	}

	if buildVersion != "" && buildVersion != loadedConfig.Meta.Version {
		return true
	}

	if buildCommit != "" && buildCommit != loadedConfig.Meta.Commit {
		return true
	}

	return false
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

// shouldUpdateToNewerVersion determines if newVersion represents a semantic upgrade over currentVersion
func shouldUpdateToNewerVersion(newVersion, currentVersion string) bool {
	if cmp, ok := compareSemanticVersions(newVersion, currentVersion); ok {
		return cmp > 0
	}

	return newVersion != currentVersion
}

// compareSemanticVersions compares two semantic version strings.
// Returns 1 if v1 > v2, -1 if v1 < v2, 0 if equal. ok is false if either version is not a semantic version.
func compareSemanticVersions(v1, v2 string) (int, bool) {
	segments1, ok1 := parseVersionSegments(v1)
	segments2, ok2 := parseVersionSegments(v2)
	if !ok1 || !ok2 {
		return 0, false
	}

	maxLen := len(segments1)
	if len(segments2) > maxLen {
		maxLen = len(segments2)
	}

	for i := 0; i < maxLen; i++ {
		s1 := segmentAt(segments1, i)
		s2 := segmentAt(segments2, i)
		if s1 > s2 {
			return 1, true
		}
		if s1 < s2 {
			return -1, true
		}
	}

	return 0, true
}

// parseVersionSegments converts semantic version string into its numeric segments.
func parseVersionSegments(version string) ([]int, bool) {
	if version == "" {
		return nil, false
	}

	trimmed := strings.TrimSpace(strings.ToLower(version))
	trimmed = strings.TrimPrefix(trimmed, "v")
	if idx := strings.IndexAny(trimmed, "-+"); idx != -1 {
		trimmed = trimmed[:idx]
	}
	if trimmed == "" {
		return nil, false
	}

	parts := strings.Split(trimmed, ".")
	segments := make([]int, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			return nil, false
		}
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, false
		}
		segments = append(segments, num)
	}

	return segments, true
}

// segmentAt fetches the segment at index i or returns 0 when out of range
func segmentAt(segments []int, i int) int {
	if i >= len(segments) {
		return 0
	}
	return segments[i]
}
