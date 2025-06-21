package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestComplete(t *testing.T) {
	// オリジナルのNewCompleterを保存
	originalNewCompleter := NewCompleter
	defer func() {
		// テスト終了後に元に戻す
		NewCompleter = originalNewCompleter
	}()

	// NewCompleterをモックに差し替える
	NewCompleter = func() *Completer {
		return &Completer{
			listLocalBranches: func() ([]string, error) {
				return []string{"feature/test-branch", "main"}, nil
			},
		}
	}

	// 標準出力をキャプチャ
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// テスト対象の関数を実行
	// "branch"サブコマンドに、"sub"という引数を与えて呼び出す
	Complete([]string{"branch", "sub"})

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}

	// 出力に期待するブランチ名が含まれているか確認
	expected := "feature/test-branch"
	if !bytes.Contains(buf.Bytes(), []byte(expected)) {
		t.Errorf("expected output to contain %q, but got %q", expected, buf.String())
	}

	expected = "main"
	if !bytes.Contains(buf.Bytes(), []byte(expected)) {
		t.Errorf("expected output to contain %q, but got %q", expected, buf.String())
	}
}
