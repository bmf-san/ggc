package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bmf-san/ggc/v4/config"
)

// Hooker handles git hook operations.
type Hooker struct {
	outputWriter io.Writer
	helper       *Helper
	execCommand  func(string, ...string) *exec.Cmd
}

// NewHooker creates a new Hooker instance.
func NewHooker() *Hooker {
	return &Hooker{
		outputWriter: os.Stdout,
		helper:       NewHelper(),
		execCommand:  exec.Command,
	}
}

// Hook executes git hook commands with the given arguments.
func (h *Hooker) Hook(args []string) {
	if len(args) == 0 {
		h.helper.ShowHookHelp()
		return
	}

	switch args[0] {
	case "list":
		h.listHooks()
	case "install":
		if len(args) < 2 {
			_, _ = fmt.Fprintf(h.outputWriter, "Error: hook name required\n")
			h.helper.ShowHookHelp()
			return
		}
		h.installHook(args[1])
	case "uninstall":
		if len(args) < 2 {
			_, _ = fmt.Fprintf(h.outputWriter, "Error: hook name required\n")
			h.helper.ShowHookHelp()
			return
		}
		h.uninstallHook(args[1])
	case "enable":
		if len(args) < 2 {
			_, _ = fmt.Fprintf(h.outputWriter, "Error: hook name required\n")
			h.helper.ShowHookHelp()
			return
		}
		h.enableHook(args[1])
	case "disable":
		if len(args) < 2 {
			_, _ = fmt.Fprintf(h.outputWriter, "Error: hook name required\n")
			h.helper.ShowHookHelp()
			return
		}
		h.disableHook(args[1])
	case "edit":
		if len(args) < 2 {
			_, _ = fmt.Fprintf(h.outputWriter, "Error: hook name required\n")
			h.helper.ShowHookHelp()
			return
		}
		h.editHook(args[1])
	default:
		h.helper.ShowHookHelp()
	}
}

// listHooks shows all available hooks and their status.
func (h *Hooker) listHooks() {
	hooksDir := filepath.Join(".git", "hooks")

	// Check if hooks directory exists
	if _, err := os.Stat(hooksDir); os.IsNotExist(err) {
		_, _ = fmt.Fprintf(h.outputWriter, "No hooks directory found\n")
		return
	}

	// Standard Git hooks
	standardHooks := []string{
		"applypatch-msg", "pre-applypatch", "post-applypatch",
		"pre-commit", "prepare-commit-msg", "commit-msg", "post-commit",
		"pre-rebase", "post-checkout", "post-merge", "pre-push",
		"pre-receive", "update", "post-receive", "post-update",
		"push-to-checkout", "pre-auto-gc", "post-rewrite",
	}

	_, _ = fmt.Fprintf(h.outputWriter, "Git Hooks Status:\n")
	_, _ = fmt.Fprintf(h.outputWriter, "------------------\n")

	for _, hook := range standardHooks {
		hookPath := filepath.Join(hooksDir, hook)
		samplePath := filepath.Join(hooksDir, hook+".sample")

		if _, err := os.Stat(hookPath); err == nil {
			// Check if it's executable
			if info, err := os.Stat(hookPath); err == nil && info.Mode()&0111 != 0 {
				_, _ = fmt.Fprintf(h.outputWriter, "✓ %s (enabled)\n", hook)
			} else {
				_, _ = fmt.Fprintf(h.outputWriter, "✗ %s (disabled)\n", hook)
			}
		} else if _, err := os.Stat(samplePath); err == nil {
			_, _ = fmt.Fprintf(h.outputWriter, "- %s (sample available)\n", hook)
		} else {
			_, _ = fmt.Fprintf(h.outputWriter, "- %s (not installed)\n", hook)
		}
	}
}

// installHook creates a new hook from sample or creates a basic template.
func (h *Hooker) installHook(hookName string) {
	hooksDir := filepath.Join(".git", "hooks")
	hookPath := filepath.Join(hooksDir, hookName)
	samplePath := filepath.Join(hooksDir, hookName+".sample")

	// Check if hook already exists
	if _, err := os.Stat(hookPath); err == nil {
		_, _ = fmt.Fprintf(h.outputWriter, "Hook '%s' already exists\n", hookName)
		return
	}

	// Try to copy from sample first
	if _, err := os.Stat(samplePath); err == nil {
		if err := h.copyFile(samplePath, hookPath); err != nil {
			_, _ = fmt.Fprintf(h.outputWriter, "Error copying sample hook: %v\n", err)
			return
		}
		_, _ = fmt.Fprintf(h.outputWriter, "Hook '%s' installed from sample\n", hookName)
	} else {
		// Create basic template
		template := h.getHookTemplate(hookName)
		if err := os.WriteFile(hookPath, []byte(template), 0755); err != nil {
			_, _ = fmt.Fprintf(h.outputWriter, "Error creating hook: %v\n", err)
			return
		}
		_, _ = fmt.Fprintf(h.outputWriter, "Hook '%s' created with basic template\n", hookName)
	}
}

// uninstallHook removes a hook.
func (h *Hooker) uninstallHook(hookName string) {
	hookPath := filepath.Join(".git", "hooks", hookName)

	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		_, _ = fmt.Fprintf(h.outputWriter, "Hook '%s' is not installed\n", hookName)
		return
	}

	if err := os.Remove(hookPath); err != nil {
		_, _ = fmt.Fprintf(h.outputWriter, "Error removing hook: %v\n", err)
		return
	}

	_, _ = fmt.Fprintf(h.outputWriter, "Hook '%s' uninstalled\n", hookName)
}

// enableHook makes a hook executable.
func (h *Hooker) enableHook(hookName string) {
	hookPath := filepath.Join(".git", "hooks", hookName)

	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		_, _ = fmt.Fprintf(h.outputWriter, "Hook '%s' is not installed\n", hookName)
		return
	}

	if err := os.Chmod(hookPath, 0755); err != nil {
		_, _ = fmt.Fprintf(h.outputWriter, "Error enabling hook: %v\n", err)
		return
	}

	_, _ = fmt.Fprintf(h.outputWriter, "Hook '%s' enabled\n", hookName)
}

// disableHook makes a hook non-executable.
func (h *Hooker) disableHook(hookName string) {
	hookPath := filepath.Join(".git", "hooks", hookName)

	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		_, _ = fmt.Fprintf(h.outputWriter, "Hook '%s' is not installed\n", hookName)
		return
	}

	if err := os.Chmod(hookPath, 0644); err != nil {
		_, _ = fmt.Fprintf(h.outputWriter, "Error disabling hook: %v\n", err)
		return
	}

	_, _ = fmt.Fprintf(h.outputWriter, "Hook '%s' disabled\n", hookName)
}

// editHook opens a hook in the default editor.
func (h *Hooker) editHook(hookName string) {
	hookPath := filepath.Join(".git", "hooks", hookName)

	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		_, _ = fmt.Fprintf(h.outputWriter, "Hook '%s' is not installed\n", hookName)
		return
	}

	val, err := config.NewConfigManager().Get("default.editor")
	editor, ok := val.(string)
	if err != nil || !ok || editor == "" {
		editor = "vi" // fallback
	}

	cmd := exec.Command(editor, hookPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(h.outputWriter, "Error opening editor: %v\n", err)
	}
}

// copyFile copies a file from src to dst.
func (h *Hooker) copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0755)
}

// getHookTemplate returns a basic template for the given hook.
func (h *Hooker) getHookTemplate(hookName string) string {
	switch hookName {
	case "pre-commit":
		return `#!/bin/sh
# Pre-commit hook
# Add your pre-commit checks here

# Example: Run tests
# npm test
# go test ./...

# Example: Run linter
# golangci-lint run

exit 0
`
	case "commit-msg":
		return `#!/bin/sh
# Commit message hook
# Validates commit message format

commit_regex='^(feat|fix|docs|style|refactor|test|chore)(\(.+\))?: .{1,50}'

if ! grep -qE "$commit_regex" "$1"; then
    echo "Invalid commit message format!"
    echo "Format: type(scope): description"
    echo "Example: feat(auth): add user authentication"
    exit 1
fi

exit 0
`
	case "pre-push":
		return `#!/bin/sh
# Pre-push hook
# Add your pre-push checks here

# Example: Run tests before push
# npm test
# go test ./...

# Example: Prevent push to main/master
protected_branch='main'
current_branch=$(git symbolic-ref HEAD | sed -e 's,.*/\(.*\),\1,')

if [ "$current_branch" = "$protected_branch" ]; then
    echo "Direct push to $protected_branch is not allowed"
    exit 1
fi

exit 0
`
	default:
		return fmt.Sprintf(`#!/bin/sh
# %s hook
# Add your %s logic here

exit 0
`, hookName, strings.ReplaceAll(hookName, "-", " "))
	}
}
