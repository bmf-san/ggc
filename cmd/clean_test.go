package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/git"
)

// mockGitClient for clean_test
type mockCleanGitClient struct {
	git.Clienter
	cleanFilesErr    error
	cleanDirsErr     error
	cleanFilesCalled bool
	cleanDirsCalled  bool
}

func (m *mockCleanGitClient) CleanFiles() error {
	m.cleanFilesCalled = true
	return m.cleanFilesErr
}

func (m *mockCleanGitClient) CleanDirs() error {
	m.cleanDirsCalled = true
	return m.cleanDirsErr
}

func (m *mockCleanGitClient) LogGraph() error {
	return nil
}

func (m *mockCleanGitClient) StashPullPop() error {
	return nil
}

func TestCleaner_Clean(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantCleanFiles bool
		wantCleanDirs  bool
	}{
		{
			name:           "clean files",
			args:           []string{"files"},
			wantCleanFiles: true,
			wantCleanDirs:  false,
		},
		{
			name:           "clean dirs",
			args:           []string{"dirs"},
			wantCleanFiles: false,
			wantCleanDirs:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockCleanGitClient{}
			var buf bytes.Buffer
			cleaner := NewCleanerWithClient(mockClient)
			cleaner.outputWriter = &buf
			cleaner.Clean(tt.args)

			if mockClient.cleanFilesCalled != tt.wantCleanFiles {
				t.Errorf("CleanFiles called = %v, want %v", mockClient.cleanFilesCalled, tt.wantCleanFiles)
			}
			if mockClient.cleanDirsCalled != tt.wantCleanDirs {
				t.Errorf("CleanDirs called = %v, want %v", mockClient.cleanDirsCalled, tt.wantCleanDirs)
			}
		})
	}
}

func TestCleaner_Clean_Files(t *testing.T) {
	var buf bytes.Buffer
	mock := &mockCleanGitClient{}
	cleaner := &Cleaner{
		gitClient:    mock,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.Clean([]string{"files"})

	if buf.Len() > 0 {
		t.Errorf("Expected no output, got %q", buf.String())
	}
}

func TestCleaner_Clean_Files_Error(t *testing.T) {
	var buf bytes.Buffer
	mock := &mockCleanGitClient{cleanFilesErr: errors.New("failed to clean files")}
	cleaner := &Cleaner{
		gitClient:    mock,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.Clean([]string{"files"})

	expected := "Error: failed to clean files\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestCleaner_Clean_Dirs(t *testing.T) {
	var buf bytes.Buffer
	mock := &mockCleanGitClient{}
	cleaner := &Cleaner{
		gitClient:    mock,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.Clean([]string{"dirs"})

	if buf.Len() > 0 {
		t.Errorf("Expected no output, got %q", buf.String())
	}
}

func TestCleaner_Clean_Dirs_Error(t *testing.T) {
	var buf bytes.Buffer
	mock := &mockCleanGitClient{cleanDirsErr: errors.New("failed to clean directories")}
	cleaner := &Cleaner{
		gitClient:    mock,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.Clean([]string{"dirs"})

	expected := "Error: failed to clean directories\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestCleaner_Clean_Help(t *testing.T) {
	var buf bytes.Buffer
	cleaner := &Cleaner{
		gitClient:    &mockCleanGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.Clean([]string{})

	output := buf.String()
	if output == "" || !bytes.Contains(buf.Bytes(), []byte("Usage")) {
		t.Errorf("Usage should be displayed, but got: %s", output)
	}
}

func TestCleaner_CleanInteractive_NoFiles(t *testing.T) {
	var buf bytes.Buffer
	cleaner := &Cleaner{
		gitClient:    &mockCleanGitClient{},
		outputWriter: &buf,
		execCommand: func(_ string, _ ...string) *exec.Cmd {
			return exec.Command("echo", "")
		},
		helper: NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.CleanInteractive()

	expected := "No files to clean.\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestCleaner_CleanInteractive_WithFiles(t *testing.T) {
	var buf bytes.Buffer
	inputBuf := strings.NewReader("all\n")
	cleaner := &Cleaner{
		gitClient:    &mockCleanGitClient{},
		outputWriter: &buf,
		execCommand: func(_ string, arg ...string) *exec.Cmd {
			if len(arg) > 0 && arg[0] == "clean" && arg[1] == "-nd" {
				return exec.Command("echo", "Would remove file1.txt\nWould remove file2.txt")
			}
			return exec.Command("echo", "")
		},
		inputReader: bufio.NewReader(inputBuf),
		helper:      NewHelper(),
	}
	cleaner.helper.outputWriter = &buf

	cleaner.CleanInteractive()

	if !strings.Contains(buf.String(), "Selected files deleted.") {
		t.Error("Expected output to contain 'Selected files deleted.'")
	}
}
