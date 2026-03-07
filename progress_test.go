package worky

import (
	"os"
	"testing"
)

func TestStateIsUnlocked(t *testing.T) {
	s := &State{Unlocked: []string{"01", "02"}}
	if !s.IsUnlocked("01") {
		t.Error("01 should be unlocked")
	}
	if !s.IsUnlocked("02") {
		t.Error("02 should be unlocked")
	}
	if s.IsUnlocked("03") {
		t.Error("03 should not be unlocked")
	}
}

func TestStateIsCompleted(t *testing.T) {
	s := &State{Completed: []string{"01"}}
	if !s.IsCompleted("01") {
		t.Error("01 should be completed")
	}
	if s.IsCompleted("02") {
		t.Error("02 should not be completed")
	}
}

func TestStateComplete(t *testing.T) {
	s := &State{Unlocked: []string{"01"}}
	s.Complete("01", "02")

	if !s.IsCompleted("01") {
		t.Error("01 should be completed")
	}
	if !s.IsUnlocked("02") {
		t.Error("02 should be unlocked after completing 01")
	}
}

func TestStateComplete_Idempotent(t *testing.T) {
	s := &State{Unlocked: []string{"01"}}
	s.Complete("01", "02")
	s.Complete("01", "02")

	count := 0
	for _, id := range s.Completed {
		if id == "01" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("01 should appear once in completed, got %d", count)
	}
}

func TestStateComplete_NoNext(t *testing.T) {
	s := &State{Unlocked: []string{"01"}}
	s.Complete("01", "")
	if !s.IsCompleted("01") {
		t.Error("01 should be completed")
	}
}

func TestStateUnlock(t *testing.T) {
	s := &State{}
	s.Unlock("01")
	if !s.IsUnlocked("01") {
		t.Error("01 should be unlocked")
	}
	// Idempotent
	s.Unlock("01")
	count := 0
	for _, id := range s.Unlocked {
		if id == "01" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("01 should appear once in unlocked, got %d", count)
	}
}

func TestLoadSaveProgress(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	w := New(Config{
		Name:    "test",
		HomeDir: ".worky-test",
		Chapters: []Chapter{
			{ID: "01", Name: "One", Slug: "01-one"},
			{ID: "02", Name: "Two", Slug: "02-two"},
		},
	})

	// loadProgress with no file returns initial state with first chapter unlocked.
	state, err := w.loadProgress()
	if err != nil {
		t.Fatalf("loadProgress: %v", err)
	}
	if !state.IsUnlocked("01") {
		t.Error("first chapter should be unlocked by default")
	}

	// Mutate and save.
	state.Complete("01", "02")
	if err := w.saveProgress(state); err != nil {
		t.Fatalf("saveProgress: %v", err)
	}

	// Reload and verify.
	loaded, err := w.loadProgress()
	if err != nil {
		t.Fatalf("loadProgress after save: %v", err)
	}
	if !loaded.IsCompleted("01") {
		t.Error("01 should be completed after reload")
	}
	if !loaded.IsUnlocked("02") {
		t.Error("02 should be unlocked after reload")
	}
}

func TestResetProgress(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	w := New(Config{
		Name:    "test",
		HomeDir: ".worky-test",
		Chapters: []Chapter{
			{ID: "01", Name: "One", Slug: "01-one"},
		},
	})

	state, _ := w.loadProgress()
	state.Complete("01", "")
	w.saveProgress(state) //nolint:errcheck

	if err := w.resetProgress(); err != nil {
		t.Fatalf("resetProgress: %v", err)
	}

	after, err := w.loadProgress()
	if err != nil {
		t.Fatalf("loadProgress after reset: %v", err)
	}
	if after.IsCompleted("01") {
		t.Error("01 should not be completed after reset")
	}
}

func TestSaveLoadCheckResults(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	w := New(Config{
		Name:    "test",
		HomeDir: ".worky-test",
		Chapters: []Chapter{
			{ID: "01", Name: "One", Slug: "01-one"},
		},
	})

	results := []CheckResult{
		{Description: "file exists", Passed: true},
		{Description: "port open", Passed: false, Error: "refused"},
	}
	if err := w.saveCheckResults("01", results); err != nil {
		t.Fatalf("saveCheckResults: %v", err)
	}

	store, err := w.loadCheckResults()
	if err != nil {
		t.Fatalf("loadCheckResults: %v", err)
	}
	got, ok := store["01"]
	if !ok {
		t.Fatal("no results for chapter 01")
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	if got[1].Error != "refused" {
		t.Errorf("unexpected error: %q", got[1].Error)
	}
}

func TestLoadCheckResults_Empty(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	w := New(Config{
		Name:    "test",
		HomeDir: ".worky-test",
		Chapters: []Chapter{
			{ID: "01", Name: "One", Slug: "01-one"},
		},
	})

	store, err := w.loadCheckResults()
	if err != nil {
		t.Fatalf("loadCheckResults on empty: %v", err)
	}
	if len(store) != 0 {
		t.Errorf("expected empty store, got %v", store)
	}
}

func TestLoadProgress_FirstChapterAlwaysUnlocked(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	w := New(Config{
		Name:    "test",
		HomeDir: ".worky-test",
		Chapters: []Chapter{
			{ID: "01", Name: "One", Slug: "01-one"},
			{ID: "02", Name: "Two", Slug: "02-two"},
		},
	})

	// Write a state that somehow doesn't have 01 unlocked.
	wDir, _ := w.workshopDir()
	os.WriteFile(wDir+"/progress.json", []byte(`{"completed":[],"unlocked":["02"]}`), 0o644)

	state, err := w.loadProgress()
	if err != nil {
		t.Fatalf("loadProgress: %v", err)
	}
	if !state.IsUnlocked("01") {
		t.Error("first chapter should always be unlocked")
	}
}
