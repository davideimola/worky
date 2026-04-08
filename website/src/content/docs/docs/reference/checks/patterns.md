---
title: Patterns
description: Retry, timeout, custom check functions, and shared state between chapters.
---

## Retry and timeout

### How it works

1. Run the first attempt.
2. If it fails and `Retries > 0`: wait `RetryDelay`, then retry — up to `Retries` additional times.
3. If `Timeout > 0`: each individual attempt is cancelled after that duration and returns `"timed out after {duration}"`.
4. If all attempts fail, the last error is returned.

### Example: wait for a service to start

```go
worky.Check{
    Description: "Server responds on :8080",
    Run:         checks.HTTPStatus("http://localhost:8080", 200),
    Timeout:     5 * time.Second,
    Retries:     5,
    RetryDelay:  2 * time.Second,
}
```

This makes up to 6 attempts (1 initial + 5 retries), waiting 2 seconds between each, with each attempt timing out after 5 seconds.

---

## Custom check function

Any `func(context.Context) error` is a valid `Run` value. Use this for validation logic that has no built-in helper. The context carries the per-attempt deadline set by `Timeout`, so pass it to any blocking calls:

```go
worky.Check{
    Description: "config.yaml is valid YAML",
    Run: func(_ context.Context) error {
        data, err := os.ReadFile("config.yaml")
        if err != nil {
            return err
        }
        var out any
        return yaml.Unmarshal(data, &out)
    },
}
```

---

## Shared state between chapters

Checks in later chapters can read values captured by checks in earlier chapters using closure variables:

```go
var deployedURL string

chapters := []worky.Chapter{
    {
        ID:   "01",
        Name: "Deploy",
        Checks: []worky.Check{
            {
                Description: "DEPLOY_URL is set",
                Run: func(_ context.Context) error {
                    deployedURL = os.Getenv("DEPLOY_URL")
                    if deployedURL == "" {
                        return errors.New("DEPLOY_URL is not set")
                    }
                    return nil
                },
            },
        },
    },
    {
        ID:   "02",
        Name: "Verify deploy",
        Checks: []worky.Check{
            {
                Description: "Deploy responds with 200",
                Run: func(ctx context.Context) error {
                    return checks.HTTPStatus(deployedURL, 200)(ctx)
                },
            },
        },
    },
}
```

`deployedURL` is populated when chapter 01's check runs and is available to chapter 02's check because both closures close over the same variable.

---

## Further reading

- [Checks](/docs/reference/checks/) — built-in check functions reference
- [Troubleshooting](/docs/troubleshooting/) — common check error messages and how to fix them
