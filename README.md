<p align="center">
  <img src="website/static/images/logo.svg" width="120" alt="worky" />
</p>

<h1 align="center">worky</h1>

<p align="center">A Go library for building interactive, self-contained workshop tools.</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/davideimola/worky"><img src="https://pkg.go.dev/badge/github.com/davideimola/worky.svg" alt="Go Reference"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-Apache%202.0-blue.svg" alt="License"></a>
</p>

**The problem:** running a workshop means participants fight setup — missing tools, broken environments, unclear progress. Worky fixes this with a single binary that locks chapters until checks pass and shows live status in the browser.

Participants download one file, run one command, and follow along in their own browser. No Docker, no cloud, no separate web server.

---

## Try the demo

The `demo/` directory is a fully working workshop you can run right now:

```sh
cd demo
go mod tidy
go run . serve --open
```

Then complete each chapter:

```sh
# Chapter 00 — Getting Started
touch done.txt
export WORKSHOP_USER=yourname
go run . check

# Chapter 01 — Hello Worky
echo "Hello, Worky!" > hello.txt
go run . check

# Chapter 02 — Finishing Up
echo "# Workshop Complete" > complete.md
go run . check
```

Watch the sidebar icons update live in your browser as each chapter unlocks. Run `go run . status` to see all chapters, `go run . reset` to start over.

---

## How it works

A workshop is a Go binary that embeds a `site/` directory of HTML and Markdown files:

```
your-workshop/
├── main.go          # imports worky, wires chapters + checks + embedded site
├── site/
│   ├── index.html           # home page
│   ├── workshop.css
│   ├── workshop-progress.js
│   └── 00-setup/
│       └── index.md         # chapter content (Markdown, rendered server-side)
└── go.mod
```

The worky server renders `.md` files on the fly and serves `.html` files as-is — no build step, no external tools.

If you prefer external slides (Google Slides, Marp, Reveal.js), skip writing chapter content — the built-in status page shows chapter progress with live 🔒/🔓/✅ icons regardless.

Participants:
1. Download the binary (or `git clone` + `go build`)
2. Run `./my-workshop serve --open`
3. Follow along in the browser, run `./my-workshop check` to unlock each chapter

---

## Create a workshop in 30 seconds

```sh
# Install the CLI
go install github.com/davideimola/worky/cmd/worky@latest

# Scaffold a new workshop
worky init "My Workshop"
cd my-workshop
go mod tidy
make serve
```

Add more chapters as your workshop grows:

```sh
worky new chapter 01 "Hello World"
# → creates site/01-hello-world/index.md
# → prints the Go snippet to add to main.go
```

Build for distribution:

```sh
worky build   # go build
```

---

## Manual setup

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
                    {Description: "Docker is running", Run: myDockerCheck},
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

Welcome to the workshop! Make sure your environment is ready.

## Steps

1. Start Docker Desktop
2. Run `./my-workshop check` to verify and unlock the next chapter
```

### 4. Build and distribute

```sh
go build -o bin/my-workshop .
GOOS=linux  GOARCH=amd64 go build -o bin/my-workshop-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -o bin/my-workshop-darwin-arm64 .
```

---

## CLI commands

| Command | Description |
|---------|-------------|
| `serve [--port N] [--open] [--detach]` | Start the local web server |
| `check [chapter-id]` | Run checks for a chapter and unlock the next one |
| `status` | Show all chapters with lock/unlock/complete icons |
| `reset` | Clear all progress |
| `stop` | Kill the background server |
| `logs [-f]` | Show (and optionally follow) server logs |

If `chapter-id` is omitted, `check` auto-detects the first incomplete unlocked chapter.

---

## Progress storage

Progress is stored in `~/<HomeDir>/progress.json`. The first chapter is always unlocked. Running `check` for a chapter marks it complete and unlocks the next one.
