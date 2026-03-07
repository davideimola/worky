package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeRunner struct {
	calls [][]string
}

func (f *fakeRunner) Run(name string, args ...string) error {
	f.calls = append(f.calls, append([]string{name}, args...))
	return nil
}

func TestBuildCmd_ErrorIfNotWorkyProject(t *testing.T) {
	dir := t.TempDir() // empty dir, no hugo.toml

	err := buildWorkshop(dir, &fakeRunner{})
	if err == nil {
		t.Fatal("expected error when hugo.toml is missing, got nil")
	}
	if !strings.Contains(err.Error(), "hugo.toml not found") {
		t.Errorf("expected 'hugo.toml not found' error, got: %v", err)
	}
}

func TestBuildCmd_RunsHugoAndGoBuild(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "hugo.toml"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	runner := &fakeRunner{}
	if err := buildWorkshop(dir, runner); err != nil {
		t.Fatalf("buildWorkshop: %v", err)
	}

	if len(runner.calls) < 2 {
		t.Fatalf("expected at least 2 command calls, got %d", len(runner.calls))
	}
	if runner.calls[0][0] != "hugo" {
		t.Errorf("expected first call to be 'hugo', got %q", runner.calls[0][0])
	}
	if runner.calls[1][0] != "go" {
		t.Errorf("expected second call to be 'go', got %q", runner.calls[1][0])
	}
}
