---
title: Getting Started
weight: 1
---

## Prerequisites

- Go 1.21+

## 1. Install the CLI

```sh
go install github.com/davideimola/worky/cmd/worky@latest
```

See [CLI](/docs/cli) for the full command reference.

## 2. Scaffold a new workshop

```sh
worky init "My Workshop"
cd my-workshop
```

Creates:

```
my-workshop/
├── main.go        # Workshop entrypoint with embed + one example chapter
├── site/
│   ├── index.html           # Home page
│   ├── workshop.css
│   ├── workshop-progress.js
│   └── 00-setup/
│       └── index.md         # Placeholder chapter content
├── Makefile       # make serve, make build, make release
└── go.mod
```

## 3. Install dependencies

`worky init` asks *"Install dependencies now?"* and runs `go mod tidy` automatically if you say yes. If you skipped it:

```sh
go mod tidy
```

## 4. Start developing

```sh
make serve
```

Your workshop is live at `http://localhost:8080`. Edit `site/00-setup/index.md` to write chapter content, or add checks to `main.go`. See [Customization](/docs/customization) for site layout and CSS options.

---

## Manual setup

<details>
<summary>Without the CLI</summary>

### 1. Create the Go project

```sh
go mod init github.com/acme/my-workshop
go get github.com/davideimola/worky
```

### 2. Write `main.go`

```go
package main

import (
    "embed"

    "github.com/davideimola/worky"
    "github.com/davideimola/worky/checks"
)

//go:embed all:site
var site embed.FS

func main() {
    worky.New(worky.Config{
        Name:    "My Workshop",
        HomeDir: ".my-workshop",
        Port:    8080,
        SiteFS:  site,
        Chapters: []worky.Chapter{
            {
                ID:   "00",
                Name: "Setup",
                Slug: "00-setup",
                Checks: []worky.Check{
                    {
                        Description: "Docker is running",
                        Run:         checks.CommandSucceeds("docker", "info"),
                    },
                },
            },
        },
    }).Run()
}
```

### 3. Write content

Create `site/00-setup/index.md`:

```markdown
# Setup

Welcome to the workshop! Let's make sure your environment is ready.

## Steps

1. Start Docker Desktop
2. Run `./my-workshop check` to verify and unlock the next chapter
```

### 4. Build and run

```sh
go run . serve --open
# or
go build -o bin/my-workshop .
```

</details>

---

## Try the demo

The `demo/` directory in the worky repository is a fully working workshop:

```sh
cd demo
go mod tidy
go run . serve --open
```

Complete each chapter to see the sidebar icons update live:

```sh
touch done.txt
export WORKSHOP_USER=yourname
go run . check   # Chapter 00

echo "Hello, Worky!" > hello.txt
go run . check   # Chapter 01

echo "# Workshop Complete" > complete.md
go run . check   # Chapter 02
```

---

## Distributing your workshop

`worky init` includes a GitHub Actions workflow (`.github/workflows/release.yml`) that cross-compiles for Linux, macOS, and Windows and publishes binaries to a GitHub Release automatically on every `v*` tag:

```sh
git tag v1.0.0
git push origin v1.0.0
```

To build locally instead:

```sh
make release   # produces bin/ with binaries for all platforms
```

Participants download the binary for their platform, run `./my-workshop serve --open`, and use `./my-workshop check` to progress through chapters. See [Runtime commands](/docs/runtime) for the full reference.

---

## What's next

- [Reference](/docs/reference) — `Config`, `Chapter`, and `Check` struct reference
- [Checks](/docs/reference/checks) — built-in validation functions
- [CLI](/docs/cli) — `worky init`, `worky new chapter`, `worky build`
- [Customization](/docs/customization) — site layout, CSS classes, Markdown rendering
