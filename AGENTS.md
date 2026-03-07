# AGENTS.md — worky

Guidelines for AI agents working on this codebase.

## Module overview

worky is a Go library (`github.com/davideimola/worky`) that workshop creators import to build self-contained workshop CLIs with an embedded Hugo site, chapter locking, and progress tracking.

It is simultaneously a **Go module** (imported in `main.go`) and a **Hugo module** (imported in `hugo.toml`) so that a single dependency covers both the Go runtime and the Hugo theme overlay.

## File layout

```
worky.go       — Config, Workshop struct, New(), Run(), all cobra commands
chapter.go     — Chapter type + lookup methods (chapterByID, chapterBySlug, nextChapter)
check.go       — Check, CheckResult, runChecks()
progress.go    — State, loadProgress/saveProgress/resetProgress, IsUnlocked/IsCompleted/Complete
server.go      — HTTP handler, /api/progress endpoint, chapter locking middleware
go.mod         — module github.com/davideimola/worky, requires cobra
go.sum         — pinned dependency hashes

layouts/partials/head/custom.html  — Hugo partial: injects workshop.css + workshop-progress.js
layouts/shortcodes/details.html    — Hugo shortcode: collapsible hint/solution block
static/workshop.css                — Styles for .ws-details and .ws-status-icon
static/workshop-progress.js        — Polls /api/progress, updates sidebar, shows completion banner
```

## Key design rules

- **No hardcoded content.** Chapters, checks, workshop name, home directory, and port all come from `Config`. There are no defaults for chapters or checks.
- **Workshop struct is the receiver.** All methods that need config or chapter/progress access are on `*Workshop`. There are no package-level variables.
- **Single dependency.** Only `github.com/spf13/cobra` is required. Everything else uses stdlib.
- **SiteFS layout.** The embedded FS passed as `Config.SiteFS` must contain a `site/` subdirectory (i.e., `//go:embed all:site` in the creator's `main.go`). The server calls `fs.Sub(siteFS, "site")` internally.
- **HomeDir defaults to `.worky`.** If `Config.HomeDir` is empty, `New()` sets it to `.worky`. Progress and PID/log files live under `~/<HomeDir>/`.

## Public API surface

Only these identifiers are part of the public API:

```go
// Types
Config
Workshop
Chapter
Check
CheckResult
State

// Functions
New(cfg Config) *Workshop

// Methods
(*Workshop).Run()

// State methods
(*State).IsUnlocked(id string) bool
(*State).IsCompleted(id string) bool
(*State).Complete(id, nextID string)
```

Everything else (chapter lookups, server handler, progress I/O, command builders) is unexported and internal.

## Cobra commands

All cobra commands are built as unexported methods (`serveCmd`, `checkCmd`, etc.) and registered in `Run()`. Do not add top-level `init()` or global flags.

## Hugo assets

The `layouts/` and `static/` directories are Hugo module assets overlaid on the creator's theme. They must not import any Go packages or reference Go code. They are pure HTML/JS/CSS.

`workshop-progress.js` uses `a.gdoc-nav__entry` selectors specific to the geekdoc theme. If selectors need to change for a different theme, that change belongs in the creator's project, not in worky.

## Verification

After any Go change:

```sh
cd worky
go build ./...
go vet ./...
```

Both must exit 0 with no output before a change is considered complete.

## CLI tool (`cmd/worky/`)

A standalone CLI installed via `go install github.com/davideimola/worky/cmd/worky@latest`.

```
cmd/worky/
  main.go                    — entrypoint
  cmd/
    root.go                  — NewRootCmd()
    init.go                  — worky init <name>
    new.go                   — worky new (parent group)
    chapter.go               — worky new chapter <id> <name>
    build.go                 — worky build
    scaffold.go              — shared scaffold() helper
  templates/
    embed.go                 — //go:embed all:files
    files/
      project/               — project scaffold templates
      chapter/               — chapter scaffold templates
```

The CLI does NOT import the worky library — it only scaffolds text files.
Templates use `<% %>` delimiters (not `{{ }}`) to avoid conflicts with Hugo shortcode syntax.

## What does NOT belong here

- Kubernetes, Docker, or any other tool-specific check helpers — those belong in the creator's project
- Example workshops or demo content
