---
title: Chapters
description: Chapter struct, locking logic, and how to add new chapters.
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

The URL path for this chapter. Must match the directory name under `site/`.

```
Slug: "01-hello-world"  →  http://localhost:8080/01-hello-world/
                        →  site/01-hello-world/index.md  (or index.html)
```

The `worky new chapter` CLI command sets this automatically.

### Checks

List of `worky.Check` values that must all pass for the chapter to be marked complete. See [Checks](/reference/checks/built-in/) for the full reference.

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

## Adding a new chapter

Use the CLI:

```sh
worky new chapter 02 "Advanced Topics"
```

This:
1. Creates `site/02-advanced-topics/index.md` with placeholder content
2. Prints the Go snippet to add to your `Chapters` slice in `main.go`

Or manually — add the entry to `Chapters` and create `site/<id>-<slug>/index.md`.

---

For the runtime commands that operate on chapters (`check`, `status`, `unlock`, `reset`) see [Runtime commands](/runtime/).
