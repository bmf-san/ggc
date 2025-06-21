package router

import (
	"testing"

	"github.com/bmf-san/ggc/cmd"
)

type mockExecuter struct {
	cmd.Executer
	addCalled                   bool
	addCommitPushCalled         bool
	branchCalled                bool
	cleanCalled                 bool
	cleanInteractiveCalled      bool
	commitCalled                bool
	commitPushInteractiveCalled bool
	completeCalled              bool
	fetchCalled                 bool
	logCalled                   bool
	pullCalled                  bool
	pullRebasePushCalled        bool
	pushCalled                  bool
	rebaseCalled                bool
	remoteCalled                bool
	resetCalled                 bool
	stashCalled                 bool
	stashPullPopCalled          bool
	showHelpCalled              bool
}

func (m *mockExecuter) Add(args []string)      { m.addCalled = true }
func (m *mockExecuter) AddCommitPush()         { m.addCommitPushCalled = true }
func (m *mockExecuter) Branch(args []string)   { m.branchCalled = true }
func (m *mockExecuter) Clean(args []string)    { m.cleanCalled = true }
func (m *mockExecuter) CleanInteractive()      { m.cleanInteractiveCalled = true }
func (m *mockExecuter) Commit(args []string)   { m.commitCalled = true }
func (m *mockExecuter) CommitPushInteractive() { m.commitPushInteractiveCalled = true }
func (m *mockExecuter) Complete(args []string) { m.completeCalled = true }
func (m *mockExecuter) Fetch(args []string)    { m.fetchCalled = true }
func (m *mockExecuter) Log(args []string)      { m.logCalled = true }
func (m *mockExecuter) Pull(args []string)     { m.pullCalled = true }
func (m *mockExecuter) PullRebasePush()        { m.pullRebasePushCalled = true }
func (m *mockExecuter) Push(args []string)     { m.pushCalled = true }
func (m *mockExecuter) Rebase(args []string)   { m.rebaseCalled = true }
func (m *mockExecuter) Remote(args []string)   { m.remoteCalled = true }
func (m *mockExecuter) Reset(args []string)    { m.resetCalled = true }
func (m *mockExecuter) Stash(args []string)    { m.stashCalled = true }
func (m *mockExecuter) StashPullPop()          { m.stashPullPopCalled = true }
func (m *mockExecuter) ShowHelp()              { m.showHelpCalled = true }

func TestRouter(t *testing.T) {
	cases := []struct {
		args     []string
		expected func(m *mockExecuter) bool
	}{
		{[]string{"ggc", "add"}, func(m *mockExecuter) bool { return m.addCalled }},
		{[]string{"ggc", "add-commit-push"}, func(m *mockExecuter) bool { return m.addCommitPushCalled }},
		{[]string{"ggc", "branch"}, func(m *mockExecuter) bool { return m.branchCalled }},
		{[]string{"ggc", "clean"}, func(m *mockExecuter) bool { return m.cleanCalled }},
		{[]string{"ggc", "clean", "interactive"}, func(m *mockExecuter) bool { return m.cleanInteractiveCalled }},
		{[]string{"ggc", "commit"}, func(m *mockExecuter) bool { return m.commitCalled }},
		{[]string{"ggc", "commit-push"}, func(m *mockExecuter) bool { return m.commitPushInteractiveCalled }},
		{[]string{"ggc", "__complete"}, func(m *mockExecuter) bool { return m.completeCalled }},
		{[]string{"ggc", "fetch"}, func(m *mockExecuter) bool { return m.fetchCalled }},
		{[]string{"ggc", "log"}, func(m *mockExecuter) bool { return m.logCalled }},
		{[]string{"ggc", "pull"}, func(m *mockExecuter) bool { return m.pullCalled }},
		{[]string{"ggc", "prp"}, func(m *mockExecuter) bool { return m.pullRebasePushCalled }},
		{[]string{"ggc", "pull-rebase-push"}, func(m *mockExecuter) bool { return m.pullRebasePushCalled }},
		{[]string{"ggc", "push"}, func(m *mockExecuter) bool { return m.pushCalled }},
		{[]string{"ggc", "rebase"}, func(m *mockExecuter) bool { return m.rebaseCalled }},
		{[]string{"ggc", "remote"}, func(m *mockExecuter) bool { return m.remoteCalled }},
		{[]string{"ggc", "reset-clean"}, func(m *mockExecuter) bool { return m.resetCalled }},
		{[]string{"ggc", "sp"}, func(m *mockExecuter) bool { return m.stashPullPopCalled }},
		{[]string{"ggc", "stash-pull-pop"}, func(m *mockExecuter) bool { return m.stashPullPopCalled }},
		{[]string{"ggc", "stash"}, func(m *mockExecuter) bool { return m.stashCalled }},
		{[]string{"ggc"}, func(m *mockExecuter) bool { return m.showHelpCalled }},
		{[]string{"ggc", "unknown"}, func(m *mockExecuter) bool { return m.showHelpCalled }},
	}

	for _, tc := range cases {
		m := &mockExecuter{}
		r := &Router{
			Executer: m,
		}
		r.Route(tc.args)
		if !tc.expected(m) {
			t.Errorf("args: %v, not called", tc.args)
		}
	}
}
