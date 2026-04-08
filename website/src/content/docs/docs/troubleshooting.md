---
title: Troubleshooting
description: Common errors and how to fix them.
---

## Server won't start

### Port 8080 already in use

```
Error: listen tcp :8080: bind: address already in use
```

Use a different port via flag or environment variable:

```sh
./myworkshop serve --port 9090
# or
PORT=9090 ./myworkshop serve
```

### PID file already present

If a previous server process crashed without cleaning up, `~/.worky/server.pid` may still exist:

```sh
worky stop          # attempts a clean shutdown
# or remove manually
rm ~/.worky/server.pid
```

---

## Checks not passing

### FileExists / DirExists

```
✗ file "done.txt" does not exist
```

Verify the path is relative to the working directory where the workshop binary runs, not the binary location.

### EnvVar checks

```
✗ environment variable KUBECONFIG is not set
```

The variable must be set in the same shell session that runs `./myworkshop check`. Variables exported in another terminal are not visible.

### Command checks

```
✗ command "kubectl" not found
```

The command must be on `PATH` in the current shell. Verify with:

```sh
which kubectl
```

### Network / HTTP checks

```
✗ dial tcp localhost:8080: connect: connection refused
```

The service is not yet running on that port. Start the service before running `check`, or add `Retries` and `RetryDelay` to the check so worky waits for it.

### Check timeout

```
✗ timed out after 5s
```

The check's `Timeout` duration was exceeded. Either increase `Timeout` in the check definition or investigate why the external process is slow to respond.

### No checks defined

```
no checks defined for chapter 02-deploy
```

The chapter struct has an empty or nil `Checks` slice. Add at least one `worky.Check` entry to the chapter, or remove the chapter if it is informational only.

---

## Progress and state issues

### Chapter stuck / unexpectedly locked

Run `./myworkshop status` to see which chapters are complete and which are unlocked. If a chapter that should be unlocked appears locked, check that the preceding chapter's checks all pass:

```sh
./myworkshop check 01-setup
```

### Unknown chapter ID

```
unknown chapter: deploy
```

The ID passed to `check` or `status` must match exactly the `ID` field in `Config.Chapters`. Use the full ID as defined (e.g. `01-setup`, not `setup`).

### Corrupted progress file

If `~/.worky/progress.json` is corrupted (e.g. truncated after a crash), worky will fail to load state. Reset progress to start fresh:

```sh
./myworkshop reset
```

Or delete the file manually:

```sh
rm ~/.worky/progress.json
```

---

## Detached mode

### Viewing logs

```sh
./myworkshop logs
```

Logs are written to `~/.worky/server.log` when the server is started with `--detach`.

### Stopping the server

```sh
./myworkshop stop
```

If `stop` fails or the process is unresponsive:

```sh
kill $(cat ~/.worky/server.pid)
rm ~/.worky/server.pid
```

---

## Further reading

- [Built-in checks](/docs/reference/checks/built-in/) — full list of built-in check functions and their error messages
- [Patterns](/docs/reference/checks/patterns/) — timeout, retry, and custom check patterns
