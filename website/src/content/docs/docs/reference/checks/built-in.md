---
title: Built-in checks
description: Pre-built check functions from the checks sub-package.
---

Import `github.com/davideimola/worky/checks` for pre-built check functions.

## File system

### `FileExists(path string) func(context.Context) error`

Passes if the file at `path` exists.

```go
checks.FileExists("done.txt")
checks.FileExists("/etc/hosts")
```

---

### `DirExists(path string) func(context.Context) error`

Passes if the directory at `path` exists.

```go
checks.DirExists("output/")
checks.DirExists("src/components")
```

---

### `FileContains(path, text string) func(context.Context) error`

Passes if the file at `path` contains the substring `text`.

```go
checks.FileContains("hello.txt", "Hello, World!")
checks.FileContains("config.yaml", "replicas: 3")
```

---

### `FileMatchesRegex(path, pattern string) func(context.Context) error`

Passes if the file at `path` matches the regular expression `pattern`. If `pattern` is not a valid regular expression, the check returns an error instead of panicking.

```go
checks.FileMatchesRegex("output.json", `"status":\s*"ok"`)
```

---

## Environment variables

### `EnvVarSet(name string) func(context.Context) error`

Passes if the environment variable `name` is set and non-empty.

```go
checks.EnvVarSet("KUBECONFIG")
checks.EnvVarSet("AWS_REGION")
```

---

### `EnvVarEquals(name, value string) func(context.Context) error`

Passes if the environment variable `name` equals `value`.

```go
checks.EnvVarEquals("STAGE", "production")
```

---

## Commands

### `CommandSucceeds(name string, args ...string) func(context.Context) error`

Passes if the command exits with code 0. Respects context cancellation and timeout.

```go
checks.CommandSucceeds("docker", "info")
checks.CommandSucceeds("kubectl", "cluster-info")
checks.CommandSucceeds("node", "--version")
```

---

### `CommandOutputContains(text, name string, args ...string) func(context.Context) error`

Passes if the command exits 0 and its combined output contains `text`. Respects context cancellation and timeout.

```go
checks.CommandOutputContains("Running", "kubectl", "get", "pods", "-n", "default")
checks.CommandOutputContains("v1.", "kubectl", "version", "--client")
```

---

## Network

### `PortOpen(host string, port int) func(context.Context) error`

Passes if a TCP connection to `host:port` succeeds within 3 seconds. Respects context cancellation.

```go
checks.PortOpen("localhost", 5432)  // PostgreSQL
checks.PortOpen("localhost", 6379)  // Redis
```

---

### `HTTPStatus(url string, expectedStatus int) func(context.Context) error`

Passes if a GET request to `url` returns `expectedStatus`. Respects context cancellation and timeout.

```go
checks.HTTPStatus("http://localhost:8080/health", 200)
checks.HTTPStatus("http://localhost:3000", 200)
```

---

### `HTTPBodyContains(url, text string) func(context.Context) error`

Passes if the body of a GET request to `url` contains `text`. Respects context cancellation and timeout.

```go
checks.HTTPBodyContains("http://localhost:8080/api/status", `"ok"`)
```
