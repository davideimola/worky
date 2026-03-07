# worky

A Go library for building interactive, self-contained workshop tools.

**The problem:** running a workshop means participants fight setup — missing tools, broken environments, unclear progress. Worky fixes this with a single binary that embeds the full documentation site, locks chapters until checks pass, and shows live status icons in the browser sidebar.

Participants download one file, run one command, and follow along in their own browser. No Docker, no cloud, no separate web server.

---

## Try the demo

The `demo/` directory is a fully working workshop you can run right now — no Hugo, no Kubernetes, no extra tools.

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

```
your-workshop/
├── main.go          # imports worky, wires chapters + checks + embedded site
├── docs/            # Hugo markdown content
├── hugo.toml        # Hugo config, imports worky as a Hugo module
└── site/            # Hugo output — embedded into the binary at build time
```

Participants:
1. Download the binary (or `git clone` + `make setup`)
2. Run `./my-workshop serve --open`
3. Follow along in the browser, run `./my-workshop check` to unlock each chapter

---

## Create a workshop in 30 seconds

```sh
# Install the CLI
go install github.com/davideimola/worky/cmd/worky@latest

# Scaffold a new workshop
worky init "My Workshop"

# Set up Hugo modules and Go dependencies
cd my-workshop
hugo mod init github.com/yourorg/my-workshop
hugo mod get github.com/davideimola/worky
hugo mod get github.com/geekdocs/geekdoc
go mod tidy

# Start developing
make serve
```

Add more chapters as your workshop grows:

```sh
worky new chapter 01 "Hello World"
# → creates docs/01-hello-world/_index.md
# → prints the Go snippet to add to main.go
```

Build for distribution:

```sh
worky build   # hugo + go build
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
            {
                ID:   "01",
                Name: "Chapter 1",
                Slug: "01-chapter1",
                Checks: []worky.Check{
                    {Description: "Service is healthy", Run: myServiceCheck},
                },
            },
        },
    }).Run()
}
```

### 3. Configure Hugo (`hugo.toml`)

```toml
baseURL    = "/"
title      = "My Workshop"
contentDir = "docs"
publishDir = "site"

[module]
  [[module.imports]]
    path = "github.com/davideimola/worky"
  [[module.imports]]
    path = "github.com/geekdocs/geekdoc"
```

worky provides the Hugo layouts, shortcodes, CSS and JS. [geekdoc](https://geekdoc.de) is the recommended theme (required for sidebar status icons).

### 4. Write content

Create markdown files in `docs/`. Use the `details` shortcode for hints and solutions:

```markdown
{{< details "Hint" >}}
Try running `kubectl get pods -A` to see what's running.
{{< /details >}}

{{< details "Solution" >}}
```yaml
apiVersion: v1
kind: Pod
...
```
{{< /details >}}
```

### 5. Build and distribute

```sh
# Build the site
hugo

# Cross-compile the binary
GOOS=linux  GOARCH=amd64 go build -o bin/my-workshop-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -o bin/my-workshop-darwin-arm64 .
```

Distribute binaries via GitHub Releases. The site is embedded — no separate web server needed.

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

## Hugo assets

worky provides these Hugo assets (overlaid on top of your theme):

| Path | Purpose |
|------|---------|
| `layouts/partials/head/custom.html` | Injects `workshop.css` and `workshop-progress.js` |
| `layouts/shortcodes/details.html` | `{{< details >}}` collapsible hint/solution |
| `static/workshop.css` | Styles for `.ws-details` and `.ws-status-icon` |
| `static/workshop-progress.js` | Polls `/api/progress`, updates sidebar icons, shows completion banner |

> `workshop-progress.js` targets `a.gdoc-nav__entry` selectors from geekdoc. If you use a different theme, you may need to override `workshop-progress.js` with your own version that matches your theme's sidebar selectors.

---

## Progress storage

Progress is stored in `~/<HomeDir>/progress.json`. The first chapter is always unlocked. Running `check` for a chapter marks it complete and unlocks the next one.

---

## Hugo module setup

Initialize the Hugo module in your project once:

```sh
hugo mod init github.com/acme/my-workshop
hugo mod get github.com/davideimola/worky
hugo mod get github.com/geekdocs/geekdoc
```
