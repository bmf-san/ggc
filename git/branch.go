package git

import (
	"os/exec"
	"strings"
)

func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	branch := strings.TrimSpace(string(out))
	return branch, nil
}

func ListLocalBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "--format", "%(refname:short)")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	return lines, nil
}

func ListRemoteBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "-r", "--format", "%(refname:short)")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	// origin/HEAD -> origin/main などのHEAD参照は除外
	filtered := []string{}
	for _, l := range lines {
		if strings.Contains(l, "->") {
			continue
		}
		filtered = append(filtered, strings.TrimSpace(l))
	}
	return filtered, nil
}
