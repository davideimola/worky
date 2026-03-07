package worky

import (
	"fmt"
	"time"
)

// Check is a single validation step.
type Check struct {
	Description string
	Run         func() error
	// Timeout is the maximum time allowed for a single attempt. Zero means no timeout.
	Timeout time.Duration
	// Retries is the number of additional attempts after the first failure. Zero means run once.
	Retries int
	// RetryDelay is the pause between retry attempts. Zero means no delay.
	RetryDelay time.Duration
}

// CheckResult holds the outcome of a single check.
type CheckResult struct {
	Description string
	Passed      bool
	Error       string
}

func (w *Workshop) runChecks(checks []Check) ([]CheckResult, bool) {
	results := make([]CheckResult, len(checks))
	allPassed := true
	for i, c := range checks {
		err := runCheck(c)
		r := CheckResult{Description: c.Description, Passed: err == nil}
		if err != nil {
			r.Error = err.Error()
			allPassed = false
		}
		results[i] = r
	}
	return results, allPassed
}

// runCheck executes a single check, honouring Timeout and Retries.
func runCheck(c Check) error {
	attempt := func() error {
		if c.Timeout <= 0 {
			return c.Run()
		}
		ch := make(chan error, 1)
		go func() { ch <- c.Run() }()
		select {
		case err := <-ch:
			return err
		case <-time.After(c.Timeout):
			return fmt.Errorf("timed out after %s", c.Timeout)
		}
	}

	var err error
	for i := 0; i <= c.Retries; i++ {
		if i > 0 && c.RetryDelay > 0 {
			time.Sleep(c.RetryDelay)
		}
		err = attempt()
		if err == nil {
			return nil
		}
	}
	return err
}
