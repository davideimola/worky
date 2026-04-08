---
title: Customization
description: Customize your workshop site layout, CSS, and Markdown rendering.
---

## SiteFS layout

A workshop's embedded site is a Go `fs.FS` — typically an `embed.FS` — that contains a `site/` subdirectory. worky calls `fs.Sub(siteFS, "site")` internally, so every path below is relative to `site/`.

```
myworkshop/
└── site/
    ├── index.html              ← custom homepage (optional)
    ├── workshop.css            ← additional styles (loaded automatically)
    ├── workshop-progress.js    ← live chapter status (do not replace)
    ├── 00-getting-started/
    │   └── index.html          ← chapter page (HTML)
    └── 01-next-chapter/
        └── index.md            ← chapter page (Markdown)
```

### Embedding in main.go

```go
//go:embed site
var siteFS embed.FS

func main() {
    worky.Run(worky.Config{
        SiteFS: siteFS,
        // ...
    })
}
```

---

## File reference

| File | Required | Purpose |
|------|----------|---------|
| `site/index.html` | No | Custom homepage. If absent, worky renders a built-in chapter-status page. |
| `site/workshop.css` | No | Additional styles. Loaded automatically by the server on every page. |
| `site/workshop-progress.js` | Yes (do not replace) | Handles live SSE updates for chapter status icons. |
| `site/{slug}/index.html` | One per chapter | Chapter page as raw HTML. |
| `site/{slug}/index.md` | One per chapter | Chapter page as Markdown (rendered server-side via Goldmark). |

Each chapter directory name must match the chapter's slug. The slug is derived from the chapter `ID` and `Name` fields — see the [Chapters](/reference/chapters/) page for the convention.

---

## CSS classes

The built-in stylesheet exposes utility classes you can use in custom HTML pages.

### Expandable hint boxes

```html
<details class="ws-details ws-details--tip">
  <summary class="ws-details__summary">
    <span class="ws-details__icon">💡</span>
    <span class="ws-details__title">Hint</span>
    <span class="ws-details__arrow">›</span>
  </summary>
  <div class="ws-details__body">Content of the hint.</div>
</details>
```

Modifier classes:

| Class | Style |
|-------|-------|
| `ws-details--tip` | Blue — hints and suggestions |
| `ws-details--info` | Grey — neutral information |

### Status icons

`.ws-status-icon` is applied to chapter status icons in the sidebar. It is managed by `workshop-progress.js` — do not manipulate it directly in CSS or JS.

---

## Markdown rendering

Chapter files ending in `.md` are rendered server-side with [Goldmark](https://github.com/yuin/goldmark). The rendered HTML is wrapped in a minimal layout with `workshop.css` applied.

Supported syntax:
- CommonMark
- Inline HTML (raw HTML blocks and inline elements pass through unchanged)

Markdown pages are a good choice for text-heavy chapters that do not need a custom layout.

---

## Worky init scaffold

Running [`worky init <name>`](/cli/) generates a ready-to-use `site/` directory with `index.html`, `workshop.css`, and `workshop-progress.js` pre-populated. Use those files as the starting point for customization rather than writing from scratch.
