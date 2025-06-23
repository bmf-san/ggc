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
	pullCalled           bool
	pullArgs             []string
	pushCalled           bool
	pushArgs             []string
	resetCalled          bool
	resetArgs            []string
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
