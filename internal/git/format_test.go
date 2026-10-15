package git

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestGoFilesWithoutBOM(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no Go files found")
	}
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			content, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			if bytes.HasPrefix(content, []byte{0xef, 0xbb, 0xbf}) {
				t.Error("UTF-8 BOM causes goimports lint failures; save as UTF-8 without BOM")
			}
		})
	}
}
