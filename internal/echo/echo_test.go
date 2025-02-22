package echo

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	oldStdout := os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	Run([]string{"Hello", "World"})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}
	r.Close()

	expected := "Hello World\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected %q but got %q", expected, got)
	}
}
