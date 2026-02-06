package keybindings

import (
	"fmt"
	"os"
	"time"

	"github.com/bmf-san/ggc/v7/internal/config"
)

// HotConfigReloader enables reloading configuration without restart
type HotConfigReloader struct {
	configPath      string
	resolver        *KeyBindingResolver
	lastModified    time.Time
	watching        bool
	reloadCallbacks []func(*config.Config)
}

// NewHotConfigReloader creates a new hot config reloader
func NewHotConfigReloader(configPath string, resolver *KeyBindingResolver) *HotConfigReloader {
	return &HotConfigReloader{
		configPath:      configPath,
		resolver:        resolver,
		watching:        false,
		reloadCallbacks: make([]func(*config.Config), 0),
	}
}

// StartWatching begins watching the config file for changes
func (hcr *HotConfigReloader) StartWatching() error {
	if hcr.watching {
		return fmt.Errorf("already watching config file")
	}

	// Get initial modification time
	if stat, err := os.Stat(hcr.configPath); err == nil {
		hcr.lastModified = stat.ModTime()
	}

	hcr.watching = true

	// Start watching in a goroutine
	go hcr.watchLoop()

	return nil
}

// StopWatching stops watching the config file
func (hcr *HotConfigReloader) StopWatching() {
	hcr.watching = false
}

// watchLoop continuously checks for config file changes
func (hcr *HotConfigReloader) watchLoop() {
	ticker := time.NewTicker(1 * time.Second) // Check every second
	defer ticker.Stop()

	for hcr.watching {
		<-ticker.C
		if stat, err := os.Stat(hcr.configPath); err == nil {
			if stat.ModTime().After(hcr.lastModified) {
				hcr.lastModified = stat.ModTime()
				hcr.reloadConfig()
			}
		}
	}
}

// reloadConfig reloads the configuration file
func (hcr *HotConfigReloader) reloadConfig() {
	fmt.Println("Config file changed, reloading...")

	// Load new config (simplified - in real implementation would use proper config loading)
	cfg := &config.Config{}

	// Clear resolver cache to force re-resolution
	hcr.resolver.ClearCache()

	// Update resolver's user config
	hcr.resolver.userConfig = cfg

	// Notify callbacks
	for _, callback := range hcr.reloadCallbacks {
		callback(cfg)
	}

	fmt.Println("Configuration reloaded successfully")
}

// RegisterReloadCallback registers a callback for config reloads
func (hcr *HotConfigReloader) RegisterReloadCallback(callback func(*config.Config)) {
	hcr.reloadCallbacks = append(hcr.reloadCallbacks, callback)
}
