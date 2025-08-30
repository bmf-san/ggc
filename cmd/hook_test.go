package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestHooker_Hook(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		expectedOutput string
		setupFiles     func(string) error
		cleanupFiles   func(string) error
	}{
		{
			name:           "no args shows help",
			args:           []string{},
			expectedOutput: "Usage: ggc hook",
		},
		{
			name:           "invalid command shows help",
			args:           []string{"invalid"},
			expectedOutput: "Usage: ggc hook",
		},
		{
			name:           "install command without hook name",
			args:           []string{"install"},
			expectedOutput: "Error: hook name required",
		},
		{
			name:           "uninstall command without hook name",
			args:           []string{"uninstall"},
			expectedOutput: "Error: hook name required",
		},
		{
			name:           "enable command without hook name",
			args:           []string{"enable"},
			expectedOutput: "Error: hook name required",
		},
		{
			name:           "disable command without hook name",
			args:           []string{"disable"},
			expectedOutput: "Error: hook name required",
		},
		{
			name:           "edit command without hook name",
			args:           []string{"edit"},
			expectedOutput: "Error: hook name required",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			h := &Hooker{
				outputWriter: &buf,
				helper:       NewHelper(),
				execCommand:  exec.Command,
			}
			h.helper.outputWriter = &buf

			h.Hook(tc.args)

			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}

func TestHooker_listHooks(t *testing.T) {
	cases := []struct {
		name           string
		setupFiles     func(string) error
		cleanupFiles   func(string) error
		expectedOutput string
	}{
		{
			name: "no hooks directory",
			setupFiles: func(_ string) error {
				return nil
			},
			cleanupFiles: func(_ string) error {
				return nil
			},
			expectedOutput: "No hooks directory found",
		},
		{
			name: "hooks directory with sample and active hooks",
			setupFiles: func(tmpDir string) error {
				hooksDir := filepath.Join(tmpDir, ".git", "hooks")
				if err := os.MkdirAll(hooksDir, 0755); err != nil {
					return err
				}

				// Create a sample hook
				samplePath := filepath.Join(hooksDir, "pre-commit.sample")
				if err := os.WriteFile(samplePath, []byte("#!/bin/sh\necho sample"), 0644); err != nil {
					return err
				}

				// Create an active executable hook
				activePath := filepath.Join(hooksDir, "post-commit")
				if err := os.WriteFile(activePath, []byte("#!/bin/sh\necho active"), 0755); err != nil {
					return err
				}

				// Create a disabled hook
				disabledPath := filepath.Join(hooksDir, "pre-push")
				if err := os.WriteFile(disabledPath, []byte("#!/bin/sh\necho disabled"), 0644); err != nil {
					return err
				}

				return nil
			},
			cleanupFiles:   os.RemoveAll,
			expectedOutput: "Git Hooks Status:",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			originalDir, _ := os.Getwd()

			if tc.setupFiles != nil {
				if err := tc.setupFiles(tmpDir); err != nil {
					t.Fatalf("failed to setup files: %v", err)
				}
			}

			err := os.Chdir(tmpDir)
			if err != nil {
				t.Fatalf("failed to chdir: %v", err)
			}
			defer func() {
				if err := os.Chdir(originalDir); err != nil {
					t.Fatalf("failed to chdir: %v", err)
				}
			}()

			var buf bytes.Buffer
			h := &Hooker{
				outputWriter: &buf,
				helper:       NewHelper(),
				execCommand:  exec.Command,
			}

			h.listHooks()

			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}

			if tc.cleanupFiles != nil {
				err := tc.cleanupFiles(tmpDir)
				if err != nil {
					t.Errorf("error during file cleanup: %v", err)
				}
			}
		})
	}
}

func TestHooker_installHook(t *testing.T) {
	cases := []struct {
		name           string
		hookName       string
		setupFiles     func(string) error
		expectedOutput string
	}{
		{
			name:     "install hook from sample",
			hookName: "pre-commit",
			setupFiles: func(tmpDir string) error {
				hooksDir := filepath.Join(tmpDir, ".git", "hooks")
				if err := os.MkdirAll(hooksDir, 0755); err != nil {
					return err
				}

				samplePath := filepath.Join(hooksDir, "pre-commit.sample")
				return os.WriteFile(samplePath, []byte("#!/bin/sh\necho sample"), 0644)
			},
			expectedOutput: "Hook 'pre-commit' installed from sample",
		},
		{
			name:     "install hook with template",
			hookName: "commit-msg",
			setupFiles: func(tmpDir string) error {
				hooksDir := filepath.Join(tmpDir, ".git", "hooks")
				return os.MkdirAll(hooksDir, 0755)
			},
			expectedOutput: "Hook 'commit-msg' created with basic template",
		},
		{
			name:     "hook already exists",
			hookName: "pre-push",
			setupFiles: func(tmpDir string) error {
				hooksDir := filepath.Join(tmpDir, ".git", "hooks")
				if err := os.MkdirAll(hooksDir, 0755); err != nil {
					return err
				}

				hookPath := filepath.Join(hooksDir, "pre-push")
				return os.WriteFile(hookPath, []byte("#!/bin/sh\necho existing"), 0755)
			},
			expectedOutput: "Hook 'pre-push' already exists",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			originalDir, _ := os.Getwd()

			if tc.setupFiles != nil {
				if err := tc.setupFiles(tmpDir); err != nil {
					t.Fatalf("failed to setup files: %v", err)
				}
			}

			err := os.Chdir(tmpDir)
			if err != nil {
				t.Fatalf("failed to chdir: %v", err)
			}
			defer func() {
				if err := os.Chdir(originalDir); err != nil {
					t.Fatalf("failed to chdir: %v", err)
				}
			}()

			var buf bytes.Buffer
			h := &Hooker{
				outputWriter: &buf,
				helper:       NewHelper(),
				execCommand:  exec.Command,
			}

			h.installHook(tc.hookName)

			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}

func TestHooker_uninstallHook(t *testing.T) {
	cases := []struct {
		name           string
		hookName       string
		setupFiles     func(string) error
		expectedOutput string
	}{
		{
			name:     "uninstall existing hook",
			hookName: "pre-commit",
			setupFiles: func(tmpDir string) error {
				hooksDir := filepath.Join(tmpDir, ".git", "hooks")
				if err := os.MkdirAll(hooksDir, 0755); err != nil {
					return err
				}

				hookPath := filepath.Join(hooksDir, "pre-commit")
				return os.WriteFile(hookPath, []byte("#!/bin/sh\necho test"), 0755)
			},
			expectedOutput: "Hook 'pre-commit' uninstalled",
		},
		{
			name:     "uninstall non-existent hook",
			hookName: "non-existent",
			setupFiles: func(tmpDir string) error {
				hooksDir := filepath.Join(tmpDir, ".git", "hooks")
				return os.MkdirAll(hooksDir, 0755)
			},
			expectedOutput: "Hook 'non-existent' is not installed",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			originalDir, _ := os.Getwd()

			if tc.setupFiles != nil {
				if err := tc.setupFiles(tmpDir); err != nil {
					t.Fatalf("failed to setup files: %v", err)
				}
			}

			err := os.Chdir(tmpDir)
			if err != nil {
				t.Fatalf("failed to chdir: %v", err)
			}
			defer func() {
				if err := os.Chdir(originalDir); err != nil {
					t.Fatalf("failed to chdir: %v", err)
				}
			}()

			var buf bytes.Buffer
			h := &Hooker{
				outputWriter: &buf,
				helper:       NewHelper(),
				execCommand:  exec.Command,
			}

			h.uninstallHook(tc.hookName)

			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}

func TestHooker_enableHook(t *testing.T) {
	cases := []struct {
		name           string
		hookName       string
		setupFiles     func(string) error
		expectedOutput string
	}{
		{
			name:     "enable existing hook",
			hookName: "pre-commit",
			setupFiles: func(tmpDir string) error {
				hooksDir := filepath.Join(tmpDir, ".git", "hooks")
				if err := os.MkdirAll(hooksDir, 0755); err != nil {
					return err
				}

				hookPath := filepath.Join(hooksDir, "pre-commit")
				return os.WriteFile(hookPath, []byte("#!/bin/sh\necho test"), 0644)
			},
			expectedOutput: "Hook 'pre-commit' enabled",
		},
		{
			name:     "enable non-existent hook",
			hookName: "non-existent",
			setupFiles: func(tmpDir string) error {
				hooksDir := filepath.Join(tmpDir, ".git", "hooks")
				return os.MkdirAll(hooksDir, 0755)
			},
			expectedOutput: "Hook 'non-existent' is not installed",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			originalDir, _ := os.Getwd()

			if tc.setupFiles != nil {
				if err := tc.setupFiles(tmpDir); err != nil {
					t.Fatalf("failed to setup files: %v", err)
				}
			}

			err := os.Chdir(tmpDir)
			if err != nil {
				t.Fatalf("failed to chdir: %v", err)
			}
			defer func() {
				if err := os.Chdir(originalDir); err != nil {
					t.Fatalf("failed to chdir: %v", err)
				}
			}()

			var buf bytes.Buffer
			h := &Hooker{
				outputWriter: &buf,
				helper:       NewHelper(),
				execCommand:  exec.Command,
			}

			h.enableHook(tc.hookName)

			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}

func TestHooker_disableHook(t *testing.T) {
	cases := []struct {
		name           string
		hookName       string
		setupFiles     func(string) error
		expectedOutput string
	}{
		{
			name:     "disable existing hook",
			hookName: "pre-commit",
			setupFiles: func(tmpDir string) error {
				hooksDir := filepath.Join(tmpDir, ".git", "hooks")
				if err := os.MkdirAll(hooksDir, 0755); err != nil {
					return err
				}

				hookPath := filepath.Join(hooksDir, "pre-commit")
				return os.WriteFile(hookPath, []byte("#!/bin/sh\necho test"), 0755)
			},
			expectedOutput: "Hook 'pre-commit' disabled",
		},
		{
			name:     "disable non-existent hook",
			hookName: "non-existent",
			setupFiles: func(tmpDir string) error {
				hooksDir := filepath.Join(tmpDir, ".git", "hooks")
				return os.MkdirAll(hooksDir, 0755)
			},
			expectedOutput: "Hook 'non-existent' is not installed",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			originalDir, _ := os.Getwd()

			if tc.setupFiles != nil {
				if err := tc.setupFiles(tmpDir); err != nil {
					t.Fatalf("failed to setup files: %v", err)
				}
			}

			err := os.Chdir(tmpDir)
			if err != nil {
				t.Fatalf("failed to chdir: %v", err)
			}
			defer func() {
				if err := os.Chdir(originalDir); err != nil {
					t.Fatalf("failed to chdir: %v", err)
				}
			}()

			var buf bytes.Buffer
			h := &Hooker{
				outputWriter: &buf,
				helper:       NewHelper(),
				execCommand:  exec.Command,
			}

			h.disableHook(tc.hookName)

			output := buf.String()
			if !strings.Contains(output, tc.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tc.expectedOutput, output)
			}
		})
	}
}

func TestHooker_getHookTemplate(t *testing.T) {
	cases := []struct {
		name           string
		hookName       string
		expectedOutput string
	}{
		{
			name:           "pre-commit template",
			hookName:       "pre-commit",
			expectedOutput: "# Pre-commit hook",
		},
		{
			name:           "commit-msg template",
			hookName:       "commit-msg",
			expectedOutput: "# Commit message hook",
		},
		{
			name:           "pre-push template",
			hookName:       "pre-push",
			expectedOutput: "# Pre-push hook",
		},
		{
			name:           "default template",
			hookName:       "custom-hook",
			expectedOutput: "# custom-hook hook",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := &Hooker{
				outputWriter: &bytes.Buffer{},
				helper:       NewHelper(),
				execCommand:  exec.Command,
			}

			template := h.getHookTemplate(tc.hookName)

			if !strings.Contains(template, tc.expectedOutput) {
				t.Errorf("expected template to contain %q, got %q", tc.expectedOutput, template)
			}
		})
	}
}

func TestHooker_copyFile(t *testing.T) {
	tmpDir := t.TempDir()

	srcPath := filepath.Join(tmpDir, "source.txt")
	dstPath := filepath.Join(tmpDir, "destination.txt")

	content := "test content"
	if err := os.WriteFile(srcPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	h := &Hooker{
		outputWriter: &bytes.Buffer{},
		helper:       NewHelper(),
		execCommand:  exec.Command,
	}

	err := h.copyFile(srcPath, dstPath)
	if err != nil {
		t.Errorf("copyFile failed: %v", err)
	}

	copiedContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Errorf("failed to read copied file: %v", err)
	}

	if string(copiedContent) != content {
		t.Errorf("expected copied content %q, got %q", content, string(copiedContent))
	}

	// Check file permissions
	info, err := os.Stat(dstPath)
	if err != nil {
		t.Errorf("failed to stat copied file: %v", err)
	}

	if info.Mode().Perm() != 0755 {
		t.Errorf("expected file permissions 0755, got %v", info.Mode().Perm())
	}
}
