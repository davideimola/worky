package worky

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// State holds the workshop progression state.
type State struct {
	Completed []string `json:"completed"`
	Unlocked  []string `json:"unlocked"`
}

func (w *Workshop) workshopDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, w.cfg.HomeDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func (w *Workshop) stateFile() (string, error) {
	dir, err := w.workshopDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "progress.json"), nil
}

func (w *Workshop) loadProgress() (*State, error) {
	path, err := w.stateFile()
	if err != nil {
		return nil, err
	}

	firstID := ""
	if len(w.cfg.Chapters) > 0 {
		firstID = w.cfg.Chapters[0].ID
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		s := &State{
			Completed: []string{},
			Unlocked:  []string{firstID},
		}
		return s, nil
	}
	if err != nil {
		return nil, err
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	// Ensure first chapter is always unlocked.
	if firstID != "" && !s.IsUnlocked(firstID) {
		s.Unlocked = append([]string{firstID}, s.Unlocked...)
	}
	return &s, nil
}

func (w *Workshop) saveProgress(s *State) error {
	path, err := w.stateFile()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (w *Workshop) resetProgress() error {
	path, err := w.stateFile()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	// Also clear saved check results.
	cpath, err := w.checkResultsFile()
	if err != nil {
		return err
	}
	if err := os.Remove(cpath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// checkResultsFile returns the path to the check results file.
func (w *Workshop) checkResultsFile() (string, error) {
	dir, err := w.workshopDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "checks.json"), nil
}

// CheckResultStore maps chapter IDs to their last check results.
type CheckResultStore map[string][]CheckResult

// saveCheckResults persists the results for a chapter to disk.
func (w *Workshop) saveCheckResults(chapterID string, results []CheckResult) error {
	path, err := w.checkResultsFile()
	if err != nil {
		return err
	}
	store := CheckResultStore{}
	if data, readErr := os.ReadFile(path); readErr == nil {
		_ = json.Unmarshal(data, &store)
	}
	store[chapterID] = results
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// loadCheckResults reads all saved check results from disk.
func (w *Workshop) loadCheckResults() (CheckResultStore, error) {
	path, err := w.checkResultsFile()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return CheckResultStore{}, nil
	}
	if err != nil {
		return nil, err
	}
	var store CheckResultStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return store, nil
}

// IsUnlocked reports whether a chapter is unlocked.
func (s *State) IsUnlocked(id string) bool {
	for _, u := range s.Unlocked {
		if u == id {
			return true
		}
	}
	return false
}

// IsCompleted reports whether a chapter is completed.
func (s *State) IsCompleted(id string) bool {
	for _, c := range s.Completed {
		if c == id {
			return true
		}
	}
	return false
}

// Complete marks a chapter as completed and unlocks the next one.
func (s *State) Complete(id, nextID string) {
	if !s.IsCompleted(id) {
		s.Completed = append(s.Completed, id)
	}
	if nextID != "" && !s.IsUnlocked(nextID) {
		s.Unlocked = append(s.Unlocked, nextID)
	}
}

// Unlock makes a chapter accessible without completing it.
func (s *State) Unlock(id string) {
	if !s.IsUnlocked(id) {
		s.Unlocked = append(s.Unlocked, id)
	}
}
