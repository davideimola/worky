package cmd

import (
	"testing"
)

type fakeRunner struct {
	calls [][]string
}

func (f *fakeRunner) Run(name string, args ...string) error {
	f.calls = append(f.calls, append([]string{name}, args...))
	return nil
}

func TestBuildCmd_RunsGoBuild(t *testing.T) {
	dir := t.TempDir()

	runner := &fakeRunner{}
	if err := buildWorkshop(dir, runner); err != nil {
		t.Fatalf("buildWorkshop: %v", err)
	}
	if len(runner.calls) != 1 {
		t.Fatalf("expected exactly 1 command call (go build), got %d", len(runner.calls))
	}
	if runner.calls[0][0] != "go" {
		t.Errorf("expected call to be 'go', got %q", runner.calls[0][0])
	}
}
