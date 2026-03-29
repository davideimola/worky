# AGENTS.md — worky

Guidelines for AI agents working on this codebase.

## Module overview

worky is a Go library (`github.com/davideimola/worky`) that workshop creators import to build self-contained workshop CLIs with an embedded site, chapter locking, and progress tracking.

A workshop embeds a `site/` directory (plain HTML and Markdown files) into its binary. The worky server serves those files directly — no build step, no external tools required.

## File layout

```
worky.go       — Config, Workshop struct, New(), Run(), all cobra commands
chapter.go     — Chapter type + lookup methods (chapterByID, chapterBySlug, nextChapter)
check.go       — Check, CheckResult, runChecks()
progress.go    — State, loadProgress/saveProgress/resetProgress, IsUnlocked/IsCompleted/Complete
server.go      — HTTP handler, SSE hub, serveMD(), serveBuiltinHome(), /api/* endpoints
go.mod         — module github.com/davideimola/worky, requires cobra + goldmark
go.sum         — pinned dependency hashes

demo/          — fully working example workshop (HTML/MD site, no build step)
  main.go      — wires 3 chapters with checks
  site/        — pre-built HTML + MD pages + workshop.css + workshop-progress.js
```

## Key design rules

- **No hardcoded content.** Chapters, checks, workshop name, home directory, and port all come from `Config`. There are no defaults for chapters or checks.
- **Workshop struct is the receiver.** All methods that need config or chapter/progress access are on `*Workshop`. There are no package-level variables.
- **Single dependency.** Only `github.com/spf13/cobra` and `github.com/yuin/goldmark` are required. Everything else uses stdlib.
- **SiteFS layout.** The embedded FS passed as `Config.SiteFS` must contain a `site/` subdirectory (i.e., `//go:embed all:site` in the creator's `main.go`). The server calls `fs.Sub(siteFS, "site")` internally.
- **HomeDir defaults to `.worky`.** If `Config.HomeDir` is empty, `New()` sets it to `.worky`. Progress and PID/log files live under `~/<HomeDir>/`.
- **Markdown rendering.** `serveMD()` looks for `site/{slug}/index.md` and renders it as HTML using goldmark, with CSS inlined in the template. HTML files are served as-is by the file server.
- **Builtin home.** When `SiteFS` is set but `site/index.html` is missing, or when `SiteFS` is nil, `serveBuiltinHome()` renders a dynamic chapter-status page.

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

## Verification

After cloning, run once to configure git hooks:

```sh
make setup
```

After any Go change:

```sh
go build ./...
go vet ./...
go test ./...
```

All must exit 0 with no output before a change is considered complete.

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
      project/               — project scaffold templates (main.go, Makefile, site/, etc.)
      chapter/               — chapter scaffold template (site/<id>-<slug>/index.md)
```

Templates use `<% %>` delimiters (not `{{ }}`) to avoid conflicts with Go template syntax in generated code.

`worky init` always generates a `site/` directory with:
- `index.html` — home page
- `workshop.css` + `workshop-progress.js` — assets
- `00-setup/index.md` — placeholder for the first chapter

`worky new chapter` creates `site/<id>-<slug>/index.md` and prints the Go snippet to add to `main.go`.

## Keeping demo and documentation in sync

**Whenever a change affects how workshops are structured, served, or scaffolded:**

1. **Update `demo/`** — the demo must always reflect the current expected workshop structure. If `site/` layout changes, chapter rendering changes, or the `main.go` pattern changes, update `demo/` accordingly.
2. **Update `website/content/docs/`** — all documentation pages must stay accurate. In particular:
   - `getting-started.md` — scaffold steps, file structure diagrams, setup commands
   - `cli.md` — `worky init` flags, `worky new chapter` output, `worky build` behavior
   - `configuration.md` — `Config` struct fields and `SiteFS` usage
   - `chapters.md` — chapter struct, slug convention, how to add chapters
3. **Update `README.md`** — the root README is the first thing creators read. Keep the quickstart, file structure, and CLI examples current.
4. **Update `AGENTS.md`** (this file) — keep the file layout, design rules, and CLI tool sections accurate when the architecture changes.

If a PR touches any of `worky.go`, `server.go`, `cmd/worky/cmd/`, or `cmd/worky/templates/`, verify that all four areas above (demo, website docs, README, AGENTS.md) are still consistent.

## What does NOT belong here

- Kubernetes, Docker, or any other tool-specific check helpers — those belong in the creator's project
- Example workshops beyond `demo/` — one reference example is enough
