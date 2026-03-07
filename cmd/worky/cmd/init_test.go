package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitProject_CreatesAllFiles(t *testing.T) {
	dir := t.TempDir()
	if err := initProject("My Workshop", "github.com/example/my-workshop", dir); err != nil {
		t.Fatalf("initProject: %v", err)
	}

	expected := []string{
		"my-workshop/main.go",
		"my-workshop/go.mod",
		"my-workshop/hugo.toml",
		"my-workshop/Makefile",
		"my-workshop/.gitignore",
		"my-workshop/README.md",
		"my-workshop/docs/_index.md",
		"my-workshop/docs/00-setup/_index.md",
		"my-workshop/.github/workflows/release.yml",
	}

	for _, rel := range expected {
		path := filepath.Join(dir, rel)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s to exist: %v", rel, err)
		}
	}
}

func TestInitProject_SlugFromName(t *testing.T) {
	dir := t.TempDir()
	if err := initProject("My Workshop", "github.com/example/my-workshop", dir); err != nil {
		t.Fatalf("initProject: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "my-workshop")); err != nil {
		t.Errorf("expected directory my-workshop to be created: %v", err)
	}
}

func TestInitProject_FileContents(t *testing.T) {
	dir := t.TempDir()
	if err := initProject("My Workshop", "github.com/example/my-workshop", dir); err != nil {
		t.Fatalf("initProject: %v", err)
	}

	mainGo, err := os.ReadFile(filepath.Join(dir, "my-workshop", "main.go"))
	if err != nil {
		t.Fatalf("read main.go: %v", err)
	}
	if !strings.Contains(string(mainGo), `Name:    "My Workshop"`) {
		t.Errorf("main.go should contain workshop name, got:\n%s", mainGo)
	}
	if !strings.Contains(string(mainGo), `HomeDir: ".my-workshop"`) {
		t.Errorf("main.go should contain HomeDir with slug, got:\n%s", mainGo)
	}

	hugoToml, err := os.ReadFile(filepath.Join(dir, "my-workshop", "hugo.toml"))
	if err != nil {
		t.Fatalf("read hugo.toml: %v", err)
	}
	if !strings.Contains(string(hugoToml), `title      = "My Workshop"`) {
		t.Errorf("hugo.toml should contain workshop name, got:\n%s", hugoToml)
	}
	if !strings.Contains(string(hugoToml), "github.com/davideimola/worky") {
		t.Errorf("hugo.toml should import worky module, got:\n%s", hugoToml)
	}
}

func TestInitProject_ErrorIfDirExists(t *testing.T) {
	dir := t.TempDir()
	// Pre-create the target directory
	if err := os.MkdirAll(filepath.Join(dir, "my-workshop"), 0o755); err != nil {
		t.Fatal(err)
	}

	err := initProject("My Workshop", "github.com/example/my-workshop", dir)
	if err == nil {
		t.Fatal("expected error when directory already exists, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' error, got: %v", err)
	}
}
