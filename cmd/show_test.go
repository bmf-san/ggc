package cmd

import (
	"bytes"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v8/internal/testutil"
)

type mockShowGitClient struct {
	testutil.MockGitClient
	called  bool
	gotArgs []string
	showErr error
}

func (m *mockShowGitClient) Show(args []string) error {
	m.called = true
	m.gotArgs = slices.Clone(args)
	return m.showErr
}

func TestShower_Show(t *testing.T) {
	cases := []struct {
		name         string
		args         []string
		expectCall   bool
		expectArgs   []string
		expectOutput string
	}{
		{
			name:       "no args shows HEAD",
			args:       []string{},
			expectCall: true,
			expectArgs: []string{},
		},
		{
			name:       "single ref",
			args:       []string{"HEAD~1"},
			expectCall: true,
			expectArgs: []string{"HEAD~1"},
		},
		{
			name:       "stat option with ref",
			args:       []string{"--stat", "abc123"},
			expectCall: true,
			expectArgs: []string{"--stat", "abc123"},
		},
		{
			name:         "help subcommand prints usage",
			args:         []string{"help"},
			expectCall:   false,
			expectOutput: "ggc show",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			mock := &mockShowGitClient{}
			s := &Shower{
				gitClient:    mock,
				outputWriter: &buf,
				helper:       NewHelper(),
			}
			s.helper.outputWriter = &buf

			s.Show(tc.args)

			if mock.called != tc.expectCall {
				t.Errorf("Show called = %v, want %v", mock.called, tc.expectCall)
			}
			if tc.expectCall && !slices.Equal(mock.gotArgs, tc.expectArgs) {
				t.Errorf("Show args = %v, want %v", mock.gotArgs, tc.expectArgs)
			}
			if tc.expectOutput != "" && !strings.Contains(buf.String(), tc.expectOutput) {
				t.Errorf("output = %q, want substring %q", buf.String(), tc.expectOutput)
			}
		})
	}
}

func TestShower_Show_Error(t *testing.T) {
	var buf bytes.Buffer
	mock := &mockShowGitClient{showErr: errors.New("bad object")}
	s := &Shower{
		gitClient:    mock,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	s.Show([]string{"deadbeef"})

	output := buf.String()
	if !strings.Contains(output, "Error") {
		t.Errorf("expected error output, got %q", output)
	}
	if !strings.Contains(output, "bad object") {
		t.Errorf("expected underlying error in output, got %q", output)
	}
}

func TestNewShower_Defaults(t *testing.T) {
	s := NewShower(&mockShowGitClient{})
	if s.outputWriter == nil {
		t.Error("outputWriter should default to os.Stdout")
	}
	if s.helper == nil {
		t.Error("helper should be initialized")
	}
}
