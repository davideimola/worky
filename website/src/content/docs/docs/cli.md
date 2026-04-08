---
title: CLI
description: worky CLI reference — init, new chapter, build, report.
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

### Flags

| Flag | Description |
|------|-------------|
| `--no-install-deps` | Skip automatic `go mod tidy` |
| `--module <path>` | Go module path (default: `github.com/example/<slug>`) |
| `-y, --yes` | Accept all defaults without prompting |

### Created files

```
kubernetes-workshop/
├── main.go        # Workshop entrypoint with embed + one example chapter
├── site/
│   ├── index.html           # Home page
│   ├── workshop.css
│   ├── workshop-progress.js
│   └── 00-setup/
│       └── index.md         # Placeholder chapter content
├── Makefile       # make serve, make build, make release
├── .gitignore
├── README.md
├── go.mod
└── .github/workflows/release.yml
```

After scaffolding:

```sh
cd kubernetes-workshop
go mod tidy
make serve
```

---

## `worky new chapter <id> <name>`

Adds a new chapter to an existing workshop.

```sh
worky new chapter 02 "Deploying Services"
# → creates site/02-deploying-services/index.md
# → prints Go snippet to add to main.go
```

**Example output:**

```
Chapter created: site/02-deploying-services/index.md

Add the following to your main.go Chapters slice:

    worky.Chapter{
        ID:   "02",
        Name: "Deploying Services",
        Slug: "02-deploying-services",
        Checks: []worky.Check{
            // TODO: add checks here
        },
    },
```

---

## `worky build`

Builds the workshop binary.

```sh
worky build
# equivalent to: go build -o bin/my-workshop .
```

For cross-compilation, use Go's standard env vars:

```sh
GOOS=linux  GOARCH=amd64 go build -o bin/my-workshop-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -o bin/my-workshop-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o bin/my-workshop-windows-amd64.exe .
```

Distribute the binaries via GitHub Releases. The entire site is embedded in the binary — no additional files needed.

---

## `worky report`

Opens a pre-filled GitHub issue to report a bug in worky itself.

```sh
worky report
```

The command asks four questions interactively:

1. Describe the problem
2. Steps to reproduce
3. Expected behavior
4. Actual behavior

Once you answer the questions, it automatically collects your environment details (worky version, OS, architecture, Go version) and opens your browser with a pre-filled GitHub issue. If no browser is available (e.g. headless environments), the URL is printed to the terminal instead.

---

For the commands available inside each workshop binary (`serve`, `check`, `status`, etc.) see [Runtime commands](/runtime/).
