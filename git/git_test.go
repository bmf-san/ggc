package git

import (
	"os/exec"
	"strings"
	"testing"
)

func TestClient_GetGitStatus(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		err     error
		want    string
		wantErr bool
	}{
		{
			name:    "正常系：変更なし",
			output:  "",
			err:     nil,
			want:    "",
			wantErr: false,
		},
		{
			name:    "正常系：変更あり",
			output:  " M file.go\n?? new.go\n",
			err:     nil,
			want:    " M file.go\n?? new.go\n",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					if name != "git" || !strings.Contains(strings.Join(arg, " "), "status --porcelain") {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					return helperCommand(t, tt.output, tt.err)
				},
			}

			got, err := c.GetGitStatus()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGitStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGitStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetBranchName(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		err     error
		want    string
		wantErr bool
	}{
		{
			name:    "正常系：mainブランチ",
			output:  "main\n",
			err:     nil,
			want:    "main",
			wantErr: false,
		},
		{
			name:    "正常系：featureブランチ",
			output:  "feature/test\n",
			err:     nil,
			want:    "feature/test",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					if name != "git" || !strings.Contains(strings.Join(arg, " "), "rev-parse --abbrev-ref HEAD") {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					return helperCommand(t, tt.output, tt.err)
				},
			}

			got, err := c.GetBranchName()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBranchName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBranchName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_CheckoutNewBranch(t *testing.T) {
	tests := []struct {
		name    string
		branch  string
		err     error
		wantErr bool
	}{
		{
			name:    "正常系：新規ブランチ作成",
			branch:  "feature/test",
			err:     nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					if name != "git" || !strings.Contains(strings.Join(arg, " "), "checkout -b "+tt.branch) {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					return helperCommand(t, "", tt.err)
				},
			}

			if err := c.CheckoutNewBranch(tt.branch); (err != nil) != tt.wantErr {
				t.Errorf("CheckoutNewBranch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_GetCurrentBranch(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		err     error
		want    string
		wantErr bool
	}{
		{
			name:    "正常系：mainブランチ",
			output:  "main\n",
			err:     nil,
			want:    "main",
			wantErr: false,
		},
		{
			name:    "正常系：featureブランチ",
			output:  "feature/test\n",
			err:     nil,
			want:    "feature/test",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					if name != "git" || !strings.Contains(strings.Join(arg, " "), "rev-parse --abbrev-ref HEAD") {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					return helperCommand(t, tt.output, tt.err)
				},
			}

			got, err := c.GetCurrentBranch()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentBranch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCurrentBranch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_LogGraph(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{
			name:    "正常系：ログ表示",
			err:     nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				execCommand: func(name string, arg ...string) *exec.Cmd {
					if name != "git" || !strings.Contains(strings.Join(arg, " "), "log --graph --oneline --decorate --all") {
						t.Errorf("unexpected command: %s %v", name, arg)
					}
					return helperCommand(t, "", tt.err)
				},
			}

			if err := c.LogGraph(); (err != nil) != tt.wantErr {
				t.Errorf("LogGraph() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// helperCommand はテスト用のモックコマンドを生成します
func helperCommand(t *testing.T, output string, err error) *exec.Cmd {
	t.Helper()
	if err != nil {
		return exec.Command("false")
	}
	return fakeExecCommand(output)
}

// fakeExecCommand は指定された出力を返すモックコマンドを生成します
func fakeExecCommand(output string) *exec.Cmd {
	return exec.Command("echo", "-n", output)
}
