package router

import (
	"testing"
)

type mockExecuter struct {
	helpCalled           bool
	branchCalled         bool
	branchArgs           []string
	commitCalled         bool
	commitArgs           []string
	logCalled            bool
	logArgs              []string
	diffCalled           bool
	diffArgs             []string
	statusCalled         bool
	statusArgs           []string
	pullCalled           bool
	pullArgs             []string
	pushCalled           bool
	pushArgs             []string
	resetCalled          bool
	resetArgs            []string
	tagCalled            bool
	tagArgs              []string
	cleanCalled          bool
	cleanArgs            []string
	pullRebasePushCalled bool
	interactiveCalled    bool
}

func (m *mockExecuter) Help() {
	m.helpCalled = true
}

func (m *mockExecuter) Branch(args []string) {
	m.branchCalled = true
	m.branchArgs = args
}

func (m *mockExecuter) Commit(args []string) {
	m.commitCalled = true
	m.commitArgs = args
}

func (m *mockExecuter) Log(args []string) {
	m.logCalled = true
	m.logArgs = args
}

func (m *mockExecuter) Status(args []string) {
	m.statusCalled = true
	m.statusArgs = args
}

func (m *mockExecuter) Diff(args []string) {
	m.diffCalled = true
	m.diffArgs = args
}

func (m *mockExecuter) Tag(args []string) {
	m.tagCalled = true
	m.tagArgs = args
}

func (m *mockExecuter) Pull(args []string) {
	m.pullCalled = true
	m.pullArgs = args
}

func (m *mockExecuter) Push(args []string) {
	m.pushCalled = true
	m.pushArgs = args
}

func (m *mockExecuter) Reset(args []string) {
	m.resetCalled = true
	m.resetArgs = args
}

func (m *mockExecuter) Clean(args []string) {
	m.cleanCalled = true
	m.cleanArgs = args
}

func (m *mockExecuter) PullRebasePush() {
	m.pullRebasePushCalled = true
}

func (m *mockExecuter) Interactive() {
	m.interactiveCalled = true
}

func TestRouter(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		validate func(t *testing.T, m *mockExecuter)
	}{
		{
			name: "help",
			args: []string{"help"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.helpCalled {
					t.Error("Help should be called")
				}
			},
		},
		{
			name: "branch",
			args: []string{"branch", "current"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.branchCalled {
					t.Error("Branch should be called")
				}
				if len(m.branchArgs) != 1 || m.branchArgs[0] != "current" {
					t.Errorf("unexpected branch args: got %v", m.branchArgs)
				}
			},
		},
		{
			name: "commit",
			args: []string{"commit", "allow-empty"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.commitCalled {
					t.Error("Commit should be called")
				}
				if len(m.commitArgs) != 1 || m.commitArgs[0] != "allow-empty" {
					t.Errorf("unexpected commit args: got %v", m.commitArgs)
				}
			},
		},
		{
			name: "log",
			args: []string{"log", "simple"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.logCalled {
					t.Error("Log should be called")
				}
				if len(m.logArgs) != 1 || m.logArgs[0] != "simple" {
					t.Errorf("unexpected log args: got %v", m.logArgs)
				}
			},
		},
		{
			name: "pull",
			args: []string{"pull", "rebase"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.pullCalled {
					t.Error("Pull should be called")
				}
				if len(m.pullArgs) != 1 || m.pullArgs[0] != "rebase" {
					t.Errorf("unexpected pull args: got %v", m.pullArgs)
				}
			},
		},
		{
			name: "push",
			args: []string{"push", "force"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.pushCalled {
					t.Error("Push should be called")
				}
				if len(m.pushArgs) != 1 || m.pushArgs[0] != "force" {
					t.Errorf("unexpected push args: got %v", m.pushArgs)
				}
			},
		},
		{
			name: "reset",
			args: []string{"reset"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.resetCalled {
					t.Error("Reset should be called")
				}
			},
		},
        {
			name: "status no args",
			args: []string{"status"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.statusCalled {
					t.Error("Status should be called")
				}
				if len(m.statusArgs) != 0 {
					t.Errorf("unexpected status args: got %v, expected empty", m.statusArgs)
				}
			},
		},
		{
			name: "status with short arg",
			args: []string{"status", "short"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.statusCalled {
					t.Error("Status should be called")
				}
				if len(m.statusArgs) != 1 || m.statusArgs[0] != "short" {
					t.Errorf("unexpected status args: got %v, expected [short]", m.statusArgs)
				}
			},
		},
        {
			name: "tag no args",
			args: []string{"tag"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.tagCalled {
					t.Error("Tag should be called")
				}
				if len(m.tagArgs) != 0 {
					t.Errorf("unexpected status args: got %v, expected empty", m.tagArgs)
				}
			},
		},
		{
			name: "tag with arg",
			args: []string{"tag", "list"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.tagCalled {
					t.Error("Tag should be called")
				}
				if len(m.tagArgs) != 1 || m.tagArgs[0] != "list" {
					t.Errorf("unexpected tag args: got %v, expected [list]", m.tagArgs)
				}
			},
		},
		{
			name: "diff no args",
			args: []string{"diff"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.diffCalled {
					t.Error("Diff should be called")
				}
				if len(m.diffArgs) != 0 {
					t.Errorf("unexpected diff args: got %v, expected empty", m.diffArgs)
				}
			},
		},
		{
			name: "diff unstaged",
			args: []string{"diff", "unstaged"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.diffCalled {
					t.Error("Diff should be called")
				}
				if len(m.diffArgs) != 1 || m.diffArgs[0] != "unstaged" {
					t.Errorf("unexpected diff args: got %v, expected [unstaged]", m.diffArgs)
				}
			},
		},
		{
			name: "diff staged",
			args: []string{"diff", "staged"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.diffCalled {
					t.Error("Diff should be called")
				}
				if len(m.diffArgs) != 1 || m.diffArgs[0] != "staged" {
					t.Errorf("unexpected diff args: got %v, expected [staged]", m.diffArgs)
				}
			},
		},
		{
			name: "clean",
			args: []string{"clean", "files"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.cleanCalled {
					t.Error("Clean should be called")
				}
				if len(m.cleanArgs) != 1 || m.cleanArgs[0] != "files" {
					t.Errorf("unexpected clean args: got %v", m.cleanArgs)
				}
			},
		},
		{
			name: "pull-rebase-push",
			args: []string{"pull-rebase-push"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.pullRebasePushCalled {
					t.Error("PullRebasePush should be called")
				}
			},
		},
		{
			name: "unknown",
			args: []string{"unknown"},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.helpCalled {
					t.Error("Help should be called")
				}
			},
		},
		{
			name: "empty",
			args: []string{},
			validate: func(t *testing.T, m *mockExecuter) {
				if !m.interactiveCalled {
					t.Error("Interactive should be called")
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockExecuter{}
			r := NewRouter(m)
			r.Route(tc.args)
			tc.validate(t, m)
		})
	}
}
