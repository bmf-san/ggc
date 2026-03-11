new_test = """
func TestBrancher_CollectDeletableBranches_FiltersCurrentBranch(t *testing.T) {
\tvar buf bytes.Buffer
\tmockClient := &mockBranchGitClient{
\t\tcurrentBranch: "main",
\t\tlistLocalBranches: func() ([]string, error) {
\t\t\treturn []string{"main", "feature/x", "feature/y"}, nil
\t\t},
\t}
\tbrancher := &Brancher{
\t\tgitClient:    mockClient,
\t\toutputWriter: &buf,
\t}
\tbranches, ok := brancher.collectDeletableBranches()
\tif !ok {
\t\tt.Error("expected ok=true")
\t}
\tfor _, br := range branches {
\t\tif br == "main" {
\t\t\tt.Error("current branch main should be filtered out")
\t\t}
\t}
\tif len(branches) != 2 {
\t\tt.Errorf("expected 2 branches after filtering, got %d: %v", len(branches), branches)
\t}
}
"""
with open("cmd/branch_test.go", "a") as f:
    f.write(new_test)
print("done")
