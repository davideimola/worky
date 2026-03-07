package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunInit_CreatesAllFiles(t *testing.T) {
	dir := t.TempDir()
	opts := InitOptions{
		Name:   "My Workshop",
		Module: "github.com/example/my-workshop",
	}
	if err := runInit(opts, dir); err != nil {
		t.Fatalf("runInit: %v", err)
	}

	expected := []string{
		"my-workshop/main.go",
		"my-workshop/go.mod",
		"my-workshop/Makefile",
		"my-workshop/.gitignore",
		"my-workshop/README.md",
		"my-workshop/.github/workflows/release.yml",
		"my-workshop/site/index.html",
		"my-workshop/site/workshop.css",
		"my-workshop/site/workshop-progress.js",
		"my-workshop/site/00-setup/index.md",
	}
	for _, rel := range expected {
		path := filepath.Join(dir, rel)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s to exist: %v", rel, err)
		}
	}
}

func TestRunInit_NoHugoFiles(t *testing.T) {
	dir := t.TempDir()
	opts := InitOptions{
		Name:   "My Workshop",
		Module: "github.com/example/my-workshop",
	}
	if err := runInit(opts, dir); err != nil {
		t.Fatalf("runInit: %v", err)
	}

	absent := []string{
		"my-workshop/hugo.toml",
		"my-workshop/docs/_index.md",
		"my-workshop/docs/00-setup/_index.md",
	}
	for _, rel := range absent {
		if _, err := os.Stat(filepath.Join(dir, rel)); err == nil {
			t.Errorf("expected file %s to NOT exist", rel)
		}
	}
}

func TestRunInit_SlugFromName(t *testing.T) {
	dir := t.TempDir()
	opts := InitOptions{Name: "My Workshop", Module: "github.com/example/my-workshop"}
	if err := runInit(opts, dir); err != nil {
		t.Fatalf("runInit: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "my-workshop")); err != nil {
		t.Errorf("expected directory my-workshop to be created: %v", err)
	}
}

func TestRunInit_FileContents(t *testing.T) {
	dir := t.TempDir()
	opts := InitOptions{
		Name:   "My Workshop",
		Module: "github.com/example/my-workshop",
	}
	if err := runInit(opts, dir); err != nil {
		t.Fatalf("runInit: %v", err)
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
}

func TestRunInit_ErrorIfDirExists(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "my-workshop"), 0o755); err != nil {
		t.Fatal(err)
	}
	opts := InitOptions{Name: "My Workshop", Module: "github.com/example/my-workshop"}
	err := runInit(opts, dir)
	if err == nil {
		t.Fatal("expected error when directory already exists, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' error, got: %v", err)
	}
}

func TestRunInit_InstallDeps_CallsGoModTidy(t *testing.T) {
	dir := t.TempDir()

	var calls [][]string
	mockRun := func(runDir, name string, args ...string) error {
		calls = append(calls, append([]string{name}, args...))
		return nil
	}

	opts := InitOptions{
		Name:        "My Workshop",
		Module:      "github.com/example/my-workshop",
		InstallDeps: true,
		depsRunner:  mockRun,
	}
	if err := runInit(opts, dir); err != nil {
		t.Fatalf("runInit: %v", err)
	}

	if len(calls) != 1 {
		t.Fatalf("expected 1 command, got %d: %v", len(calls), calls)
	}
	if strings.Join(calls[0], " ") != "go mod tidy" {
		t.Errorf("expected 'go mod tidy', got %v", calls[0])
	}
}

func TestRunInit_InstallDeps_SkippedWhenFalse(t *testing.T) {
	dir := t.TempDir()

	var calls [][]string
	mockRun := func(runDir, name string, args ...string) error {
		calls = append(calls, append([]string{name}, args...))
		return nil
	}

	opts := InitOptions{
		Name:        "My Workshop",
		Module:      "github.com/example/my-workshop",
		InstallDeps: false,
		depsRunner:  mockRun,
	}
	if err := runInit(opts, dir); err != nil {
		t.Fatalf("runInit: %v", err)
	}
	if len(calls) != 0 {
		t.Errorf("expected no commands when InstallDeps=false, got: %v", calls)
	}
}
