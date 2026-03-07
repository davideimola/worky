---
title: Checks
weight: 4
---

## Check struct

```go
type Check struct {
    Description string
    Run         func() error
    Timeout     time.Duration // Zero means no timeout
    Retries     int           // Additional attempts after first failure (zero = run once)
    RetryDelay  time.Duration // Pause between retries (zero = no delay)
}
```

- `Description` is shown in CLI output next to the pass/fail icon.
- `Run` is any function returning `error`. `nil` means pass; non-nil means fail with that error message.
- `Timeout`, `Retries`, `RetryDelay` are optional — useful for checks that need to wait for services to start.

---

## Built-in checks

Import `github.com/davideimola/worky/checks` for pre-built check functions.

### File system

#### `FileExists(path string) func() error`

Passes if the file at `path` exists.

```go
checks.FileExists("done.txt")
checks.FileExists("/etc/hosts")
```

---

#### `DirExists(path string) func() error`

Passes if the directory at `path` exists.

```go
checks.DirExists("output/")
checks.DirExists("src/components")
```

---

#### `FileContains(path, text string) func() error`

Passes if the file at `path` contains the substring `text`.

```go
checks.FileContains("hello.txt", "Hello, World!")
checks.FileContains("config.yaml", "replicas: 3")
```

---

#### `FileMatchesRegex(path, pattern string) func() error`

Passes if the file at `path` matches the regular expression `pattern`.

```go
checks.FileMatchesRegex("output.json", `"status":\s*"ok"`)
```

---

### Environment variables

#### `EnvVarSet(name string) func() error`

Passes if the environment variable `name` is set and non-empty.

```go
checks.EnvVarSet("KUBECONFIG")
checks.EnvVarSet("AWS_REGION")
```

---

#### `EnvVarEquals(name, value string) func() error`

Passes if the environment variable `name` equals `value`.

```go
checks.EnvVarEquals("STAGE", "production")
```

---

### Commands

#### `CommandSucceeds(name string, args ...string) func() error`

Passes if the command exits with code 0.

```go
checks.CommandSucceeds("docker", "info")
checks.CommandSucceeds("kubectl", "cluster-info")
checks.CommandSucceeds("node", "--version")
```

---

#### `CommandOutputContains(text, name string, args ...string) func() error`

Passes if the command exits 0 and its combined output contains `text`.

```go
checks.CommandOutputContains("Running", "kubectl", "get", "pods", "-n", "default")
checks.CommandOutputContains("v1.", "kubectl", "version", "--client")
```

---

### Network

#### `PortOpen(host string, port int) func() error`

Passes if a TCP connection to `host:port` succeeds within 3 seconds.

```go
checks.PortOpen("localhost", 5432)  // PostgreSQL
checks.PortOpen("localhost", 6379)  // Redis
```

---

#### `HTTPStatus(url string, expectedStatus int) func() error`

Passes if a GET request to `url` returns `expectedStatus`.

```go
checks.HTTPStatus("http://localhost:8080/health", 200)
checks.HTTPStatus("http://localhost:3000", 200)
```

---

#### `HTTPBodyContains(url, text string) func() error`

Passes if the body of a GET request to `url` contains `text`.

```go
checks.HTTPBodyContains("http://localhost:8080/api/status", `"ok"`)
```

---

## Custom checks

Any `func() error` works:

```go
worky.Check{
    Description: "Namespace 'workshop' exists in Kubernetes",
    Run: func() error {
        out, err := exec.Command("kubectl", "get", "ns", "workshop").CombinedOutput()
        if err != nil {
            return fmt.Errorf("namespace not found: %s", out)
        }
        return nil
    },
}
```

---

## Retry example

For checks that need to wait for a service:

```go
worky.Check{
    Description: "API server is healthy",
    Run:         checks.HTTPStatus("http://localhost:8080/health", 200),
    Timeout:     5 * time.Second,
    Retries:     6,
    RetryDelay:  5 * time.Second,
}
```

This will attempt the check up to 7 times (1 initial + 6 retries), waiting 5 seconds between each, with each attempt timing out after 5 seconds.
