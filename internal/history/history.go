package history

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// historyFile returns the path to the history file using XDG or home dir.
func historyFile() (string, error) {
	cfg, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(cfg, "ggc")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "history"), nil
}

// AppendCommand appends a timestamped command line to the history file.
func AppendCommand(args []string) error {
	if len(args) == 0 {
		return nil
	}
	hf, err := historyFile()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(hf, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	line := time.Now().UTC().Format(time.RFC3339) + "\t" + strings.Join(args, " ") + "\n"
	_, err = f.WriteString(line)
	return err
}

// ReadAll returns all history lines (most recent last).
func ReadAll() ([]string, error) {
	hf, err := historyFile()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(hf)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	defer f.Close()
	var out []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		out = append(out, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// ReadLast returns up to n most recent history entries (most recent last).
func ReadLast(n int) ([]string, error) {
	all, err := ReadAll()
	if err != nil {
		return nil, err
	}
	if n <= 0 || n >= len(all) {
		return all, nil
	}
	return all[len(all)-n:], nil
}

// Search returns lines that contain pattern (case-insensitive).
func Search(pattern string) ([]string, error) {
	all, err := ReadAll()
	if err != nil {
		return nil, err
	}
	var out []string
	low := strings.ToLower(pattern)
	for _, l := range all {
		if strings.Contains(strings.ToLower(l), low) {
			out = append(out, l)
		}
	}
	return out, nil
}
