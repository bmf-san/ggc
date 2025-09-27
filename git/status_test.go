package git

import (
	"os/exec"
	"slices"
	"testing"
)

func TestClient_Status(t *testing.T) {
	var gotArgs []string
	expectedOutput := "On branch main\nnothing to commit, working tree clean"

	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo", "-n", expectedOutput)
		},
	}

	result, err := client.Status()
	if err != nil {
		t.Errorf("Status() error = %v", err)
	}

	wantArgs := []string{"git", "status"}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("Status() gotArgs = %v, want %v", gotArgs, wantArgs)
	}

	if result != expectedOutput {
		t.Errorf("Status() result = %v, want %v", result, expectedOutput)
	}
}

func TestClient_StatusShort(t *testing.T) {
	var gotArgs []string
	expectedOutput := " M file.go\n?? new_file.go"

	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo", "-n", expectedOutput)
		},
	}

	result, err := client.StatusShort()
	if err != nil {
		t.Errorf("StatusShort() error = %v", err)
	}

	wantArgs := []string{"git", "status", "--short"}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("StatusShort() gotArgs = %v, want %v", gotArgs, wantArgs)
	}

	if result != expectedOutput {
		t.Errorf("StatusShort() result = %v, want %v", result, expectedOutput)
	}
}

func TestClient_StatusWithColor(t *testing.T) {
	var gotArgs []string
	expectedOutput := "On branch main\nnothing to commit, working tree clean"

	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo", "-n", expectedOutput)
		},
	}

	result, err := client.StatusWithColor()
	if err != nil {
		t.Errorf("StatusWithColor() error = %v", err)
	}

	wantArgs := []string{"git", "-c", "color.status=always", "status"}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("StatusWithColor() gotArgs = %v, want %v", gotArgs, wantArgs)
	}

	if result != expectedOutput {
		t.Errorf("StatusWithColor() result = %v, want %v", result, expectedOutput)
	}
}

func TestClient_StatusShortWithColor(t *testing.T) {
	var gotArgs []string
	expectedOutput := " M file.go\n?? new_file.go"

	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			gotArgs = append([]string{name}, args...)
			return exec.Command("echo", "-n", expectedOutput)
		},
	}

	result, err := client.StatusShortWithColor()
	if err != nil {
		t.Errorf("StatusShortWithColor() error = %v", err)
	}

	wantArgs := []string{"git", "-c", "color.status=always", "status", "--short"}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("StatusShortWithColor() gotArgs = %v, want %v", gotArgs, wantArgs)
	}

	if result != expectedOutput {
		t.Errorf("StatusShortWithColor() result = %v, want %v", result, expectedOutput)
	}
}

// Error case tests for better coverage
func TestClient_Status_Error(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			return exec.Command("false") // Command that always fails
		},
	}

	_, err := client.Status()
	if err == nil {
		t.Error("Expected Status to return an error")
	}
}

func TestClient_StatusShort_Error(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			return exec.Command("false") // Command that always fails
		},
	}

	_, err := client.StatusShort()
	if err == nil {
		t.Error("Expected StatusShort to return an error")
	}
}

func TestClient_StatusWithColor_Error(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			return exec.Command("false") // Command that always fails
		},
	}

	_, err := client.StatusWithColor()
	if err == nil {
		t.Error("Expected StatusWithColor to return an error")
	}
}

func TestClient_StatusShortWithColor_Error(t *testing.T) {
	client := &Client{
		execCommand: func(name string, args ...string) *exec.Cmd {
			return exec.Command("false") // Command that always fails
		},
	}

	_, err := client.StatusShortWithColor()
	if err == nil {
		t.Error("Expected StatusShortWithColor to return an error")
	}
}
