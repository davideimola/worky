package worky

import (
	"context"
	"fmt"
	"time"
)

// Check is a single validation step.
type Check struct {
	Description string
	Run         func(context.Context) error
	// Timeout is the maximum time allowed for a single attempt. Zero means no timeout.
	Timeout time.Duration
	// Retries is the number of additional attempts after the first failure. Zero means run once.
	Retries int
	// RetryDelay is the pause between retry attempts. Zero means no delay.
	RetryDelay time.Duration
}

// CheckResult holds the outcome of a single check.
type CheckResult struct {
	Description string `json:"description"`
	Passed      bool   `json:"passed"`
	Error       string `json:"error,omitempty"`
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
// If Timeout is set, a context with deadline is passed to Run; if the context
// expires the attempt reports "timed out after <duration>".
func runCheck(c Check) error {
	attempt := func() error {
		ctx := context.Background()
		var cancel context.CancelFunc
		if c.Timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, c.Timeout)
			defer cancel()
		}
		err := c.Run(ctx)
		if err != nil && ctx.Err() != nil {
			return fmt.Errorf("timed out after %s", c.Timeout)
		}
		return err
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
