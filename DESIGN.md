# Worky Design System

> Design by [Fabrizio Lamanna](https://www.linkedin.com/in/fabrizio-lamanna/).

This document is the single source of truth for Worky's visual identity.
Use it when implementing styled CLI output, TUI components, or any new UI surface.

---

## Colors

| Token | Hex | Usage |
|-------|-----|-------|
| `bg` | `#071122` | Page / terminal background |
| `surface` | `#0B1828` | Cards, panels, elevated elements |
| `teal` | `#00C2A0` | Primary accent — CTAs, success, active state, icons |
| `teal-deep` | `#00A086` | Hover state, icon fill in SVG assets |
| `teal-border` | `rgba(0,194,160,0.25)` | Borders on teal-accented elements |
| `border` | `rgba(255,255,255,0.07)` | Subtle dividers and card borders |
| `border-strong` | `rgba(255,255,255,0.12)` | More visible dividers |
| `text` | `rgba(255,255,255,0.88)` | Body text |
| `text-muted` | `rgba(255,255,255,0.45)` | Secondary text, placeholders |
| `warning` | `#FEBC2E` | In-progress, caution |
| `error` | `#FF5757` | Failures, destructive actions |

> Replaces the **Catppuccin Mocha** placeholder used before the brand was finalised.

### Catppuccin → Worky mapping

| Old (Catppuccin Mocha) | New (Worky) | Role |
|------------------------|-------------|------|
| `#cdd6f4` | `rgba(255,255,255,0.88)` | Primary text |
| `#a6e3a1` | `#00C2A0` | Success / accent |
| `#f38ba8` | `#FF5757` | Error |
| `#fab387` | `#FEBC2E` | Warning |
| `#313244` | `#0B1828` | Surface |
| `#1e1e2e` | `#071122` | Background |

---

## Typography

| Role | Font | Notes |
|------|------|-------|
| UI / prose | **Inter** | All body text, labels, prompts |
| Code / terminal | **JetBrains Mono** | Code blocks, command output, monospace UI |

In terminal environments (Lipgloss / Bubbletea) font rendering is handled by the user's terminal emulator. No font configuration is needed — use the color tokens above and rely on the system monospace font for code.

---

## Lipgloss palette (Go)

```go
// internal/theme/theme.go
package theme

import "github.com/charmbracelet/lipgloss"

var (
    Bg      = lipgloss.Color("#071122")
    Surface = lipgloss.Color("#0B1828")
    Teal    = lipgloss.Color("#00C2A0")
    TealDeep = lipgloss.Color("#00A086")
    Text    = lipgloss.Color("#E2E8F0") // closest to rgba(255,255,255,0.88) on dark bg
    Muted   = lipgloss.Color("#64748B")
    Warning = lipgloss.Color("#FEBC2E")
    Error   = lipgloss.Color("#FF5757")
    Border  = lipgloss.Color("#0F2236")  // approximate rgba(255,255,255,0.07) on #071122
)

// Ready-made styles
var (
    Success = lipgloss.NewStyle().Foreground(Teal).Bold(true)
    Failure = lipgloss.NewStyle().Foreground(Error).Bold(true)
    Warn    = lipgloss.NewStyle().Foreground(Warning).Bold(true)
    Info    = lipgloss.NewStyle().Foreground(Muted)
    Code    = lipgloss.NewStyle().Foreground(Teal)
    Header  = lipgloss.NewStyle().Foreground(Text).Bold(true)
)
```

---

## Icons / status indicators

| State | Symbol | Color |
|-------|--------|-------|
| Success / unlocked | `✓` | `Teal` (`#00C2A0`) |
| In progress / checking | `○` | `Warning` (`#FEBC2E`) |
| Locked / pending | `○` | `Muted` |
| Error / failed | `✗` | `Error` (`#FF5757`) |
| Bullet / list item | `•` | `Muted` |

---

## Assets

| File | Usage |
|------|-------|
| `website/public/logo.svg` | Full wordmark (pixel icon + "worky" text). Use in web/print. |
| `website/public/icon.svg` | Pixel icon only. Use as favicon, app icon, small contexts. |
| `.assets/img/worky-lightmode.svg` | Light-mode variant for external docs / GitHub. |
| `.assets/img/worky-darkmode.svg` | Dark-mode variant for external docs / GitHub. |

---

## Voice & tone

- **Terse and direct** — output only what the user needs to act on.
- **No filler** — avoid "Please wait…", "Done!", or emoji unless the context calls for it.
- **Consistent icons** — use the symbol set above; don't mix `✔`, `✓`, `☑` etc.
- **Credits** — "Powered by Worky" appears by default in styled output; suppressible via `Config.HideCredits`.
