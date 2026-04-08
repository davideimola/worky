---
title: Configuration
description: worky.Config struct — all fields explained.
---

The `worky.Config` struct is the only configuration surface. Pass it to `worky.New()`.

```go
type Config struct {
    Name     string    // Workshop display name
    HomeDir  string    // State directory under home (default: ".worky")
    Port     int       // Default HTTP port (default: 8080)
    SiteFS   fs.FS     // Embedded site (optional; if nil, a built-in UI is served)
    Chapters []Chapter
}
```

## Fields

### Name

```go
Name: "My Workshop"
```

Displayed as the CLI tool name in help output and printed by `status`. Required.

---

### HomeDir

```go
HomeDir: ".my-workshop"
```

The subdirectory under `$HOME` where worky stores state files:

| File | Purpose |
|------|---------|
| `progress.json` | Chapter unlock/complete state |
| `server.pid` | PID of the background server (when `--detach` is used) |
| `server.log` | Log output of the background server |

Defaults to `.worky` if empty.

---

### Port

```go
Port: 8080
```

The default HTTP port for `serve`. Participants can override with `--port N`. Defaults to `8080` if zero.

---

### SiteFS

`SiteFS` is **optional**.

**With embedded site:**

```go
//go:embed all:site
var site embed.FS

worky.New(worky.Config{
    SiteFS: site,
    // ...
})
```

The embedded FS must contain a `site/` subdirectory. See [Customization](/customization/) for the full site layout, file reference, and CSS classes.

**Without embedded site (built-in UI):**

```go
worky.New(worky.Config{
    // SiteFS omitted — worky serves a built-in progress page
    // ...
})
```

Use this when your slides live elsewhere (Google Slides, Marp, Reveal.js) and you only need the check/progress layer.

---

### Chapters

```go
Chapters: []worky.Chapter{
    {ID: "00", Name: "Setup",     Slug: "00-setup",     Checks: [...]},
    {ID: "01", Name: "Chapter 1", Slug: "01-chapter1",  Checks: [...]},
}
```

Ordered list of chapters. The first chapter is always unlocked. See [Chapters](/reference/chapters/) for the full `Chapter` struct reference.

---

## Full example

```go
package main

import (
    "embed"
    "time"

    "github.com/davideimola/worky"
    "github.com/davideimola/worky/checks"
)

//go:embed all:site
var site embed.FS

func main() {
    worky.New(worky.Config{
        Name:    "Kubernetes Workshop",
        HomeDir: ".k8s-workshop",
        Port:    9090,
        SiteFS:  site,
        Chapters: []worky.Chapter{
            {
                ID:   "00",
                Name: "Prerequisites",
                Slug: "00-prerequisites",
                Checks: []worky.Check{
                    {
                        Description: "kubectl is installed",
                        Run:         checks.CommandSucceeds("kubectl", "version", "--client"),
                    },
                    {
                        Description: "Cluster is reachable",
                        Run:         checks.CommandSucceeds("kubectl", "cluster-info"),
                        Timeout:     10 * time.Second,
                        Retries:     3,
                        RetryDelay:  2 * time.Second,
                    },
                },
            },
        },
    }).Run()
}
```
