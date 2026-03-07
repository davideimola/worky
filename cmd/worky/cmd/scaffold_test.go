package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestScaffold_CreatesFile(t *testing.T) {
	tmplFS := fstest.MapFS{
		"hello.txt.tmpl": &fstest.MapFile{Data: []byte("hello")},
	}
	dest := filepath.Join(t.TempDir(), "hello.txt")

	if err := scaffold(tmplFS, "hello.txt.tmpl", dest, nil); err != nil {
		t.Fatalf("scaffold: %v", err)
	}

	if _, err := os.Stat(dest); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestScaffold_RendersTemplateVars(t *testing.T) {
	tmplFS := fstest.MapFS{
		"tmpl": &fstest.MapFile{Data: []byte("Hello <% .Name %>!")},
	}
	dest := filepath.Join(t.TempDir(), "out.txt")

	if err := scaffold(tmplFS, "tmpl", dest, map[string]string{"Name": "World"}); err != nil {
		t.Fatalf("scaffold: %v", err)
	}

	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "Hello World!" {
		t.Errorf("got %q, want %q", string(got), "Hello World!")
	}
}

func TestScaffold_CreatesParentDirs(t *testing.T) {
	tmplFS := fstest.MapFS{
		"tmpl": &fstest.MapFile{Data: []byte("data")},
	}
	dest := filepath.Join(t.TempDir(), "a", "b", "c", "out.txt")

	if err := scaffold(tmplFS, "tmpl", dest, nil); err != nil {
		t.Fatalf("scaffold: %v", err)
	}

	if _, err := os.Stat(dest); err != nil {
		t.Fatalf("expected file at nested path: %v", err)
	}
}

func TestScaffold_ErrorOnInvalidTemplate(t *testing.T) {
	tmplFS := fstest.MapFS{
		"bad.tmpl": &fstest.MapFile{Data: []byte("<% .Unclosed ")},
	}
	dest := filepath.Join(t.TempDir(), "out.txt")

	if err := scaffold(tmplFS, "bad.tmpl", dest, nil); err == nil {
		t.Fatal("expected error on invalid template, got nil")
	}
}
