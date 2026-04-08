---
title: Runtime commands
description: Commands built into every workshop binary — serve, check, status, reset, unlock, logs, stop.
---

Every workshop binary built with worky includes a set of runtime commands for both participants and facilitators.

## Command reference

| Command | Description |
|---------|-------------|
| `serve [--port N] [--open] [--detach] [--preview]` | Start the local web server |
| `check [chapter-id]` | Run checks for a chapter and unlock the next one |
| `status` | Show all chapters with lock/unlock/complete icons |
| `reset` | Clear all progress |
| `unlock <chapter-id>` | Manually unlock a chapter (facilitator use) |
| `stop` | Kill the background server |
| `logs [-f]` | Show (and optionally follow) server logs |

---

## `serve`

Starts the local web server and serves the embedded site.

| Flag | Default | Description |
|------|---------|-------------|
| `--port N` | `Config.Port` | Port to listen on |
| `--open` | false | Open browser automatically after start |
| `--detach` | false | Run server in background (writes PID + log file) |
| `--preview` | false | Disable chapter locking for content review |

```sh
./my-workshop serve --open
./my-workshop serve --port 9090 --detach
```

---

## `check [chapter-id]`

Runs all checks for a chapter. If `chapter-id` is omitted, auto-detects the first unlocked, incomplete chapter.

```sh
./my-workshop check        # auto-detect
./my-workshop check 01     # explicit
```

If all checks pass the chapter is marked complete and the next one is unlocked. See [Checks](/docs/reference/checks/built-in/) for the full list of built-in validation functions.

---

## `status`

Shows all chapters with icons:

```
  ✅  Chapter 00: Getting Started
  🔓  Chapter 01: Hello World
  🔒  Chapter 02: Finishing Up
```

---

## `unlock <chapter-id>`

Manually unlocks a chapter without running checks. Intended for facilitators who need to unblock a participant.

```sh
./my-workshop unlock 02
```

---

## `reset`

Clears all progress. The first chapter is unlocked, all others return to locked state.

```sh
./my-workshop reset
```

---

## `logs`

```sh
./my-workshop logs         # print log file
./my-workshop logs -f      # follow (tail -f style)
```

Logs are written to `~/<HomeDir>/server.log` when the server is started with `--detach`.

---

## `stop`

Stops a server started with `--detach`.

```sh
./my-workshop stop
```

If the process is unresponsive, see [Troubleshooting](/docs/troubleshooting/).
