---
title: Documentation
weight: 1
---

worky is a Go library for building self-contained interactive workshop tools. A workshop creator imports worky, defines chapters with validation checks, embeds a site directory, and distributes a single binary that participants download and run.

## What worky provides

- **A CLI runtime** — `serve`, `check`, `status`, `reset`, `stop`, `logs` commands built with Cobra
- **Chapter locking** — chapters unlock progressively as checks pass
- **Progress tracking** — JSON state file + SSE push to update the browser in real-time
- **Markdown rendering** — `.md` chapter files rendered server-side with inline CSS (no build step)
- **Pre-built checks** — `checks` sub-package with common validation functions

## Navigation

{{< cards >}}
  {{< card link="getting-started" title="Getting Started" subtitle="Install, scaffold, and run your first workshop" icon="play" >}}
  {{< card link="reference" title="Reference" subtitle="Config, Chapter, and Check struct — full Go API" icon="book-open" >}}
  {{< card link="cli" title="CLI" subtitle="worky init, worky new chapter, worky build" icon="terminal" >}}
  {{< card link="runtime" title="Runtime commands" subtitle="serve, check, status, reset, unlock, logs" icon="terminal" >}}
  {{< card link="customization" title="Customization" subtitle="SiteFS layout, CSS classes, and Markdown rendering" icon="color-swatch" >}}
  {{< card link="troubleshooting" title="Troubleshooting" subtitle="Common errors and how to fix them" icon="exclamation-circle" >}}
{{< /cards >}}
