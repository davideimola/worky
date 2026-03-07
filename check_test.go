package worky

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRunCheck_Success(t *testing.T) {
	c := Check{
		Description: "always passes",
		Run:         func(_ context.Context) error { return nil },
	}
	if err := runCheck(c); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRunCheck_Failure(t *testing.T) {
	want := errors.New("boom")
	c := Check{
		Description: "always fails",
		Run:         func(_ context.Context) error { return want },
	}
	if err := runCheck(c); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRunCheck_Retries(t *testing.T) {
	attempts := 0
	c := Check{
		Description: "fails twice then passes",
		Retries:     2,
		Run: func(_ context.Context) error {
			attempts++
			if attempts < 3 {
				return errors.New("not yet")
			}
			return nil
		},
	}
	if err := runCheck(c); err != nil {
		t.Fatalf("expected nil after retries, got %v", err)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestRunCheck_RetriesExhausted(t *testing.T) {
	c := Check{
		Description: "always fails",
		Retries:     2,
		Run:         func(_ context.Context) error { return errors.New("fail") },
	}
	if err := runCheck(c); err == nil {
		t.Fatal("expected error after retries exhausted")
	}
}

func TestRunCheck_Timeout(t *testing.T) {
	c := Check{
		Description: "hangs",
		Timeout:     50 * time.Millisecond,
		Run: func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(10 * time.Second):
				return nil
			}
		},
	}
	err := runCheck(c)
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if err.Error() != "timed out after 50ms" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunCheck_TimeoutSucceeds(t *testing.T) {
	c := Check{
		Description: "fast enough",
		Timeout:     500 * time.Millisecond,
		Run:         func(_ context.Context) error { return nil },
	}
	if err := runCheck(c); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRunChecks_AllPass(t *testing.T) {
	w := New(Config{
		Name:     "test",
		Chapters: []Chapter{{ID: "01", Name: "One", Slug: "01-one"}},
	})
	checks := []Check{
		{Description: "a", Run: func(_ context.Context) error { return nil }},
		{Description: "b", Run: func(_ context.Context) error { return nil }},
	}
	results, passed := w.runChecks(checks)
	if !passed {
		t.Fatal("expected all checks to pass")
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Passed {
			t.Errorf("check %q should have passed", r.Description)
		}
	}
}

func TestRunChecks_SomeFail(t *testing.T) {
	w := New(Config{
		Name:     "test",
		Chapters: []Chapter{{ID: "01", Name: "One", Slug: "01-one"}},
	})
	checks := []Check{
		{Description: "pass", Run: func(_ context.Context) error { return nil }},
		{Description: "fail", Run: func(_ context.Context) error { return errors.New("nope") }},
	}
	results, passed := w.runChecks(checks)
	if passed {
		t.Fatal("expected not all checks to pass")
	}
	if results[0].Passed != true {
		t.Error("first check should pass")
	}
	if results[1].Passed != false {
		t.Error("second check should fail")
	}
	if results[1].Error != "nope" {
		t.Errorf("unexpected error message: %q", results[1].Error)
	}
}
