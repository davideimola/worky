---
title: Getting Started
weight: 1
---

## Prerequisites

- Go 1.21+
- Hugo extended v0.120+ (for module support)

## Option A: CLI scaffold (recommended)

### 1. Install the CLI

```sh
go install github.com/davideimola/worky/cmd/worky@latest
```

### 2. Scaffold a new workshop

```sh
worky init "My Workshop"
cd my-workshop
```

This creates:

```
my-workshop/
├── main.go        # Workshop entrypoint with one example chapter
├── docs/
│   └── 00-setup/
│       └── _index.md
├── hugo.toml      # Hugo config importing worky + geekdoc modules
├── Makefile       # make serve, make build, make release
└── go.mod
```

### 3. Set up dependencies

```sh
hugo mod init github.com/yourorg/my-workshop
hugo mod get github.com/davideimola/worky
hugo mod get github.com/geekdocs/geekdoc
go mod tidy
```

### 4. Start developing

```sh
make serve
# → hugo server (watch mode) + go run . serve --open
```

---

## Option B: Manual setup

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

### 4. Write content

Create `docs/00-setup/_index.md`:

```markdown
---
title: Setup
---

Welcome to the workshop! Let's make sure your environment is ready.

## Steps

1. Start Docker Desktop
2. Run `./my-workshop check` to verify and unlock the next chapter
```

Use the `details` shortcode for hints:

```markdown
{{</* details "Hint" */>}}
Try running `docker info` to see if Docker is responding.
{{</* /details */>}}
```

### 5. Build the docs and run

```sh
# Build Hugo site into site/
hugo

# Run the workshop server
go run . serve --open
```

---

## Participant experience

Once distributed, participants:

1. Download the binary (or `git clone` + `make setup`)
2. Run `./my-workshop serve --open`
3. Follow along in the browser
4. Run `./my-workshop check` after each chapter to unlock the next one

```sh
./my-workshop status   # show all chapters with lock/unlock/complete icons
./my-workshop check    # auto-detect current chapter and validate
./my-workshop reset    # clear all progress and start over
```

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
