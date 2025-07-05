package git

import (
	"errors"
	"testing"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "コマンドありのエラー",
			err: &Error{
				Op:      "test operation",
				Command: "git test",
				Err:     errors.New("test error"),
			},
			want: "git: test operation failed: test error (command: git test)",
		},
		{
			name: "コマンドなしのエラー",
			err: &Error{
				Op:      "test operation",
				Command: "",
				Err:     errors.New("test error"),
			},
			want: "git: test operation failed: test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewError(t *testing.T) {
	op := "test operation"
	command := "git test"
	err := errors.New("test error")

	result := NewError(op, command, err)

	gitErr, ok := result.(*Error)
	if !ok {
		t.Fatalf("NewError should return *Error, got %T", result)
	}

	if gitErr.Op != op {
		t.Errorf("Op = %v, want %v", gitErr.Op, op)
	}
	if gitErr.Command != command {
		t.Errorf("Command = %v, want %v", gitErr.Command, command)
	}
	if gitErr.Err != err {
		t.Errorf("Err = %v, want %v", gitErr.Err, err)
	}
}
