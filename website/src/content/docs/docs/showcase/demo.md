---
title: Worky Demo
description: The demo workshop included in the worky repository — three chapters, file and environment checks.
---

The demo workshop is included in the [worky repository](https://github.com/davideimola/worky/tree/main/demo) under `demo/`. It's the fastest way to see worky in action.

## Run it now

Just Go:

```sh
git clone https://github.com/davideimola/worky
cd worky/demo
go mod tidy
go run . serve --open
```

The workshop server starts on `http://localhost:8080` and opens in your browser.

## Chapters

### Chapter 00 — Getting Started

Verifies your environment is ready.

**Checks:**
- `done.txt` exists in the current directory
- `WORKSHOP_USER` environment variable is set

```sh
touch done.txt
export WORKSHOP_USER=yourname
go run . check
```

---

### Chapter 01 — Hello Worky

A simple file creation and content check.

**Checks:**
- `hello.txt` exists
- `hello.txt` contains the text `"Hello, Worky!"`

```sh
echo "Hello, Worky!" > hello.txt
go run . check
```

---

### Chapter 02 — Finishing Up

The final chapter.

**Checks:**
- `complete.md` exists
- `complete.md` contains `"# Workshop Complete"`

```sh
echo "# Workshop Complete" > complete.md
go run . check
```

---

## What you'll see

As each `check` passes, the sidebar icons in your browser update live — no refresh needed. The status command shows your overall progress:

```sh
go run . status

  ✅  Chapter 00: Getting Started
  ✅  Chapter 01: Hello Worky
  🔓  Chapter 02: Finishing Up
```

When all chapters are complete, a banner appears in the browser.

## Reset and replay

```sh
go run . reset
go run . status

  🔓  Chapter 00: Getting Started
  🔒  Chapter 01: Hello Worky
  🔒  Chapter 02: Finishing Up
```

## Source code

The demo `main.go` shows a minimal but complete worky setup using the `checks` sub-package:

```go
worky.New(worky.Config{
    Name:    "Worky Demo",
    HomeDir: ".worky-demo",
    Port:    8080,
    SiteFS:  site,
    Chapters: []worky.Chapter{
        {
            ID:   "00",
            Name: "Getting Started",
            Slug: "00-getting-started",
            Checks: []worky.Check{
                {
                    Description: "done.txt exists in current directory",
                    Run:         checks.FileExists("done.txt"),
                },
                {
                    Description: "WORKSHOP_USER environment variable is set",
                    Run:         checks.EnvVarSet("WORKSHOP_USER"),
                },
            },
        },
        // ...
    },
}).Run()
```
