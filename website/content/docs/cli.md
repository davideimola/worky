---
title: CLI
weight: 5
---

worky ships a standalone CLI for scaffolding workshops. Install it once, use it in any project.

```sh
go install github.com/davideimola/worky/cmd/worky@latest
```

---

## `worky init <name>`

Scaffolds a complete new workshop project in a new directory.

```sh
worky init "Kubernetes Workshop"
# → creates kubernetes-workshop/
```

**Created files:**

```
kubernetes-workshop/
├── main.go        # Workshop entrypoint with one example chapter
├── docs/
│   └── 00-setup/
│       └── _index.md
├── hugo.toml      # Hugo config with worky + geekdoc module imports
├── Makefile       # make serve, make build, make release
└── go.mod
```

After scaffolding, initialize the Hugo modules:

```sh
cd kubernetes-workshop
hugo mod init github.com/yourorg/kubernetes-workshop
hugo mod get github.com/davideimola/worky
hugo mod get github.com/geekdocs/geekdoc
go mod tidy
```

---

## `worky new chapter <id> <name>`

Adds a new chapter to an existing workshop.

```sh
worky new chapter 02 "Deploying Services"
# → creates docs/02-deploying-services/_index.md
# → prints Go snippet to add to main.go
```

**Example output:**

```
Created: docs/02-deploying-services/_index.md

Add this to your Chapters slice in main.go:

    {
        ID:   "02",
        Name: "Deploying Services",
        Slug: "02-deploying-services",
        Checks: []worky.Check{
            // TODO: add checks
        },
    },
```

---

## `worky build`

Builds the Hugo site and compiles the Go binary.

```sh
worky build
```

Equivalent to:

```sh
hugo
go build -o bin/my-workshop .
```

For cross-compilation, use Go's standard env vars:

```sh
GOOS=linux  GOARCH=amd64 go build -o bin/my-workshop-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -o bin/my-workshop-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o bin/my-workshop-windows-amd64.exe .
```

Distribute the binaries via GitHub Releases. The entire Hugo site is embedded — no additional files needed.

---

## Workshop runtime commands

These are the commands built into every workshop binary (not the `worky` CLI):

| Command | Description |
|---------|-------------|
| `serve [--port N] [--open] [--detach] [--preview]` | Start the local web server |
| `check [chapter-id]` | Run checks for a chapter and unlock the next one |
| `status` | Show all chapters with lock/unlock/complete icons |
| `reset` | Clear all progress |
| `unlock <chapter-id>` | Manually unlock a chapter (facilitator use) |
| `stop` | Kill the background server |
| `logs [-f]` | Show (and optionally follow) server logs |

### `serve` flags

| Flag | Default | Description |
|------|---------|-------------|
| `--port N` | Config.Port | Port to listen on |
| `--open` | false | Open browser automatically after start |
| `--detach` | false | Run server in background (writes PID + log file) |
| `--preview` | false | Disable chapter locking for content review |

### `check`

If `chapter-id` is omitted, auto-detects the first unlocked, incomplete chapter.

```sh
./my-workshop check        # auto-detect
./my-workshop check 01     # explicit
```

### `logs`

```sh
./my-workshop logs         # print log file
./my-workshop logs -f      # follow (tail -f style)
```
