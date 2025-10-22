package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/bmf-san/ggc/v7/internal/prompt"
	"github.com/bmf-san/ggc/v7/pkg/git"
)

// mockBranchGitClient is a mock implementation focused on branch operations for tests
type mockBranchGitClient struct {
	getCurrentBranchCalled bool
	currentBranch          string
	err                    error
	listLocalBranches      func() ([]string, error)
	listRemoteBranches     func() ([]string, error)
	mergedBranches         []string
	checkoutNewBranchError bool
	createdBranches        []string
	deletedBranches        []string
	ops                    *mockBranchOperations
}

func (m *mockBranchGitClient) GetCurrentBranch() (string, error) {
	m.getCurrentBranchCalled = true
	return m.currentBranch, m.err
}

func (m *mockBranchGitClient) ListLocalBranches() ([]string, error) {
	if m.listLocalBranches != nil {
		return m.listLocalBranches()
	}
	return []string{"main", "feature/test"}, nil
}

func (m *mockBranchGitClient) ListRemoteBranches() ([]string, error) {
	if m.listRemoteBranches != nil {
		return m.listRemoteBranches()
	}
	return []string{"origin/main", "origin/feature/test"}, nil
}

// Branch Operations methods
func (m *mockBranchGitClient) CheckoutNewBranch(branchName string) error {
	m.createdBranches = append(m.createdBranches, branchName)
	if m.checkoutNewBranchError {
		return errors.New("exit status 1")
	}
	// Check if branch already exists
	if branchName == "main" {
		return errors.New("fatal: a branch named 'main' already exists")
	}
	return nil
}
func (m *mockBranchGitClient) CheckoutBranch(_ string) error { return nil }
func (m *mockBranchGitClient) CheckoutNewBranchFromRemote(_, _ string) error {
	return nil
}
func (m *mockBranchGitClient) DeleteBranch(name string) error {
	m.deletedBranches = append(m.deletedBranches, name)
	return nil
}
func (m *mockBranchGitClient) ListMergedBranches() ([]string, error) {
	if m.mergedBranches != nil {
		return m.mergedBranches, nil
	}
	return []string{}, nil
}

// Track calls for better testing
type mockBranchOperations struct {
	renameBranchCalls       []struct{ old, new string }
	moveBranchCalls         []struct{ branch, commit string }
	setUpstreamBranchCalls  []struct{ branch, upstream string }
	renameBranchError       error
	moveBranchError         error
	setUpstreamError        error
	revParseVerifyResult    bool
	sortBranchesCalls       []string
	sortBranchesError       error
	branchesContainingCalls []string
	branchesContainingError error
}

func (m *mockBranchGitClient) RenameBranch(old, new string) error {
	if m.ops == nil {
		m.ops = &mockBranchOperations{}
	}
	m.ops.renameBranchCalls = append(m.ops.renameBranchCalls, struct{ old, new string }{old, new})
	return m.ops.renameBranchError
}

func (m *mockBranchGitClient) MoveBranch(branch, commit string) error {
	if m.ops == nil {
		m.ops = &mockBranchOperations{}
	}
	m.ops.moveBranchCalls = append(m.ops.moveBranchCalls, struct{ branch, commit string }{branch, commit})
	return m.ops.moveBranchError
}

func (m *mockBranchGitClient) SetUpstreamBranch(branch, upstream string) error {
	if m.ops == nil {
		m.ops = &mockBranchOperations{}
	}
	m.ops.setUpstreamBranchCalls = append(m.ops.setUpstreamBranchCalls, struct{ branch, upstream string }{branch, upstream})
	return m.ops.setUpstreamError
}

func (m *mockBranchGitClient) RevParseVerify(ref string) bool {
	if m.ops != nil {
		return m.ops.revParseVerifyResult
	}
	return true // Default to valid
}
func (m *mockBranchGitClient) GetBranchInfo(branch string) (*git.BranchInfo, error) {
	// Provide simple consistent info
	bi := &git.BranchInfo{
		Name:            branch,
		IsCurrentBranch: branch == "main",
		Upstream:        "origin/" + branch,
		AheadBehind:     "ahead 1",
		LastCommitSHA:   "abc1234",
		LastCommitMsg:   "Test commit",
	}
	return bi, nil
}
func (m *mockBranchGitClient) ListBranchesVerbose() ([]git.BranchInfo, error) {
	return []git.BranchInfo{
		{Name: "main", IsCurrentBranch: true, Upstream: "origin/main", LastCommitSHA: "1234567", LastCommitMsg: "msg"},
		{Name: "feature", IsCurrentBranch: false, Upstream: "origin/feature", LastCommitSHA: "89abcde", LastCommitMsg: "msg2"},
	}, nil
}
func (m *mockBranchGitClient) SortBranches(by string) ([]string, error) {
	if m.ops == nil {
		m.ops = &mockBranchOperations{}
	}
	m.ops.sortBranchesCalls = append(m.ops.sortBranchesCalls, by)
	if m.ops.sortBranchesError != nil {
		return nil, m.ops.sortBranchesError
	}
	if by == "date" {
		return []string{"feature/test", "main"}, nil
	}
	return []string{"feature/test", "main"}, nil
}
func (m *mockBranchGitClient) BranchesContaining(commit string) ([]string, error) {
	if m.ops == nil {
		m.ops = &mockBranchOperations{}
	}
	m.ops.branchesContainingCalls = append(m.ops.branchesContainingCalls, commit)
	if m.ops.branchesContainingError != nil {
		return nil, m.ops.branchesContainingError
	}
	return []string{"main", "feature"}, nil
}

func TestBrancher_Branch_Current(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		currentBranch: "feature/test",
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
	}
	brancher.Branch([]string{"current"})

	if !mockClient.getCurrentBranchCalled {
		t.Error("GetCurrentBranch should be called")
	}

	output := buf.String()
	if output != "feature/test\n" {
		t.Errorf("unexpected output: got %q, want %q", output, "feature/test\n")
	}
}

func TestBrancher_Branch_Current_Error(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		err: errors.New("failed to get current branch"),
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
	}
	brancher.Branch([]string{"current"})

	if !mockClient.getCurrentBranchCalled {
		t.Error("GetCurrentBranch should be called")
	}

	output := buf.String()
	if output != "Error: failed to get current branch\n" {
		t.Errorf("unexpected output: got %q, want %q", output, "Error: failed to get current branch\n")
	}
}

func TestBrancher_Branch_Checkout(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		prompter:     prompt.New(strings.NewReader("1\n"), &buf),
	}
	brancher.Branch([]string{"checkout"})

	output := buf.String()
	expected := "Local branches:\n[1] main\n[2] feature/test\nEnter the number to checkout: "
	if !strings.Contains(output, expected) {
		t.Errorf("unexpected output: got %q, want %q", output, expected)
	}
}

func TestBrancher_Branch_CheckoutRemote(t *testing.T) {
	var buf bytes.Buffer
	brancher := &Brancher{
		gitClient:    &mockBranchGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
		prompter:     prompt.New(strings.NewReader("1\n"), &buf),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"checkout", "remote"})

	output := buf.String()
	if !strings.Contains(output, "Remote branches:") {
		t.Error("Expected remote branches list")
	}
}

func TestBrancher_Branch_Delete(t *testing.T) {
	var buf bytes.Buffer
	brancher := &Brancher{
		gitClient:    &mockBranchGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
		prompter:     prompt.New(strings.NewReader("1\n"), &buf),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"delete"})

	output := buf.String()
	if !strings.Contains(output, "Select local branches to delete") {
		t.Error("Expected delete prompt")
	}
}

func TestBrancher_Branch_DeleteMerged(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		currentBranch:  "main",
		mergedBranches: []string{"feature/old", "feature/completed"},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
		prompter:     prompt.New(strings.NewReader("1\n"), &buf),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"delete", "merged"})

	output := buf.String()
	if !strings.Contains(output, "Select merged local branches to delete") {
		t.Error("Expected delete merged prompt")
	}
}

func TestBrancher_Branch_Rename_WithArgs(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"rename", "feature/old", "feature/new"})

	if mockClient.ops == nil || len(mockClient.ops.renameBranchCalls) != 1 {
		t.Fatalf("expected one rename call, got %v", mockClient.ops)
	}
	call := mockClient.ops.renameBranchCalls[0]
	if call.old != "feature/old" || call.new != "feature/new" {
		t.Errorf("unexpected rename args: got %+v", call)
	}
}

func TestBrancher_Branch_Rename_WithArgsInvalid(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"rename", "feature/old", ""})

	if mockClient.ops != nil && len(mockClient.ops.renameBranchCalls) != 0 {
		t.Fatalf("expected no rename calls, got %v", mockClient.ops.renameBranchCalls)
	}
	if !strings.Contains(buf.String(), "Error: new branch name cannot be empty.") {
		t.Errorf("expected empty-name error, got %q", buf.String())
	}
}

func TestBrancher_Branch_Rename_WithArgsGitError(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		ops: &mockBranchOperations{renameBranchError: errors.New("rename failed")},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"rename", "feature/old", "feature/new"})

	if mockClient.ops == nil || len(mockClient.ops.renameBranchCalls) != 1 {
		t.Fatalf("expected one rename call, got %v", mockClient.ops)
	}
	if !strings.Contains(buf.String(), "Error: rename failed") {
		t.Errorf("expected git error output, got %q", buf.String())
	}
}

func TestBrancher_Branch_Move_WithArgs(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"move", "feature/work", "abc123"})

	if mockClient.ops == nil || len(mockClient.ops.moveBranchCalls) != 1 {
		t.Fatalf("expected one move call, got %v", mockClient.ops)
	}
	call := mockClient.ops.moveBranchCalls[0]
	if call.branch != "feature/work" || call.commit != "abc123" {
		t.Errorf("unexpected move args: got %+v", call)
	}
}

func TestBrancher_Branch_Move_WithArgsInvalidCommit(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		ops: &mockBranchOperations{revParseVerifyResult: false},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"move", "feature/work", "abc123"})

	if len(mockClient.ops.moveBranchCalls) != 0 {
		t.Fatalf("expected no move call, got %v", mockClient.ops.moveBranchCalls)
	}
	if !strings.Contains(buf.String(), "Invalid commit or ref.") {
		t.Errorf("expected invalid commit message, got %q", buf.String())
	}
}

func TestBrancher_Branch_Move_WithArgsGitError(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		ops: &mockBranchOperations{revParseVerifyResult: true, moveBranchError: errors.New("move failed")},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"move", "feature/work", "abc123"})

	if mockClient.ops == nil || len(mockClient.ops.moveBranchCalls) != 1 {
		t.Fatalf("expected one move call, got %v", mockClient.ops)
	}
	if !strings.Contains(buf.String(), "Error: move failed") {
		t.Errorf("expected move error output, got %q", buf.String())
	}
}

func TestBrancher_Branch_SetUpstream_WithArgs(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"set", "upstream", "feature/work", "origin/feature"})

	if mockClient.ops == nil || len(mockClient.ops.setUpstreamBranchCalls) != 1 {
		t.Fatalf("expected one upstream call, got %v", mockClient.ops)
	}
	call := mockClient.ops.setUpstreamBranchCalls[0]
	if call.branch != "feature/work" || call.upstream != "origin/feature" {
		t.Errorf("unexpected upstream args: got %+v", call)
	}
}

func TestBrancher_Branch_SetUpstream_WithNumericArg(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listRemoteBranches: func() ([]string, error) {
			return []string{"origin/main", "origin/feature"}, nil
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"set", "upstream", "feature/work", "2"})

	if mockClient.ops == nil || len(mockClient.ops.setUpstreamBranchCalls) != 1 {
		t.Fatalf("expected one upstream call, got %v", mockClient.ops)
	}
	call := mockClient.ops.setUpstreamBranchCalls[0]
	if call.upstream != "origin/feature" {
		t.Errorf("expected numeric mapping, got %+v", call)
	}
}

func TestBrancher_Branch_SetUpstream_WithInvalidSelection(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listRemoteBranches: func() ([]string, error) {
			return []string{"origin/main"}, nil
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"set", "upstream", "feature/work", "5"})

	if mockClient.ops != nil && len(mockClient.ops.setUpstreamBranchCalls) != 0 {
		t.Fatalf("expected no upstream call, got %v", mockClient.ops.setUpstreamBranchCalls)
	}
	if !strings.Contains(buf.String(), "Error: invalid remote selection") {
		t.Errorf("expected invalid selection error, got %q", buf.String())
	}
}

func TestBrancher_Branch_SetUpstream_WithArgsGitError(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		ops: &mockBranchOperations{setUpstreamError: errors.New("set upstream failed")},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"set", "upstream", "feature/work", "origin/feature"})

	if mockClient.ops == nil || len(mockClient.ops.setUpstreamBranchCalls) != 1 {
		t.Fatalf("expected one upstream call, got %v", mockClient.ops)
	}
	if !strings.Contains(buf.String(), "Error: set upstream failed") {
		t.Errorf("expected upstream error, got %q", buf.String())
	}
}

func TestBrancher_Branch_SetUpstream_WithInvalidArity(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"set", "upstream", "feature/work"})

	if mockClient.ops != nil && len(mockClient.ops.setUpstreamBranchCalls) != 0 {
		t.Fatalf("expected no upstream call, got %v", mockClient.ops.setUpstreamBranchCalls)
	}
	if !strings.Contains(buf.String(), "expects <branch> <upstream>") {
		t.Errorf("expected arity error, got %q", buf.String())
	}
}

func TestBrancher_Branch_Info_WithArgs(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"info", "feature/test"})

	if !strings.Contains(buf.String(), "Name: feature/test") {
		t.Errorf("expected branch info output, got %q", buf.String())
	}
}

func TestBrancher_Branch_Info_WithEmptyArg(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"info", ""})

	if !strings.Contains(buf.String(), "branch name cannot be empty") {
		t.Errorf("expected empty-name error, got %q", buf.String())
	}
}

func TestBrancher_Branch_Info_WithExtraArgs(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"info", "one", "two"})

	if !strings.Contains(buf.String(), "at most one branch") {
		t.Errorf("expected extra-arg error, got %q", buf.String())
	}
}

func TestBrancher_Branch_Sort_WithArgs(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"sort", "date"})

	if mockClient.ops == nil || len(mockClient.ops.sortBranchesCalls) != 1 {
		t.Fatalf("expected one sort call, got %v", mockClient.ops)
	}
	if mockClient.ops.sortBranchesCalls[0] != "date" {
		t.Errorf("expected sort by date, got %v", mockClient.ops.sortBranchesCalls)
	}
	if !strings.Contains(buf.String(), "feature/test") {
		t.Errorf("expected sorted output, got %q", buf.String())
	}
}

func TestBrancher_Branch_Sort_InvalidOption(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"sort", "semantic"})

	if mockClient.ops != nil && len(mockClient.ops.sortBranchesCalls) != 0 {
		t.Fatalf("expected no sort call, got %v", mockClient.ops.sortBranchesCalls)
	}
	if !strings.Contains(buf.String(), "invalid sort option") {
		t.Errorf("expected invalid option error, got %q", buf.String())
	}
}

func TestBrancher_Branch_Sort_EmptyOption(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"sort", ""})

	if mockClient.ops != nil && len(mockClient.ops.sortBranchesCalls) != 0 {
		t.Fatalf("expected no sort call, got %v", mockClient.ops.sortBranchesCalls)
	}
	if !strings.Contains(buf.String(), "sort option cannot be empty") {
		t.Errorf("expected empty option error, got %q", buf.String())
	}
}

func TestBrancher_Branch_Sort_WithError(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		ops: &mockBranchOperations{sortBranchesError: errors.New("sort failed")},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"sort", "name"})

	if mockClient.ops == nil || len(mockClient.ops.sortBranchesCalls) != 1 {
		t.Fatalf("expected one sort call, got %v", mockClient.ops)
	}
	if !strings.Contains(buf.String(), "Error: sort failed") {
		t.Errorf("expected sort error output, got %q", buf.String())
	}
}

func TestBrancher_Branch_Contains_WithArgs(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"contains", "abc123"})

	if mockClient.ops == nil || len(mockClient.ops.branchesContainingCalls) != 1 {
		t.Fatalf("expected one contains call, got %v", mockClient.ops)
	}
	if mockClient.ops.branchesContainingCalls[0] != "abc123" {
		t.Errorf("expected commit 'abc123', got %v", mockClient.ops.branchesContainingCalls)
	}
	if !strings.Contains(buf.String(), "main") {
		t.Errorf("expected branches output, got %q", buf.String())
	}
}

func TestBrancher_Branch_Contains_InvalidCommit(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		ops: &mockBranchOperations{revParseVerifyResult: false},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"contains", "abc123"})

	if len(mockClient.ops.branchesContainingCalls) != 0 {
		t.Fatalf("expected no contains call, got %v", mockClient.ops.branchesContainingCalls)
	}
	if !strings.Contains(buf.String(), "Invalid commit or ref.") {
		t.Errorf("expected invalid commit error, got %q", buf.String())
	}
}

func TestBrancher_Branch_Contains_EmptyArg(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"contains", ""})

	if !strings.Contains(buf.String(), "commit or ref cannot be empty") {
		t.Errorf("expected empty commit error, got %q", buf.String())
	}
}

func TestBrancher_Branch_Contains_TooManyArgs(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"contains", "one", "two"})

	if !strings.Contains(buf.String(), "at most one commit") {
		t.Errorf("expected arity error, got %q", buf.String())
	}
}

func TestBrancher_Branch_Contains_WithError(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		ops: &mockBranchOperations{revParseVerifyResult: true, branchesContainingError: errors.New("lookup failed")},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	brancher.Branch([]string{"contains", "abc123"})

	if mockClient.ops == nil || len(mockClient.ops.branchesContainingCalls) != 1 {
		t.Fatalf("expected one contains call, got %v", mockClient.ops)
	}
	if !strings.Contains(buf.String(), "Error: lookup failed") {
		t.Errorf("expected lookup error, got %q", buf.String())
	}
}

func TestBrancher_Branch_Help(t *testing.T) {
	var buf bytes.Buffer
	brancher := &Brancher{
		gitClient:    &mockBranchGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	// Show help when no arguments provided
	brancher.Branch([]string{})

	output := buf.String()
	if !strings.Contains(output, "Usage") {
		t.Error("Expected help message to contain 'Usage'")
	}
}

func TestBrancher_Branch_UnknownCommand(t *testing.T) {
	var buf bytes.Buffer
	brancher := &Brancher{
		gitClient:    &mockBranchGitClient{},
		outputWriter: &buf,
		helper:       NewHelper(),
	}
	brancher.helper.outputWriter = &buf

	// Show help for unknown commands
	brancher.Branch([]string{"unknown"})

	output := buf.String()
	if !strings.Contains(output, "Usage") {
		t.Error("Expected help message to contain 'Usage'")
	}
}

func TestBrancher_branchCheckout_Error(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listLocalBranches: func() ([]string, error) {
			return nil, errors.New("failed to list branches")
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
	}
	brancher.branchCheckout()

	output := buf.String()
	if !strings.Contains(output, "Error: failed to list branches") {
		t.Errorf("Expected error message, got: %s", output)
	}
}

func TestBrancher_branchDelete_Success(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listLocalBranches: func() ([]string, error) {
			return []string{"feature/test", "bugfix/issue"}, nil
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		prompter:     prompt.New(strings.NewReader("1 2\n"), &buf),
	}

	brancher.branchDeleteArgs(nil)

	output := buf.String()
	if !strings.Contains(output, "feature/test") {
		t.Error("Expected branch name in output")
	}
	if !strings.Contains(output, "Selected branches deleted.") {
		t.Error("Expected success message")
	}
}

func TestBrancher_branchDelete_All(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listLocalBranches: func() ([]string, error) {
			return []string{"feature/test", "bugfix/issue"}, nil
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		prompter:     prompt.New(strings.NewReader("all\n"), &buf),
	}

	brancher.branchDeleteArgs(nil)

	output := buf.String()
	if !strings.Contains(output, "All branches deleted.") {
		t.Error("Expected all branches deleted message")
	}
}

func TestBrancher_branchDelete_Cancel(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listLocalBranches: func() ([]string, error) {
			return []string{"feature/test"}, nil
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		prompter:     prompt.New(strings.NewReader("\n"), &buf),
	}

	brancher.branchDeleteArgs(nil)

	output := buf.String()
	if !strings.Contains(output, "Canceled.") {
		t.Error("Expected canceled message")
	}
}

func TestBrancher_branchDelete_Error(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listLocalBranches: func() ([]string, error) {
			return nil, errors.New("failed to list branches")
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
	}

	brancher.branchDeleteArgs(nil)

	output := buf.String()
	if !strings.Contains(output, "Error: failed to list branches") {
		t.Error("Expected error message")
	}
}

func TestBrancher_branchDelete_NoBranches(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listLocalBranches: func() ([]string, error) {
			return []string{}, nil
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
	}

	brancher.branchDeleteArgs(nil)

	output := buf.String()
	if !strings.Contains(output, "No local branches found.") {
		t.Error("Expected no branches message")
	}
}

func TestBrancher_branchDeleteMerged_Success(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		currentBranch:  "main",
		mergedBranches: []string{"feature/old", "feature/completed"},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		prompter:     prompt.New(strings.NewReader("1 2\n"), &buf),
	}

	brancher.branchDeleteMerged()

	output := buf.String()
	if !strings.Contains(output, "Selected merged branches deleted.") {
		t.Error("Expected success message")
	}
}

func TestBrancher_branchDeleteMerged_All(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		currentBranch:  "main",
		mergedBranches: []string{"feature/old", "feature/completed", "hotfix/bug"},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		prompter:     prompt.New(strings.NewReader("all\n"), &buf),
	}

	brancher.branchDeleteMerged()

	output := buf.String()
	if !strings.Contains(output, "All merged branches deleted.") {
		t.Error("Expected all merged branches deleted message")
	}
}

func TestBrancher_branchDeleteMerged_Cancel(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		currentBranch:  "main",
		mergedBranches: []string{"feature/old", "feature/completed"},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		prompter:     prompt.New(strings.NewReader("\n"), &buf),
	}

	brancher.branchDeleteMerged()

	output := buf.String()
	if !strings.Contains(output, "Canceled.") {
		t.Error("Expected canceled message")
	}
}

func TestBrancher_branchDeleteMerged_CurrentBranchError(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		err: errors.New("failed to get current branch"),
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
	}

	brancher.branchDeleteMerged()

	output := buf.String()
	if !strings.Contains(output, "Error: failed to get current branch") {
		t.Error("Expected current branch error message")
	}
}

func TestBrancher_branchDeleteMerged_NoMergedBranches(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		currentBranch: "main",
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
	}

	brancher.branchDeleteMerged()

	output := buf.String()
	if !strings.Contains(output, "No merged local branches.") {
		t.Error("Expected no merged branches message")
	}
}

func TestBrancher_branchCheckoutRemote_Success(t *testing.T) {
	var buf bytes.Buffer
	brancher := &Brancher{
		gitClient:    &mockBranchGitClient{},
		outputWriter: &buf,
		prompter:     prompt.New(strings.NewReader("2\n"), &buf),
	}

	brancher.branchCheckoutRemote()

	output := buf.String()
	if !strings.Contains(output, "Remote branches:") {
		t.Error("Expected remote branches list")
	}
	if !strings.Contains(output, "origin/feature/test") {
		t.Error("Expected remote branch name in output")
	}
}

func TestBrancher_branchCheckoutRemote_InvalidNumber(t *testing.T) {
	var buf bytes.Buffer
	brancher := &Brancher{
		gitClient:    &mockBranchGitClient{},
		outputWriter: &buf,
		prompter:     prompt.New(strings.NewReader("invalid\n"), &buf),
	}

	brancher.branchCheckoutRemote()

	output := buf.String()
	if !strings.Contains(output, "Invalid number.") {
		t.Error("Expected invalid number message")
	}
}

func TestBrancher_branchCheckoutRemote_InvalidBranchName(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listRemoteBranches: func() ([]string, error) {
			return []string{"invalidbranch"}, nil
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		prompter:     prompt.New(strings.NewReader("1\n"), &buf),
	}

	brancher.branchCheckoutRemote()

	output := buf.String()
	if !strings.Contains(output, "Invalid remote branch name.") {
		t.Error("Expected invalid branch name message")
	}
}

func TestBrancher_branchCheckoutRemote_EmptyLocalFromRemote(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listRemoteBranches: func() ([]string, error) {
			return []string{"remote/"}, nil
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		prompter:     prompt.New(strings.NewReader("1\n"), &buf),
	}

	brancher.branchCheckoutRemote()

	output := buf.String()
	if !strings.Contains(output, "Invalid remote branch name.") {
		t.Error("Expected invalid remote branch name message for empty local name")
	}
}

func TestBrancher_Branch_Create(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput string
		cmdOutput      string
		cmdError       bool
	}{
		{
			name:           "Success: Create new branch",
			input:          "feature/test\n",
			expectedOutput: "Enter new branch name: ",
			cmdOutput:      "",
			cmdError:       false,
		},
		{
			name:           "Error: Empty branch name",
			input:          "\n",
			expectedOutput: "Enter new branch name: Canceled.\n",
			cmdOutput:      "",
			cmdError:       false,
		},
		{
			name:           "Error: Branch creation failed",
			input:          "feature/test\n",
			expectedOutput: "Enter new branch name: Error: failed to create and checkout branch: exit status 1\n",
			cmdOutput:      "",
			cmdError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			brancher := &Brancher{
				gitClient:    &mockBranchGitClient{checkoutNewBranchError: tt.cmdError},
				outputWriter: &buf,
				prompter:     prompt.New(strings.NewReader(tt.input), &buf),
			}

			brancher.Branch([]string{"create"})

			output := buf.String()
			if output != tt.expectedOutput {
				t.Errorf("unexpected output:\ngot:  %q\nwant: %q", output, tt.expectedOutput)
			}
		})
	}
}

func TestBrancher_branchCreate_ExistingBranch(t *testing.T) {
	var buf bytes.Buffer
	brancher := &Brancher{
		gitClient:    &mockBranchGitClient{},
		outputWriter: &buf,
		prompter:     prompt.New(strings.NewReader("main\n"), &buf),
	}

	brancher.Branch([]string{"create"})

	output := buf.String()
	expectedOutput := "Enter new branch name: Error: failed to create and checkout branch: fatal: a branch named 'main' already exists\n"
	if output != expectedOutput {
		t.Errorf("unexpected output:\ngot:  %q\nwant: %q", output, expectedOutput)
	}
}

func TestBrancher_branchCreate_WithArgument(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
	}

	brancher.Branch([]string{"create", "feature/new"})

	if len(mockClient.createdBranches) != 1 || mockClient.createdBranches[0] != "feature/new" {
		t.Fatalf("expected branch creation with name %q, got %#v", "feature/new", mockClient.createdBranches)
	}
	if output := buf.String(); output != "" {
		t.Fatalf("expected no output, got %q", output)
	}
}

func TestBrancher_branchCreate_WithEmptyArgument(t *testing.T) {
	var buf bytes.Buffer
	brancher := &Brancher{
		gitClient:    &mockBranchGitClient{},
		outputWriter: &buf,
	}

	brancher.Branch([]string{"create", "   "})

	output := buf.String()
	expected := "Error: invalid branch name: branch name cannot be empty\n"
	if output != expected {
		t.Errorf("unexpected output for empty argument:\ngot:  %q\nwant: %q", output, expected)
	}
}

// Consolidated boundary tests for branch checkout input to avoid repetitive single-case tests.
// Covers various numeric/non-numeric forms and validates expected outputs for each.
func TestBrancher_Branch_BoundaryInputValues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Single character input",
			input:    "1\n",
			expected: "main",
		},
		{
			name:     "Maximum integer input",
			input:    "999999\n",
			expected: "Invalid number",
		},
		{
			name:     "Negative number input",
			input:    "-1\n",
			expected: "Invalid number",
		},
		{
			name:     "Zero input",
			input:    "0\n",
			expected: "Invalid number",
		},
		{
			name:     "Leading zeros",
			input:    "001\n",
			expected: "main",
		},
		{
			name:     "Floating point number",
			input:    "1.5\n",
			expected: "Invalid number",
		},
		{
			name:     "Scientific notation",
			input:    "1e2\n",
			expected: "Invalid number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			mockClient := &mockBranchGitClient{}
			brancher := &Brancher{
				gitClient:    mockClient,
				outputWriter: &buf,
				prompter:     prompt.New(strings.NewReader(tt.input), &buf),
			}

			brancher.branchCheckout()

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected %q in output, got: %s", tt.expected, output)
			}
		})
	}
}

func TestBrancher_Branch_BoundaryBranchNames(t *testing.T) {
	tests := []struct {
		name         string
		branchName   string
		expectedPass bool
		description  string
	}{
		{
			name:         "Single character branch name",
			branchName:   "a",
			expectedPass: true,
			description:  "Minimum length branch name",
		},
		{
			name:         "Maximum length branch name",
			branchName:   strings.Repeat("a", 255),
			expectedPass: true,
			description:  "Long branch name allowed by git",
		},
		{
			name:         "Branch name with dots",
			branchName:   "feature.test",
			expectedPass: true,
			description:  "Branch name with dots",
		},
		{
			name:         "Branch name with hyphens",
			branchName:   "feature-test",
			expectedPass: true,
			description:  "Branch name with hyphens",
		},
		{
			name:         "Branch name with underscores",
			branchName:   "feature_test",
			expectedPass: true,
			description:  "Branch name with underscores",
		},
		{
			name:         "Branch name with slashes",
			branchName:   "feature/test",
			expectedPass: true,
			description:  "Branch name with slashes",
		},
		{
			name:         "Branch name starting with dot",
			branchName:   ".test",
			expectedPass: false,
			description:  "Leading dot is rejected by git",
		},
		{
			name:         "Branch name ending with dot",
			branchName:   "test.",
			expectedPass: false,
			description:  "Invalid: ends with dot",
		},
		{
			name:         "Branch name with consecutive dots",
			branchName:   "test..branch",
			expectedPass: false,
			description:  "Invalid: consecutive dots",
		},
		{
			name:         "Branch name with spaces",
			branchName:   "test branch",
			expectedPass: false,
			description:  "Invalid: contains spaces",
		},
		{
			name:         "Branch name with special characters",
			branchName:   "test@#$%",
			expectedPass: true,
			description:  "Allowed by git check-ref-format (no forbidden sequences)",
		},
		{
			name:         "Branch name with control characters",
			branchName:   "test\x00branch",
			expectedPass: false,
			description:  "Invalid: control characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			brancher := &Brancher{
				gitClient:    &mockBranchGitClient{},
				outputWriter: &buf,
				prompter:     prompt.New(strings.NewReader(tt.branchName+"\n"), &buf),
			}

			brancher.branchCreate(nil)

			output := buf.String()
			if tt.expectedPass {
				if strings.Contains(output, "Error:") {
					t.Errorf("Expected success for %s, but got error: %s", tt.description, output)
				}
			} else {
				if !strings.Contains(output, "Error:") && !strings.Contains(output, "fatal:") {
					t.Errorf("Expected error for %s, but got success: %s", tt.description, output)
				}
			}
		})
	}
}

func TestBrancher_Branch_BoundaryUserInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty input",
			input:    "\n",
			expected: "Canceled",
		},
		{
			name:     "Whitespace only",
			input:    "   \n",
			expected: "Canceled",
		},
		{
			name:     "Tab characters",
			input:    "\t\t\n",
			expected: "Canceled",
		},
		{
			name:     "Multiple newlines",
			input:    "\n\n\n",
			expected: "Canceled",
		},
		{
			name:     "Very long input",
			input:    strings.Repeat("a", 1000) + "\n",
			expected: "", // Git allows long branch names
		},
		{
			name:     "Input with trailing spaces",
			input:    "branch   \n",
			expected: "", // Should be trimmed and work
		},
		{
			name:     "Input with leading spaces",
			input:    "   branch\n",
			expected: "", // Should be trimmed and work
		},
		{
			name:     "Unicode characters",
			input:    "brαnch\n",
			expected: "", // Git allows UTF-8 branch names
		},
		{
			name:     "Mixed case input",
			input:    "BrAnCh\n",
			expected: "", // Should work
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			brancher := &Brancher{
				gitClient:    &mockBranchGitClient{},
				outputWriter: &buf,
				prompter:     prompt.New(strings.NewReader(tt.input), &buf),
			}

			brancher.branchCreate(nil)

			output := buf.String()
			if tt.expected != "" {
				if !strings.Contains(output, tt.expected) {
					t.Errorf("Expected %q in output, got: %s", tt.expected, output)
				}
			}
		})
	}
}

// Consolidated boundary tests for listing/selection behavior (e.g., empty lists, long names),
// reducing the need for separate, repetitive tests for each edge case.
func TestBrancher_Branch_BoundaryListOperations(t *testing.T) {
	tests := []struct {
		name          string
		localBranches []string
		input         string
		expected      string
	}{
		{
			name:          "Empty branch list",
			localBranches: []string{},
			input:         "1\n",
			expected:      "No local branches found",
		},
		{
			name:          "Single branch",
			localBranches: []string{"main"},
			input:         "1\n",
			expected:      "main",
		},
		{
			name:          "Branch with very long name",
			localBranches: []string{"main", strings.Repeat("feature-", 30)},
			input:         "2\n",
			expected:      strings.Repeat("feature-", 30),
		},
		{
			name:          "Branches with special naming patterns",
			localBranches: []string{"main", "feature/ABC-123", "bugfix/fix-issue-456", "release/v1.0.0"},
			input:         "3\n",
			expected:      "bugfix/fix-issue-456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			mockClient := &mockBranchGitClient{
				listLocalBranches: func() ([]string, error) {
					return tt.localBranches, nil
				},
			}
			brancher := &Brancher{
				gitClient:    mockClient,
				outputWriter: &buf,
				prompter:     prompt.New(strings.NewReader(tt.input), &buf),
			}

			brancher.branchCheckout()

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected %q in output, got: %s", tt.expected, output)
			}
		})
	}
}

func TestBrancher_Branch_BoundaryDeleteOperations(t *testing.T) {
	tests := []struct {
		name          string
		localBranches []string
		input         string
		expected      string
	}{
		{
			name:          "Delete single branch with boundary input",
			localBranches: []string{"main", "feature", "bugfix"},
			input:         "1\n",
			expected:      "Selected branches deleted",
		},
		{
			name:          "Delete multiple branches",
			localBranches: []string{"main", "feature", "bugfix", "hotfix"},
			input:         "1 2 3\n",
			expected:      "Selected branches deleted",
		},
		{
			name:          "Delete all branches",
			localBranches: []string{"feature", "bugfix", "hotfix"},
			input:         "all\n",
			expected:      "All branches deleted",
		},
		{
			name:          "Delete with out-of-range numbers",
			localBranches: []string{"main", "feature"},
			input:         "1 5 10\n",
			expected:      "Invalid number:",
		},
		{
			name:          "Delete with mixed valid/invalid input",
			localBranches: []string{"main", "feature", "bugfix"},
			input:         "1 invalid 2\n",
			expected:      "Invalid number:",
		},
		{
			name:          "Delete with very large numbers",
			localBranches: []string{"main", "feature"},
			input:         "999999 1000000\n",
			expected:      "Invalid number:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			mockClient := &mockBranchGitClient{
				listLocalBranches: func() ([]string, error) {
					return tt.localBranches, nil
				},
			}
			brancher := &Brancher{
				gitClient:    mockClient,
				outputWriter: &buf,
				prompter:     prompt.New(strings.NewReader(tt.input), &buf),
			}

			brancher.branchDeleteArgs(nil)

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected %q in output, got: %s", tt.expected, output)
			}
		})
	}
}

func TestBrancher_Branch_BoundaryRemoteOperations(t *testing.T) {
	tests := []struct {
		name           string
		remoteBranches []string
		input          string
		expected       string
	}{
		{
			name:           "Empty remote branches",
			remoteBranches: []string{},
			input:          "1\n",
			expected:       "Selected branches deleted.",
		},
		{
			name:           "Single remote branch",
			remoteBranches: []string{"origin/main"},
			input:          "1\n",
			expected:       "main",
		},
		{
			name:           "Multiple remote origins",
			remoteBranches: []string{"origin/main", "upstream/main", "fork/main"},
			input:          "2\n",
			expected:       "main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			mockClient := &mockBranchGitClient{
				listRemoteBranches: func() ([]string, error) {
					return tt.remoteBranches, nil
				},
			}
			brancher := &Brancher{
				gitClient:    mockClient,
				outputWriter: &buf,
				prompter:     prompt.New(strings.NewReader(tt.input), &buf),
			}

			brancher.branchDeleteArgs(nil)

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected %q in output, got: %s", tt.expected, output)
			}
		})
	}
}

func TestBrancher_branchDeleteArgs_WithArgs(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{currentBranch: "main"}
	brancher := &Brancher{gitClient: mockClient, outputWriter: &buf}

	brancher.branchDeleteArgs([]string{"feature/test"})

	if len(mockClient.deletedBranches) != 1 || mockClient.deletedBranches[0] != "feature/test" {
		t.Fatalf("expected to delete feature/test, got %v", mockClient.deletedBranches)
	}
	if got := buf.String(); got != "" {
		t.Errorf("expected no output on successful deletion, got %q", got)
	}
}

func TestBrancher_branchDeleteArgs_MultipleAndSkipCurrent(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{currentBranch: "main"}
	brancher := &Brancher{gitClient: mockClient, outputWriter: &buf}

	brancher.branchDeleteArgs([]string{"main", "feature/test", "foo"})

	// Should skip current branch 'main'
	if got := buf.String(); !strings.Contains(got, "Skipping current branch: main") {
		t.Errorf("expected skip message, got %q", got)
	}
	if len(mockClient.deletedBranches) != 2 {
		t.Fatalf("expected 2 deletions, got %d (%v)", len(mockClient.deletedBranches), mockClient.deletedBranches)
	}
	if mockClient.deletedBranches[0] != "feature/test" || mockClient.deletedBranches[1] != "foo" {
		t.Errorf("unexpected deleted branches: %v", mockClient.deletedBranches)
	}
}

func TestBrancher_branchDeleteArgs_NoArgsFallsBackInteractive(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{gitClient: mockClient, outputWriter: &buf, prompter: prompt.New(strings.NewReader("1\n"), &buf)}

	brancher.branchDeleteArgs(nil)

	if got := buf.String(); !strings.Contains(got, "Selected branches deleted.") {
		t.Errorf("expected interactive deletion confirmation, got %q", got)
	}
	if len(mockClient.deletedBranches) == 0 {
		t.Error("expected at least one branch deletion call")
	}
}

func TestBrancher_branchInfo_WithArgsLegacy(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{currentBranch: "main"}
	brancher := &Brancher{gitClient: mockClient, outputWriter: &buf}

	brancher.branchInfo([]string{"feature/test"})

	out := buf.String()
	if !strings.Contains(out, "Name: feature/test") {
		t.Errorf("expected branch info to be printed; got %q", out)
	}
}

func TestBrancher_branchInfo_MultipleArgsShowsError(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{gitClient: mockClient, outputWriter: &buf}

	brancher.branchInfo([]string{"main", "feature/test"})

	out := buf.String()
	if !strings.Contains(out, "at most one branch") {
		t.Errorf("expected error message for multiple args; got %q", out)
	}
}

func TestBrancher_branchInfo_NoArgsFallsBackInteractive(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{}
	brancher := &Brancher{gitClient: mockClient, outputWriter: &buf, prompter: prompt.New(strings.NewReader("2\n"), &buf)}

	brancher.branchInfo(nil)

	out := buf.String()
	if !strings.Contains(out, "Name: feature/test") {
		t.Errorf("expected info printed for selected branch; got %q", out)
	}
}

func TestBrancher_selectUpstreamBranch_ErrorHandling(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listRemoteBranches: func() ([]string, error) {
			return nil, errors.New("network error: connection refused")
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		prompter:     prompt.New(strings.NewReader("origin/main\n"), &buf),
	}

	result := brancher.selectUpstreamBranch()

	if result != "" {
		t.Errorf("expected empty string when ListRemoteBranches fails, got %q", result)
	}
	output := buf.String()
	if !strings.Contains(output, "Error listing remote branches") {
		t.Errorf("expected error message in output, got: %s", output)
	}
	if !strings.Contains(output, "network error") {
		t.Errorf("expected original error message in output, got: %s", output)
	}
}

func TestBrancher_selectUpstreamBranch_EmptyBranchesFiltered(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listRemoteBranches: func() ([]string, error) {
			return []string{"origin/main", "", "  ", "origin/feature"}, nil
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		prompter:     prompt.New(strings.NewReader("2\n"), &buf),
	}

	result := brancher.selectUpstreamBranch()

	// After filtering empty strings, index 2 should map to "origin/feature"
	if result != "origin/feature" {
		t.Errorf("expected 'origin/feature', got %q", result)
	}
	output := buf.String()
	// Should only show 2 valid branches
	if strings.Count(output, "[1]") != 1 || strings.Count(output, "[2]") != 1 {
		t.Errorf("expected exactly 2 numbered branches in output, got: %s", output)
	}
	if strings.Count(output, "[3]") > 0 {
		t.Errorf("expected no third branch (empty ones should be filtered), got: %s", output)
	}
}

func TestBrancher_selectUpstreamBranch_Success(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listRemoteBranches: func() ([]string, error) {
			return []string{"origin/main", "origin/develop"}, nil
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		prompter:     prompt.New(strings.NewReader("1\n"), &buf),
	}

	result := brancher.selectUpstreamBranch()

	if result != "origin/main" {
		t.Errorf("expected 'origin/main', got %q", result)
	}
	output := buf.String()
	if !strings.Contains(output, "Remote branches:") {
		t.Errorf("expected 'Remote branches:' in output, got: %s", output)
	}
}

func TestBrancher_selectUpstreamBranch_IndexOutOfRange(t *testing.T) {
	var buf bytes.Buffer
	mockClient := &mockBranchGitClient{
		listRemoteBranches: func() ([]string, error) {
			return []string{"origin/main"}, nil
		},
	}
	brancher := &Brancher{
		gitClient:    mockClient,
		outputWriter: &buf,
		prompter:     prompt.New(strings.NewReader("999\n"), &buf),
	}

	result := brancher.selectUpstreamBranch()

	// When index is out of range, it should return the input as-is (not index into array)
	if result != "999" {
		t.Errorf("expected '999' (out of range index returned as-is), got %q", result)
	}
}
