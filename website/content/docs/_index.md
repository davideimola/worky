---
title: Documentation
weight: 1
---

worky is a Go library for building self-contained interactive workshop tools. A workshop creator imports worky, defines chapters with validation checks, embeds a Hugo-built documentation site, and distributes a single binary that participants download and run.

## What worky provides

- **A CLI runtime** — `serve`, `check`, `status`, `reset`, `stop`, `logs` commands built with Cobra
- **Chapter locking** — chapters unlock progressively as checks pass
- **Progress tracking** — JSON state file + SSE push to update the browser in real-time
- **Hugo theme overlay** — layouts, shortcodes, CSS and JS for the embedded docs site
- **Pre-built checks** — `checks` sub-package with common validation functions

## Navigation

{{< cards >}}
  {{< card link="getting-started" title="Getting Started" subtitle="Install, scaffold, and run your first workshop" icon="play" >}}
  {{< card link="configuration" title="Configuration" subtitle="Config struct reference — all fields explained" icon="cog" >}}
  {{< card link="chapters" title="Chapters" subtitle="Chapter struct, locking logic, status commands" icon="book-open" >}}
  {{< card link="checks" title="Checks" subtitle="Built-in check functions with signatures" icon="check-circle" >}}
  {{< card link="cli" title="CLI" subtitle="worky init, worky new chapter, worky build" icon="terminal" >}}
{{< /cards >}}
