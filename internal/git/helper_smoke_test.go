package git

import (
	"bytes"
	"testing"
)

func TestEchoCommandHelper(t *testing.T) {
	cmd := echoCommand("hello world")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	if buf.String() != "hello world" {
		t.Fatalf("got %q", buf.String())
	}
	if failingCommand().Run() == nil {
		t.Fatal("failingCommand should fail")
	}
	if err := successCommand().Run(); err != nil {
		t.Fatalf("successCommand: %v", err)
	}
}
