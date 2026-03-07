---
title: Checks
weight: 4
---

A check is a named validation step that a participant must pass to unlock the next chapter.

## Check struct

```go
type Check struct {
    Description string
    Run         func(context.Context) error
    Timeout     time.Duration // 0 = no limit per attempt
    Retries     int           // additional attempts after first failure (0 = run once)
    RetryDelay  time.Duration // pause between retries (0 = no delay)
}
```

- `Description` is shown in CLI output next to the pass/fail icon.
- `Run` is any `func(context.Context) error`. Return `nil` to pass; return a non-nil error to fail with that message. When `Timeout` is set, the context passed to `Run` carries the per-attempt deadline — pass it to any blocking calls so they can be cancelled.
- `Timeout`, `Retries`, `RetryDelay` are optional — useful for checks that need to wait for services to start.

{{< cards >}}
  {{< card link="built-in" title="Built-in checks" subtitle="FileExists, EnvVarSet, CommandSucceeds, HTTPStatus and more" icon="check-circle" >}}
  {{< card link="patterns" title="Patterns" subtitle="Timeout, retry, custom functions, shared state between chapters" icon="code" >}}
{{< /cards >}}
