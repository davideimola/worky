---
title: Chapters
weight: 3
---

## Chapter struct

```go
type Chapter struct {
    ID     string  // "00", "01", etc.
    Name   string  // Display name shown in CLI output
    Slug   string  // URL path, e.g. "00-setup"
    Checks []Check // Validation steps for this chapter
}
```

### ID

Short string identifier, conventionally zero-padded: `"00"`, `"01"`, `"02"`. Used by `check`, `status`, and `unlock` commands. Must be unique across chapters.

### Name

Human-readable name shown in `status` output and CLI messages. Example: `"Getting Started"`.

### Slug

The URL path for this chapter in the embedded Hugo site. Must match the directory name in your `docs/` folder.

```
Slug: "01-hello-world"  →  http://localhost:8080/01-hello-world/
```

The `worky new chapter` CLI command sets this automatically.

### Checks

List of `worky.Check` values that must all pass for the chapter to be marked complete. See [Checks](../checks) for the full reference.

---

## Locking logic

- The **first chapter** is always unlocked when progress is initialized (or after `reset`).
- Running `check` for a chapter validates all its checks. If all pass, the chapter is marked **complete** and the next chapter is **unlocked**.
- Locked chapters return HTTP 403 from the embedded server — participants cannot browse ahead.

State transitions:

```
locked → unlocked → complete
```

Only forward transitions are possible via normal flow. `unlock` can manually unlock any chapter (facilitator use).

---

## CLI commands affecting chapters

### `check [chapter-id]`

Runs all checks for a chapter. If omitted, auto-detects the first unlocked, incomplete chapter.

```sh
./my-workshop check      # auto-detect
./my-workshop check 01   # explicit
```

### `status`

Shows all chapters with icons:

```
  ✅  Chapter 00: Getting Started
  🔓  Chapter 01: Hello World
  🔒  Chapter 02: Finishing Up
```

### `unlock <chapter-id>`

Manually unlocks a chapter without running checks. Intended for facilitators who need to unblock a participant.

```sh
./my-workshop unlock 02
```

### `reset`

Clears all progress. The first chapter is unlocked, all others return to locked state.

```sh
./my-workshop reset
```

---

## Adding a new chapter

Use the CLI:

```sh
worky new chapter 02 "Advanced Topics"
```

This:
1. Creates `docs/02-advanced-topics/_index.md` with front matter
2. Prints the Go snippet to add to your `Chapters` slice in `main.go`

Or manually — add the entry to `Chapters` and create the matching docs directory.
